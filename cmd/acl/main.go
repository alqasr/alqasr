package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/hashicorp/boundary/api"

	"github.com/alqasr/alqasr/internal/service"
	"github.com/alqasr/alqasr/internal/squid"
)

func main() {
	logger := log.New(os.Stderr, "alqasr_acl: ", log.LstdFlags|log.Lmsgprefix)
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
	matchService := service.NewMatchService(ctx, client)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		request, err := squid.NewExtAclRequest(scanner.Text())
		if err != nil {
			logger.Println(err)
			squid.SendERR()
			continue
		}

		ok, err := matchService.Match(request)
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
