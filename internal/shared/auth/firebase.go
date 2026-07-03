// Package auth はモジュラーモノリス全体で共有する認証・認可ユーティリティ。
// Firebase Authentication / App Check の検証と Gin ハンドラ向けの
// gate ヘルパを提供する。
package auth

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/appcheck"
	"firebase.google.com/go/v4/auth"
)

// Clients は Extract / RequireUserUnlessAllowed が依存する Firebase クライアントの束。
// Auth / AppCheck の両方が必須。NewClients が初期化に失敗した場合は呼び出し側で
// 起動を中断する想定 (fail-loud)。
type Clients struct {
	Auth     *auth.Client
	AppCheck *appcheck.Client
}

// NewClients は Application Default Credentials を用いて Firebase を初期化し、
// Auth / AppCheck クライアントを返す。Cloud Run 上ではサービスアカウントが
// そのまま ADC として利用される。
func NewClients(ctx context.Context) (*Clients, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("firebase.NewApp: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("app.Auth: %w", err)
	}

	appCheckClient, err := app.AppCheck(ctx)
	if err != nil {
		return nil, fmt.Errorf("app.AppCheck: %w", err)
	}

	return &Clients{
		Auth:     authClient,
		AppCheck: appCheckClient,
	}, nil
}
