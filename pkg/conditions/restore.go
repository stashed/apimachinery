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

func SetRestoreTargetFoundConditionToTrue(inv invoker.RestoreInvoker, index int) error {
	target := inv.GetTargetInfo()[index].Target
	return inv.SetCondition(&target.Ref, kmapi.Condition{
		Type:   v1beta1.RestoreTargetFound,
		Status: core.ConditionTrue,
		Reason: v1beta1.TargetAvailable,
		Message: fmt.Sprintf("Restore target %s %s/%s found.",
			target.Ref.APIVersion,
			strings.ToLower(target.Ref.Kind),
			target.Ref.Name,
		),
	})
}

func SetRestoreTargetFoundConditionToFalse(inv invoker.RestoreInvoker, index int) error {
	target := inv.GetTargetInfo()[index].Target
	return inv.SetCondition(&target.Ref, kmapi.Condition{
		Type:   v1beta1.RestoreTargetFound,
		Status: core.ConditionFalse,
		Reason: v1beta1.TargetNotAvailable,
		Message: fmt.Sprintf("Restore target %s %s/%s does not exist.",
			target.Ref.APIVersion,
			strings.ToLower(target.Ref.Kind),
			target.Ref.Name,
		),
	})
}

func SetRestoreTargetFoundConditionToUnknown(inv invoker.RestoreInvoker, index int, err error) error {
	target := inv.GetTargetInfo()[index].Target
	return inv.SetCondition(&target.Ref, kmapi.Condition{
		Type:   v1beta1.RestoreTargetFound,
		Status: core.ConditionUnknown,
		Reason: v1beta1.UnableToCheckTargetAvailability,
		Message: fmt.Sprintf("Failed to check whether restore target %s %s/%s exist or not. Reason: %v",
			target.Ref.APIVersion,
			strings.ToLower(target.Ref.Kind),
			target.Ref.Name,
			err,
		),
	})
}

func SetRestoreJobCreatedConditionToTrue(inv invoker.RestoreInvoker, tref *v1beta1.TargetRef) error {
	return inv.SetCondition(tref, kmapi.Condition{
		Type:    v1beta1.RestoreJobCreated,
		Status:  core.ConditionTrue,
		Reason:  v1beta1.RestoreJobCreationSucceeded,
		Message: "Successfully created restore job.",
	})
}

func SetRestoreJobCreatedConditionToFalse(inv invoker.RestoreInvoker, tref *v1beta1.TargetRef, err error) error {
	return inv.SetCondition(tref, kmapi.Condition{
		Type:    v1beta1.RestoreJobCreated,
		Status:  core.ConditionFalse,
		Reason:  v1beta1.RestoreJobCreationFailed,
		Message: fmt.Sprintf("Failed to create restore job. Reason: %v", err.Error()),
	})
}

func SetInitContainerInjectedConditionToTrue(inv invoker.RestoreInvoker, tref *v1beta1.TargetRef) error {
	return inv.SetCondition(tref, kmapi.Condition{
		Type:    v1beta1.StashInitContainerInjected,
		Status:  core.ConditionTrue,
		Reason:  v1beta1.InitContainerInjectionSucceeded,
		Message: "Successfully injected stash init-container.",
	})
}

func SetInitContainerInjectedConditionToFalse(inv invoker.RestoreInvoker, tref *v1beta1.TargetRef, err error) error {
	return inv.SetCondition(tref, kmapi.Condition{
		Type:    v1beta1.StashInitContainerInjected,
		Status:  core.ConditionFalse,
		Reason:  v1beta1.InitContainerInjectionFailed,
		Message: fmt.Sprintf("Failed to inject Stash init-container. Reason: %v", err.Error()),
	})
}

func SetRestoreCompletedConditionToTrue(inv invoker.RestoreInvoker, tref *v1beta1.TargetRef, msg string) error {
	return inv.SetCondition(tref, kmapi.Condition{
		Type:    v1beta1.RestoreCompleted,
		Status:  core.ConditionTrue,
		Reason:  "PostRestoreTasksExecuted",
		Message: msg,
	})
}

