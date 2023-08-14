package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/muhammadnurbasari/onesmile-test-protobuffer/proto/generate"
	"github.com/muhammadnurbasari/onesmile-test-service-transaction/dialService"
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

type items struct {
	HistoryId uint64 `json:"history_id"`
	Name      string `json:"name"`
	Quantity  uint32 `json:"quantity"`
	Subtotal  uint64 `json:"subtotal"`
}

type transactionResponse struct {
	Status string `json:"status"`
	Total  uint64 `json:"total"`
}

type transactionResponseError struct {
	Error error `json:"error"`
}

type historyResponse struct {
	Id         int32    `json:"id"`
	Items      []*items `json:"items"`
	GrandTotal int64    `json:"grand_total"`
	CreditCard string   `json:"credit_card"`
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

	var myItems []*generate.Item

	for _, v := range req.Items {
		each := &generate.Item{
			Name:     v.Name,
			Quantity: int32(v.Quantity),
			SubTotal: int64(v.Subtotal),
		}
		myItems = append(myItems, each)
	}

	// validate credit card number
	reqCreditCard := generate.CreditCard{CreditCard: req.CreditCard}
	trxCC := dialService.ServiceCreditCard()
	validate, err := trxCC.ValidateCreditCard(context.Background(), &reqCreditCard)

	if err != nil {
		return nil, err
	}

	if !validate.IsValidate {
		return nil, errors.New("not a valid credit card")
	}

	// create a new transaction
	myRequest := generate.Transaction{
		Items:      myItems,
		GrandTotal: int64(req.GrandTotal),
		CreditCard: req.CreditCard,
	}
	trx := dialService.ServiceHistory()
	_, err = trx.Create(context.Background(), &myRequest)

	if err != nil {
		return nil, err
	}

	return &transactionResponse{Status: "Berhasil", Total: req.GrandTotal}, nil
}

func (s *service) HistoryTransaction(ctx context.Context) (*[]historyResponse, error) {
	trx := dialService.ServiceHistory()
	result, err := trx.Histories(context.Background(), new(empty.Empty))

	if err != nil {
		fmt.Println("error 1 : " + err.Error())
		return nil, err
	}

	var response []historyResponse

	for _, v := range result.List {
		var itemData []*items

		for _, item := range v.Items {
			var eachItem = items{
				HistoryId: uint64(item.HistoryId),
				Name:      item.Name,
				Quantity:  uint32(item.Quantity),
				Subtotal:  uint64(item.SubTotal),
			}
			itemData = append(itemData, &eachItem)
		}

		var each = historyResponse{
			Id:         v.Id,
			Items:      itemData,
			GrandTotal: v.GrandTotal,
			CreditCard: v.CreditCard,
		}

		response = append(response, each)
	}
	return &response, nil
}

type ServiceTransaction interface {
	CreateTransaction(ctx context.Context, req *transactionRequest) (*transactionResponse, error)
	HistoryTransaction(ctx context.Context) (*[]historyResponse, error)
}

func NewService(logger log.Logger) ServiceTransaction {
	return &service{logger: log.With(logger, "service", "transaction")}
}
