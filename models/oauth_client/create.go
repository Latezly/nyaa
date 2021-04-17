package oauth_client

import (
	"strings"

	"github.com/Latezly/nyaa_go/models"
	"github.com/Latezly/nyaa_go/utils/sanitize"
	"github.com/Latezly/nyaa_go/utils/validator/api"
)

func Create(form *apiValidator.CreateForm) (*models.OauthClient, error) {
	client := &models.OauthClient{
		Name:              form.Name,
		RedirectURIs:      strings.Join(sanitize.ClearEmpty(form.RedirectURI), "|"),
		GrantTypes:        strings.Join(sanitize.ClearEmpty(form.GrantTypes), "|"),
		ResponseTypes:     strings.Join(sanitize.ClearEmpty(form.ResponseTypes), "|"),
		Scope:             form.Scope,
		Owner:             form.Owner,
		PolicyURI:         form.PolicyURI,
		TermsOfServiceURI: form.TermsOfServiceURI,
		ClientURI:         form.ClientURI,
		LogoURI:           form.LogoURI,
		Contacts:          strings.Join(sanitize.ClearEmpty(form.Contacts), "|"),
	}

	err := models.ORM.Create(client).Error
	if err != nil {
		return client, err
	}

	return client, nil
}
