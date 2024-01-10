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
	"strings"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type mlagPorts struct {
	Disabled      float64
	Configured    float64
	Inactive      float64
	ActivePartial float64 `json:"Active-partial"`
	ActiveFull    float64 `json:"Active-full"`
}

type mlagDetail struct {
	MlagState             string
	PeerMlagState         string
	StateChanges          float64
	LastStateChangeTime   float64
	MlagHwReady           bool
	Failover              bool
	FailoverCauseList     []string
	FailoverInitiated     bool
	SecondaryFromFailover bool
	UdpHeartbeatAlive     bool
}

type MLAGProber struct {
	NegStatus       string
	State           string
	ConfigSanity    string
	PeerLinkStatus  string
	LocalIntfStatus string
	Ports           mlagPorts `json:"mlagPorts"`
	Detail          mlagDetail

	detailGuage               *prometheus.GaugeVec
	stateGuage                *prometheus.GaugeVec
	stateChangesGuage         prometheus.Gauge
	lastStateChangeGuage      prometheus.Gauge
	configSanityGuage         *prometheus.GaugeVec
	negStatusGuage            *prometheus.GaugeVec
	peerLinkStatusGuage       *prometheus.GaugeVec
	localInfStatusStatusGuage *prometheus.GaugeVec
	failoverGuage             *prometheus.GaugeVec
	portGuage                 *prometheus.GaugeVec
}

func (m *MLAGProber) GetCmd() string {
	return "show mlag detail"
}

func (m *MLAGProber) Register(registry *prometheus.Registry) {
	m.detailGuage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "arista_mlag_detail",
		Help: "Contains MLAG state detail",
	}, []string{"mlagState", "peerMlagState"})
	registry.MustRegister(m.detailGuage)

	m.stateGuage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "arista_mlag_state",
		Help: "Contains MLAG state",
	}, []string{"state"})
	registry.MustRegister(m.stateGuage)

	m.stateChangesGuage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "arista_mlag_state_changes",
		Help: "Contains MLAG state changes",
	})
	registry.MustRegister(m.stateChangesGuage)

	m.lastStateChangeGuage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "arista_mlag_last_state_change",
		Help: "Contains MLAG last state change",
	})
	registry.MustRegister(m.lastStateChangeGuage)

	m.configSanityGuage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "arista_mlag_config_sanity",
		Help: "Contains MLAG last state change",
	}, []string{"status"})
	registry.MustRegister(m.configSanityGuage)

	m.negStatusGuage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "arista_mlag_neg_status",
		Help: "Contains MLAG neg status",
	}, []string{"status"})
	registry.MustRegister(m.negStatusGuage)

	m.peerLinkStatusGuage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "arista_mlag_peer_link",
		Help: "Contains MLAG peer link",
	}, []string{"status"})
	registry.MustRegister(m.peerLinkStatusGuage)

	m.localInfStatusStatusGuage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "arista_mlag_local_inf_status",
		Help: "Contains MLAG local inf status",
	}, []string{"status"})
	registry.MustRegister(m.localInfStatusStatusGuage)

	m.failoverGuage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "arista_mlag_failover",
		Help: "Contains MLAG failover",
	}, []string{"cause"})
	registry.MustRegister(m.failoverGuage)

	m.portGuage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "arista_mlag_ports",
		Help: "Contains MLAG port information by state",
	}, []string{"state"})

	for _, lv := range []string{"disabled", "configured", "inactive", "activePartial", "activeFull"} {
		m.portGuage.WithLabelValues(lv)
	}

	registry.MustRegister(m.portGuage)
}

func (m *MLAGProber) Handler(logger log.Logger) {
	m.detailGuage.WithLabelValues(m.Detail.MlagState, m.Detail.PeerMlagState).Set(1)
	if m.State == "active" {
		m.stateGuage.WithLabelValues(m.State).Set(1)
	} else {
		m.stateGuage.WithLabelValues(m.State).Set(0)
	}
	m.stateChangesGuage.Set(m.Detail.StateChanges)
	m.lastStateChangeGuage.Set(m.Detail.LastStateChangeTime)

	if m.ConfigSanity == "consistent" {
		m.configSanityGuage.WithLabelValues(m.ConfigSanity).Set(1)
	} else {
		m.configSanityGuage.WithLabelValues(m.ConfigSanity).Set(0)
	}

	if m.NegStatus == "connected" {
		m.negStatusGuage.WithLabelValues(m.NegStatus).Set(1)
	} else {
		m.negStatusGuage.WithLabelValues(m.NegStatus).Set(0)
	}

	if m.PeerLinkStatus == "up" {
		m.peerLinkStatusGuage.WithLabelValues(m.PeerLinkStatus).Set(1)
	} else {
		m.peerLinkStatusGuage.WithLabelValues(m.PeerLinkStatus).Set(0)
	}

	if m.LocalIntfStatus == "up" {
		m.localInfStatusStatusGuage.WithLabelValues(m.LocalIntfStatus).Set(1)
	} else {
		m.localInfStatusStatusGuage.WithLabelValues(m.LocalIntfStatus).Set(0)
	}

	if m.Detail.Failover {
		m.failoverGuage.WithLabelValues(strings.Join(m.Detail.FailoverCauseList, ",")).Set(1)
	} else {
		m.failoverGuage.WithLabelValues("").Set(0)
	}

	m.portGuage.WithLabelValues("disabled").Set(m.Ports.Disabled)
	m.portGuage.WithLabelValues("configured").Set(m.Ports.Configured)
	m.portGuage.WithLabelValues("inactive").Set(m.Ports.Inactive)
	m.portGuage.WithLabelValues("activePartial").Set(m.Ports.ActivePartial)
	m.portGuage.WithLabelValues("activeFull").Set(m.Ports.ActiveFull)
}
