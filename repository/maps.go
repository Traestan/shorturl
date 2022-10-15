package repository

import "github.com/dgrijalva/jwt-go"

//request
type FindAllRequest struct {
	Query string
}
type FindUrlDataRequest struct {
	Uuid      string
	UHeader   string
	UIp       string
	Ureferer  string
	UAHeaders interface{}
	UDateTime string
}
type CreateRequest struct {
	SUrl string `json:"surl"`
	Ip   string `json:"ip"`
}
type IndexRequest struct {
	User string `json:"user"`
}

//response
type ResponseUserLogin struct {
	Result interface{}
	Err    error
}

type CreateResponse struct {
	Hurl interface{} `json:"hurl"`
}

type IndexResponse struct {
	Docs *PageContent
}
type ResponseFindUrls struct {
	Result []ResultUrlData
	Err    error
}
type ResponseFindUrl struct {
	Result *UrlData
	Err    error
}
type ResponseStatFindUrl struct {
	Result *ResultStatUrlData
	Err    error
}

//types
type CreateUrlData struct {
	SourceUrl string `json:"surl"`
	Ip        string `json:"ip"`
}

//User
/*type UserRegistrationData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Date_Add string `json:"date_add"`
}*/

/*
type UserData struct {
	ID    string `json:"_id,omitempty"`
	Email string `json:"email,omitempty"`
	//Password string      `json:"password"`
	Date_Add string `json:"date_add,omitempty"`
	Token    string `json:"token,omitempty"`
}*/

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type ResultUrlData struct {
	ID        interface{}
	ShortUrl  string
	SourceUrl string
	Ip        string
	Uuid      string
	Date_Add  string
}
type ResultUrlDataAll struct {
	Result []ResultUrlData
}

type ResultStatUrlData struct {
	Shorturl interface{}
	Statist  []UrlLog
}

type PageContent struct {
	Content string
}
type UrlLog struct {
	Uip       string
	Ureferer  string
	Uaheaders interface{}
}

type UrlData struct {
	ID        interface{}
	SourceUrl string `json:"slug,omitempty"`
	ShortUrl  string `json:"hurl,omitempty"`
	Ip        string `json:"ip,omitempty"`
	Uuid      string `json:"uuid,omitempty"`
	Date_Add  string `json:"date_add,omitempty"`
}

//stinger
/*func (p UrlData) String() string {
	return fmt.Sprintf("source -%v (%s ShortUrl) uuid - %s", p.SourceUrl, p.ShortUrl, p.Uuid)
}*/

type User struct {
	//Username       string `json:"email,omitempty"`
	ID             string `json:"_id,omitempty"`
	Email          string `json:"email,omitempty"`
	Password       string `json:"password"`
	NewPassword    string `json:"newpass,omitempty"`
	NewPassCompare string `json:"newpasscompare,omitempty"`
	Date_Add       string `json:"date_add,omitempty"`
	Token          string `json:"token,omitempty"`
}
