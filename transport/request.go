package transport

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Traestan/shorturl/helper"
	"github.com/Traestan/shorturl/repository"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

// request
func decodeCreateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

	jsonNews, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logger.Log(
			"panic", r,
			"req", r.URL.Query(),
			"took", time.Now())
	}

	news := &repository.UrlData{}
	err = json.Unmarshal(jsonNews, &news)
	if err != nil {
		logger.Log(
			"panic", r,
			"req", r.URL.Query(),
			"took", time.Now())
	}

	return news, nil
}

func decodeIndexRequest(_ context.Context, r *http.Request) (request interface{}, err error) {

	keys := mux.Vars(r)
	if len(keys) < 1 {
		return &repository.IndexRequest{User: ""}, nil
	}

	querySting := &repository.IndexRequest{User: ""}
	return querySting, nil
}
func decodeProfileRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

	jsonNews, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Log(
			"panic", r,
			"req", r.URL.Query(),
			"took", time.Now())
	}

	user := &repository.User{}
	err = json.Unmarshal(jsonNews, &user)
	if err != nil {
		logger.Log(
			"panic", r,
			"req", r.URL.Query(),
			"took", time.Now())
	}
	//logger.Log("msg", user)
	return user, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

	jsonNews, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Log(
			"panic", r,
			"req", r.URL.Query(),
			"took", time.Now())
	}

	user := &repository.User{}
	err = json.Unmarshal(jsonNews, &user)
	if err != nil {
		logger.Log(
			"panic", r,
			"req", r.URL.Query(),
			"took", time.Now())
	}
	return user, nil
}

func decodeFindOneRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	keys := mux.Vars(r)
	if len(keys) < 1 {
		return &repository.FindUrlDataRequest{Uuid: ""}, nil
	}
	querySting := &repository.FindUrlDataRequest{Uuid: keys["id"]}

	return querySting, nil
}

func decodeEditUrlRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	keys := mux.Vars(r)
	jsonNews, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Log(
			"panic", r,
			"req", r.URL.Query(),
			"took", time.Now())
	}

	news := &repository.UrlData{}
	err = json.Unmarshal(jsonNews, &news)
	if err != nil {
		logger.Log(
			"panic", r,
			"req", r.URL.Query(),
			"took", time.Now())
	}
	news.ShortUrl = keys["id"]
	return news, nil
}

func decodeRootRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	keys := mux.Vars(r)
	if len(keys) < 1 {
		return &repository.FindUrlDataRequest{Uuid: ""}, nil
	}

	headers := make(map[string]interface{})

	ua := r.Header.Get("User-Agent")
	referer := r.Header.Get("Referer")

	for k, v := range r.Header {
		headers[strings.ToLower(k)] = string(v[0])
	}

	querySting := &repository.FindUrlDataRequest{Uuid: keys["shorturl"],
		UHeader:   ua,
		UIp:       helper.GetIPAdress(r),
		Ureferer:  referer,
		UAHeaders: headers,
		UDateTime: time.Now().String(),
	}

	return querySting, nil
}

//qrcode
func decodeGerQrcodeRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	keys := mux.Vars(r)
	if len(keys) < 1 {
		return &repository.FindUrlDataRequest{Uuid: ""}, nil
	}
	querySting := &repository.FindUrlDataRequest{Uuid: keys["id"]}

	return querySting, nil
}
