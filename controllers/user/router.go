package userController

import "github.com/Latezly/nyaa_go/controllers/router"
import "github.com/Latezly/nyaa_go/controllers/feed"
import "github.com/Latezly/nyaa_go/controllers/search"

func init() {

	// Login
	router.Get().POST("/login", UserLoginPostHandler)
	router.Get().GET("/login", UserLoginFormHandler)

	// Register
	router.Get().GET("/register", UserRegisterFormHandler)
	router.Get().POST("/register", UserRegisterPostHandler)

	// Logout
	router.Get().POST("/logout", UserLogoutHandler)

	// Notifications
	router.Get().GET("/notifications", UserNotificationsHandler)

	// Verify Email
	router.Get().Any("/verify/email/:token", UserVerifyEmailHandler)

	// User Profile specific routes
	userRoutes := router.Get().Group("/user")
	{
		userRoutes.GET("", RedirectToUserSearch)
		userRoutes.GET("/:id", UserProfileHandler)
		userRoutes.GET("/:id/:username", UserProfileHandler)
		userRoutes.GET("/:id/:username/follow", UserFollowHandler)
		userRoutes.GET("/:id/:username/edit", UserDetailsHandler)
		userRoutes.POST("/:id/:username/edit", UserProfileFormHandler)
		userRoutes.GET("/:id/:username/apireset", UserAPIKeyResetHandler)
		userRoutes.GET("/:id/:username/search", searchController.SearchHandler)
		userRoutes.GET("/:id/:username/search/:page", searchController.SearchHandler)
		userRoutes.GET("/:id/:username/feed", feedController.RSSHandler)
		userRoutes.GET("/:id/:username/feed/:page", feedController.RSSHandler)
		userRoutes.POST("/:id/:username/delete", UserProfileDelete)
		userRoutes.POST("/:id/:username/ban", UserProfileBan)
	}
	
	router.Get().Any("/username", RedirectToUserSearch)
	router.Get().Any("/username/:username", UserGetFromName)
	router.Get().Any("/username/:username/search", searchController.SearchHandler)
	router.Get().Any("/username/:username/search:page", searchController.SearchHandler)
}
