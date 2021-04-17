package themeToggleController

import "github.com/Latezly/nyaa_go/controllers/router"

func init() {
	router.Get().Any("/dark", toggleThemeHandler)
	router.Get().Any("/dark/*redirect", toggleThemeHandler)
}
