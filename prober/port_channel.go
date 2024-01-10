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
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type lagPort struct {
	LacpMode           string
	Protocol           string
	TimeBecameActive   float64
	Weight             float64
	TimeBecameInactive float64
	ReasonUnconfigured string
}

type portChannel struct {
	ActivePorts   map[string]lagPort
	InactivePorts map[string]lagPort
}

type PortChannelProber struct {
	PortChannels map[string]*portChannel

	metrics *prometheus.GaugeVec
}

func (p *PortChannelProber) GetCmd() string {
	return "show port-channel detailed"
}

func (p *PortChannelProber) Register(registry *prometheus.Registry) {
	p.metrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "arista_portchannel_ports",
		Help: "Contains port channel ports by state",
	}, []string{"interface", "state"})

	registry.MustRegister(p.metrics)
}

func (p *PortChannelProber) Handler(logger log.Logger) {
	for id, pc := range p.PortChannels {
		p.metrics.WithLabelValues(id, "active").Set(float64(len(pc.ActivePorts)))
		p.metrics.WithLabelValues(id, "inactive").Set(float64(len(pc.InactivePorts)))
	}
}
