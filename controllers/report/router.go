package reportController

import "github.com/Latezly/nyaa_go/controllers/router"

func init() {
	reportRoutes := router.Get().Group("/report")
	{
		//reporting a torrent
		reportRoutes.GET("/:id", ReportViewTorrentHandler)
		reportRoutes.POST("/:id", ReportTorrentHandler)
	}
}
