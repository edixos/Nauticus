// Copyright 2022-2023 Edixos
// SPDX-License-Identifier: Apache-2.0
package spacetemplate

import (
	"github.com/edixos/nauticus/api/v1alpha1"
	"github.com/edixos/nauticus/pkg/controller/constants"
	"github.com/edixos/nauticus/pkg/metrics"
)

func (r *Reconciler) setMetrics(spacetpl *v1alpha1.SpaceTemplate, conditionType v1alpha1.ConditionType) {
	switch conditionType {
	case v1alpha1.ConditionType(constants.SpaceTplConditionReady):
		metrics.ReadySpaceTemplates.WithLabelValues(spacetpl.Name).Set(0)
		metrics.InProgresSpaceTemplates.WithLabelValues(spacetpl.Name).Set(0)
		metrics.FailedSpaceTemplates.WithLabelValues(spacetpl.Name).Set(0)
	case v1alpha1.ConditionType(constants.SpaceTplConditionCreating):
		metrics.ReadySpaceTemplates.WithLabelValues(spacetpl.Name).Set(0)
		metrics.InProgresSpaceTemplates.WithLabelValues(spacetpl.Name).Set(0)
		metrics.FailedSpaceTemplates.WithLabelValues(spacetpl.Name).Set(0)
	case v1alpha1.ConditionType(constants.SpaceTplConditionFailed):
		metrics.ReadySpaceTemplates.WithLabelValues(spacetpl.Name).Set(0)
		metrics.InProgresSpaceTemplates.WithLabelValues(spacetpl.Name).Set(0)
		metrics.FailedSpaceTemplates.WithLabelValues(spacetpl.Name).Set(0)
	}
}
