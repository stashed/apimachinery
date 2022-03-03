/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package conditions

import (
	"fmt"
	"strings"

	"stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	"stash.appscode.dev/apimachinery/pkg/invoker"

	core "k8s.io/api/core/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

func SetBackupTargetFoundConditionToUnknown(invoker invoker.BackupInvoker, tref v1beta1.TargetRef, err error) error {
	return invoker.SetCondition(&tref, kmapi.Condition{
		Type:   v1beta1.BackupTargetFound,
		Status: core.ConditionUnknown,
		Reason: v1beta1.UnableToCheckTargetAvailability,
		Message: fmt.Sprintf("Failed to check whether backup target %s %s/%s exist or not. Reason: %v",
			tref.APIVersion,
			strings.ToLower(tref.Kind),
			tref.Name,
			err.Error(),
		),
	})
}

func SetBackupTargetFoundConditionToFalse(invoker invoker.BackupInvoker, tref v1beta1.TargetRef) error {
	return invoker.SetCondition(&tref, kmapi.Condition{
		// Set the "BackupTargetFound" condition to "False"
		Type:   v1beta1.BackupTargetFound,
		Status: core.ConditionFalse,
		Reason: v1beta1.TargetNotAvailable,
		Message: fmt.Sprintf("Backup target %s %s/%s does not exist.",
			tref.APIVersion,
			strings.ToLower(tref.Kind),
			tref.Name,
		),
	})
}

func SetBackupTargetFoundConditionToTrue(invoker invoker.BackupInvoker, tref v1beta1.TargetRef) error {
	return invoker.SetCondition(&tref, kmapi.Condition{
		Type:   v1beta1.BackupTargetFound,
		Status: core.ConditionTrue,
		Reason: v1beta1.TargetAvailable,
		Message: fmt.Sprintf("Backup target %s %s/%s found.",
			tref.APIVersion,
			strings.ToLower(tref.Kind),
			tref.Name,
		),
	})
}

func SetCronJobCreatedConditionToFalse(invoker invoker.BackupInvoker, err error) error {
	return invoker.SetCondition(nil, kmapi.Condition{
		Type:    v1beta1.CronJobCreated,
		Status:  core.ConditionFalse,
		Reason:  v1beta1.CronJobCreationFailed,
		Message: fmt.Sprintf("Failed to create backup triggering CronJob. Reason: %v", err.Error()),
	})
}

func SetCronJobCreatedConditionToTrue(invoker invoker.BackupInvoker) error {
	return invoker.SetCondition(nil, kmapi.Condition{
		Type:    v1beta1.CronJobCreated,
		Status:  core.ConditionTrue,
		Reason:  v1beta1.CronJobCreationSucceeded,
		Message: "Successfully created backup triggering CronJob.",
	})
}

func SetSidecarInjectedConditionToTrue(invoker invoker.BackupInvoker, tref v1beta1.TargetRef) error {
	return invoker.SetCondition(&tref, kmapi.Condition{
		Type:   v1beta1.StashSidecarInjected,
		Status: core.ConditionTrue,
		Reason: v1beta1.SidecarInjectionSucceeded,
		Message: fmt.Sprintf("Successfully injected stash sidecar into %s %s/%s",
			tref.APIVersion,
			strings.ToLower(tref.Kind),
			tref.Name,
		),
	})
}

func SetSidecarInjectedConditionToFalse(invoker invoker.BackupInvoker, tref v1beta1.TargetRef, err error) error {
	return invoker.SetCondition(&tref, kmapi.Condition{
		Type:   v1beta1.StashSidecarInjected,
		Status: core.ConditionFalse,
		Reason: v1beta1.SidecarInjectionFailed,
		Message: fmt.Sprintf("Failed to inject stash sidecar into %s %s/%s. Reason: %v",
			tref.APIVersion,
			strings.ToLower(tref.Kind),
			tref.Name,
			err.Error(),
		),
	})
}
