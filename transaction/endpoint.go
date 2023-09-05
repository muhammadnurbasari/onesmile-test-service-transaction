package transaction

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateTransaction  endpoint.Endpoint
	HistoryTransaction endpoint.Endpoint
}

func MakeEndpoints(s ServiceTransaction) Endpoints {
	return Endpoints{
		CreateTransaction:  makeCreateTransactionEndpoint(s),
		HistoryTransaction: makeHistoryTransactionEndpoint(s),
	}
}

// Create Transaction godoc
// @Summary Create Transaction
// @Description API for Create transaction
// @Tags TRANSACTION
// @Accept  json
// @Param request body transactionRequest true "Request Body Raw"
// @Produce  json
// @Success 200 {object} transactionResponse
// @Failure 400 {object} transactionResponseError
// @Router /transaction [post]
func makeCreateTransactionEndpoint(svc ServiceTransaction) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(transactionRequest)
		res, err := svc.CreateTransaction(ctx, &req)

		if err != nil {
			return nil, err
		}

		return res, err
	}
}

// History godoc
// @Summary History
// @Description API for History
// @Tags TRANSACTION
// @Accept  json
// @Produce  json
// @Success 200 {object} []historyResponse
// @Failure 400 {object} transactionResponseError
// @Router /history [get]
func makeHistoryTransactionEndpoint(svc ServiceTransaction) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (response interface{}, err error) {
		res, err := svc.HistoryTransaction(ctx)

		if err != nil {
			return HistoryResponse{Message: err.Error(), Data: []historyResponse{}, Error: err}, err
		}

		if len(*res) == 0 {
			return HistoryResponse{Message: "success", Data: []historyResponse{}, Error: nil}, err
		}

		return HistoryResponse{Message: "success", Data: *res, Error: nil}, err
	}
}

type HistoryResponse struct {
	Message string            `json:"message"`
	Data    []historyResponse `json:"data"`
	Error   error             `json:"error"`
}

func (r HistoryResponse) Failed() error { return r.Error }
