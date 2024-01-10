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

package prober

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/aristanetworks/goeapi"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prober interface {
	GetCmd() string
	Register(*prometheus.Registry)
	Handler(log.Logger)
}

func Handler(w http.ResponseWriter, r *http.Request, logger log.Logger, params url.Values) {
	if params == nil {
		params = r.URL.Query()
	}

	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", http.StatusBadRequest)
		return
	}

	modules := params.Get("module")
	if modules == "" {
		modules = "power"
	}

	probers := make([]Prober, 0)

	for _, moduleName := range strings.Split(modules, ",") {
		switch moduleName {
		case "mlag":
			probers = append(probers, &MLAGProber{})
		case "power":
			probers = append(probers, &PowerProber{})
		case "portchannel":
			probers = append(probers, &PortChannelProber{})
		case "redundancy":
			probers = append(probers, &RedundancyProber{})
		case "switchover":
			probers = append(probers, &SwitchoverProber{})
		default:
			http.Error(w, fmt.Sprintf("Unknown module %q", moduleName), http.StatusBadRequest)
			level.Debug(logger).Log("msg", "Unknown module", "module", moduleName)
			return
		}
	}

	node, err := goeapi.ConnectTo(target)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to %q", target), http.StatusBadRequest)
		level.Error(logger).Log("msg", "Failed to connect to target", "err", err)
		return
	}

	handle, _ := node.GetHandle("json")
	registry := prometheus.NewRegistry()

	for _, p := range probers {
		p.Register(registry)
		handle.AddCommand(p)
	}

	if err := handle.Call(); err != nil {
		http.Error(w, "Failed to run arista command", http.StatusInternalServerError)
		level.Error(logger).Log("msg", "Failed to run arista command", "err", err)
		return
	}

	for _, p := range probers {
		p.Handler(logger)
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
