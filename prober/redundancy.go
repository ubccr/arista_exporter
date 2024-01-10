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

type RedundancyProber struct {
	SlotID                         float64
	MyMode                         string
	PeerMode                       string
	UnitDesc                       string
	CommunicationDesc              string
	PeerState                      string
	SwitchoverReady                bool
	AllAgentSsoReady               bool
	LastRedundancyModeChangeTime   float64
	LastRedundancyModeChangeReason string

	slotIDGuage                       *prometheus.GaugeVec
	myModeGuage                       *prometheus.GaugeVec
	peerModeGuage                     *prometheus.GaugeVec
	communicationDescGuage            *prometheus.GaugeVec
	switchoverReadyGuage              prometheus.Gauge
	allAgentSsoReadyGuage             prometheus.Gauge
	lastRedundancyModeChangeTimeGuage *prometheus.GaugeVec
}

func (r *RedundancyProber) GetCmd() string {
	return "show redundancy status"
}

func (r *RedundancyProber) Register(registry *prometheus.Registry) {
	r.slotIDGuage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "arista_redundancy_slot_id",
		Help: "Contains redundancy slot id",
	}, []string{"unitDesc"})
	registry.MustRegister(r.slotIDGuage)

	r.myModeGuage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "arista_redundancy_mode",
		Help: "Contains redundancy mode",
	}, []string{"status"})
	registry.MustRegister(r.myModeGuage)

	r.peerModeGuage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "arista_redundancy_peer_mode",
		Help: "Contains redundancy peer mode",
	}, []string{"status"})
	registry.MustRegister(r.peerModeGuage)

	r.communicationDescGuage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "arista_redundancy_communication_desc",
		Help: "Contains redundancy communication desc",
	}, []string{"status"})
	registry.MustRegister(r.communicationDescGuage)

	r.switchoverReadyGuage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "arista_redundancy_switchover_ready",
		Help: "Contains redundancy switchover ready",
	})
	registry.MustRegister(r.switchoverReadyGuage)

	r.allAgentSsoReadyGuage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "arista_redundancy_all_agent_sso_ready",
		Help: "Contains redundancy all Agent SSO Ready",
	})
	registry.MustRegister(r.allAgentSsoReadyGuage)

	r.lastRedundancyModeChangeTimeGuage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "arista_redundancy_last_mode_change_time",
		Help: "Contains redundancy last mode change time",
	}, []string{"reason"})
	registry.MustRegister(r.lastRedundancyModeChangeTimeGuage)
}

func (r *RedundancyProber) Handler(logger log.Logger) {
	r.slotIDGuage.WithLabelValues(r.UnitDesc).Set(r.SlotID)
	if r.MyMode == "active" {
		r.myModeGuage.WithLabelValues(r.MyMode).Set(1)
	} else {
		r.myModeGuage.WithLabelValues(r.MyMode).Set(0)
	}
	if r.PeerMode == "standby" {
		r.peerModeGuage.WithLabelValues(r.PeerMode).Set(1)
	} else {
		r.peerModeGuage.WithLabelValues(r.PeerMode).Set(0)
	}
	if r.CommunicationDesc == "Up" {
		r.communicationDescGuage.WithLabelValues(r.CommunicationDesc).Set(1)
	} else {
		r.communicationDescGuage.WithLabelValues(r.CommunicationDesc).Set(0)
	}
	if r.SwitchoverReady {
		r.switchoverReadyGuage.Set(1)
	}
	if r.AllAgentSsoReady {
		r.allAgentSsoReadyGuage.Set(1)
	}

	r.lastRedundancyModeChangeTimeGuage.WithLabelValues(r.LastRedundancyModeChangeReason).Set(r.LastRedundancyModeChangeTime)
}
