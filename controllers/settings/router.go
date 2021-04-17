package settingsController

import "github.com/Latezly/nyaa_go/controllers/router"

func init() {
	router.Get().GET("/settings", SeePublicSettingsHandler)
	router.Get().POST("/settings", ChangePublicSettingsHandler)
}
