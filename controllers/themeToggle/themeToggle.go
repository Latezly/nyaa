package themeToggleController

import (
	"net/http"

	"github.com/Latezly/nyaa_go/config"
	"github.com/Latezly/nyaa_go/controllers/router"
	"github.com/Latezly/nyaa_go/utils/timeHelper"
	"github.com/gin-gonic/gin"
	"github.com/Latezly/nyaa_go/utils/publicSettings"
	
)

// toggleThemeHandler : Controller to switch between theme1 & theme2
func toggleThemeHandler(c *gin.Context) {

	DefaultTheme := config.DefaultTheme(false) 
	DefaultDarkTheme := config.DefaultTheme(true)

	theme := publicSettings.GetThemeFromRequest(c)
	theme2, err := c.Cookie("theme2")
	
	if err != nil {
		theme2 = publicSettings.GetDarkThemeFromRequest(c)
	}
	
	if theme != DefaultDarkTheme && theme2 != DefaultDarkTheme {
		//None of the themes are dark ones, force the second one as the dark one
		theme2 = DefaultDarkTheme
	} else if  theme == theme2 {
		//Both theme are dark ones, force the second one as the default (light) theme
		theme2 = DefaultTheme
	}
	//Get theme1 & theme2 value
	
	// If logged in, update user theme (will not work otherwise)
	user := router.GetUser(c)
	if user.ID > 0 {
		user.Theme = theme2
		user.UpdateRaw()
	}
	
	//Switch theme & theme2 value
	http.SetCookie(c.Writer, &http.Cookie{Name: "theme", Value: theme2, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "theme2", Value: theme, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})	
	
	//Redirect user to page he was in beforehand
	if  c.Request.URL.Query()["no_redirect"] == nil {
		c.Redirect(http.StatusSeeOther, c.Param("redirect") + "#footer")
	}
	return
}

func getDomainName() string {
	domain := config.Get().Cookies.DomainName
	if config.Get().Environment == "DEVELOPMENT" {
		domain = ""
	}
	return domain
}
