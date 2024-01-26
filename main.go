package main

import (
	"net/http"

	. "github.com/jwtly10/simple-site-tracker/api/router"
	"github.com/jwtly10/simple-site-tracker/utils/logger"
)

func main() {
	l := logger.Get()

	router := NewRouter()

	l.Info().Msg("Starting server on port 8080")
	http.ListenAndServe(":8080", router)
}
