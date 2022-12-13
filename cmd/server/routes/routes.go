package routes

import (
	"database/sql"
	"github.com/Nerzal/gocloak/v12"
	"github.com/gin-gonic/gin"
	"gitlab.com/leorodriguez/grupo-04/cmd/server/handler"
	"gitlab.com/leorodriguez/grupo-04/docs"
	"gitlab.com/leorodriguez/grupo-04/internal/accounts"
	"gitlab.com/leorodriguez/grupo-04/internal/auth"
	"gitlab.com/leorodriguez/grupo-04/internal/cards"
	"gitlab.com/leorodriguez/grupo-04/internal/transactions"
	"gitlab.com/leorodriguez/grupo-04/internal/users"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router interface {
	MapRoutes()
}

type router struct {
	r  *gin.Engine
	rg *gin.RouterGroup
	db *sql.DB

	aliasWords []string
}

func NewRouter(r *gin.Engine, db *sql.DB, aliasWords []string) Router {
	return &router{r: r, db: db, aliasWords: aliasWords}
}

func (r *router) MapRoutes() {
	keycloakSettings := auth.KeycloakSettings{
		GoCloak:      gocloak.NewClient(os.Getenv("KEYCLOAK_URL")),
		ClientId:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		Realm:        os.Getenv("KEYCLOAK_REALM"),
	}

	keycloakService := auth.NewAuth(keycloakSettings)
	authRepository := users.NewRepository(r.db)
	accountsRepository := accounts.NewRepository(r.db)
	transactionsRepository := transactions.NewRepository(r.db)
	cardsRepository := cards.NewRepository(r.db)

	authService := users.NewUsers(keycloakService, authRepository, r.aliasWords)
	accountsService := accounts.NewService(authService, accountsRepository, transactionsRepository, keycloakService, r.aliasWords)
	cardService := cards.NewService(cardsRepository)

	authHandler := handler.NewAuthHandler(authService, accountsService)
	accountsHandler := handler.NewAccountsHandler(accountsService)
	cardsHandler := handler.NewCardHandler(cardService, accountsService)
	middlewares := handler.NewMiddlewares(accountsService)

	r.rg = r.r.Group("/api")

	accountsGroup := r.rg.Group("/accounts")
	accountsGroup.GET("/:accountID", middlewares.IsAuthorized, accountsHandler.GetAccount)
	accountsGroup.PATCH("/:accountID", middlewares.IsAuthorized, accountsHandler.ChangeAlias())
	accountsGroup.GET("/:accountID/transactions", middlewares.IsAuthorized, accountsHandler.GetTransactionsLastFive)

	cardsGroup := r.rg.Group("/accounts")
	cardsGroup.POST("/:accountID/cards", middlewares.IsAuthorized, cardsHandler.NewCard())
	cardsGroup.GET("/:accountID/cards", middlewares.IsAuthorized, cardsHandler.GetAll)
	cardsGroup.GET("/:accountID/cards/:cardID", middlewares.IsAuthorized, cardsHandler.GetByCardID)
	cardsGroup.DELETE("/:accountID/cards/:cardID", middlewares.IsAuthorized, cardsHandler.DeleteByCardID)

	usersGroup := r.rg.Group("/users")
	usersGroup.POST("/", authHandler.Register())
	usersGroup.GET("/:userID", middlewares.IsAuthorized, accountsHandler.GetUser)
	usersGroup.PATCH("/:accountID", middlewares.IsAuthorized, accountsHandler.UpdateAccount())
	usersGroup.POST("/login", authHandler.Login())
	usersGroup.GET("/logout", authHandler.HasToken, authHandler.Logout())
	usersGroup.POST("/forgot", authHandler.ForgotPassword())

	docs.SwaggerInfo.Host = "localhost:8080"
	r.rg.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
