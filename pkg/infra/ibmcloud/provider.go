package ibmcloud

import (
	"errors"
	"net/http"

	"github.com/moficodes/ibmcloud-kubernetes-admin/pkg/infra"
)

type Provider struct {
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) GetSessionWithToken(token *infra.AuthToken) (infra.CloudSession, error) {
	if token == nil {
		return nil, errors.New("token can not be nil")
	}
	t := &Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiration:   int(token.Expiration),
	}
	session := &Session{Token: t}
	return session, nil
}

func (p *Provider) GetSessionWithCookie(r *http.Request) (infra.CloudSession, error) {
	return GetSession(r)
}
