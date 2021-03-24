package vault

import (
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/github"
	"github.com/pkg/errors"
)

func GetVaultClient() (*api.Client, error) {
	config := api.DefaultConfig()

	if err := config.ReadEnvironment(); err != nil {
		return nil, errors.Wrap(err, "failed to read env vars")
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "can not get client")
	}

	handler := github.CLIHandler{}
	var m map[string]string
	secret, err := handler.Auth(client, m)
	if err != nil {
		return nil, errors.Wrap(err, "could not get secret")
	}
	token, err := secret.TokenID()
	if err != nil {
		return nil, errors.Wrap(err, "error getting client token")
	}

	client.SetToken(token)
	return client, nil
}
