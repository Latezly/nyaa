package apiController

import (
	"github.com/Latezly/nyaa_go/controllers/middlewares"
	"github.com/Latezly/nyaa_go/controllers/router"
)

func init() {
	api := router.Get().Group("/api")
	{
		api.GET("", APIHandler)
		api.GET("/p/:page", APIHandler)
		api.GET("/view/:id", APIViewHandler)
		api.HEAD("/view/:id", APIViewHeadHandler)
		api.POST("/upload", APIUploadHandler)
		api.POST("/login", APILoginHandler)
		api.GET("/profile", APIProfileHandler)
		api.GET("/user", middlewares.ScopesRequired("user"), APIOwnProfile)
		api.GET("/token/check", APICheckTokenHandler)
		api.GET("/token/refresh", APIRefreshTokenHandler)
		api.Any("/search", APISearchHandler)
		api.Any("/search/p/:page", APISearchHandler)
		api.PUT("/update", APIUpdateHandler)
	}
}
