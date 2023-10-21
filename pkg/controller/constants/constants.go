// Copyright 2022-2023 Edixos
// SPDX-License-Identifier: Apache-2.0
package constants

import (
	"time"
)

const (
	NauticusSpaceFinalizer = "nauticus.io/finalizer"
	RequeueAfter           = time.Minute * 3

	SpaceConditionReady    string = "Ready"
	SpaceConditionCreating string = "Creating"
	SpaceConditionFailed   string = "Failed"

	SpaceSyncSuccessReason string = "SpaceSyncedSuccessfully"
	SpaceCreatingReason    string = "SpaceCreating"
	SpaceFailedReason      string = "SpaceSyncFailed"

	SpaceSyncSuccessMessage string = "Space synced successfully"
	SpaceSyncFailMessage    string = "Space failed to sync"
	SpaceCreatingMessage    string = "Creating Space in progress"

	SpaceTplConditionReady string = "Ready"

	SpaceTplSyncSuccessReason string = "SpaceTemplateSyncedSuccessfully"

	SpaceTplSyncSuccessMessage string = "SpaceTemplate synced successfully"
	SpaceTplCreatingMessage    string = "Creating SpaceTemplate in progress"
)
