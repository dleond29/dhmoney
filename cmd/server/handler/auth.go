package handler

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/leorodriguez/grupo-04/internal/accounts"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
	"gitlab.com/leorodriguez/grupo-04/internal/users"
	"gitlab.com/leorodriguez/grupo-04/pkg/logger"
	"gitlab.com/leorodriguez/grupo-04/pkg/web"
)

type AuthHandler struct {
	usersService    users.Service
	accountsService accounts.Service
}

func NewAuthHandler(service users.Service, accountsService accounts.Service) AuthHandler {
	return AuthHandler{
		usersService:    service,
		accountsService: accountsService,
	}
}

// Authorization godoc
// @Summary      Create a new users
// @Description  New users registration
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        RegisterRequest   body  domain.RegisterRequest  true  "RegisterRequest"
// @Success      200  {object}  users.UserDto
// @Failure      400  {string} string  "Bad json, Requiered fields or Email already registered"
// @Failure      500  {string} string  "Internal error"
// @Router       /users/ [post]
func (ah *AuthHandler) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var rq domain.RegisterRequest
		if err := ctx.ShouldBindJSON(&rq); err != nil {
			logger.Error(err.Error())
			web.Error(ctx, http.StatusBadRequest, "Bad json")

			return
		}

		rqRefValue := reflect.ValueOf(rq)
		rqRefType := reflect.TypeOf(rq)
		var valuesNil []string
		for i := 0; i < rqRefValue.NumField(); i++ {
			if e := rqRefValue.Field(i); e.IsZero() {
				valuesNil = append(valuesNil, rqRefType.Field(i).Tag.Get("json"))
			}
		}

		if len(valuesNil) > 0 {
			web.Error(ctx, http.StatusBadRequest, "Required fields: %s", strings.Join(valuesNil, ", "))
			return
		}

		user, err := ah.accountsService.Register(ctx, rq)
		if err != nil {
			logger.Error(err.Error())
			switch err {
			case users.ErrEmailAlreadyRegistered:
				web.Error(ctx, http.StatusBadRequest, "Email already registered")
			default:
				web.Error(ctx, http.StatusInternalServerError, "Internal error")
			}

			return
		}

		web.Response(ctx, http.StatusOK, user)
		return
	}
}

// Authorization godoc
// @Summary      Login session
// @Description  user login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        LoginRequest   body   domain.LoginRequest  true  "LoginRequest"
// @Success      200  {object} domain.LoginResponse
// @Failure      400  {string} string  "All fields are required or Invalid user credentials"
// @Failure      401  {string} string  "Email not verified"
// @Failure      404  {string} string  "User not exists"
// @Failure      500  {string} string  "Internal error"
// @Router       /users/login [post]
func (ah *AuthHandler) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var rq domain.LoginRequest
		if err := ctx.ShouldBindJSON(&rq); err != nil {
			logger.Error(err.Error())
			web.Error(ctx, http.StatusBadRequest, "Internal error")

			return
		}

		if rq.Email == "" || rq.Password == "" {
			//logger.Error("Nil values in login request")
			web.Error(ctx, http.StatusBadRequest, "All fields are required")

			return
		}

		response, err := ah.usersService.Login(ctx, rq)
		if err != nil {
			logger.Error(err.Error())
			switch err {
			case users.ErrInvalidUserCredentials:
				web.Error(ctx, http.StatusBadRequest, "Invalid user credentials")
			case users.ErrUserNotExists:
				web.Error(ctx, http.StatusNotFound, "User not exists")
			case users.ErrEmailNotVerified:
				web.Error(ctx, http.StatusUnauthorized, "Email not verified")

			default:
				web.Error(ctx, http.StatusInternalServerError, "Internal error")
			}

			return
		}

		web.Response(ctx, http.StatusOK, response)
		return
	}
}

// Authorization godoc
// @Summary      Logout session
// @Description  user logout
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization  header   string  true  "Authorization"
// @Success      200  {string} string  "ok"
// @Failure      500  {string} string  "Internal error"
// @Router       /users/logout [get]
func (ah *AuthHandler) Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, ok := ctx.Request.Header["Authorization"]
		if !ok {
			web.Error(ctx, http.StatusBadRequest, "Internal error")

			return
		}

		err := ah.usersService.Logout(ctx, strings.Split(token[0], " ")[1])
		if err != nil {
			logger.Error(err.Error())
			web.Error(ctx, http.StatusInternalServerError, "Internal error")

			return
		}

		web.Response(ctx, http.StatusOK, "ok")
		return
	}
}

// Authorization godoc
// @Summary      Recover password through email
// @Description  user forgot credentials
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        ForgotRequest   body   domain.ForgotRequest  true  "ForgotRequest"
// @Success      200  {string} string  "ok"
// @Failure      400  {string} string  "Bad request"
// @Failure      500  {string} string  "Internal error"
// @Router       /users/forgot [post]
func (ah *AuthHandler) ForgotPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var rq domain.ForgotRequest
		if err := ctx.ShouldBindJSON(&rq); err != nil {
			logger.Error(err.Error())
			web.Error(ctx, http.StatusBadRequest, "Bad request")

			return
		}

		err := ah.usersService.ForgotPassword(ctx, rq.Email)
		if err != nil {
			logger.Error(err.Error())
			web.Error(ctx, http.StatusInternalServerError, "Internal error")

			return
		}

		web.Response(ctx, http.StatusOK, "ok")
		return
	}
}

func (ah *AuthHandler) HasToken(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		web.Error(ctx, http.StatusUnauthorized, "Token not sent")

		ctx.Abort()
		return
	}

	if split := strings.Split(token, " "); len(split) != 2 || split[0] != "Bearer" {
		web.Error(ctx, http.StatusUnauthorized, "Bad token")

		ctx.Abort()
		return
	}

	ctx.Next()
}
