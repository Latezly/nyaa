package moderatorController

import (
	"fmt"
	"html"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/Latezly/nyaa_go/models"
	"github.com/Latezly/nyaa_go/models/notifications"
	"github.com/Latezly/nyaa_go/templates"
	"github.com/Latezly/nyaa_go/utils/log"
	msg "github.com/Latezly/nyaa_go/utils/messages"
	"github.com/Latezly/nyaa_go/utils/validator"
	"github.com/Latezly/nyaa_go/utils/validator/announcement"
	"github.com/gin-gonic/gin"
)

func listAnnouncements(c *gin.Context) {
	page := c.Param("page")
	pagenum := 1
	offset := 100
	var err error
	messages := msg.GetMessages(c)
	deleted := c.Request.URL.Query()["deleted"]
	if deleted != nil {
		messages.AddInfoTf("infos", "annoucement_deleted")
	}
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	var conditions string
	var values []interface{}

	announcements, nbAnnouncements := notifications.FindAll(offset, (pagenum-1)*offset, conditions, values...)
	nav := templates.Navigation{nbAnnouncements, offset, pagenum, "mod/announcements/p"}
	templates.ModelList(c, "admin/announcements.jet.html", announcements, nav, templates.NewSearchForm(c))
}

func addAnnouncement(c *gin.Context) {
	announcement := &models.Notification{}
	messages := msg.GetMessages(c)

	id := c.Query("id")
	if id == "" && len(messages.GetInfos("ID_ANNOUNCEMENT")) > 0 {
		id = messages.GetInfos("ID_ANNOUNCEMENT")[0]
	}
	idInt, _ := strconv.Atoi(id)
	if idInt > 0 {
		var err error
		announcement, _ = notifications.FindByID(uint(idInt))
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
	duration := int(math.Ceil(math.Max(1, float64(announcement.Expire.Sub(time.Now())/(24*time.Hour)))))
	form := &announcementValidator.CreateForm{
		ID:      announcement.ID,
		Message: announcement.Content,
		Duration:   duration,
	}
	c.Bind(form)
	if form.Duration == 0 {
		form.Duration = duration
	}
	templates.Form(c, "admin/announcement_form.jet.html", form)
}

func postAnnouncement(c *gin.Context) {
	messages := msg.GetMessages(c)
	announcement := &models.Notification{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	if id > 0 {
		var err error
		announcement, err = notifications.FindByID(uint(id))
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
	form := &announcementValidator.CreateForm{}
	// We bind the request to the form
	c.Bind(form)
	// We try to validate the form
	validator.ValidateForm(form, messages)
	// If validation has failed, errors are added in messages variable
	if !messages.HasErrors() {
		// No errors, check if we update or create
		if id > 0 { // announcement exists we update
			err := notifications.UpdateAnnouncement(announcement, form) // Making the update query
			if err != nil {
				// Error, we add it to the messages variable
				messages.AddErrorT("errors", "update_failed")
			} else {
				// Success, we add a notice to the messages variable
				messages.AddInfoT("infos", "update_success")
			}
		} else { // announcement doesn't exist, we create it
			var err error
			currentTime := time.Now()
			announcement, err := notifications.NotifyAll(form.Message, currentTime.Add(time.Hour * time.Duration(form.Duration)))
			if err != nil {
				// Error, we add it as a message
				messages.AddErrorT("errors", "create_failed")
			} else {
				// Success, we redirect to the edit form
				messages.AddInfoT("infos", "create_anouncement_success")
				id := fmt.Sprintf("%d", announcement.ID)
				messages.AddInfo("ID_ANNOUNCEMENT", id)
			}
		}
	}
	// If we are still here, we show the form
	addAnnouncement(c)
}

// deleteAnnouncement : Controller for deleting an announcement
func deleteAnnouncement(c *gin.Context) {
	id, _ := strconv.ParseInt(c.PostForm("id"), 10, 32)
	announcement, err := notifications.FindByID(uint(id))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	err = announcement.Delete()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/mod/announcement?deleted")
}
