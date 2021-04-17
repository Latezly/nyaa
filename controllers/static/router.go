package staticController

import (
	"net/http"

	"github.com/Latezly/nyaa_go/controllers/databasedumps"
	"github.com/Latezly/nyaa_go/controllers/router"
)

func init() {
	// Static file handlers
	// TODO Use config from cli
	// TODO Make sure the directory exists
	router.Get().StaticFS("/css/", http.Dir("./public/css/"))
	router.Get().StaticFS("/js/", http.Dir("./public/js/"))
	router.Get().StaticFS("/img/", http.Dir("./public/img/"))
	router.Get().StaticFS("/apidoc/", http.Dir("./apidoc/"))
	router.Get().StaticFS("/dbdumps/", http.Dir(databasedumpsController.DatabaseDumpPath))
	router.Get().StaticFS("/gpg/", http.Dir(databasedumpsController.GPGPublicKeyPath))
}
