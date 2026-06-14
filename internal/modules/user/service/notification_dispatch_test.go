package service

import (
	"context"
	"errors"
	"sort"
	"testing"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/fun-dotto/user-api/internal/domain"
)

type stubNotificationRepo struct {
	getByIDs           func(ctx context.Context, ids []string) ([]domain.Notification, error)
	dispatch           func(ctx context.Context, deliveries map[string][]string) ([]domain.Notification, error)
	dispatchCalled     bool
	lastDeliveriesArg  map[string][]string
	getByIDsCalledWith []string
}

func (s *stubNotificationRepo) ListNotifications(context.Context, domain.NotificationListFilter) ([]domain.Notification, error) {
	return nil, errors.New("not implemented")
}
func (s *stubNotificationRepo) CreateNotification(context.Context, domain.Notification) (domain.Notification, error) {
	return domain.Notification{}, errors.New("not implemented")
}
func (s *stubNotificationRepo) UpdateNotification(context.Context, domain.Notification) (domain.Notification, error) {
	return domain.Notification{}, errors.New("not implemented")
}
func (s *stubNotificationRepo) DeleteNotification(context.Context, string) error {
	return errors.New("not implemented")
}
func (s *stubNotificationRepo) GetNotificationsByIDs(ctx context.Context, ids []string) ([]domain.Notification, error) {
	s.getByIDsCalledWith = ids
	if s.getByIDs == nil {
		return nil, nil
	}
	return s.getByIDs(ctx, ids)
}
func (s *stubNotificationRepo) DispatchNotifications(ctx context.Context, deliveries map[string][]string) ([]domain.Notification, error) {
	s.dispatchCalled = true
	s.lastDeliveriesArg = deliveries
	if s.dispatch == nil {
		return nil, nil
	}
	return s.dispatch(ctx, deliveries)
}

type stubFCMTokenRepo struct {
	list       func(ctx context.Context, filter domain.FCMTokenListFilter) ([]domain.FCMToken, error)
	called     bool
	lastFilter domain.FCMTokenListFilter
}

func (s *stubFCMTokenRepo) ListFCMTokens(ctx context.Context, filter domain.FCMTokenListFilter) ([]domain.FCMToken, error) {
	s.called = true
	s.lastFilter = filter
	if s.list == nil {
		return nil, nil
	}
	return s.list(ctx, filter)
}

type stubMessagingClient struct {
	send       func(ctx context.Context, msg *messaging.MulticastMessage) (*messaging.BatchResponse, error)
	calls      int
	tokenCalls [][]string
}

func (s *stubMessagingClient) SendEachForMulticast(ctx context.Context, msg *messaging.MulticastMessage) (*messaging.BatchResponse, error) {
	s.calls++
	tokens := append([]string{}, msg.Tokens...)
	s.tokenCalls = append(s.tokenCalls, tokens)
	if s.send == nil {
		responses := make([]*messaging.SendResponse, len(msg.Tokens))
		for i := range responses {
			responses[i] = &messaging.SendResponse{Success: true, MessageID: "msg"}
		}
		return &messaging.BatchResponse{SuccessCount: len(responses), Responses: responses}, nil
	}
	return s.send(ctx, msg)
}

func newServiceWithStubs(repo *stubNotificationRepo, tokenRepo *stubFCMTokenRepo, msg *stubMessagingClient) *NotificationService {
	return &NotificationService{repo: repo, fcmTokenRepo: tokenRepo, messagingClient: msg}
}

func notifiedAt(t time.Time) *time.Time { return &t }

func TestDispatchNotifications_NoNotifications(t *testing.T) {
	repo := &stubNotificationRepo{getByIDs: func(context.Context, []string) ([]domain.Notification, error) {
		return []domain.Notification{}, nil
	}}
	tokenRepo := &stubFCMTokenRepo{}
	msg := &stubMessagingClient{}
	svc := newServiceWithStubs(repo, tokenRepo, msg)

	got, err := svc.DispatchNotifications(context.Background(), []string{"n1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty result, got %d", len(got))
	}
	if tokenRepo.called {
		t.Errorf("expected ListFCMTokens not called")
	}
	if msg.calls != 0 {
		t.Errorf("expected no FCM calls, got %d", msg.calls)
	}
	if repo.dispatchCalled {
		t.Errorf("expected DispatchNotifications not called")
	}
}

