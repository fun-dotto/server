package handler

import (
	"context"

	api "github.com/fun-dotto/subject-api/generated"
)

func (h *Handler) DayOfWeekTimetableSlotsV1List(ctx context.Context, _ api.DayOfWeekTimetableSlotsV1ListRequestObject) (api.DayOfWeekTimetableSlotsV1ListResponseObject, error) {
	slots, err := h.slotSvc.List(ctx)
	if err != nil {
		return nil, err
	}
	return api.DayOfWeekTimetableSlotsV1List200JSONResponse{DayOfWeekTimetableSlots: slotsToAPI(slots)}, nil
}

func (h *Handler) DayOfWeekTimetableSlotsV1Create(ctx context.Context, request api.DayOfWeekTimetableSlotsV1CreateRequestObject) (api.DayOfWeekTimetableSlotsV1CreateResponseObject, error) {
	slot, err := h.slotSvc.Create(ctx, slotRequestToDomain(*request.Body))
	if err != nil {
		return nil, err
	}
	return api.DayOfWeekTimetableSlotsV1Create201JSONResponse{DayOfWeekTimetableSlot: slotToAPI(slot)}, nil
}

func (h *Handler) DayOfWeekTimetableSlotsV1Detail(ctx context.Context, request api.DayOfWeekTimetableSlotsV1DetailRequestObject) (api.DayOfWeekTimetableSlotsV1DetailResponseObject, error) {
	slot, err := h.slotSvc.GetByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return api.DayOfWeekTimetableSlotsV1Detail200JSONResponse{DayOfWeekTimetableSlot: slotToAPI(slot)}, nil
}

func (h *Handler) DayOfWeekTimetableSlotsV1Update(ctx context.Context, request api.DayOfWeekTimetableSlotsV1UpdateRequestObject) (api.DayOfWeekTimetableSlotsV1UpdateResponseObject, error) {
	slot, err := h.slotSvc.Update(ctx, request.Id, slotRequestToDomain(*request.Body))
	if err != nil {
		return nil, err
	}
	return api.DayOfWeekTimetableSlotsV1Update200JSONResponse{DayOfWeekTimetableSlot: slotToAPI(slot)}, nil
}

func (h *Handler) DayOfWeekTimetableSlotsV1Delete(ctx context.Context, request api.DayOfWeekTimetableSlotsV1DeleteRequestObject) (api.DayOfWeekTimetableSlotsV1DeleteResponseObject, error) {
	if err := h.slotSvc.Delete(ctx, request.Id); err != nil {
		return nil, err
	}
	return api.DayOfWeekTimetableSlotsV1Delete204Response{}, nil
}
