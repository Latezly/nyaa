package moderatorController

import (
	"github.com/Latezly/nyaa_go/models"
	"github.com/Latezly/nyaa_go/models/comments"
	"github.com/Latezly/nyaa_go/models/reports"
	"github.com/Latezly/nyaa_go/models/torrents"
	"github.com/Latezly/nyaa_go/models/users"
	"github.com/Latezly/nyaa_go/templates"
	"github.com/gin-gonic/gin"
)

// IndexModPanel : Controller for showing index page of Mod Panel
func IndexModPanel(c *gin.Context) {
	offset := 10
	torrents, _, _ := torrents.FindAllForAdminsOrderBy("torrent_id DESC", offset, 0)
	users, _ := users.FindUsersForAdmin(offset, 0)
	comments, _ := comments.FindAll(offset, 0, "", "")
	torrentReports, _, _ := reports.GetAll(offset, 0)

	templates.PanelAdmin(c, torrents, models.TorrentReportsToJSON(torrentReports), users, comments)
}

func GuidelinesModPanel(c *gin.Context) {
	templates.Static(c, "admin/guidelines.jet.html")
}
