package endpoint

import (
	"context"
	"fmt"

	//"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	//"github.com/go-kit/kit/auth/jwt"
	"github.com/Traestan/shorturl/repository"
	"github.com/go-kit/kit/endpoint"
	"github.com/skip2/go-qrcode"
)

// Endpoints holds all Go kit endpoints for the Order service.
type Endpoints struct {
	Register   endpoint.Endpoint
	Login      endpoint.Endpoint
	Forgot     endpoint.Endpoint
	Changepass endpoint.Endpoint
	//Work url
	Create    endpoint.Endpoint
	Index     endpoint.Endpoint
	GetUrl    endpoint.Endpoint
	EditUrl   endpoint.Endpoint
	DeleteUrl endpoint.Endpoint
	GetStat   endpoint.Endpoint

	//Work ShortUrl
	Root      endpoint.Endpoint
	GerQrcode endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the Order service.
func MakeEndpoints(s *repository.Model) *Endpoints {
	return &Endpoints{
		Register:   makeRegisterEndpoint(s),
		Login:      makeLoginEndpoint(s),
		Forgot:     makeForgotEndpoint(s),
		Changepass: ChangepassEndpoint(s),

		Index:     makeIndexEndpoint(s),
		Create:    makeCreateEndpoint(s),
		GetUrl:    makeGetUrlEndpoint(s),
		EditUrl:   makeEditUrlEndpoint(s),
		DeleteUrl: makeDeleteUrlEndpoint(s),
		GetStat:   makeGetStatEndpoint(s),

		//Work Shorturl
		Root:      makeRootEndpoint(s),
		GerQrcode: makeGetQrcodeEndpoint(s),
	}
}

//register
func makeRegisterEndpoint(s *repository.Model) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*repository.User) // type assertion

		id, err := s.Registration(ctx, req)
		if err != nil {
			return nil, err
		}
		if req.Email != "" {
			//helper.SendEmail(user.Email, "Login", "Yare login")
			MailSendRegister(req, "Регистрация", "Спасибо за регистрацию.")
		}
		return id, nil
	}
}

//login
func makeLoginEndpoint(s *repository.Model) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(*repository.User) // type assertion
		user, err := s.Login(ctx, req)

		if err != nil {
			return nil, err
		}
		expirationTime := time.Now().Add(200 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time
		claims := &repository.Claims{
			Username: user.Email,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}

		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Create the JWT string
		tokenString, err := token.SignedString(repository.JwtKey)
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			return nil, err
		}

		// Finally, we set the client cookie for "token" as the JWT we just generated
		// we also set an expiry time which is the same as the token itself
		/*http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
			Domain:  ".localhost",
			Path:    "/",
		})*/

		user.Token = tokenString
		//fmt.Print(user.Email != "")
		if user.Email != "" {
			//helper.SendEmail(user.Email, "Login", "Yare login")
			MailSend(user.Email, "Login", "Yare login")
		}

		return user, nil
	}
}
func makeForgotEndpoint(s *repository.Model) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(*repository.User) // type assertion
		user, err := s.Forgot(ctx, req)

		if err != nil {
			return nil, err
		}
		if user.Email != "" {
			MailSendForgot(user, "Forgot", user.Password)
		}
		user.Password = ""
		return user, nil
	}
}

func ChangepassEndpoint(s *repository.Model) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		cnts := ctx.Value("Username")
		username := "null"
		if cnts != nil {
			//fmt.Print(cnts)
			username = cnts.(jwt.MapClaims)["username"].(string)
		}
		req := request.(*repository.User) // type assertion
		req.Email = username
		ResultFind, err := s.Changepass(ctx, req)
		return ResultFind, err
	}
}

//user work
//add
func makeCreateEndpoint(s *repository.Model) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*repository.UrlData) // type assertion
		cnts := ctx.Value("Username")
		username := "null"
		if cnts != nil {
			fmt.Print(cnts)
			username = cnts.(jwt.MapClaims)["username"].(string)
		} else {
			fmt.Print("usernma3e")
		}
		id, err := s.Create(ctx, req, username)
		if err != nil {
			return nil, err
		}

		return repository.CreateResponse{Hurl: id}, nil
	}
}

//index
func makeIndexEndpoint(s *repository.Model) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		cnts := ctx.Value("Username")
		username := "null"
		if cnts != nil {
			//fmt.Print(cnts)
			username = cnts.(jwt.MapClaims)["username"].(string)
		} else {
			fmt.Print(cnts)
		}

		//req := request.(*repository.IndexRequest) // type assertion
		ResultFind, err := s.GetUrls(ctx, username)
		return repository.ResponseFindUrls{Result: ResultFind.Result, Err: err}, err
	}
}

//get
func makeGetUrlEndpoint(s *repository.Model) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(*repository.FindUrlDataRequest) // type assertion
		ResultFind, err := s.GetOne(ctx, req.Uuid)

		return repository.ResponseFindUrl{Result: ResultFind, Err: err}, nil
	}
}

//edit
func makeEditUrlEndpoint(s *repository.Model) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(*repository.UrlData) // type assertion
		ResultFind, err := s.EditOne(ctx, req)

		return repository.ResponseFindUrl{Result: ResultFind, Err: err}, nil
	}
}

//edit

func makeDeleteUrlEndpoint(s *repository.Model) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(*repository.UrlData) // type assertion
		ResultFind, err := s.DeleteOne(ctx, req)
		if err != nil {
			return nil, err
		}
		return ResultFind, nil
	}
}

//stat
func makeGetStatEndpoint(s *repository.Model) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//fmt.Print("STAT")
		cnts := ctx.Value("Username")
		username := "null"
		if cnts != nil {
			username = cnts.(jwt.MapClaims)["username"].(string)
		}

		req := request.(*repository.FindUrlDataRequest) // type assertion
		ResultFind, err := s.GetStats(ctx, req.Uuid, username)
		return ResultFind, err
	}
}

//work shorturl
//Переходи по sourceurl
func makeRootEndpoint(s *repository.Model) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*repository.FindUrlDataRequest) // type assertion
		//fmt.Print(req.Uuid)
		ResultFind, err := s.GetOne(ctx, req.Uuid)
		if err != nil {
			return "", err
		}
		s.SetLog(ctx, req)
		return ResultFind.SourceUrl, nil
	}
}

//возвращаем qrcode
func makeGetQrcodeEndpoint(s *repository.Model) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//flag.Parse()
		req := request.(*repository.FindUrlDataRequest) // type assertion
		ResultFind, err := s.GetOne(ctx, req.Uuid)

		s.SetLog(ctx, req)

		if err != nil {
			return nil, err
		}
		var png []byte
		png, _ = qrcode.Encode("http://localhost/"+ResultFind.ShortUrl, qrcode.Medium, 256)

		return png, nil
	}
}
