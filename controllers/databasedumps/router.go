package databasedumpsController

import "github.com/Latezly/nyaa_go/controllers/router"

func init() {
	router.Get().Any("/dumps", DatabaseDumpHandler)
}