func TestDispatchNotifications_GetByIDsError(t *testing.T) {
	wantErr := errors.New("db down")
	repo := &stubNotificationRepo{getByIDs: func(context.Context, []string) ([]domain.Notification, error) {
		return nil, wantErr
	}}
	svc := newServiceWithStubs(repo, &stubFCMTokenRepo{}, &stubMessagingClient{})

	if _, err := svc.DispatchNotifications(context.Background(), []string{"n1"}); !errors.Is(err, wantErr) {
		t.Errorf("expected error %v, got %v", wantErr, err)
	}
}

func TestDispatchNotifications_NoTargetUsers(t *testing.T) {
	repo := &stubNotificationRepo{getByIDs: func(context.Context, []string) ([]domain.Notification, error) {
		return []domain.Notification{{ID: "n1", Title: "t", Body: "b"}}, nil
	}}
	tokenRepo := &stubFCMTokenRepo{}
	msg := &stubMessagingClient{}
	svc := newServiceWithStubs(repo, tokenRepo, msg)

	got, err := svc.DispatchNotifications(context.Background(), []string{"n1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty result, got %d", len(got))
	}
	if tokenRepo.called {
		t.Errorf("expected ListFCMTokens not called when no target users")
	}
	if repo.dispatchCalled {
		t.Errorf("expected DispatchNotifications not called")
	}
}

