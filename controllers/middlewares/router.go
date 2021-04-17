package middlewares

import "github.com/Latezly/nyaa_go/controllers/router"

func init() {
	router.Get().Use(CSP(), ErrorMiddleware())
}