func SetRestoreCompletedConditionToFalse(inv invoker.RestoreInvoker, tref *v1beta1.TargetRef, msg string) error {
	return inv.SetCondition(tref, kmapi.Condition{
		Type:    v1beta1.RestoreCompleted,
		Status:  core.ConditionFalse,
		Reason:  "PostRestoreTasksNotExecuted",
		Message: msg,
	})
}

func SetRestoreExecutorEnsuredToTrue(inv invoker.RestoreInvoker, tref *v1beta1.TargetRef, msg string) error {
	return inv.SetCondition(tref, kmapi.Condition{
		Type:    v1beta1.RestoreExecutorEnsured,
		Status:  core.ConditionTrue,
		Reason:  v1beta1.SuccessfullyEnsuredRestoreExecutor,
		Message: msg,
	})
}

func SetRestoreExecutorEnsuredToFalse(inv invoker.RestoreInvoker, tref *v1beta1.TargetRef, msg string) error {
	return inv.SetCondition(tref, kmapi.Condition{
		Type:    v1beta1.RestoreExecutorEnsured,
		Status:  core.ConditionFalse,
		Reason:  v1beta1.FailedToEnsureRestoreExecutor,
		Message: msg,
	})
}

func SetRestoreMetricsPushedConditionToFalse(inv invoker.RestoreInvoker, err error) error {
	return inv.SetCondition(nil, kmapi.Condition{
		Type:    v1beta1.MetricsPushed,
		Status:  core.ConditionFalse,
		Reason:  v1beta1.FailedToPushMetrics,
		Message: fmt.Sprintf("Failed to push metrics. Reason: %v", err.Error()),
	})
}

func SetRestoreMetricsPushedConditionToTrue(inv invoker.RestoreInvoker) error {
	return inv.SetCondition(nil, kmapi.Condition{
		Type:    v1beta1.MetricsPushed,
		Status:  core.ConditionTrue,
		Reason:  v1beta1.SuccessfullyPushedMetrics,
		Message: "Successfully pushed metrics.",
	})
}

func SetPreRestoreHookExecutionSucceededToFalse(inv invoker.RestoreInvoker, err error) error {
	return inv.SetCondition(nil, kmapi.Condition{
		Type:    v1beta1.PreRestoreHookExecutionSucceeded,
		Status:  core.ConditionFalse,
		Reason:  v1beta1.FailedToExecutePreRestoreHook,
		Message: fmt.Sprintf("Failed to execute preRestore hook. Reason: %v", err.Error()),
	})
}

func SetPreRestoreHookExecutionSucceededToTrue(inv invoker.RestoreInvoker) error {
	return inv.SetCondition(nil, kmapi.Condition{
		Type:    v1beta1.PreRestoreHookExecutionSucceeded,
		Status:  core.ConditionTrue,
		Reason:  v1beta1.SuccessfullyExecutedPreRestoreHook,
		Message: "Successfully executed preRestore hook.",
	})
}

func SetPostRestoreHookExecutionSucceededToFalse(inv invoker.RestoreInvoker, err error) error {
	return inv.SetCondition(nil, kmapi.Condition{
		Type:    v1beta1.PostRestoreHookExecutionSucceeded,
		Status:  core.ConditionFalse,
		Reason:  v1beta1.FailedToExecutePostRestoreHook,
		Message: fmt.Sprintf("Failed to execute postRestore hook. Reason: %v", err.Error()),
	})
}

func SetPostRestoreHookExecutionSucceededToTrue(inv invoker.RestoreInvoker) error {
	return inv.SetCondition(nil, kmapi.Condition{
		Type:    v1beta1.PostRestoreHookExecutionSucceeded,
		Status:  core.ConditionTrue,
		Reason:  v1beta1.SuccessfullyExecutedPostRestoreHook,
		Message: "Successfully executed postRestore hook.",
	})
}
