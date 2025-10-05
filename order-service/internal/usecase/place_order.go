package usecase

import (
	"context"
	"errors"
	"time"

	invport "github.com/adopabianko/commerce/order-service/internal/domain/inventory"
	"github.com/adopabianko/commerce/order-service/internal/domain/order"
	inventoryv1 "github.com/adopabianko/commerce/proto/gen/inventory/v1"
	"github.com/google/uuid"
)

type PlaceOrder struct {
	repo    order.Repository
	inv     invport.Client
	timeout time.Duration
}

func NewPlaceOrder(repo order.Repository, inv invport.Client, timeout time.Duration) *PlaceOrder {
	return &PlaceOrder{repo: repo, inv: inv, timeout: timeout}
}

type Item struct {
	SKU string
	Qty int32
}
type Request struct{ Items []Item }
type Response struct{ OrderID, Status string }

func (uc *PlaceOrder) Exec(ctx context.Context, req Request) (*Response, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("no items")
	}

	id := uuid.New().String()
	o := &order.Order{ID: id, Status: "PENDING"}
	for _, it := range req.Items {
		o.Items = append(o.Items, order.OrderItem{SKU: it.SKU, Qty: it.Qty})
	}
	if err := uc.repo.CreateOrder(o); err != nil {
		return nil, err
	}

	ctxTO, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	chk := &inventoryv1.CheckStockRequest{}
	for _, it := range req.Items {
		chk.Items = append(chk.Items, &inventoryv1.Item{Sku: it.SKU, Qty: it.Qty})
	}
	resChk, err := uc.inv.CheckStock(ctxTO, chk)
	if err != nil || !resChk.Ok {
		_ = uc.repo.UpdateStatus(id, "REJECTED")
		return &Response{OrderID: id, Status: "REJECTED"}, nil
	}

	resReq := &inventoryv1.ReserveStockRequest{}
	for _, it := range req.Items {
		resReq.Items = append(resReq.Items, &inventoryv1.Item{Sku: it.SKU, Qty: it.Qty})
	}
	resRes, err := uc.inv.ReserveStock(ctxTO, resReq)
	if err != nil || !resRes.Ok {
		_ = uc.repo.UpdateStatus(id, "REJECTED")
		return &Response{OrderID: id, Status: "REJECTED"}, nil
	}

	_ = uc.repo.UpdateStatus(id, "CONFIRMED")
	return &Response{OrderID: id, Status: "CONFIRMED"}, nil
}
