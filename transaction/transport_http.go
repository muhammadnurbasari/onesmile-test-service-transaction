package transaction

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHttpServerTransaction(ctx context.Context, endpoints Endpoints, router *gin.Engine) {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeErrorResponse),
	}

	//definition of a handler
	validateTransactionHandler := httptransport.NewServer(
		endpoints.CreateTransaction, //use the endpoint
		decodeTransactionRequest,    //converts the parameters received via the request body into the struct expected by the endpoint
		encodeTransactionResponse,   //converts the struct returned by the endpoint to a json response
		options...,
	)

	router.POST("/transaction", gin.WrapH(validateTransactionHandler))

}

//converts the struct returned by the endpoint to a json response
func encodeTransactionResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

//converts the parameters received via the request body into the struct expected by the endpoint
func decodeTransactionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request transactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeErrorResponse(ctx context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
