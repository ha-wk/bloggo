package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/ha-wk/bloggo/apis/models"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to Database %s", Dbdriver)
			log.Fatal("Error is", err)
		} else {
			fmt.Printf("Connected to database %s", Dbdriver)
		}
	}

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{}) //LOGIC FOR DATA MIGRATION
	server.Router = mux.NewRouter()
	server.initializeRoutes()

}
func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8090")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
