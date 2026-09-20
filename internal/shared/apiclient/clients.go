// Package apiclient は BFF (app / admin) が呼び出す各ドメイン API の
// HTTP クライアントをまとめて初期化するパッケージ。
// Cloud Run 間呼び出しのため ID トークン付きの HTTP クライアントを利用する。
package apiclient

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	academic_api "github.com/fun-dotto/server/gen/academic"
	announcement_api "github.com/fun-dotto/server/gen/announcement"
	funch_api "github.com/fun-dotto/server/gen/funch"
	user_api "github.com/fun-dotto/server/gen/user"
	"google.golang.org/api/idtoken"
)

const httpClientTimeout = 30 * time.Second

// ExternalClients 外部APIクライアントをまとめて管理
type ExternalClients struct {
	Academic     *academic_api.ClientWithResponses
	Announcement *announcement_api.ClientWithResponses
	Funch        *funch_api.ClientWithResponses
	User         *user_api.ClientWithResponses
}

// NewExternalClients 全ての外部APIクライアントを初期化
func NewExternalClients(ctx context.Context) (*ExternalClients, error) {
	academicURL, err := requireURL("ACADEMIC_API_URL")
	if err != nil {
		return nil, err
	}
	academicHTTP, err := newAuthHTTPClient(ctx, academicURL)
	if err != nil {
		return nil, fmt.Errorf("academic client: %w", err)
	}
	academic, err := academic_api.NewClientWithResponses(academicURL, academic_api.WithHTTPClient(academicHTTP))
	if err != nil {
		return nil, fmt.Errorf("academic client: %w", err)
	}

	announcementURL, err := requireURL("ANNOUNCEMENT_API_URL")
	if err != nil {
		return nil, err
	}
	announcementHTTP, err := newAuthHTTPClient(ctx, announcementURL)
	if err != nil {
		return nil, fmt.Errorf("announcement client: %w", err)
	}
	announcement, err := announcement_api.NewClientWithResponses(announcementURL, announcement_api.WithHTTPClient(announcementHTTP))
	if err != nil {
		return nil, fmt.Errorf("announcement client: %w", err)
	}

	funchURL, err := requireURL("FUNCH_API_URL")
	if err != nil {
		return nil, err
	}
	funchHTTP, err := newAuthHTTPClient(ctx, funchURL)
	if err != nil {
		return nil, fmt.Errorf("funch client: %w", err)
	}
	funch, err := funch_api.NewClientWithResponses(funchURL, funch_api.WithHTTPClient(funchHTTP))
	if err != nil {
		return nil, fmt.Errorf("funch client: %w", err)
	}

	userURL, err := requireURL("USER_API_URL")
	if err != nil {
		return nil, err
	}
	userHTTP, err := newAuthHTTPClient(ctx, userURL)
	if err != nil {
		return nil, fmt.Errorf("user client: %w", err)
	}
	user, err := user_api.NewClientWithResponses(userURL, user_api.WithHTTPClient(userHTTP))
	if err != nil {
		return nil, fmt.Errorf("user client: %w", err)
	}

	return &ExternalClients{
		Academic:     academic,
		Announcement: announcement,
		Funch:        funch,
		User:         user,
	}, nil
}

func requireURL(envKey string) (string, error) {
	url := os.Getenv(envKey)
	if url == "" {
		return "", fmt.Errorf("%s is required", envKey)
	}
	return url, nil
}

// newAuthHTTPClient Google Cloud認証付きHTTPクライアントを作成
func newAuthHTTPClient(ctx context.Context, targetURL string) (*http.Client, error) {
	client, err := idtoken.NewClient(ctx, targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth client: %w", err)
	}
	client.Timeout = httpClientTimeout
	return client, nil
}
