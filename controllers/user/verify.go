package userController

import (
	"github.com/Latezly/nyaa_go/templates"
	"github.com/Latezly/nyaa_go/utils/email"
	msg "github.com/Latezly/nyaa_go/utils/messages"
	"github.com/gin-gonic/gin"
)

// UserVerifyEmailHandler : Controller when verifying email, needs a token
func UserVerifyEmailHandler(c *gin.Context) {
	token := c.Param("token")
	messages := msg.GetMessages(c)

	_, errEmail := email.EmailVerification(token, c)
	if errEmail != nil {
		messages.ImportFromError("errors", errEmail)
	}
	templates.Static(c, "site/static/verify_success.jet.html")
}
