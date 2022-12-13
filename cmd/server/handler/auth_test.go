package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
	"gitlab.com/leorodriguez/grupo-04/internal/users"
	"net/http"
	"net/http/httptest"
	"testing"
)

type usersMock struct {
	mock.Mock
	users.Service
}

func (u *usersMock) Register(ctx context.Context, rq domain.RegisterRequest) (users.UserDto, error) {
	args := u.Called(rq)
	return args.Get(0).(users.UserDto), args.Error(1)
}

func (u *usersMock) Login(ctx context.Context, rq domain.LoginRequest) (domain.LoginResponse, error) {
	args := u.Called(rq)
	return args.Get(0).(domain.LoginResponse), args.Error(1)
}

func (u *usersMock) Logout(ctx context.Context, token string) error {
	args := u.Called(token)
	return args.Error(0)
}

func (u *usersMock) ForgotPassword(ctx context.Context, email string) error {
	args := u.Called(email)
	return args.Error(0)
}

func createRequest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")

	return req, httptest.NewRecorder()
}

func TestRegister(t *testing.T) {
	var tests = []struct {
		name           string
		body           string
		handler        AuthHandler
		responseStatus int
		responseBody   string
	}{
		{
			name: "ok",
			body: `{
				"name": "name",
				"last_name": "lastname",
				"dni": 21312731,
				"phone": 123123,
				"email": "email@gmail.com",
				"password": "password"
			}`,
			handler: func() AuthHandler {
				serviceMock := new(usersMock)

				rq := domain.RegisterRequest{
					Name:     "name",
					LastName: "lastname",
					DNI:      21312731,
					Phone:    123123,
					Email:    "email@gmail.com",
					Password: "password",
				}
				user := users.UserDto{
					ID:       1,
					Name:     "name",
					LastName: "lastname",
					DNI:      21312731,
					Phone:    123123,
					Email:    "email@gmail.com",
					CVU:      "1234567890123456789012",
					Alias:    "casa.perro.hola",
				}
				serviceMock.On("Register", rq).
					Return(user, nil)

				handler := NewAuthHandler(serviceMock)
				return handler
			}(),
			responseStatus: http.StatusOK,
			responseBody:   `{"id":1,"name":"name","last_name":"lastname","dni":21312731,"phone":123123,"email":"email@gmail.com","cvu":"1234567890123456789012","alias":"casa.perro.hola"}`,
		},
		{
			name: "bad body request",
			body: `{
				"name": name,
			}`,
			handler: func() AuthHandler {
				serviceMock := new(usersMock)
				handler := NewAuthHandler(serviceMock)
				return handler
			}(),
			responseStatus: http.StatusBadRequest,
			responseBody:   `{"code":"bad_request","message":"Bad json"}`,
		},
		{
			name: "missing required fields",
			body: `{
				"name": "name",
				"last_name": "lastname",
				"dni": 21312731
			}`,
			handler: func() AuthHandler {
				serviceMock := new(usersMock)
				handler := NewAuthHandler(serviceMock)
				return handler
			}(),
			responseStatus: http.StatusBadRequest,
			responseBody:   `{"code":"bad_request","message":"Required fields: phone, email, password"}`,
		},
		{
			name: "email already registered",
			body: `{
				"name": "name",
				"last_name": "lastname",
				"dni": 21312731,
				"phone": 123123,
				"email": "email@gmail.com",
				"password": "password"
			}`,
			handler: func() AuthHandler {
				serviceMock := new(usersMock)

				rq := domain.RegisterRequest{
					Name:     "name",
					LastName: "lastname",
					DNI:      21312731,
					Phone:    123123,
					Email:    "email@gmail.com",
					Password: "password",
				}
				serviceMock.On("Register", rq).
					Return(users.UserDto{}, users.ErrEmailAlreadyRegistered)

				handler := NewAuthHandler(serviceMock)
				return handler
			}(),
			responseStatus: http.StatusBadRequest,
			responseBody:   `{"code":"bad_request","message":"Email already registered"}`,
		},
		{
			name: "internal service error",
			body: `{
				"name": "name",
				"last_name": "lastname",
				"dni": 21312731,
				"phone": 123123,
				"email": "email@gmail.com",
				"password": "password"
			}`,
			handler: func() AuthHandler {
				serviceMock := new(usersMock)

				rq := domain.RegisterRequest{
					Name:     "name",
					LastName: "lastname",
					DNI:      21312731,
					Phone:    123123,
					Email:    "email@gmail.com",
					Password: "password",
				}
				serviceMock.On("Register", rq).
					Return(users.UserDto{}, errors.New("internal error"))

				handler := NewAuthHandler(serviceMock)
				return handler
			}(),
			responseStatus: http.StatusInternalServerError,
			responseBody:   `{"code":"internal_server_error","message":"Internal error"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := gin.Default()
			rg := r.Group("/api")
			rg.POST("/users", tt.handler.Register())

			req, rr := createRequest(http.MethodPost, "http://localhost:8080/api/users", tt.body)

			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.responseStatus, rr.Code)
			assert.Equal(t, tt.responseBody, rr.Body.String())
		})
	}
}

func TestLogin(t *testing.T) {
	var tests = []struct {
		name           string
		body           string
		handler        AuthHandler
		responseStatus int
		responseBody   string
	}{
		{
			name: "login - ok",
			body: `{
				"email": "email@gmail.com",
				"password": "password"
			}`,
			handler: func() AuthHandler {
				serviceMock := new(usersMock)

				rq := domain.LoginRequest{
					Email:    "email@gmail.com",
					Password: "password",
				}
				response := domain.LoginResponse{
					Token: "asdasd",
				}
				serviceMock.On("Login", rq).
					Return(response, nil)

				handler := NewAuthHandler(serviceMock)
				return handler
			}(),
			responseStatus: http.StatusOK,
			responseBody:   `{"token":"asdasd"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := gin.Default()
			rg := r.Group("/api")
			rg.POST("/users/login", tt.handler.Login())

			req, rr := createRequest(http.MethodPost, "http://localhost:8080/api/users/login", tt.body)

			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.responseStatus, rr.Code)
			assert.Equal(t, tt.responseBody, rr.Body.String())
		})
	}
}

