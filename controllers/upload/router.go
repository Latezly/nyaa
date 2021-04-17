package uploadController

import "github.com/Latezly/nyaa_go/controllers/router"

func init() {
	router.Get().Any("/upload", UploadHandler)
	router.Get().Any("/upload/status/:id", multiUploadStatus)
}
