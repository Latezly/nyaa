package searchController

import (
	"html"
	"net/http"
	"strconv"
	"strings"
	"fmt"

	"math"

	"github.com/Latezly/nyaa_go/controllers/router"
	"github.com/Latezly/nyaa_go/models"
	"github.com/Latezly/nyaa_go/templates"
	"github.com/Latezly/nyaa_go/models/torrents"
	"github.com/Latezly/nyaa_go/utils/search"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)
// SearchHandler : Controller for displaying search result page, accepting common search arguments
func SearchHandler(c *gin.Context) {
	var err error
	// TODO Don't create a new client for each request
	// TODO Fallback to postgres search if es is down

	page := c.Param("page")
	currentUser := router.GetUser(c)
	// db params url
	pagenum := 1
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		if pagenum <= 0 {
			c.AbortWithError(http.StatusNotFound, errors.New("Can't find a page with negative value"))
			return
		}
	}
	
	searchForm := templates.NewSearchForm(c)
	
	if c.Param("id") != "" {
		query := c.Request.URL.Query()
		query.Set("userID", c.Param("id"))
		c.Request.URL.RawQuery = query.Encode()
		searchForm.SearchURL = fmt.Sprintf("/user/%s/%s/search", c.Param("id"), c.Param("username"))
		searchForm.UserName = c.Param("username") //Only add username if user search route
	}
	
	userID, err := strconv.ParseUint(c.Query("userID"), 10, 32)
	if err != nil {
		userID = 0
	}
	
	if userID == 0 && c.Param("id") != "" && c.Param("id") != "0" {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/user/%s/%s", c.Param("id"), c.Param("username")))
		//User is trying to use the user search route with an inexisting user
		//Must redirect him to user search instead of simply showing "no torrents found!"
	}
	
	if c.Query("hash") != "" {
		torrent, err := torrents.FindRawByHash(strings.TrimSpace(c.Query("hash")))
		//Wanna make sure to remove spaces because user copy-pasting hashes might include spaces at time
		if err == nil {
			templates.ModelList(c, "site/torrents/listing.jet.html", models.TorrentsToJSON([]models.Torrent{torrent}), templates.Navigation{1, 1, 0, "/search"}, searchForm)
			//We already fetched the torrent so we can directly show the template without needing to do a search.AuthorizedQuery
		} else {
			variables := templates.Commonvariables(c)
			searchForm.ShowRefine = true
			variables.Set("Search", searchForm)
			templates.Render(c, "errors/no_results.jet.html", variables)
			//The hash hasn't been found so no point in showing any result, but we do want to show the refine so that the user can search another hash or remove it from the search
		}
		return
	}
	
	searchParam, torrents, nbTorrents, err := search.AuthorizedQuery(c, pagenum, currentUser.CurrentOrJanitor(uint(userID)), currentUser.IsJanitor())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	
	// Convert back to strings for now.
	category := ""
	if len(searchParam.Category) > 0 {
		category = searchParam.Category[0].String()
	}
	nav := templates.Navigation{int(nbTorrents), int(searchParam.Max), int(searchParam.Offset), "search"}
	
	searchForm.TorrentParam, searchForm.Category = searchParam, category

	if c.Query("refine") == "1" || nbTorrents == 0 {
		searchForm.ShowRefine = true
	}

	maxPages := math.Ceil(float64(nbTorrents) / float64(searchParam.Max))
	if pagenum > int(maxPages) {
		variables := templates.Commonvariables(c)
		variables.Set("Search", searchForm)
		templates.Render(c, "errors/no_results.jet.html", variables)
		return
	}

	templates.ModelList(c, "site/torrents/listing.jet.html", models.TorrentsToJSON(torrents), nav, searchForm)
}
