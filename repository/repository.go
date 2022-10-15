package repository

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"

	"fmt"

	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var JwtKey = []byte("LOLO")

type Models interface {
	//FindAll([]string) string
	Registration(map[string]string) string
	GetOne(string) string
}

type Model struct {
	mongo   *mongo.Client
	logs    log.Logger
	encoder *urlEncoder
}

// NewModel creates object of Models
func NewModel(connection *mongo.Client, logger log.Logger, urlencoder *urlEncoder) Model {
	return Model{
		mongo:   connection,
		logs:    logger,
		encoder: urlencoder,
	}
}

func (r Model) Registration(ctx context.Context, user *User) (interface{}, error) {
	r.logs.Log("msg", "Add")

	user.Password = GetMD5Hash(user.Password)
	user.Date_Add = time.Now().String()
	if user.Email != "" {
		collection := r.mongo.Database("shorturl").Collection("users_shop")

		//проверка на наличие email в базе
		var result User
		filter := bson.D{{"email", user.Email}}
		err := collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			r.logs.Log("error", err)
		}
		if err == nil {
			r.logs.Log("msg", "Email isset")
			return nil, NewError("Email isset") //errors.New("Email isset")
		} else {
			insertResult, err := collection.InsertOne(context.TODO(), user)
			if err != nil {
				r.logs.Log("Error on inserting new Hero", err)
			}

			return insertResult.InsertedID, nil
		}

	} else {
		return nil, nil
	}

}

func (r Model) Login(ctx context.Context, user *User) (*User, error) {
	r.logs.Log("msg", "Login")

	//var result User
	result := &User{}
	if user.Email != "" {
		password := GetMD5Hash(user.Password)
		filter := bson.D{{"email", user.Email}, {"password", password}}

		connect := r.mongo.Database("shorturl")

		urlShorted := connect.Collection("users_shop")
		err := urlShorted.FindOne(context.TODO(), filter).Decode(&result)

		if err != nil {
			r.logs.Log("error", err)
			return nil, err
		}
	}
	return result, nil
}

func (r Model) Forgot(ctx context.Context, user *User) (*User, error) {
	r.logs.Log("msg", "Forgot")

	var result User
	if user.Email != "" {

		filter := bson.D{{"email", user.Email}}

		connect := r.mongo.Database("shorturl")

		urlShorted := connect.Collection("users_shop")
		err := urlShorted.FindOne(context.TODO(), filter).Decode(&result)

		if err != nil {
			r.logs.Log("error", err)
			return nil, NewError("Email not isset")
		}

		rand.Seed(time.Now().UnixNano())

		chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
			"abcdefghijklmnopqrstuvwxyz" +
			"0123456789")
		length := 8
		var b strings.Builder
		for i := 0; i < length; i++ {
			b.WriteRune(chars[rand.Intn(len(chars))])
		}
		newPass := b.String()
		user.Password = GetMD5Hash(newPass)

		update := bson.D{
			{"$set", bson.D{
				{"password", user.Password},
			}},
		}
		//err := urlShorted.FindOne(context.TODO(), filter).Decode(&result)
		_, err = urlShorted.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			r.logs.Log("error", err)
			return nil, err
		}
		user.Password = newPass
	}
	return user, nil
}

func (r Model) Changepass(ctx context.Context, user *User) (string, error) {
	r.logs.Log("msg", "Changepass")
	if user.Email != "" {

		filter := bson.D{{"email", user.Email}}
		connect := r.mongo.Database("shorturl")

		urlShorted := connect.Collection("urls_shorted")
		password := GetMD5Hash(user.NewPassword)
		update := bson.D{
			{"$set", bson.D{
				{"password", password},
			}},
		}
		_, err := urlShorted.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return "", err
		}
	} else {
		r.logs.Log("err", "Username empty")
		return "", nil
	}
	return "yes", nil
}
func (r Model) Create(ctx context.Context, url *UrlData, user string) (interface{}, error) {
	r.logs.Log("msg", "Add")

	url.Uuid = GetMD5Hash(url.SourceUrl + "-" + url.Ip)
	url.ShortUrl = r.encoder.EncodeURL(r.encoder.DecodeURL(url.SourceUrl + time.Now().String())) //GetMD5HashShort(url.Uuid + "-" + time.Now().String())
	url.Date_Add = time.Now().String()
	url.Ip = user

	r.logs.Log("msg", url)

	if url.SourceUrl != "" {
		collection := r.mongo.Database("shorturl").Collection("urls_shorted")
		_, err := collection.InsertOne(context.TODO(), url)
		if err != nil {
			r.logs.Log("Error on inserting", err)
		}
		return url.ShortUrl, nil
	}
	return nil, nil
}

