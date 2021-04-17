package activitiesController

import "github.com/Latezly/nyaa_go/controllers/router"

func init() {
	router.Get().Any("/activities", ActivityListHandler)
	router.Get().Any("/activities/p/:page", ActivityListHandler)
}
