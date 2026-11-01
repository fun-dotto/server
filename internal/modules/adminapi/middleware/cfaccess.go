// Package middleware は admin-api 用の HTTP ミドルウェアを提供する。
package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

// CFAccessJWTHeader は Cloudflare Access が付与する JWT のヘッダー名。
const CFAccessJWTHeader = "Cf-Access-Jwt-Assertion"

const (
	jwksRefreshInterval = 10 * time.Minute
	jwksMinRefetch      = 30 * time.Second
	clockLeeway         = time.Minute
)

var allowedAlgs = []jose.SignatureAlgorithm{jose.RS256, jose.ES256}

// CFAccessVerifier は Cloudflare Access の JWT を検証する。
type CFAccessVerifier struct {
	issuer   string
	audience string
	certsURL string
	client   *http.Client
	now      func() time.Time

	mu        sync.Mutex
	keys      jose.JSONWebKeySet
	fetchedAt time.Time
}

// NewCFAccessVerifier は teamDomain (例: "example.cloudflareaccess.com") と
// Access アプリケーションの AUD タグから Verifier を作る。
func NewCFAccessVerifier(teamDomain, audience string) (*CFAccessVerifier, error) {
	teamDomain = strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(teamDomain), "https://"), "/")
	audience = strings.TrimSpace(audience)
	if teamDomain == "" {
		return nil, errors.New("cloudflare access team domain is empty")
	}
	if audience == "" {
		return nil, errors.New("cloudflare access audience is empty")
	}
	issuer := "https://" + teamDomain
	return &CFAccessVerifier{
		issuer:   issuer,
		audience: audience,
		certsURL: issuer + "/cdn-cgi/access/certs",
		client:   &http.Client{Timeout: 10 * time.Second},
		now:      time.Now,
	}, nil
}

// Middleware は Cloudflare Access を経由していない (有効な JWT を持たない) リクエストを 403 で拒否する。
func (v *CFAccessVerifier) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(CFAccessJWTHeader)
		if token == "" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if err := v.Verify(r.Context(), token); err != nil {
			slog.WarnContext(r.Context(), "cloudflare access jwt verification failed", "error", err)
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Verify は JWT の署名・issuer・audience・有効期限を検証する。
func (v *CFAccessVerifier) Verify(ctx context.Context, token string) error {
	parsed, err := jwt.ParseSigned(token, allowedAlgs)
	if err != nil {
		return fmt.Errorf("parse jwt: %w", err)
	}
	if len(parsed.Headers) != 1 {
		return errors.New("unexpected number of signatures")
	}
	kid := parsed.Headers[0].KeyID

	key, err := v.key(ctx, kid)
	if err != nil {
		return err
	}

	var claims jwt.Claims
	if err := parsed.Claims(key, &claims); err != nil {
		return fmt.Errorf("verify signature: %w", err)
	}
	if err := claims.ValidateWithLeeway(jwt.Expected{
		Issuer:      v.issuer,
		AnyAudience: jwt.Audience{v.audience},
		Time:        v.now(),
	}, clockLeeway); err != nil {
		return fmt.Errorf("validate claims: %w", err)
	}
	if claims.Expiry == nil {
		return errors.New("jwt has no exp")
	}
	return nil
}

func (v *CFAccessVerifier) key(ctx context.Context, kid string) (*jose.JSONWebKey, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	now := v.now()
	stale := now.Sub(v.fetchedAt) > jwksRefreshInterval
	if k := v.lookup(kid); k != nil && !stale {
		return k, nil
	}
	// 未知の kid (鍵ローテーション) または期限切れの場合は再取得する。ただし連続取得は抑制する。
	if stale || now.Sub(v.fetchedAt) > jwksMinRefetch {
		if err := v.fetch(ctx); err != nil {
			if k := v.lookup(kid); k != nil {
				return k, nil
			}
			return nil, err
		}
	}
	if k := v.lookup(kid); k != nil {
		return k, nil
	}
	return nil, fmt.Errorf("unknown key id %q", kid)
}

func (v *CFAccessVerifier) lookup(kid string) *jose.JSONWebKey {
	keys := v.keys.Key(kid)
	idx := slices.IndexFunc(keys, func(k jose.JSONWebKey) bool { return k.Use == "" || k.Use == "sig" })
	if idx < 0 {
		return nil
	}
	return &keys[idx]
}

func (v *CFAccessVerifier) fetch(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.certsURL, nil)
	if err != nil {
		return err
	}
	res, err := v.client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch cloudflare access certs: %w", err)
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch cloudflare access certs: status %d", res.StatusCode)
	}
	var set jose.JSONWebKeySet
	if err := json.NewDecoder(res.Body).Decode(&set); err != nil {
		return fmt.Errorf("decode cloudflare access certs: %w", err)
	}
	v.keys = set
	v.fetchedAt = v.now()
	return nil
}
