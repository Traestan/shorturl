package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Traestan/shorturl/endpoint"
	"github.com/Traestan/shorturl/repository"
	"github.com/Traestan/shorturl/transport"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMakeHTTPHandler(t *testing.T) {
	// в этот массив можно добавлять всех типов запросы, которые должны выполняться
	reqURL := "http://127.0.0.1:5000/add"
	var requestsArr = []string{
		"http://yandex.ru",
		"http://google.ru",
		"http://mail.ru",
		"http://ya.ru",
	}

	for _, req := range requestsArr {
		values := map[string]string{"surl": req, "ip": "127.0.0.1"}
		jsonValue, err := json.Marshal(values)
		if err != nil {
			t.Errorf("ERRORS MArshals: %v", err)
			return
		}
		t.Errorf("%s\n", jsonValue)
		resp, err := http.Post(reqURL, "application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			t.Errorf("ERRORS: %v", err)
			return
		}

		if resp == nil {
			t.Errorf("Request ERROR: null response")
			return
		}

		respData := make(map[string]string)
		err = json.NewDecoder(resp.Body).Decode(&respData)
		if err != nil {
			t.Errorf("Decode ERROR: %v", err)
			return
		}

		if _, ok := respData["hurl"]; ok {
			t.Errorf("Response ERROR: %v", respData["hurl"])
		}
	}
}
func TestGetresult(t *testing.T) {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

	clientOptions := options.Client().ApplyURI("mongodb://shorturl:123123@localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		t.Error("Init ERROR:")
	}
	urlEncoder := repository.NewURLEncoder(&repository.URLEncoderConfig{})
	mng := repository.NewModel(client, logger, urlEncoder) //connect to mongo

	if err != nil {
		t.Errorf("Init ERROR: %v", err)
		return
	}

	//endpoint init
	sep := endpoint.MakeEndpoints(&mng)

	go http.ListenAndServe(":5000", transport.MakeHTTPHandler(sep, &mng, logger))

	var requestsArr = []string{
		"3aa754e4ed",
		"3705ce87a6",
	}
	reqURL := "http://127.0.0.1:5000/get/"
	for _, req := range requestsArr {
		t.Errorf(reqURL + req)
		resp, err := http.Get(reqURL + req)
		if resp == nil {
			t.Errorf("Request ERROR: null response")
			return
		}

		respData := make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&respData)
		if err != nil {
			t.Errorf("Decode ERROR: %v", err)
			return
		}

	}
}

func TestHTTP(t *testing.T) {
	/*zkt, _ := zipkin.NewTracer(nil, zipkin.WithNoopTracer(true))
	mt := service.NewMetricsTest()
	svc := personservice.New(mt)
	ep := personendpoints.New(svc, log.NewNopLogger(), mt, zkt)*/
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))

	clientOptions := options.Client().ApplyURI("mongodb://shorturl:123123@localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		t.Error("Init ERROR:")
	}
	urlEncoder := repository.NewURLEncoder(&repository.URLEncoderConfig{})
	mng := repository.NewModel(client, logger, urlEncoder) //connect to mongo

	if err != nil {
		t.Errorf("Init ERROR: %v", err)
		return
	}

	//endpoint init
	sep := endpoint.MakeEndpoints(&mng)

	mux := transport.MakeHTTPHandler(sep, &mng, log.NewNopLogger())
	srv := httptest.NewServer(mux)
	defer srv.Close()

	for _, testcase := range []struct {
		method, url, body, want string
	}{
		{"POST", srv.URL + "/add", `{"ip":"127.0.0.1","surl":"https://yandex.ru"}`, `{"error":"","response":{"message":"test","err":null}}`},
		{"POST", srv.URL + "/add", `{"ip":"127.0.0.1","surl":"https://google.ru"}`, `{"error":"","response":{"message":"test","err":null}}`},
		{"POST", srv.URL + "/add", `{"ip":"127.0.0.1","surl":"https://rambler.ru"}`, `{"error":"","response":{"message":"test","err":null}}`},
	} {
		t.Errorf("%s", testcase.body)
		req, _ := http.NewRequest(testcase.method, testcase.url, strings.NewReader(testcase.body))
		resp, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		if want, have := testcase.want, strings.TrimSpace(string(body)); want != have {
			t.Errorf("%s %s %s: want %q, have %q", testcase.method, testcase.url, testcase.body, want, have)
		}
	}
}
