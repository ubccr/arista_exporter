// Copyright 2024 Andrew E. Bruno
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/aristanetworks/goeapi"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/ubccr/arista_exporter/prober"
)

const (
	aristaEndpoint  = "/arista"
	metricsEndpoint = "/metrics"
)

var (
	configFile    = kingpin.Flag("config.file", "Arista exporter config file").Default("~/.eapi.conf").String()
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9465").String()
)

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("arista_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting arista_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())
	level.Info(logger).Log("msg", "Starting Server", "address", *listenAddress)

	if *configFile != "" {
		level.Info(logger).Log("msg", "Using config file", "path", *configFile)
		goeapi.LoadConfig(*configFile)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Arista Exporter</title></head>
             <body>
             <h1>Arista Exporter</h1>
             <p><a href='` + aristaEndpoint + `'>Arista Metrics</a></p>
             <p><a href='` + metricsEndpoint + `'>Exporter Metrics</a></p>
             </body>
             </html>`))
	})
	http.HandleFunc(aristaEndpoint, func(w http.ResponseWriter, r *http.Request) {
		prober.Handler(w, r, logger, nil)
	})
	http.Handle(metricsEndpoint, promhttp.Handler())
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}