func TestLogout(t *testing.T) {
	var tests = []struct {
		name           string
		token          string
		handler        AuthHandler
		responseStatus int
		responseBody   string
	}{
		{
			name:  "login - ok",
			token: "Bearer tokenAsd",
			handler: func() AuthHandler {
				serviceMock := new(usersMock)

				serviceMock.On("Logout", "tokenAsd").
					Return(nil)

				handler := NewAuthHandler(serviceMock)
				return handler
			}(),
			responseStatus: http.StatusOK,
			responseBody:   `"ok"`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := gin.Default()
			rg := r.Group("/api")
			rg.GET("/users/logout", tt.handler.HasToken, tt.handler.Logout())

			req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/users/logout", nil)
			req.Header.Add("Authorization", tt.token)

			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.responseStatus, rr.Code)
			assert.Equal(t, tt.responseBody, rr.Body.String())
		})
	}
}

func TestForgotPassword(t *testing.T) {
	var tests = []struct {
		name           string
		body           string
		handler        AuthHandler
		responseStatus int
		responseBody   string
	}{
		{
			name: "login - ok",
			body: `{
				"email": "email@c.com"
			}`,
			handler: func() AuthHandler {
				serviceMock := new(usersMock)

				serviceMock.On("ForgotPassword", "email@c.com").
					Return(nil)

				handler := NewAuthHandler(serviceMock)
				return handler
			}(),
			responseStatus: http.StatusOK,
			responseBody:   `"ok"`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := gin.Default()
			rg := r.Group("/api")
			rg.POST("/users/forgot", tt.handler.ForgotPassword())

			req, rr := createRequest(http.MethodPost, "http://localhost:8080/api/users/forgot", tt.body)

			rr = httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.responseStatus, rr.Code)
			assert.Equal(t, tt.responseBody, rr.Body.String())
		})
	}
}
