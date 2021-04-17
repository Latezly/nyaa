package captchaController

import (
	"github.com/Latezly/nyaa_go/controllers/router"
	"github.com/Latezly/nyaa_go/utils/captcha"
)

func init() {
	router.Get().Any("/captcha/*hash", captcha.ServeFiles)
}
