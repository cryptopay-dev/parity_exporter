package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	exporterPort, etherscanUrl, etherscanKey, parityUrl, endpoint string
)

func init() {
	exporterPort = os.Getenv("EXPORTER_PORT")
	if exporterPort == "" {
		log.Fatal("EXPORTER_PORT is not set")
	}

	etherscanUrl = os.Getenv("ETHERSCAN_URL")
	if etherscanUrl == "" {
		log.Fatal("ETHERSCAN_URL is not set")
	}

	etherscanKey = os.Getenv("ETHERSCAN_KEY")
	if etherscanKey == "" {
		log.Fatal("ETHERSCAN_KEY is not set")
	}

	parityUrl = os.Getenv("PARITY_URL")
	if parityUrl == "" {
		log.Fatal("PARITY_URL is not set")
	}

	endpoint = os.Getenv("ENDPOINT")
	if endpoint == "" {
		log.Println("ENDPOINT not provided. Setting ENDPOINT to '/'")
		endpoint = "/"
	}
	if match, _ := regexp.MatchString(`(?:^[^/]|[\s])`, endpoint); match {
		log.Fatal("ENDPOINT is invalid")
	}
}

func main() {
	registry := prometheus.NewPedanticRegistry()

	registry.MustRegister(NewParityLastBlock(parityUrl))
	registry.MustRegister(NewEtherscanLastBlock(etherscanUrl, etherscanKey))

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog:      log.New(os.Stderr, log.Prefix(), log.Flags()),
		ErrorHandling: promhttp.ContinueOnError,
	})

	http.Handle(endpoint, handler)
	fmt.Printf("Parity/Etherscan prometheus exporter started on port %s\n", exporterPort)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", exporterPort), nil))
}
