package main

import (
	"bufio"
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/boundary/api"
	"github.com/hashicorp/boundary/api/authtokens"
	"github.com/hashicorp/boundary/api/users"

	"github.com/alqasr/alqasr/internal/boundary"
	"github.com/alqasr/alqasr/internal/squid"
)

const (
	globalScope = "global"
	addr        = "http://10.5.0.5:9200/"
)

func auth(ctx context.Context, client *api.Client, username string, password string) (bool, error) {
	client.SetToken(password)

	id, err := boundary.TokenIdFromToken(password)
	if err != nil {
		return false, err
	}

	token, err := authtokens.NewClient(client).Read(ctx, id)
	if err != nil {
		return false, err
	}

	user, err := users.NewClient(client).Read(ctx, token.Item.UserId)
	if err != nil {
		return false, err
	}

	if username != user.Item.LoginName {
		return false, errors.New("security warning: token belongs to another user")
	}

	return true, nil
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
		args := strings.Split(scanner.Text(), " ")
		if len(args) < 2 {
			squid.SendERR()
			continue
		}

		if ok, _ := auth(ctx, client, args[0], args[1]); !ok {
			squid.SendERR()
			continue
		}

		squid.SendOK()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
