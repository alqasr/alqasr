package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/boundary/api"
	"github.com/hashicorp/boundary/api/hostcatalogs"
	"github.com/hashicorp/boundary/api/hosts"
	"github.com/hashicorp/boundary/api/hostsets"

	"github.com/alqasr/alqasr/internal/squid"
)

const (
	globalScope = "global"
	addr        = "http://10.5.0.5:9200/"
)

type Request struct {
	Login   string
	Token   string
	Src     string
	SrcPort string
	Dst     string
	DstPort string
}

func check(ctx context.Context, client *api.Client, request *Request) (bool, error) {
	allowList := map[string]struct{}{}

	client.SetToken(request.Token)

	var (
		hostCatalogsClient = hostcatalogs.NewClient(client)
		hostSetsClient     = hostsets.NewClient(client)
		hostsClient        = hosts.NewClient(client)
	)

	hostCatalogs, err := hostCatalogsClient.List(ctx, globalScope, hostcatalogs.WithRecursive(true))
	if err != nil {
		return false, err
	}

	for _, hostCatalog := range hostCatalogs.Items {
		hostSets, err := hostSetsClient.List(ctx, hostCatalog.Id, hostsets.WithFilter(fmt.Sprintf(`"/item/description"==%s`, "HTTP")))
		if err != nil {
			return false, err
		}

		for _, hostSet := range hostSets.Items {
			hosts, err := hostsClient.List(ctx, hostSet.HostCatalogId)
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

func main() {
	ctx := context.Background()

	// The default Addr is http://127.0.0.1:9200, but this can be overridden by
	// setting the `BOUNDARY_ADDR` environment variable.
	os.Setenv("BOUNDARY_ADDR", addr)

	clientConfig, err := api.DefaultConfig()
	if err != nil {
		log.Fatal(err)
	}

	// change default value (60 seconds)
	clientConfig.Timeout = time.Second * 1

	client, err := api.NewClient(clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// format: %LOGIN %>{Proxy-Authorization} %SRC %SRCPORT %DST %PORT
		args := strings.Split(line, " ")
		if len(args) < 7 {
			squid.SendERR()
			continue
		}

		token, err := squid.PasswordFromProxyAuthorization(args[1])
		if err != nil {
			squid.SendERR()
			continue
		}

		// TODO: check token validity; is really needed?

		request := Request{
			Login:   args[0],
			Token:   token,
			Src:     args[2],
			SrcPort: args[3],
			Dst:     args[4],
			DstPort: args[5],
		}

		if ok, _ := check(ctx, client, &request); !ok {
			// TODO: log error
			squid.SendERR()
			continue
		}

		squid.SendOK()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
