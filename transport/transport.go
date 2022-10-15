package transport

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-kit/kit/ratelimit"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/Traestan/shorturl/endpoint"
	"github.com/Traestan/shorturl/repository"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"

	"golang.org/x/time/rate"
)

var (
	ErrBadRouting    = errors.New("bad routing")
	ErrOrderNotFound = errors.New("url not found")
)

// NewService wires Go kit endpoints to the HTTP transport.

func MakeHTTPHandler(svcEndpoints *endpoint.Endpoints, db *repository.Model, logger log.Logger) http.Handler {
	// set-up router and initialize http endpoints
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeErrorResponse),
		kithttp.ServerBefore(jwt.HTTPToContext()),
	}

	r.Use(mux.CORSMethodMiddleware(r))

	//user work
	eRegister := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))(svcEndpoints.Register)
	r.Methods("POST", "OPTIONS").Path("/user/register").Handler(kithttp.NewServer(
		eRegister,
		decodeLoginRequest,
		encodeResponse,
		options...,
	))

	eLogin := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))(svcEndpoints.Login)

	r.Methods("POST", "OPTIONS").Path("/user/login").Handler(kithttp.NewServer(
		eLogin,
		decodeLoginRequest,
		encodeResponseClaim,
		options...,
	))

	r.Methods("POST", "OPTIONS").Path("/user/forgot").Handler(kithttp.NewServer(
		svcEndpoints.Forgot,
		decodeLoginRequest,
		encodeResponseForgot,
		options...,
	))
	r.Methods("POST", "OPTIONS").Path("/user/changepass").Handler(MyMiddleware(logger, kithttp.NewServer(
		svcEndpoints.Changepass,
		decodeProfileRequest,
		encodeResponse,
		options...,
	)))

	//urls work
	//index
	r.Methods("GET", "OPTIONS").Path("/urls").Handler(MyMiddleware(logger, kithttp.NewServer(
		svcEndpoints.Index,
		decodeIndexRequest,
		encodeResponse,
		options...,
	)))
	//
	r.Methods("POST", "OPTIONS").Path("/add").Handler(MyMiddleware(logger, kithttp.NewServer(
		svcEndpoints.Create,
		decodeCreateRequest,
		encodeResponse,
		options...,
	)))
	r.Methods("GET", "OPTIONS").Path("/get/{id}").Handler(kithttp.NewServer(
		svcEndpoints.GetUrl,
		decodeFindOneRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT", "OPTIONS").Path("/edit/{id}").Handler(MyMiddleware(logger, kithttp.NewServer(
		svcEndpoints.EditUrl,
		decodeEditUrlRequest,
		encodeResponse,
		options...,
	)))
	r.Methods("POST", "OPTIONS").Path("/del/{id}").Handler(kithttp.NewServer(
		svcEndpoints.DeleteUrl,
		decodeEditUrlRequest,
		encodeResponse,
		options...,
	))
	//stat
	r.Methods("GET", "OPTIONS").Path("/stat/{id}").Handler(MyMiddleware(logger, kithttp.NewServer(
		svcEndpoints.GetStat,
		decodeFindOneRequest,
		encodeResponse,
		options...,
	)))

	//work shorturl
	//переход по ссылке
	r.Methods("GET").Path("/{shorturl}").Handler(kithttp.NewServer(
		svcEndpoints.Root,
		decodeRootRequest,
		redirectResponse,
		options...,
	))
	r.Methods("GET").Path("/qrcode/{id}").Handler(kithttp.NewServer(
		svcEndpoints.GerQrcode,
		decodeGerQrcodeRequest,
		qrCodeResponse,
		options...,
	))

	return r
}
