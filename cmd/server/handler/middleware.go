package handler

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/leorodriguez/grupo-04/internal/accounts"
	"gitlab.com/leorodriguez/grupo-04/pkg/logger"
	"gitlab.com/leorodriguez/grupo-04/pkg/web"
	"net/http"
	"strconv"
)

type Middlewares struct {
	accountsService accounts.Service
}

func NewMiddlewares(accountsService accounts.Service) Middlewares {
	return Middlewares{
		accountsService: accountsService,
	}
}

func (m *Middlewares) IsAuthorized(ctx *gin.Context) {
	idParam := ctx.Param("accountID")
	var isUserID bool
	if idParam == "" {
		idParam = ctx.Param("userID")
		isUserID = true
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		web.Error(ctx, http.StatusBadRequest, "invalid id")
		ctx.Abort()
		return
	}

	token := ctx.GetHeader("Authorization")
	if token == "" {
		web.Error(ctx, http.StatusBadRequest, "Token not sent")
		ctx.Abort()
		return
	}
	isAuthorized, err := m.accountsService.IsAuthorized(ctx, id, isUserID, token)
	if err != nil {
		switch err {
		case accounts.ErrTokenExpired:
			web.Error(ctx, http.StatusUnauthorized, "Session expired. Please login again")
		default:
			web.Error(ctx, http.StatusInternalServerError, "Internal error")
		}

		logger.Error(err.Error())
		ctx.Abort()
		return
	}

	if !isAuthorized {
		web.Error(ctx, http.StatusForbidden, "Not authorized")
		ctx.Abort()
		return
	}

	ctx.Next()
}