func (r Model) GetUrls(ctx context.Context, user string) (*ResultUrlDataAll, error) {
	r.logs.Log("msg", "Get Urls")

	var err error

	filter := bson.D{{"ip", user}}
	connect := r.mongo.Database("shorturl")

	collect := connect.Collection("urls_shorted")
	cur, err := collect.Find(context.TODO(), filter)

	urlShorted := &ResultUrlDataAll{}
	if err != nil {
		r.logs.Log("error", err)
		result := &ResultUrlDataAll{}
		return result, nil
	}

	for cur.Next(context.TODO()) {
		statist := ResultUrlData{}
		err := cur.Decode(&statist)
		if err != nil {
			r.logs.Log("error", err)
			result := &ResultUrlDataAll{}
			return result, nil
		}
		urlShorted.Result = append(urlShorted.Result, statist)
	}
	//r.logs.Log("msg", urlShorted)
	return urlShorted, nil
}

func (r Model) GetOne(ctx context.Context, uuid string) (*UrlData, error) {
	r.logs.Log("msg", "Find One"+uuid)
	//r.logs.Log("msg")
	var result UrlData
	filter := bson.D{{"shorturl", uuid}}
	connect := r.mongo.Database("shorturl")

	urlShorted := connect.Collection("urls_shorted")
	err := urlShorted.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		r.logs.Log("error", err)
		return &result, nil
	}

	return &result, nil
}
func (r Model) GetOneUrl(ctx context.Context, uuid, user string) (*UrlData, error) {
	r.logs.Log("msg", "Find One "+uuid)
	//r.logs.Log("msg")
	var result UrlData
	filter := bson.D{{"shorturl", uuid}, {"ip", user}}
	connect := r.mongo.Database("shorturl")

	urlShorted := connect.Collection("urls_shorted")
	err := urlShorted.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		r.logs.Log("error", err)
		return &result, nil
	}

	return &result, nil
}

func (r Model) GetPage(ctx context.Context, page string) (*PageContent, error) {
	r.logs.Log("msg", "GetPage")

	result := &PageContent{Content: "Hello"}
	return result, nil
}

func (r Model) EditOne(ctx context.Context, url *UrlData) (*UrlData, error) {
	r.logs.Log("msg", "Edit One")

	var result UrlData
	filter := bson.D{{"shorturl", url.ShortUrl}}
	connect := r.mongo.Database("shorturl")

	urlShorted := connect.Collection("urls_shorted")
	update := bson.D{
		{"$set", bson.D{
			{"sourceurl", url.SourceUrl},
		}},
	}
	//err := urlShorted.FindOne(context.TODO(), filter).Decode(&result)
	_, err := urlShorted.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		r.logs.Log("err", err)
		return &result, nil
	}

	return &result, nil
}
func (r Model) DeleteOne(ctx context.Context, url *UrlData) (string, error) {
	r.logs.Log("msg", "Delete One")

	connect := r.mongo.Database("shorturl")

	filter := bson.D{{"shorturl", url.ShortUrl}}
	filterMany := bson.D{{"uuid", url.ShortUrl}}

	// delete url_shorted
	res, err := connect.Collection("urls_shorted").DeleteOne(ctx, filter)

	if err != nil {
		r.logs.Log("err", err)
		return "err", nil
	}
	r.logs.Log("msg", fmt.Sprintf("DeleteOne Result TYPE: %v", res.DeletedCount))

	//delete logs for url
	resm, errm := connect.Collection("urls_shorted_log").DeleteMany(ctx, filterMany)

	if errm != nil {
		return "err", errm
	}
	r.logs.Log("msg", fmt.Sprintf("DeleteMany Result TYPE: %v", resm.DeletedCount))

	return "success", nil
}
func (r Model) GetStats(ctx context.Context, uuid, user string) (*ResultStatUrlData, error) {
	r.logs.Log("msg", "Get Stat")

	//var urlShorted *ResultStatUrlData
	var ShortUrl *UrlData
	var err error

	r.logs.Log("msg", uuid)

	ShortUrl, err = r.GetOneUrl(ctx, uuid, user)
	if err != nil {
		r.logs.Log("err", err)
		return nil, err
	}
	r.logs.Log("msg", "CHECK STAT")
	urlShorted := &ResultStatUrlData{Shorturl: ShortUrl}
	//urlShorted.Shorturl = ShortUrl

	filter := bson.D{{"uuid", uuid}}
	connect := r.mongo.Database("shorturl")

	collect := connect.Collection("urls_shorted_log")
	cur, err := collect.Find(context.TODO(), filter)

	if err != nil {
		r.logs.Log("err", err)
		result := &ResultStatUrlData{}
		return result, nil
	}

	for cur.Next(context.TODO()) {
		statist := UrlLog{}
		err := cur.Decode(&statist)
		if err != nil {
			r.logs.Log("err", err)
			result := &ResultStatUrlData{}
			return result, nil
		}
		urlShorted.Statist = append(urlShorted.Statist, statist)
	}

	return urlShorted, nil
}

func (r Model) SetLog(ctx context.Context, uuid *FindUrlDataRequest) (string, error) {
	collection := r.mongo.Database("shorturl").Collection("urls_shorted_log")

	_, err := collection.InsertOne(context.TODO(), uuid)
	if err != nil {
		r.logs.Log("err", err)
		return "error", err
	}
	return "yes", nil
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
func GetMD5HashShort(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[2:7])
}
