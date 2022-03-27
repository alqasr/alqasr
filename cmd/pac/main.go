package main

import (
	"bytes"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alqasr/alqasr/internal/config"
)

const pac = `function FindProxyForURL(url, host) {
  {{- $port := .Port -}}
  {{- if not .AllowedDomains }}
    return "PROXY 127.0.0.1:{{ $port }}";
  {{- else }}
  {{- range .AllowedDomains }}
  if (dnsDomainIs(host, "{{.}}")) 
    return "PROXY 127.0.0.1:{{ $port }}";
  {{ end }}

  return "DIRECT";
  {{- end }}
}
`

func render(cfg config.Proxy, pac string) ([]byte, error) {
	tmpl, err := template.New("pac").Parse(pac)
	if err != nil {
		return nil, err
	}

	output := bytes.NewBuffer(nil)

	err = tmpl.Execute(output, cfg)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func main() {
	logger := log.New(os.Stdout, "alqasr_pac: ", log.LstdFlags|log.Lmsgprefix)
	logger.Println("starting service")

	var configFile string
	flag.StringVar(&configFile, "config.file", "config.yml", "alqasr configuration file path")
	flag.Parse()

	logger.Println("loading config file")

	cfg, err := config.Load(configFile)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("render proxy.pac using config")

	pac, err := render(cfg.Proxy, pac)
	if err != nil {
		logger.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/proxy.pac", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/x-ns-proxy-autoconfig")
		w.Write(pac)
	})

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 120,
		ErrorLog:     logger,
	}

	logger.Fatal(server.ListenAndServe())
}
