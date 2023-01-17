package controllers

import (
	"context"

	nauticusiov1alpha1 "github.com/edixos/nauticus/api/v1alpha1"
	"github.com/go-logr/logr"
)

func (s *SpaceReconciler) reconcileSpace(ctx context.Context, space *nauticusiov1alpha1.Space, log logr.Logger) error {
	log.Info("Reconciling Namespace for space.", "Space", space.Name)
	err := s.reconcileNamespace(ctx, space, log)
	if err != nil {
		return err
	}
	log.Info("Reconciling Resource Quota for space", "Space", space.Name)
	err = s.reconcileResourceQuota(ctx, space, log)
	if err != nil {
		return err
	}
	return nil
}
