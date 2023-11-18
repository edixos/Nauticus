// Copyright 2022-2023 Edixos
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	FailedSpaces = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "not_ready_spaces",
			Help: "Not Ready Spaces",
		}, []string{"name"},
	)
	ReadySpaces = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ready_spaces",
			Help: "Ready Spaces",
		}, []string{"name"},
	)
	InProgressSpaces = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "in_progress_spaces",
			Help: "In Progress Spaces",
		}, []string{"name"},
	)

	ReadySpaceTemplates = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ready_spacetemplates",
			Help: "ready SpaceTemplates",
		}, []string{"name"},
	)

	FailedSpaceTemplates = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "not_ready_spacetemplates",
			Help: "not ready SpaceTemplates",
		}, []string{"name"},
	)
	InProgresSpaceTemplates = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "in_progress_spacetemplates",
			Help: "In Progress SpaceTemplates",
		}, []string{"name"},
	)
)

func init() {
	metrics.Registry.MustRegister(FailedSpaces, ReadySpaces, InProgressSpaces, ReadySpaceTemplates, FailedSpaceTemplates, InProgresSpaceTemplates)
}
