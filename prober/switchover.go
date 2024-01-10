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

type SwitchoverProber struct {
	SwitchoverCount      float64
	switchoverCountGuage prometheus.Gauge
}

func (s *SwitchoverProber) GetCmd() string {
	return "show redundancy switchover sso"
}

func (s *SwitchoverProber) Register(registry *prometheus.Registry) {
	s.switchoverCountGuage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "arista_redundancy_switchover_count",
		Help: "Contains redundancy switchover count",
	})
	registry.MustRegister(s.switchoverCountGuage)
}

func (s *SwitchoverProber) Handler(logger log.Logger) {
	s.switchoverCountGuage.Set(s.SwitchoverCount)
}
