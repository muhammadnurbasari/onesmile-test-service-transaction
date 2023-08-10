package transaction

import (
	"context"
)

type transaction struct {
	CreditCard string `json:"credit_card"`
	GrandTotal uint64 `json:"grand_total"`
	Items      []item `json:"items"`
}

type item struct {
	Name     string `json:"name"`
	Quantity uint32 `json:"quantity"`
	Subtotal uint64 `json:"subtotal"`
}

type transactionResponse struct {
	Status string `json:"status"`
	Total  uint64 `json:"total"`
}

func (t *transaction) CreateTransaction(ctx context.Context, req *transaction) (*transactionResponse, error) {
	return &transactionResponse{Status: "Berhasil", Total: 10000}, nil
}

type ServiceTransaction interface {
	CreateTransaction(ctx context.Context, req *transaction) (*transactionResponse, error)
}

func NewService(CreditCard string, GrandTotal uint64, Items []item) ServiceTransaction {
	return &transaction{CreditCard, GrandTotal, Items}
}
