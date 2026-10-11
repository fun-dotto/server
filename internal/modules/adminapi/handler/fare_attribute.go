package handler

import (
	"context"

	"connectrpc.com/connect"
	adminv1 "github.com/fun-dotto/server/gen/admin/v1"
	"github.com/fun-dotto/server/internal/shared/model"
	"github.com/shopspring/decimal"
	"google.golang.org/genproto/googleapis/type/money"
)

type FareAttributeHandler struct {
	c crud[model.FareAttribute]
}

func NewFareAttributeHandler(repo Repository[model.FareAttribute]) *FareAttributeHandler {
	return &FareAttributeHandler{c: crud[model.FareAttribute]{
		name: "fare_attribute",
		repo: repo,
		keys: func(m model.FareAttribute) []keyField { return []keyField{{"fare_id", m.FareID}} },
	}}
}

// fareCurrencyCode は運賃の通貨。DB は通貨を持たないため日本円に固定する。
const fareCurrencyCode = "JPY"

var fareAttributeRiderCategoryEnum = enumMap[adminv1.FareAttributeRiderCategory]{
	adminv1.FareAttributeRiderCategory_FARE_ATTRIBUTE_RIDER_CATEGORY_ADULT: string(model.RiderCategoryAdult),
}

type fareAttributeInput interface {
	GetFareId() string
	GetRiderCategory() adminv1.FareAttributeRiderCategory
	GetPrice() *money.Money
}

func fareAttributeFromProto(in fareAttributeInput, s fieldSet) (model.FareAttribute, error) {
	var v violations
	var m model.FareAttribute
	if s.has("fare_id") {
		m.FareID = requireString(&v, "fare_id", in.GetFareId())
	}
	if s.has("rider_category") {
		m.RiderCategory = model.RiderCategory(fareAttributeRiderCategoryEnum.toDB(&v, "rider_category", in.GetRiderCategory(), true))
	}
	if s.has("price") {
		p := in.GetPrice()
		switch {
		case p == nil:
			v.add("price", "must be set")
		case p.GetCurrencyCode() != fareCurrencyCode:
			v.add("price.currency_code", "must be "+fareCurrencyCode)
		case p.GetUnits() < 0 || p.GetNanos() < 0:
			v.add("price", "must not be negative")
		case p.GetNanos() >= 1_000_000_000:
			v.add("price.nanos", "must be less than 1,000,000,000")
		}
		m.Price = decimal.New(p.GetUnits(), 0).Add(decimal.New(int64(p.GetNanos()), -9))
	}
	return m, v.err()
}

func fareAttributeToProto(m model.FareAttribute) *adminv1.FareAttribute {
	units := m.Price.Truncate(0)
	return &adminv1.FareAttribute{
		FareId:        m.FareID,
		RiderCategory: fareAttributeRiderCategoryEnum.fromDB(string(m.RiderCategory)),
		Price: &money.Money{
			CurrencyCode: fareCurrencyCode,
			Units:        units.IntPart(),
			Nanos:        int32(m.Price.Sub(units).Shift(9).IntPart()),
		},
	}
}

func (h *FareAttributeHandler) ListFareAttributes(ctx context.Context, req *connect.Request[adminv1.ListFareAttributesRequest]) (*connect.Response[adminv1.ListFareAttributesResponse], error) {
	items, next, err := h.c.list(ctx, req.Msg.GetPageSize(), req.Msg.GetPageToken())
	if err != nil {
		return nil, err
	}
	res := &adminv1.ListFareAttributesResponse{NextPageToken: next}
	for _, m := range items {
		res.FareAttributes = append(res.FareAttributes, fareAttributeToProto(m))
	}
	return connect.NewResponse(res), nil
}

func (h *FareAttributeHandler) GetFareAttribute(ctx context.Context, req *connect.Request[adminv1.GetFareAttributeRequest]) (*connect.Response[adminv1.GetFareAttributeResponse], error) {
	fareID, err := requireKeyString("fare_id", req.Msg.GetFareId())
	if err != nil {
		return nil, err
	}
	m, err := h.c.get(ctx, []keyField{{"fare_id", fareID}})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.GetFareAttributeResponse{FareAttribute: fareAttributeToProto(m)}), nil
}

func (h *FareAttributeHandler) CreateFareAttribute(ctx context.Context, req *connect.Request[adminv1.CreateFareAttributeRequest]) (*connect.Response[adminv1.CreateFareAttributeResponse], error) {
	m, err := fareAttributeFromProto(req.Msg, nil)
	if err != nil {
		return nil, err
	}
	created, err := h.c.create(ctx, &m)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.CreateFareAttributeResponse{FareAttribute: fareAttributeToProto(created)}), nil
}

func (h *FareAttributeHandler) UpdateFareAttribute(ctx context.Context, req *connect.Request[adminv1.UpdateFareAttributeRequest]) (*connect.Response[adminv1.UpdateFareAttributeResponse], error) {
	fareID, err := requireKeyString("fare_id", req.Msg.GetFareId())
	if err != nil {
		return nil, err
	}
	s, err := updateFields(req.Msg, req.Msg.GetUpdateMask(), "fare_id")
	if err != nil {
		return nil, err
	}
	m, err := fareAttributeFromProto(req.Msg, s)
	if err != nil {
		return nil, err
	}
	updated, err := h.c.update(ctx, []keyField{{"fare_id", fareID}}, &m, s.columns(nil))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.UpdateFareAttributeResponse{FareAttribute: fareAttributeToProto(updated)}), nil
}

func (h *FareAttributeHandler) DeleteFareAttribute(ctx context.Context, req *connect.Request[adminv1.DeleteFareAttributeRequest]) (*connect.Response[adminv1.DeleteFareAttributeResponse], error) {
	fareID, err := requireKeyString("fare_id", req.Msg.GetFareId())
	if err != nil {
		return nil, err
	}
	if err := h.c.delete(ctx, []keyField{{"fare_id", fareID}}); err != nil {
		return nil, err
	}
	return connect.NewResponse(&adminv1.DeleteFareAttributeResponse{}), nil
}
