package main

import (
	"database/sql"
	"fmt"
	"gitlab.com/leorodriguez/grupo-04/pkg/logger"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gitlab.com/leorodriguez/grupo-04/cmd/server/routes"
)

// @title           Grupo 4 Swagger
// @version         1.0
// @description     Project develop for group 4 CTD backend specialist.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Grupo 4
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api
func main() {
	logger.Init()
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}

	aliasWordsFilePath := os.Getenv("ALIAS_WORDS_FILE_PATH")
	aliasWordsRaw, err := os.ReadFile(aliasWordsFilePath)
	if err != nil {
		panic(err)

	}
	aliasWords := strings.Split(string(aliasWordsRaw), "\n")

	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbName)

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	router := routes.NewRouter(r, db, aliasWords)
	router.MapRoutes()

	if err = r.Run(); err != nil {
		panic(err)
	}

}
