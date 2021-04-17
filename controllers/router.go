package controllers

import (
	"net/http"

	_ "github.com/Latezly/nyaa_go/controllers/activities"    // activities controller
	_ "github.com/Latezly/nyaa_go/controllers/api"           // api controller
	_ "github.com/Latezly/nyaa_go/controllers/captcha"       // captcha controller
	_ "github.com/Latezly/nyaa_go/controllers/databasedumps" // databasedumps controller
	_ "github.com/Latezly/nyaa_go/controllers/faq"           // faq controller
	_ "github.com/Latezly/nyaa_go/controllers/feed"          // feed controller
	_ "github.com/Latezly/nyaa_go/controllers/middlewares"   // middlewares
	_ "github.com/Latezly/nyaa_go/controllers/moderator"     // moderator controller
	_ "github.com/Latezly/nyaa_go/controllers/oauth"         // oauth2 controller
	_ "github.com/Latezly/nyaa_go/controllers/pprof"         // pprof controller
	_ "github.com/Latezly/nyaa_go/controllers/report"        // report controller
	"github.com/Latezly/nyaa_go/controllers/router"
	_ "github.com/Latezly/nyaa_go/controllers/search"   // search controller
	_ "github.com/Latezly/nyaa_go/controllers/settings" // settings controller
	_ "github.com/Latezly/nyaa_go/controllers/static"   // static files
	_ "github.com/Latezly/nyaa_go/controllers/themeToggle" // themeToggle controller
	_ "github.com/Latezly/nyaa_go/controllers/torrent"  // torrent controller
	_ "github.com/Latezly/nyaa_go/controllers/upload"   // upload controller
	_ "github.com/Latezly/nyaa_go/controllers/user"     // user controller
	"github.com/justinas/nosurf"
)

// CSRFRouter : CSRF protection for Router variable for exporting the route configuration
var CSRFRouter *nosurf.CSRFHandler

func init() {
	CSRFRouter = nosurf.New(router.Get())
	CSRFRouter.ExemptRegexp("/api(?:/.+)*")
	CSRFRouter.ExemptRegexp("/mod(?:/.+)*")
	CSRFRouter.ExemptPath("/upload")
	CSRFRouter.ExemptPath("/user/login")
	CSRFRouter.ExemptPath("/oauth2/token")
	CSRFRouter.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid CSRF tokens", http.StatusBadRequest)
	}))
	CSRFRouter.SetBaseCookie(http.Cookie{
		Path:   "/",
		MaxAge: nosurf.MaxAge,
	})

}
