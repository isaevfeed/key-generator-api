package main

import (
	"fmt"
	serv "isaevfeed/internal/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Generator struct{}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env file isn't found")
	}
}

func main() {
	host, _ := os.LookupEnv("SERVER_HOST")
	port, _ := os.LookupEnv("SERVER_PORT")

	srv := serv.New(fmt.Sprintf("%s:%s", host, port))
	srv.Listen()
}
