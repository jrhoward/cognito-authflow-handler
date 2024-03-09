package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jrhoward/cognito-authflow-handler/auth"
	"github.com/jrhoward/cognito-authflow-handler/config"
)

func main() {
	logInfo := log.New(os.Stdout, "INFO ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
	logError := log.New(os.Stdout, "ERROR ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
	if len(os.Args) != 2 {
		logError.Fatal("Usage: cognito-auth-flow-handler /path/to/config.yaml")
	}
	configPath := os.Args[1]
	_, err := os.Stat(configPath)
	if err != nil {
		logError.Fatal(err)
	}
	err = config.Init(configPath)
	if err != nil {
		logError.Fatal(err)
	}
	err = auth.Init(configPath)
	if err != nil {
		logError.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/auth", auth.AuthWrapper).Methods("GET")
	router.HandleFunc("/logout", auth.LogoutHandler).Methods("GET")
	listening := config.GetServerHost()
	logInfo.Println("starting the server at " + listening)
	err = http.ListenAndServe(listening, router)
	if err != nil {
		logError.Fatal("could not start the server", err)
	}
}
