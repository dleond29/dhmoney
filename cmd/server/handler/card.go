package handler

import "C"
import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/leorodriguez/grupo-04/internal/accounts"
	"gitlab.com/leorodriguez/grupo-04/internal/cards"
	"gitlab.com/leorodriguez/grupo-04/internal/domain"
	"gitlab.com/leorodriguez/grupo-04/pkg/logger"
	"gitlab.com/leorodriguez/grupo-04/pkg/web"
)

type CardHandler struct {
	cardsService    cards.Service
	accountsService accounts.Service
}

func NewCardHandler(cardsService cards.Service, accountsService accounts.Service) CardHandler {
	return CardHandler{
		cardsService:    cardsService,
		accountsService: accountsService,
	}
}

// Card godoc
// @Summary      Create a new card
// @Description  New card
// @Tags         card
// @Accept       json
// @Produce      json
// @Param        accountID   path   int   true  "accountID"
// @Param        NewCard   body  domain.CardDto  true  "NewCard"
// @Success      200  {string} string  "Ok"
// @Failure      400  {string} string  "invalid id, bad json, Required fields, Card already associated to another account"
// @Failure      500  {string} string  "Internal error"
// @Router       /accounts/{accountID}/cards [post]
func (c *CardHandler) NewCard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("accountID")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Error(ctx, http.StatusBadRequest, "invalid id")
			return
		}

		var rq domain.CardDto
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

		err = c.cardsService.Save(ctx, id, rq)
		if err != nil {
			logger.Error(err.Error())
			switch err {
			case cards.ErrCardAlreadyAssociated:
				web.Error(ctx, http.StatusBadRequest, "Card already associated to another account")
			default:
				web.Error(ctx, http.StatusInternalServerError, "Internal error")
			}

			return
		}

		web.Response(ctx, http.StatusOK, "OK")
		return
	}
}

// Card godoc
// @Summary      Get all cards
// @Description  Get all cards
// @Tags         card
// @Accept       json
// @Produce      json
// @Param        accountID   path   int   true  "accountID"
// @Success      200  {array}  domain.Card
// @Failure      400  {string} string  "invalid account id"
// @Failure      500  {string} string  "Internal error"
// @Router       /accounts/{accountID}/cards [get]
func (c *CardHandler) GetAll(ctx *gin.Context) {
	accountIDParam := ctx.Param("accountID")
	accountID, err := strconv.Atoi(accountIDParam)
	if err != nil {
		web.Error(ctx, http.StatusBadRequest, "invalid account id")
		return
	}

	cards, err := c.cardsService.GetAll(ctx, accountID)
	if err != nil {

		web.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	web.Response(ctx, http.StatusOK, cards)
}

// Card godoc
// @Summary      Get card by id
// @Description  Get card by id
// @Tags         card
// @Accept       json
// @Produce      json
// @Param        accountID   path   int   true  "accountID"
// @Param        cardID   path   int   true  "cardID"
// @Success      200  {object}  domain.Card
// @Failure      400  {string} string  "invalid account id, invalid card id"
// @Failure      404  {string} string  "Card not found"
// @Failure      500  {string} string  "Internal error"
// @Router       /accounts/{accountID}/cards/{cardID} [get]
func (c *CardHandler) GetByCardID(ctx *gin.Context) {
	accountIDParam := ctx.Param("accountID")
	accountID, err := strconv.Atoi(accountIDParam)
	if err != nil {
		web.Error(ctx, http.StatusBadRequest, "invalid account id")
		return
	}

	cardIDParam := ctx.Param("cardID")
	cardID, err := strconv.Atoi(cardIDParam)
	if err != nil {
		web.Error(ctx, http.StatusBadRequest, "invalid card id")
		return
	}

	card, err := c.cardsService.GetByCardID(ctx, accountID, cardID)
	if err != nil {
		switch err {
		case cards.ErrCardNotFound:
			web.Error(ctx, http.StatusNotFound, "Card not found")
			return
		}
		web.Error(ctx, http.StatusInternalServerError, "Internal error")
		return
	}
	web.Response(ctx, http.StatusOK, card)
}

// Card godoc
// @Summary      Delete card by id
// @Description  Delete card by id
// @Tags         card
// @Accept       json
// @Produce      json
// @Param        accountID   path   int   true  "accountID"
// @Param        cardID   path   int   true  "cardID"
// @Success      200  {string} string  "ok"
// @Failure      400  {string} string  "invalid card id"
// @Failure      404  {string} string  "Card not found"
// @Failure      500  {string} string  "Internal error"
// @Router       /accounts/{accountID}/cards/{cardID} [delete]
func (c *CardHandler) DeleteByCardID(ctx *gin.Context) {
	cardIDParam := ctx.Param("cardID")
	cardID, err := strconv.Atoi(cardIDParam)
	if err != nil {
		web.Error(ctx, http.StatusBadRequest, "invalid card id")
		return
	}

	err = c.cardsService.DeleteByCardID(ctx, cardID)
	if err != nil {
		switch err {
		case cards.ErrCardNotFound:
			web.Error(ctx, http.StatusNotFound, "Card not found")
		default:
			web.Error(ctx, http.StatusInternalServerError, "Internal error")
		}
		return
	}
	web.Response(ctx, http.StatusOK, "OK")
}
