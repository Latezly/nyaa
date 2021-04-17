package torrentController

import (
	"net/http"
	"strconv"

	"github.com/Latezly/nyaa_go/controllers/router"
	"github.com/Latezly/nyaa_go/models"
	"github.com/Latezly/nyaa_go/models/activities"
	"github.com/Latezly/nyaa_go/models/reports"
	"github.com/Latezly/nyaa_go/models/torrents"
	"github.com/Latezly/nyaa_go/utils/search"
	"github.com/gin-gonic/gin"
)

// TorrentDeleteUserPanel : Controller for deleting a user torrent by a user
func TorrentDeleteUserPanel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.PostForm("id"), 10, 32)
	currentUser := router.GetUser(c)
	torrent, _ := torrents.FindByID(uint(id))
	if currentUser.CurrentOrAdmin(torrent.UploaderID) && torrent.ID > 0 {
		_, _, err := torrent.Delete(false)
		if err == nil {
			if torrent.Uploader == nil {
				torrent.Uploader = &models.User{}
			}
			_, username := torrents.HideUser(torrent.UploaderID, torrent.Uploader.Username, torrent.Hidden)
			if currentUser.IsModerator() { // We hide username on log activity if user is not admin and torrent is hidden
				activities.Log(&models.User{}, torrent.Identifier(), "delete", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), username, currentUser.Username)
			} else {
				activities.Log(&models.User{}, torrent.Identifier(), "delete", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), username, username)
			}
			//delete reports of torrent
			query := &search.Query{}
			query.Append("torrent_id", id)
			torrentReports, _, _ := reports.FindOrderBy(query, "", 0, 0)
			for _, report := range torrentReports {
				report.Delete()
			}
		}
		c.Redirect(http.StatusSeeOther, "/?deleted")
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}
