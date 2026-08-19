package repository

import (
	"context"
	"fmt"
	"net/http"

	user_api "github.com/fun-dotto/server/gen/user"
	"github.com/fun-dotto/server/internal/modules/app/domain"
	"github.com/fun-dotto/server/internal/modules/app/external"
)

// UserRepository は外部APIからユーザーを取得する
type UserRepository struct {
	client *user_api.ClientWithResponses
}

func NewUserRepository(client *user_api.ClientWithResponses) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) GetUser(id string) (*domain.User, error) {
	response, err := r.client.UsersV1DetailWithResponse(context.Background(), id)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() == http.StatusNotFound {
		return nil, domain.ErrUserNotFound
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get user: status %d", response.StatusCode())
	}

	u := external.ToDomainUser(response.JSON200.User)
	return &u, nil
}

func (r *UserRepository) UpsertUser(id string, req domain.UserRequest) (*domain.User, error) {
	body := external.ToExternalUserRequest(req)
	response, err := r.client.UsersV1UpsertWithResponse(context.Background(), id, body)
	if err != nil {
		return nil, err
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to upsert user: status %d", response.StatusCode())
	}

	u := external.ToDomainUser(response.JSON200.User)
	return &u, nil
}

func (r *UserRepository) UpsertFCMToken(userID string, req domain.FCMTokenRequest) (*domain.FCMToken, error) {
	body := external.ToExternalFCMTokenRequest(userID, req)
	response, err := r.client.FCMTokenV1UpsertWithResponse(context.Background(), body)
	if err != nil {
		return nil, err
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to upsert fcm token: status %d", response.StatusCode())
	}

	token := external.ToDomainFCMToken(response.JSON200.FcmToken)
	return &token, nil
}
