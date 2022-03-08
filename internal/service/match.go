package service

import (
	"context"
	"fmt"

	"github.com/hashicorp/boundary/api"
	"github.com/hashicorp/boundary/api/hostcatalogs"
	"github.com/hashicorp/boundary/api/hosts"
	"github.com/hashicorp/boundary/api/hostsets"

	"github.com/alqasr/alqasr/internal/boundary"
	"github.com/alqasr/alqasr/internal/squid"
)

type MatchService struct {
	ctx    context.Context
	client *api.Client
}

func NewMatchService(ctx context.Context, client *api.Client) *MatchService {
	return &MatchService{
		ctx:    ctx,
		client: client,
	}
}

func (srv *MatchService) Match(request *squid.ExtAclRequest) (bool, error) {
	allowList := map[string]struct{}{}

	client := srv.client.Clone()
	client.SetToken(request.Token)

	var (
		hostCatalogsClient = hostcatalogs.NewClient(client)
		hostSetsClient     = hostsets.NewClient(client)
		hostsClient        = hosts.NewClient(client)
	)

	hostCatalogs, err := hostCatalogsClient.List(srv.ctx, boundary.GlobalScope, hostcatalogs.WithRecursive(true))
	if err != nil {
		return false, err
	}

	for _, hostCatalog := range hostCatalogs.Items {
		hostSets, err := hostSetsClient.List(srv.ctx, hostCatalog.Id, hostsets.WithFilter(fmt.Sprintf(`"/item/description"==%s`, boundary.HostSetDescription)))
		if err != nil {
			return false, err
		}

		for _, hostSet := range hostSets.Items {
			hosts, err := hostsClient.List(srv.ctx, hostSet.HostCatalogId)
			if err != nil {
				return false, err
			}

			for _, host := range hosts.Items {
				allowList[host.Attributes["address"].(string)] = struct{}{}
			}
		}
	}

	if _, ok := allowList[request.Dst]; ok {
		return true, nil
	}

	return false, nil
}
