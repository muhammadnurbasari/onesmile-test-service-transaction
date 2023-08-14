package transaction

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
)

type transactionRequest struct {
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

type service struct {
	logger log.Logger
}

func (s *service) CreateTransaction(ctx context.Context, req *transactionRequest) (*transactionResponse, error) {
	var err error
	if req.CreditCard == "" {
		err = errors.New("Error, credit card cant be empty")
	}

	if req.GrandTotal <= 0 {
		err = errors.New("Error, GrandTotal cannot be less than 0")
	}

	if len(req.Items) < 1 {
		err = errors.New("Error, items cannot be empty")
	}

	if err != nil {
		return nil, err
	}

	return &transactionResponse{Status: "Berhasil", Total: req.GrandTotal}, nil
}

type ServiceTransaction interface {
	CreateTransaction(ctx context.Context, req *transactionRequest) (*transactionResponse, error)
}

func NewService(logger log.Logger) ServiceTransaction {
	return &service{logger: log.With(logger, "service", "transaction")}
}
