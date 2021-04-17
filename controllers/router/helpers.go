package router

import (
	"github.com/Latezly/nyaa_go/models"
	"github.com/Latezly/nyaa_go/utils/cookies"
	"github.com/gin-gonic/gin"
)

// GetUser return the current user from the context
func GetUser(c *gin.Context) *models.User {
	user, _, _ := cookies.CurrentUser(c)
	if user == nil {
		return &models.User{}
	}
	return user
}
