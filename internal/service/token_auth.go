package service

import (
	"context"

	"github.com/hashicorp/boundary/api"
	"github.com/hashicorp/boundary/api/authtokens"

	"github.com/alqasr/alqasr/internal/boundary"
)

type TokenAuthService struct {
	ctx    context.Context
	client *api.Client
}

func NewTokenAuthService(ctx context.Context, client *api.Client) *TokenAuthService {
	return &TokenAuthService{
		ctx:    ctx,
		client: client,
	}
}

func (srv *TokenAuthService) Auth(username string, password string) (bool, error) {
	client := srv.client.Clone()
	client.SetToken(password)

	id, err := boundary.TokenIdFromToken(password)
	if err != nil {
		return false, err
	}

	_, err = authtokens.NewClient(client).Read(srv.ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}
