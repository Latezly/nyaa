package oauth_client

import (
	"github.com/Latezly/nyaa_go/models"
	"github.com/pkg/errors"
)

func Delete(id string) error {
	err := models.ORM.Where("id = ?", id).Delete(&models.OauthClient{}).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
