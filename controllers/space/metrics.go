// Copyright 2023-2024 Edixos
// SPDX-License-Identifier: Apache-2.0

package space

import (
	"github.com/edixos/nauticus/api/v1alpha1"
	"github.com/edixos/nauticus/pkg/controller/constants"
	"github.com/edixos/nauticus/pkg/metrics"
)

func (r *Reconciler) setMetrics(space *v1alpha1.Space, conditionType v1alpha1.ConditionType) {
	switch conditionType {
	case v1alpha1.ConditionType(constants.SpaceConditionCreating):
		metrics.ReadySpaces.WithLabelValues(space.Name).Set(0)
		metrics.InProgressSpaces.WithLabelValues(space.Name).Set(1)
		metrics.FailedSpaces.WithLabelValues(space.Name).Set(0)
	case v1alpha1.ConditionType(constants.SpaceConditionReady):
		metrics.ReadySpaces.WithLabelValues(space.Name).Set(1)
		metrics.InProgressSpaces.WithLabelValues(space.Name).Set(0)
		metrics.FailedSpaces.WithLabelValues(space.Name).Set(0)
	case v1alpha1.ConditionType(constants.SpaceConditionFailed):
		metrics.ReadySpaces.WithLabelValues(space.Name).Set(0)
		metrics.InProgressSpaces.WithLabelValues(space.Name).Set(0)
		metrics.FailedSpaces.WithLabelValues(space.Name).Set(1)
	}
}
