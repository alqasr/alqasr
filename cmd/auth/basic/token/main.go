package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/boundary/api"

	"github.com/alqasr/alqasr/internal/service"
	"github.com/alqasr/alqasr/internal/squid"
)

func main() {
	logger := log.New(os.Stderr, "alqasr_auth: ", log.LstdFlags|log.Lmsgprefix)
	logger.Println("starting service")

	var boundaryController string
	flag.StringVar(&boundaryController, "boundary-controller", "http://localhost:9200/", "Boundary Controller address as URL.")
	flag.Parse()

	client, err := api.NewClient(nil)
	if err != nil {
		log.Fatal(err)
	}

	client.SetAddr(boundaryController)
	client.SetClientTimeout(time.Second * 10)

	ctx := context.Background()
	tokenAuthService := service.NewTokenAuthService(ctx, client)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		values := strings.Split(scanner.Text(), " ")
		if len(values) < 2 {
			squid.SendERR()
			continue
		}

		ok, err := tokenAuthService.Auth(values[0], values[1])
		if err != nil {
			logger.Println(err)
		}

		if !ok {
			squid.SendERR()
			continue
		}

		squid.SendOK()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
