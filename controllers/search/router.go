package searchController

import "github.com/Latezly/nyaa_go/controllers/router"

func init() {
	router.Get().Any("/", SearchHandler)
	router.Get().Any("/p/:page", SearchHandler)
	router.Get().Any("/search", SearchHandler)
	router.Get().Any("/search/:page", SearchHandler)
}
