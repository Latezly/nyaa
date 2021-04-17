package userController

import (
	"net/http"

	"github.com/Latezly/nyaa_go/templates"
	"github.com/Latezly/nyaa_go/utils/cookies"
	msg "github.com/Latezly/nyaa_go/utils/messages"
	"github.com/Latezly/nyaa_go/utils/validator"
	"github.com/Latezly/nyaa_go/utils/validator/user"
	"github.com/gin-gonic/gin"
)

// UserLoginFormHandler : Getting View User Login
func UserLoginFormHandler(c *gin.Context) {
	_, _, errorUser := cookies.CurrentUser(c)
	// User is already connected, redirect to home
	if errorUser == nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	loginForm := userValidator.LoginForm{
		RedirectTo: c.DefaultQuery("redirectTo", ""),
	}
	templates.Form(c, "site/user/login.jet.html", loginForm)
}

// UserLoginPostHandler : Post Login controller
func UserLoginPostHandler(c *gin.Context) {
	b := userValidator.LoginForm{}
	c.Bind(&b)
	messages := msg.GetMessages(c)

	validator.ValidateForm(&b, messages)
	if !messages.HasErrors() {
		_, _, errorUser := cookies.CreateUserAuthentication(c, &b)
		if errorUser == nil {
			url := c.DefaultPostForm("redirectTo", "/")
			c.Redirect(http.StatusSeeOther, url)
			return
		}
		messages.ErrorT(errorUser)
	}
	UserLoginFormHandler(c)
}