func TestDispatchNotifications_NoTokens_StillRecordsDelivery(t *testing.T) {
	notification := domain.Notification{
		ID:    "n1",
		Title: "t",
		Body:  "b",
		TargetUsers: []domain.NotificationTargetUser{
			{UserID: "u1"},
			{UserID: "u2", NotifiedAt: notifiedAt(time.Now())},
		},
	}
	repo := &stubNotificationRepo{
		getByIDs: func(context.Context, []string) ([]domain.Notification, error) {
			return []domain.Notification{notification}, nil
		},
		dispatch: func(_ context.Context, _ map[string][]string) ([]domain.Notification, error) {
			return []domain.Notification{notification}, nil
		},
	}
	tokenRepo := &stubFCMTokenRepo{list: func(context.Context, domain.FCMTokenListFilter) ([]domain.FCMToken, error) {
		return []domain.FCMToken{}, nil
	}}
	msg := &stubMessagingClient{}
	svc := newServiceWithStubs(repo, tokenRepo, msg)

	got, err := svc.DispatchNotifications(context.Background(), []string{"n1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(got))
	}
	if msg.calls != 0 {
		t.Errorf("expected no FCM calls when no tokens, got %d", msg.calls)
	}
	gotUsers := repo.lastDeliveriesArg["n1"]
	sort.Strings(gotUsers)
	if len(gotUsers) != 2 || gotUsers[0] != "u1" || gotUsers[1] != "u2" {
		t.Errorf("expected both u1 and u2 in deliveries (force re-send), got %v", gotUsers)
	}
}

func TestDispatchNotifications_AllTokensSucceed(t *testing.T) {
	notification := domain.Notification{
		ID:    "n1",
		Title: "t",
		Body:  "b",
		TargetUsers: []domain.NotificationTargetUser{
			{UserID: "u1"},
			{UserID: "u2", NotifiedAt: notifiedAt(time.Now())},
		},
	}
	repo := &stubNotificationRepo{
		getByIDs: func(context.Context, []string) ([]domain.Notification, error) {
			return []domain.Notification{notification}, nil
		},
		dispatch: func(_ context.Context, _ map[string][]string) ([]domain.Notification, error) {
			return []domain.Notification{notification}, nil
		},
	}
	tokenRepo := &stubFCMTokenRepo{list: func(context.Context, domain.FCMTokenListFilter) ([]domain.FCMToken, error) {
		return []domain.FCMToken{
			{UserID: "u1", Token: "tok-u1"},
			{UserID: "u2", Token: "tok-u2"},
		}, nil
	}}
	msg := &stubMessagingClient{}
	svc := newServiceWithStubs(repo, tokenRepo, msg)

	if _, err := svc.DispatchNotifications(context.Background(), []string{"n1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.calls != 1 {
		t.Errorf("expected 1 FCM batch call, got %d", msg.calls)
	}
	gotUsers := repo.lastDeliveriesArg["n1"]
	sort.Strings(gotUsers)
	if len(gotUsers) != 2 || gotUsers[0] != "u1" || gotUsers[1] != "u2" {
		t.Errorf("expected both users delivered, got %v", gotUsers)
	}
}

func TestDispatchNotifications_PartialFCMFailure(t *testing.T) {
	notification := domain.Notification{
		ID:    "n1",
		Title: "t",
		Body:  "b",
		TargetUsers: []domain.NotificationTargetUser{
			{UserID: "u1"},
			{UserID: "u2"},
			{UserID: "u3"},
		},
	}
	repo := &stubNotificationRepo{
		getByIDs: func(context.Context, []string) ([]domain.Notification, error) {
			return []domain.Notification{notification}, nil
		},
		dispatch: func(_ context.Context, _ map[string][]string) ([]domain.Notification, error) {
			return []domain.Notification{notification}, nil
		},
	}
	tokenRepo := &stubFCMTokenRepo{list: func(context.Context, domain.FCMTokenListFilter) ([]domain.FCMToken, error) {
		return []domain.FCMToken{
			{UserID: "u1", Token: "tok-u1"},
			{UserID: "u2", Token: "tok-u2"},
		}, nil
	}}
	msg := &stubMessagingClient{send: func(_ context.Context, m *messaging.MulticastMessage) (*messaging.BatchResponse, error) {
		responses := make([]*messaging.SendResponse, len(m.Tokens))
		for i, tk := range m.Tokens {
			if tk == "tok-u2" {
				responses[i] = &messaging.SendResponse{Error: errors.New("invalid token")}
			} else {
				responses[i] = &messaging.SendResponse{Success: true, MessageID: "msg"}
			}
		}
		return &messaging.BatchResponse{Responses: responses}, nil
	}}
	svc := newServiceWithStubs(repo, tokenRepo, msg)

	if _, err := svc.DispatchNotifications(context.Background(), []string{"n1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gotUsers := repo.lastDeliveriesArg["n1"]
	sort.Strings(gotUsers)
	want := []string{"u1", "u3"}
	if len(gotUsers) != len(want) || gotUsers[0] != want[0] || gotUsers[1] != want[1] {
		t.Errorf("expected delivered=%v (u1 succeeded; u3 had no token), got %v", want, gotUsers)
	}
}

func TestDispatchNotifications_FCMBatchError_KeepsPartialSuccess(t *testing.T) {
	notification := domain.Notification{
		ID:    "n1",
		Title: "t",
		Body:  "b",
		TargetUsers: []domain.NotificationTargetUser{
			{UserID: "u1"},
		},
	}
	repo := &stubNotificationRepo{
		getByIDs: func(context.Context, []string) ([]domain.Notification, error) {
			return []domain.Notification{notification}, nil
		},
		dispatch: func(_ context.Context, _ map[string][]string) ([]domain.Notification, error) {
			return []domain.Notification{notification}, nil
		},
	}
	tokenRepo := &stubFCMTokenRepo{list: func(context.Context, domain.FCMTokenListFilter) ([]domain.FCMToken, error) {
		return []domain.FCMToken{{UserID: "u1", Token: "tok-u1"}}, nil
	}}
	msg := &stubMessagingClient{send: func(context.Context, *messaging.MulticastMessage) (*messaging.BatchResponse, error) {
		return nil, errors.New("fcm down")
	}}
	svc := newServiceWithStubs(repo, tokenRepo, msg)

	got, err := svc.DispatchNotifications(context.Background(), []string{"n1"})
	if err != nil {
		t.Fatalf("expected nil error (FCM error is logged, not returned), got %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty result when no users delivered, got %d", len(got))
	}
	if repo.dispatchCalled {
		t.Errorf("expected DispatchNotifications not called when no successful deliveries")
	}
}

func TestDispatchNotifications_FCMTokenListError(t *testing.T) {
	notification := domain.Notification{
		ID:    "n1",
		Title: "t",
		Body:  "b",
		TargetUsers: []domain.NotificationTargetUser{
			{UserID: "u1"},
		},
	}
	repo := &stubNotificationRepo{getByIDs: func(context.Context, []string) ([]domain.Notification, error) {
		return []domain.Notification{notification}, nil
	}}
	wantErr := errors.New("token db down")
	tokenRepo := &stubFCMTokenRepo{list: func(context.Context, domain.FCMTokenListFilter) ([]domain.FCMToken, error) {
		return nil, wantErr
	}}
	svc := newServiceWithStubs(repo, tokenRepo, &stubMessagingClient{})

	if _, err := svc.DispatchNotifications(context.Background(), []string{"n1"}); !errors.Is(err, wantErr) {
		t.Errorf("expected error %v, got %v", wantErr, err)
	}
}
