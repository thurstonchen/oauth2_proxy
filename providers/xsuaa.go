package providers

import (
	"errors"
	//"fmt"
	//"github.com/bitly/go-simplejson"
	"github.com/thurstonchen/oauth2_proxy/api"
	"log"
	"net/http"
	"net/url"
)

type XSUAAProvider struct {
	*ProviderData
	Tenant string
}

func NewXSUAAProvider(p *ProviderData) *XSUAAProvider {
	p.ProviderName = "xsuaa"

	if p.LoginURL == nil || p.LoginURL.String() == "" {
		p.LoginURL = &url.URL{
			Scheme: "https",
			Host:   "gtt-newdevsandbox.authentication.sap.hana.ondemand.com",
			Path:   "/login",
		}
	}
	if p.RedeemURL == nil || p.RedeemURL.String() == "" {
		p.RedeemURL = &url.URL{
			Scheme: "https",
			Host:   "gtt-newdevsandbox.authentication.sap.hana.ondemand.com",
			Path:   "/oauth/token?grant_type=client_credentials",
		}
	}
	// ValidationURL is the API Base URL
	if p.ValidateURL == nil || p.ValidateURL.String() == "" {
		p.ValidateURL = &url.URL{
			Scheme: "https",
			Host:   "gtt-newdevsandbox.authentication.sap.hana.ondemand.com",
			Path:   "/config?action=who",
		}
	}
	if p.Scope == "" {
		p.Scope = "user:email"
	}

	return &XSUAAProvider{ProviderData: p}
}


func (p *XSUAAProvider) GetEmailAddress(s *SessionState) (string, error) {
	var email string
	var err error

	if s.AccessToken == "" {
		return "", errors.New("missing access token")
	}
	req, err := http.NewRequest("GET", p.ValidateURL.String(), nil)
	if err != nil {
		return "", err
	}

	json, err := api.Request(req)

	if err != nil {
		return "", err
	}

	email, err = getEmailFromJSON(json)

	if err == nil && email != "" {
		return email, err
	}
	email, err = json.Get("userPrincipalName").String()

	if err != nil {
		log.Printf("failed making request %s", err)
		return "", err
	}

	if email == "" {
		log.Printf("failed to get email address")
		return "", err
	}

	return email, err
}
