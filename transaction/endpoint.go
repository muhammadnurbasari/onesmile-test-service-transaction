package transaction

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func makeCreateTransactionEndpoint(svc ServiceTransaction) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(transaction)
		res, err := svc.CreateTransaction(ctx, &req)

		if err != nil {
			return nil, err
		}

		return res, err
	}
}
