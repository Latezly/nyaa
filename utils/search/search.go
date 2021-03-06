package search

import (
	"time"

	"github.com/Latezly/nyaa_go/config"
	"github.com/Latezly/nyaa_go/models"
	"github.com/Latezly/nyaa_go/utils/cache"
	"github.com/Latezly/nyaa_go/utils/log"
	"github.com/gin-gonic/gin"
)

var searchOperator string
var useTSQuery bool

// Configure : initialize search
func Configure(conf *config.SearchConfig) (err error) {
	useTSQuery = false
	// Postgres needs ILIKE for case-insensitivity
	if models.ORM.Dialect().GetName() == "postgres" {
		searchOperator = "ILIKE ?"
		//useTSQuery = true
		// !!DISABLED!! because this makes search a lot stricter
		// (only matches at word borders)
	} else {
		searchOperator = "LIKE ?"
	}
	return
}

func stringIsASCII(input string) bool {
	for _, char := range input {
		if char > 127 {
			return false
		}
	}
	return true
}

// ByQueryNoUser : search torrents according to request without user
func ByQueryNoUser(c *gin.Context, pagenum int) (search TorrentParam, tor []models.Torrent, count int, err error) {
	search, tor, count, err = ByQuery(c, pagenum, false, false, false, false)
	return
}

// ByQueryWithUser : search torrents according to request with user
func ByQueryWithUser(c *gin.Context, pagenum int) (search TorrentParam, tor []models.Torrent, count int, err error) {
	search, tor, count, err = ByQuery(c, pagenum, true, false, false, true)
	return
}

// ByQueryDeleted : search deleted torrents according to request with user and count
func ByQueryDeleted(c *gin.Context, pagenum int) (search TorrentParam, tor []models.Torrent, count int, err error) {
	search, tor, count, err = ByQuery(c, pagenum, true, true, false, true)
	return
}

// ByQueryNoHidden : search torrents and filter those hidden
func ByQueryNoHidden(c *gin.Context, pagenum int) (search TorrentParam, tor []models.Torrent, count int, err error) {
	search, tor, count, err = ByQuery(c, pagenum, false, false, true, false)
	return
}

// TODO Clean this up
// Some fields are postgres specific (countAll, withUser)
// elasticsearch always provide a count to how many hits
// ES doesn't store users
// deleted is unused because es doesn't index deleted torrents
func ByQuery(c *gin.Context, pagenum int, withUser bool, deleted bool, hidden bool, locked bool) (TorrentParam, []models.Torrent, int, error) {
	var torrentParam TorrentParam
	torrentParam.FromRequest(c)
	torrentParam.Offset = uint32(pagenum)
	torrentParam.Hidden = hidden
	torrentParam.Locked = locked
	torrentParam.Full = withUser
	torrentParam.Deleted = deleted

	if torrentParam.Abort {
		return torrentParam, []models.Torrent{}, 0, nil
	}

	if found, ok := cache.C.Get(torrentParam.Identifier()); ok {
		log.Infof("Retrieve results from Cache in %s", torrentParam.Identifier())
		torrentCache := found.(*TorrentCache)
		return torrentParam, torrentCache.Torrents, torrentCache.Count, nil
	}
	if config.Get().Search.EnableElasticSearch && models.ElasticSearchClient != nil && !deleted {
		tor, totalHits, err := torrentParam.FindES(c, models.ElasticSearchClient)
		// If there are results no errors from ES search we use the ES client results
		if totalHits > 0 && err == nil {
			// Since we have results, we cache them so we don't ask everytime ES for the same results
			cache.C.Set(torrentParam.Identifier(), &TorrentCache{tor, int(totalHits)}, 5*time.Minute)
			// we return the results
			// Convert back to non-json torrents
			return torrentParam, tor, int(totalHits), nil
		}
		// Errors from ES should be managed in the if condition. Log is triggered only if err != nil (checkError behaviour)
		log.CheckErrorWithMessage(err, "ES_ERROR_MSG: Seems like ES was not reachable whereas it was when starting the app.")
	}
	// We fallback to PG, if ES gives error or no results or if ES is disabled in config or if deleted search is enabled
	log.Errorf("Falling back to postgresql query")
	tor, totalHits, err := torrentParam.FindDB(c)
	if totalHits > 0 && err == nil {
		cache.C.Set(torrentParam.Identifier(), &TorrentCache{tor, int(totalHits)}, 5*time.Minute)
	}
	return torrentParam, tor, int(totalHits), err
}

// AuthorizedQuery return a seach byquery according to the bool. If false, it doesn't look for hidden torrents, else it looks for every torrents
func AuthorizedQuery(c *gin.Context, pagenum int, authorized bool, locked bool) (TorrentParam, []models.Torrent, int, error) {
	return ByQuery(c, pagenum, true, false, !authorized, locked)
}
