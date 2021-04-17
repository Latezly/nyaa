package faqController

import "github.com/Latezly/nyaa_go/controllers/router"

func init() {
	router.Get().Any("/faq", FaqHandler)
}
