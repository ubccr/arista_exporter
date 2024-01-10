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

type powerSupply struct {
	ModelName     string
	Capacity      float64
	Dominant      bool
	InputCurrent  float64
	OutputCurrent float64
	InputVoltage  float64
	OutputPower   float64
	State         string
}

type PowerProber struct {
	PowerSupplies map[string]*powerSupply

	powerGuage *prometheus.GaugeVec
}

func (p *PowerProber) GetCmd() string {
	return "show system environment power"
}

func (p *PowerProber) Register(registry *prometheus.Registry) {
	p.powerGuage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "arista_power_supply_state",
		Help: "Contains Power Supply state",
	}, []string{"powerSupply", "state"})

	registry.MustRegister(p.powerGuage)
}

func (p *PowerProber) Handler(logger log.Logger) {
	for id, ps := range p.PowerSupplies {
		if ps.State == "ok" {
			p.powerGuage.WithLabelValues(id, ps.State).Set(1)
		} else {
			p.powerGuage.WithLabelValues(id, ps.State).Set(0)
		}
	}
}
