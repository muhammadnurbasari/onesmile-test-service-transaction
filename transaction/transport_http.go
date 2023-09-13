package transaction

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHttpServerTransaction(_ context.Context, endpoints Endpoints, router *gin.Engine) {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeErrorResponse),
	}

	//definition of a handler
	validateTransactionHandler := httptransport.NewServer(
		endpoints.CreateTransaction, //use the endpoint
		decodeTransactionRequest,    //converts the parameters received via the request body into the struct expected by the endpoint
		encodeTransactionResponse,   //converts the struct returned by the endpoint to a json response
		append(options, httptransport.ServerBefore(jwt.HTTPToContext()))...,
	)

	validateHistoryHandler := httptransport.NewServer(
		endpoints.HistoryTransaction, //use the endpoint
		func(ctx context.Context, r *http.Request) (request interface{}, err error) { return nil, nil }, //converts the parameters received via the request body into the struct expected by the endpoint
		encodeTransactionResponse, //converts the struct returned by the endpoint to a json response
		append(options, httptransport.ServerBefore(jwt.HTTPToContext()))...,
	)

	router.POST("/transaction", gin.WrapH(validateTransactionHandler))
	router.GET("/history", gin.WrapH(validateHistoryHandler))

}

//converts the struct returned by the endpoint to a json response
func encodeTransactionResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	// return json.NewEncoder(w).Encode(response)

	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func err2code(err error) int {
	return http.StatusInternalServerError
}

type errorWrapper struct {
	Error string `json:"error"`
}

//converts the parameters received via the request body into the struct expected by the endpoint
func decodeTransactionRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request transactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
