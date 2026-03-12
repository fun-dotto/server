package service

import (
	"context"

	"github.com/fun-dotto/academic-api/internal/domain"
	"github.com/google/uuid"
)

type dayOfWeekTimetableSlotRepository interface {
	List(ctx context.Context) ([]domain.DayOfWeekTimetableSlot, error)
	GetByID(ctx context.Context, id string) (domain.DayOfWeekTimetableSlot, error)
	Create(ctx context.Context, slot domain.DayOfWeekTimetableSlot) (domain.DayOfWeekTimetableSlot, error)
	Update(ctx context.Context, slot domain.DayOfWeekTimetableSlot) (domain.DayOfWeekTimetableSlot, error)
	Delete(ctx context.Context, id string) error
}

type DayOfWeekTimetableSlotService struct {
	repo dayOfWeekTimetableSlotRepository
}

func NewDayOfWeekTimetableSlotService(repo dayOfWeekTimetableSlotRepository) *DayOfWeekTimetableSlotService {
	return &DayOfWeekTimetableSlotService{repo: repo}
}

func (s *DayOfWeekTimetableSlotService) List(ctx context.Context) ([]domain.DayOfWeekTimetableSlot, error) {
	return s.repo.List(ctx)
}

func (s *DayOfWeekTimetableSlotService) GetByID(ctx context.Context, id string) (domain.DayOfWeekTimetableSlot, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DayOfWeekTimetableSlotService) Create(ctx context.Context, slot domain.DayOfWeekTimetableSlot) (domain.DayOfWeekTimetableSlot, error) {
	slot.ID = uuid.New().String()
	return s.repo.Create(ctx, slot)
}

func (s *DayOfWeekTimetableSlotService) Update(ctx context.Context, id string, slot domain.DayOfWeekTimetableSlot) (domain.DayOfWeekTimetableSlot, error) {
	slot.ID = id
	return s.repo.Update(ctx, slot)
}

func (s *DayOfWeekTimetableSlotService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
