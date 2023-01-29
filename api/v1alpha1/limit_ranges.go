package v1alpha1

import corev1 "k8s.io/api/core/v1"

type LimitRangesSpec struct {
	Items []corev1.LimitRangeSpec `json:"items,omitempty"`
}
