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
	cs "stash.appscode.dev/apimachinery/client/clientset/versioned"
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

func SetBackupSkippedConditionToTrue(stashClient cs.Interface, backupSession *v1beta1.BackupSession, msg string) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.BackupSkipped,
				Status:  core.ConditionTrue,
				Reason:  v1beta1.SkippedTakingNewBackup,
				Message: msg,
			},
		},
	})
}

func SetBackupMetricsPushedConditionToFalse(stashClient cs.Interface, backupSession *v1beta1.BackupSession, err error) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.MetricsPushed,
				Status:  core.ConditionFalse,
				Reason:  v1beta1.FailedToPushMetrics,
				Message: fmt.Sprintf("Failed to push metrics. Reason: %v", err.Error()),
			},
		},
	})
}

func SetBackupMetricsPushedConditionToTrue(stashClient cs.Interface, backupSession *v1beta1.BackupSession) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.MetricsPushed,
				Status:  core.ConditionTrue,
				Reason:  v1beta1.SuccessfullyPushedMetrics,
				Message: "Successfully pushed metrics.",
			},
		},
	})
}

func SetBackupHistoryCleanedConditionToFalse(stashClient cs.Interface, backupSession *v1beta1.BackupSession, err error) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.BackupHistoryCleaned,
				Status:  core.ConditionFalse,
				Reason:  v1beta1.FailedToCleanBackupHistory,
				Message: fmt.Sprintf("Failed to cleanup old BackupSessions. Reason: %v", err.Error()),
			},
		},
	})
}

func SetBackupHistoryCleanedConditionToTrue(stashClient cs.Interface, backupSession *v1beta1.BackupSession) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.BackupHistoryCleaned,
				Status:  core.ConditionTrue,
				Reason:  v1beta1.SuccessfullyCleanedBackupHistory,
				Message: "Successfully cleaned up backup history according to backupHistoryLimit.",
			},
		},
	})
}

func SetBackupExecutorEnsuredToFalse(stashClient cs.Interface, backupSession *v1beta1.BackupSession, err error) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.BackupExecutorEnsured,
				Status:  core.ConditionFalse,
				Reason:  v1beta1.FailedToEnsureBackupExecutor,
				Message: fmt.Sprintf("Failed to ensure backup executor. Reason: %v", err.Error()),
			},
		},
	})
}

func SetBackupExecutorEnsuredToTrue(stashClient cs.Interface, backupSession *v1beta1.BackupSession) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.BackupExecutorEnsured,
				Status:  core.ConditionTrue,
				Reason:  v1beta1.SuccessfullyEnsuredBackupExecutor,
				Message: "Successfully ensured backup executor.",
			},
		},
	})
}

func SetPreBackupHookExecutionSucceededToFalse(stashClient cs.Interface, backupSession *v1beta1.BackupSession, err error) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.PreBackupHookExecutionSucceeded,
				Status:  core.ConditionFalse,
				Reason:  v1beta1.FailedToExecutePreBackupHook,
				Message: fmt.Sprintf("Failed to execute preBackup hook. Reason: %v", err.Error()),
			},
		},
	})
}

func SetPreBackupHookExecutionSucceededToTrue(stashClient cs.Interface, backupSession *v1beta1.BackupSession) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.PreBackupHookExecutionSucceeded,
				Status:  core.ConditionTrue,
				Reason:  v1beta1.SuccessfullyExecutedPreBackupHook,
				Message: "Successfully executed preBackup hook.",
			},
		},
	})
}

func SetPostBackupHookExecutionSucceededToFalse(stashClient cs.Interface, backupSession *v1beta1.BackupSession, err error) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.PostBackupHookExecutionSucceeded,
				Status:  core.ConditionFalse,
				Reason:  v1beta1.FailedToExecutePostBackupHook,
				Message: fmt.Sprintf("Failed to execute postBackup hook. Reason: %v", err.Error()),
			},
		},
	})
}

func SetPostBackupHookExecutionSucceededToTrue(stashClient cs.Interface, backupSession *v1beta1.BackupSession) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.PostBackupHookExecutionSucceeded,
				Status:  core.ConditionTrue,
				Reason:  v1beta1.SuccessfullyExecutedPostBackupHook,
				Message: "Successfully executed postBackup hook.",
			},
		},
	})
}
