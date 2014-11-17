// Copyright 2013-2014 Bowery, Inc.
// Contains the main entry point

package main

import (
	"os"

	"github.com/Bowery/gopackages/config"
	"github.com/Bowery/gopackages/web"
)

func main() {
	port := os.Getenv("PORT")
	if os.Getenv("ENV") == "production" {
		port = ":80"
	}
	if port == "" {
		port = ":5000"
	}

	server := web.NewServer(port, []web.Handler{
		new(web.SlashHandler),
		new(web.CorsHandler),
		&web.StatHandler{Key: config.StatHatKey, Name: "mercer"},
	}, routes)
	server.Router.NotFoundHandler = &web.NotFoundHandler{renderer}
	server.ListenAndServe()
}
