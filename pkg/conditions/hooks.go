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

	"stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	cs "stash.appscode.dev/apimachinery/client/clientset/versioned"
	"stash.appscode.dev/apimachinery/pkg/invoker"

	core "k8s.io/api/core/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

func SetGlobalPreBackupHookSucceededConditionToFalse(stashClient cs.Interface, backupSession *v1beta1.BackupSession, hookErr error) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.GlobalPreBackupHookSucceeded,
				Status:  core.ConditionFalse,
				Reason:  v1beta1.GlobalPreBackupHookExecutionFailed,
				Message: fmt.Sprintf("Failed to execute global PreBackup Hook. Reason: %v.", hookErr),
			},
		},
	})
}

func SetGlobalPreBackupHookSucceededConditionToTrue(stashClient cs.Interface, backupSession *v1beta1.BackupSession) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.GlobalPreBackupHookSucceeded,
				Status:  core.ConditionTrue,
				Reason:  v1beta1.GlobalPreBackupHookExecutedSuccessfully,
				Message: "Global PreBackup hook has been executed successfully",
			},
		},
	})
}

func SetGlobalPostBackupHookSucceededConditionToFalse(stashClient cs.Interface, backupSession *v1beta1.BackupSession, hookErr error) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.GlobalPostBackupHookSucceeded,
				Status:  core.ConditionFalse,
				Reason:  v1beta1.GlobalPostBackupHookExecutionFailed,
				Message: fmt.Sprintf("Failed to execute global PostBackup Hook. Reason: %v.", hookErr),
			},
		},
	})
}

func SetGlobalPostBackupHookSucceededConditionToTrue(stashClient cs.Interface, backupSession *v1beta1.BackupSession) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.GlobalPostBackupHookSucceeded,
				Status:  core.ConditionTrue,
				Reason:  v1beta1.GlobalPostBackupHookExecutedSuccessfully,
				Message: "Global PostBackup hook has been executed successfully",
			},
		},
	})
}

func SetGlobalPreRestoreHookSucceededConditionToFalse(invoker invoker.RestoreInvoker, hookErr error) error {
	return invoker.SetCondition(nil, kmapi.Condition{
		Type:    v1beta1.GlobalPreRestoreHookSucceeded,
		Status:  core.ConditionFalse,
		Reason:  v1beta1.GlobalPreRestoreHookExecutionFailed,
		Message: fmt.Sprintf("Failed to execute global PreRestore Hook. Reason: %v.", hookErr),
	})
}

func SetGlobalPreRestoreHookSucceededConditionToTrue(invoker invoker.RestoreInvoker) error {
	return invoker.SetCondition(nil, kmapi.Condition{
		Type:    v1beta1.GlobalPreRestoreHookSucceeded,
		Status:  core.ConditionTrue,
		Reason:  v1beta1.GlobalPreRestoreHookExecutedSuccessfully,
		Message: "Global PreRestore hook has been executed successfully",
	})
}

func SetGlobalPostRestoreHookSucceededConditionToFalse(invoker invoker.RestoreInvoker, hookErr error) error {
	return invoker.SetCondition(nil, kmapi.Condition{
		Type:    v1beta1.GlobalPostRestoreHookSucceeded,
		Status:  core.ConditionFalse,
		Reason:  v1beta1.GlobalPostRestoreHookExecutionFailed,
		Message: fmt.Sprintf("Failed to execute global PostRestore Hook. Reason: %v.", hookErr),
	})
}

func SetGlobalPostRestoreHookSucceededConditionToTrue(invoker invoker.RestoreInvoker) error {
	return invoker.SetCondition(nil, kmapi.Condition{
		Type:    v1beta1.GlobalPostRestoreHookSucceeded,
		Status:  core.ConditionTrue,
		Reason:  v1beta1.GlobalPostRestoreHookExecutedSuccessfully,
		Message: "Global PostRestore hook has been executed successfully",
	})
}
