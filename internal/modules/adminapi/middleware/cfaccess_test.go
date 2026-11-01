package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

const testAud = "test-aud"

func setup(t *testing.T) (*CFAccessVerifier, func(jwt.Claims) string) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	pub := jose.JSONWebKey{Key: &priv.PublicKey, KeyID: "kid1", Algorithm: "RS256", Use: "sig"}
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{pub}})
	}))
	t.Cleanup(srv.Close)

	v, err := NewCFAccessVerifier(srv.URL, testAud)
	if err != nil {
		t.Fatal(err)
	}
	v.client = srv.Client()

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: jose.JSONWebKey{Key: priv, KeyID: "kid1"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	sign := func(c jwt.Claims) string {
		s, err := jwt.Signed(signer).Claims(c).Serialize()
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	return v, sign
}

func TestCFAccessMiddleware(t *testing.T) {
	v, sign := setup(t)
	now := time.Now()
	valid := jwt.Claims{Issuer: v.issuer, Audience: jwt.Audience{testAud}, Expiry: jwt.NewNumericDate(now.Add(time.Hour)), IssuedAt: jwt.NewNumericDate(now)}

	expired := valid
	expired.Expiry = jwt.NewNumericDate(now.Add(-time.Hour))
	wrongAud := valid
	wrongAud.Audience = jwt.Audience{"other"}
	wrongIss := valid
	wrongIss.Issuer = "https://evil.example.com"
	noExp := valid
	noExp.Expiry = nil

	h := v.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))

	tests := []struct {
		name  string
		token string
		want  int
	}{
		{"valid", sign(valid), http.StatusOK},
		{"missing", "", http.StatusForbidden},
		{"garbage", "abc.def.ghi", http.StatusForbidden},
		{"expired", sign(expired), http.StatusForbidden},
		{"wrong audience", sign(wrongAud), http.StatusForbidden},
		{"wrong issuer", sign(wrongIss), http.StatusForbidden},
		{"no exp", sign(noExp), http.StatusForbidden},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/admin.v1.UserService/ListUsers", nil)
			if tt.token != "" {
				req.Header.Set(CFAccessJWTHeader, tt.token)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tt.want {
				t.Errorf("status = %d, want %d", rec.Code, tt.want)
			}
		})
	}
}

func TestNewCFAccessVerifierRequiresConfig(t *testing.T) {
	if _, err := NewCFAccessVerifier("", testAud); err == nil {
		t.Error("expected error for empty issuer")
	}
	if _, err := NewCFAccessVerifier("team.cloudflareaccess.com", testAud); err == nil {
		t.Error("expected error for issuer without https://")
	}
	if _, err := NewCFAccessVerifier("https://team.cloudflareaccess.com", " "); err == nil {
		t.Error("expected error for empty audience")
	}
	v, _ := NewCFAccessVerifier("https://team.cloudflareaccess.com/", testAud)
	if !strings.HasSuffix(v.certsURL, "team.cloudflareaccess.com/cdn-cgi/access/certs") || v.issuer != "https://team.cloudflareaccess.com" {
		t.Errorf("unexpected urls: %s %s", v.issuer, v.certsURL)
	}
}
