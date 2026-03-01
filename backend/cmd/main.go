package main

import (
	"log"
	"net/http"
	"song-match-backend/api/route"
	"song-match-backend/bootstrap"
	"time"
)

func main() {
	app := bootstrap.App()
	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	srv := &http.Server{
		Addr:    env.ServerAddress,
		Handler: route.Setup(env, timeout, db),
	}

	log.Printf("Starting HTTP server on port %s\n", env.ServerAddress)

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
