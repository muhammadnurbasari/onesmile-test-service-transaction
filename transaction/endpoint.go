package transaction

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateTransaction endpoint.Endpoint
}

func MakeEndpoints(s ServiceTransaction) Endpoints {
	return Endpoints{
		CreateTransaction: makeCreateTransactionEndpoint(s),
	}
}

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
