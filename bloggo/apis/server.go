package apis

import (
	"fmt"
	"log"
	"os"

	"github.com/ha-wk/bloggo/apis/controllers"
	//"github.com/ha-wk/bloggo/apis/seed"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error in getting env,not coming through %v", err)
	} else {
		fmt.Println("WE ARE GETTING THE ENV VALUES")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	//seed.Load(server.DB)
	server.Run(":8090")
}
