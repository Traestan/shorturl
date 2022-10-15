package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Traestan/shorturl/repository"
	"github.com/go-kit/kit/log"
)

type errorer interface {
	error() error
}

func codeFrom(err error) int {
	switch err {
	case ErrOrderNotFound:
		return http.StatusBadRequest
	default:
		return http.StatusOK
	}
}

// response
func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func encodeResponseClaim(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeErrorResponse(ctx, e.error(), w)
		return nil
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	ltoken, ok := response.(*repository.User)
	if !ok {
		return nil
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: ltoken.Token,
		//Expires: expirationTime,
	})

	return json.NewEncoder(w).Encode(response)
}

func encodeResponseForgot(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeErrorResponse(ctx, e.error(), w)
		return nil
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	//fmt.Print(respons)
	return json.NewEncoder(w).Encode(response)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
		logger.Log("msg", "Error in encode response")

		//encodeErrorResponse(ctx, e.error(), w)
		return nil
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	return json.NewEncoder(w).Encode(response)
}

func redirectResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeErrorResponse(ctx, e.error(), w)
		return nil
	}

	resp := fmt.Sprintf("%s", response)
	//подсмотрено https://golang.org/src/net/http/server.go redirect
	w.Header().Set("Location", resp)
	w.WriteHeader(http.StatusMovedPermanently)

	return json.NewEncoder(w).Encode(response)
}

func qrCodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "image/png")
	w.Write(response.([]byte))

	return nil
}
