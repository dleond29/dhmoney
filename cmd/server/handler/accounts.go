package handler

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/leorodriguez/grupo-04/internal/accounts"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
	"gitlab.com/leorodriguez/grupo-04/internal/users"
	"gitlab.com/leorodriguez/grupo-04/pkg/logger"
	"gitlab.com/leorodriguez/grupo-04/pkg/web"
)

type AccountsHandler struct {
	service accounts.Service
}

func NewAccountsHandler(service accounts.Service) *AccountsHandler {
	return &AccountsHandler{service: service}
}

// Accounts godoc
// @Summary      Get account info
// @Description  Get account info
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        Authorization  header   string  true  "Authorization"
// @Param        accountID   path   int   true  "accountID"
// @Success      200  {object}  domain.AccountInfo
// @Failure      500  {string} string  "Internal error"
// @Router       /accounts/{accountID} [get]
func (t *AccountsHandler) GetAccount(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	idParam := ctx.Param("accountID")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		web.Error(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	account, err := t.service.GetAccountInfo(ctx, id, token)
	if err != nil {

		web.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	web.Response(ctx, http.StatusOK, account)
}

// Users godoc
// @Summary      Get users info
// @Description  Get users info
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userID   path   int   true  "userID"
// @Success      200  {object}  domain.UserInfo
// @Failure      400  {string} string  "invalid id"
// @Failure      500  {string} string  "Internal error"
// @Router       /users/{userID} [get]
func (t *AccountsHandler) GetUser(ctx *gin.Context) {
	idParam := ctx.Param("userID")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		web.Error(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	account, err := t.service.GetUserInfo(ctx, id)
	if err != nil {
		web.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	web.Response(ctx, http.StatusOK, account)
}

// Transactions  godoc
// @Summary      Get last five transactions info
// @Description  Get last five transactions info
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        Authorization  header   string  true  "Authorization"
// @Param        accountID   path   int   true  "accountID"
// @Success      200  {object}  []domain.TransactionInfo
// @Failure      400  {string} string  "invalid id"
// @Failure      500  {string} string  "Internal error"
// @Router       /accounts/{accountID}/transactions [get]
func (t *AccountsHandler) GetTransactionsLastFive(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	idParam := ctx.Param("accountID")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		web.Error(ctx, http.StatusBadRequest, "invalid id")
		return
	}

	trxs, err := t.service.GetTransactionsLastFive(ctx, id, token)
	if err != nil {

		web.Error(ctx, http.StatusInternalServerError, "Internal error")
		return
	}

	if trxs == nil {
		web.Response(ctx, http.StatusOK, []domain.Transaction{})
		return
	}

	web.Response(ctx, http.StatusOK, trxs)
}

// Accounts godoc
// @Summary      Change alias account info
// @Description  Change alias account info
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        accountID   path   int   true  "accountID"
// @Success      200  {string} string  "Ok"
// @Failure      400  {string} string  "invalid id, Bad json, Required field or Alias already in use"
// @Failure      500  {string} string  "Internal error"
// @Router       /accounts/{accountID} [patch]
func (t *AccountsHandler) ChangeAlias() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("accountID")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			logger.Error(err.Error())
			web.Error(ctx, http.StatusBadRequest, "invalid id")
			return
		}

		var rq domain.Alias
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

		err = t.service.UpdateAlias(ctx, id, rq.Alias)
		if err != nil {
			logger.Error(err.Error())
			switch err {
			case accounts.ErrAliasAlreadyExists:
				web.Error(ctx, http.StatusBadRequest, "Alias already in use")
			default:
				web.Error(ctx, http.StatusInternalServerError, "Internal error")
			}

			return
		}

		web.Response(ctx, http.StatusOK, "OK")
		return
	}
}

// Users godoc
// @Summary      Update user info
// @Description  Update user info
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        accountID   path   int   true  "accountID"
// @Success      200  {string} string  "Ok"
// @Failure      400  {string} string  "invalid id, Bad json, Email already registered"
// @Failure      500  {string} string  "Internal error"
// @Router       /users/{accountID} [patch]
func (t *AccountsHandler) UpdateAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("accountID")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Error(ctx, http.StatusBadRequest, "invalid id")
			return
		}

		var rq domain.RegisterRequest
		if err := ctx.ShouldBindJSON(&rq); err != nil {
			logger.Error(err.Error())
			web.Error(ctx, http.StatusBadRequest, "Bad json")

			return
		}

		_, err = t.service.UpdateAccount(ctx, rq, id)
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

		web.Response(ctx, http.StatusOK, "OK")
		return
	}
}
