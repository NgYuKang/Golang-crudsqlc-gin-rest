package main

import (
	"context"
	"fmt"
	"golang-crudsqlc-gin-rest/controllers"
	db "golang-crudsqlc-gin-rest/db/sqlc"
	"golang-crudsqlc-gin-rest/routes"
	"golang-crudsqlc-gin-rest/util"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(config.DbSource)
	conn, err := pgx.Connect(context.Background(), config.DbSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	pdb := db.New(conn)
	contactController := controllers.NewContactController(pdb, context.Background())
	contactRouter := routes.NewContactRoute(*contactController)

	server := gin.Default()
	router := server.Group("/api")

	contactRouter.ContactRouter(router)

	log.Fatal(server.Run(":8080"))

}
