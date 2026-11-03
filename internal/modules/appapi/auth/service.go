package auth

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"google.golang.org/api/idtoken"
)

// GoogleServiceVerifier は Google の ID トークンを検証し、許可されたサービスアカウントからの呼び出しかを確かめる。
type GoogleServiceVerifier struct {
	audience        string
	allowedAccounts []string
	validate        func(ctx context.Context, token, audience string) (*idtoken.Payload, error)
}

// NewGoogleServiceVerifier は audience（App API の URL）と、カンマ区切りの許可するサービスアカウントから Verifier を作る。
func NewGoogleServiceVerifier(audience, allowedAccounts string) (*GoogleServiceVerifier, error) {
	audience = strings.TrimSpace(audience)
	if audience == "" {
		return nil, errors.New("id token audience is empty")
	}
	var accounts []string
	for a := range strings.SplitSeq(allowedAccounts, ",") {
		if a = strings.TrimSpace(a); a != "" {
			accounts = append(accounts, a)
		}
	}
	if len(accounts) == 0 {
		return nil, errors.New("id token allowed service accounts is empty")
	}
	return &GoogleServiceVerifier{audience: audience, allowedAccounts: accounts, validate: idtoken.Validate}, nil
}

// Verify はトークンの署名・有効期限・audience を検証し、発行元が許可されたサービスアカウントかを確かめる。
func (v *GoogleServiceVerifier) Verify(ctx context.Context, token string) error {
	payload, err := v.validate(ctx, token, v.audience)
	if err != nil {
		return err
	}
	email, _ := payload.Claims["email"].(string)
	verified, _ := payload.Claims["email_verified"].(bool)
	if !verified {
		return fmt.Errorf("email is not verified: %q", email)
	}
	if !slices.Contains(v.allowedAccounts, email) {
		return fmt.Errorf("service account is not allowed: %q", email)
	}
	return nil
}
