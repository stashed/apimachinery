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

func SetBackendRepositoryInitializedConditionToFalse(stashClient cs.Interface, backupSession *v1beta1.BackupSession, err error) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.BackendRepositoryInitialized,
				Status:  core.ConditionFalse,
				Reason:  v1beta1.FailedToInitializeBackendRepository,
				Message: fmt.Sprintf("Failed to initialize backend repository. Reason: %v", err.Error()),
			},
		},
	})
}

func SetBackendRepositoryInitializedConditionToTrue(stashClient cs.Interface, backupSession *v1beta1.BackupSession) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.BackendRepositoryInitialized,
				Status:  core.ConditionTrue,
				Reason:  v1beta1.BackendRepositoryFound,
				Message: "Repository exist in the backend.",
			},
		},
	})
}

func SetBackupExecutorEnsuredToFalse(stashClient cs.Interface, backupSession *v1beta1.BackupSession, target v1beta1.TargetRef, err error) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Targets: []v1beta1.BackupTargetStatus{
			{
				Ref: target,
				Conditions: []kmapi.Condition{
					{
						Type:    v1beta1.BackupExecutorEnsured,
						Status:  core.ConditionFalse,
						Reason:  v1beta1.FailedToEnsureBackupExecutor,
						Message: fmt.Sprintf("Failed to ensure backup executor. Reason: %v", err.Error()),
					},
				},
			},
		},
	})
}

func SetBackupExecutorEnsuredToTrue(stashClient cs.Interface, backupSession *v1beta1.BackupSession, target v1beta1.TargetRef) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Targets: []v1beta1.BackupTargetStatus{
			{
				Ref: target,
				Conditions: []kmapi.Condition{
					{
						Type:    v1beta1.BackupExecutorEnsured,
						Status:  core.ConditionTrue,
						Reason:  v1beta1.SuccessfullyEnsuredBackupExecutor,
						Message: "Successfully ensured backup executor.",
					},
				},
			},
		},
	})
}

func SetPreBackupHookExecutionSucceededToFalse(stashClient cs.Interface, backupSession *v1beta1.BackupSession, target v1beta1.TargetRef, err error) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Targets: []v1beta1.BackupTargetStatus{
			{
				Ref: target,
				Conditions: []kmapi.Condition{
					{
						Type:    v1beta1.PreBackupHookExecutionSucceeded,
						Status:  core.ConditionFalse,
						Reason:  v1beta1.FailedToExecutePreBackupHook,
						Message: fmt.Sprintf("Failed to execute preBackup hook. Reason: %v", err.Error()),
					},
				},
			},
		},
	})
}

func SetPreBackupHookExecutionSucceededToTrue(stashClient cs.Interface, backupSession *v1beta1.BackupSession, target v1beta1.TargetRef) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Targets: []v1beta1.BackupTargetStatus{
			{
				Ref: target,
				Conditions: []kmapi.Condition{
					{
						Type:    v1beta1.PreBackupHookExecutionSucceeded,
						Status:  core.ConditionTrue,
						Reason:  v1beta1.SuccessfullyExecutedPreBackupHook,
						Message: "Successfully executed preBackup hook.",
					},
				},
			},
		},
	})
}

func SetPostBackupHookExecutionSucceededToFalse(stashClient cs.Interface, backupSession *v1beta1.BackupSession, target v1beta1.TargetRef, err error) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Targets: []v1beta1.BackupTargetStatus{
			{
				Ref: target,
				Conditions: []kmapi.Condition{
					{
						Type:    v1beta1.PostBackupHookExecutionSucceeded,
						Status:  core.ConditionFalse,
						Reason:  v1beta1.FailedToExecutePostBackupHook,
						Message: fmt.Sprintf("Failed to execute postBackup hook. Reason: %v", err.Error()),
					},
				},
			},
		},
	})
}

func SetPostBackupHookExecutionSucceededToTrue(stashClient cs.Interface, backupSession *v1beta1.BackupSession, target v1beta1.TargetRef) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Targets: []v1beta1.BackupTargetStatus{
			{
				Ref: target,
				Conditions: []kmapi.Condition{
					{
						Type:    v1beta1.PostBackupHookExecutionSucceeded,
						Status:  core.ConditionTrue,
						Reason:  v1beta1.SuccessfullyExecutedPostBackupHook,
						Message: "Successfully executed postBackup hook.",
					},
				},
			},
		},
	})
}

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

func SetRetentionPolicyAppliedConditionToFalse(stashClient cs.Interface, backupSession *v1beta1.BackupSession, err error) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.RetentionPolicyApplied,
				Status:  core.ConditionFalse,
				Reason:  v1beta1.FailedToApplyRetentionPolicy,
				Message: fmt.Sprintf("Failed to apply retention policy. Reason: %v", err.Error()),
			},
		},
	})
}

func SetRetentionPolicyAppliedConditionToTrue(stashClient cs.Interface, backupSession *v1beta1.BackupSession) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.RetentionPolicyApplied,
				Status:  core.ConditionTrue,
				Reason:  v1beta1.SuccessfullyAppliedRetentionPolicy,
				Message: "Successfully applied retention policy.",
			},
		},
	})
}

func SetRepositoryIntegrityVerifiedConditionToFalse(stashClient cs.Interface, backupSession *v1beta1.BackupSession, err error) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.RepositoryIntegrityVerified,
				Status:  core.ConditionFalse,
				Reason:  v1beta1.FailedToVerifyRepositoryIntegrity,
				Message: fmt.Sprintf("Repository integrity verification failed. Reason: %v", err.Error()),
			},
		},
	})
}

func SetRepositoryIntegrityVerifiedConditionToTrue(stashClient cs.Interface, backupSession *v1beta1.BackupSession) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.RepositoryIntegrityVerified,
				Status:  core.ConditionTrue,
				Reason:  v1beta1.SuccessfullyVerifiedRepositoryIntegrity,
				Message: "Repository integrity verification succeeded.",
			},
		},
	})
}

func SetRepositoryMetricsPushedConditionToFalse(stashClient cs.Interface, backupSession *v1beta1.BackupSession, err error) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.RepositoryMetricsPushed,
				Status:  core.ConditionFalse,
				Reason:  v1beta1.FailedToPushRepositoryMetrics,
				Message: fmt.Sprintf("Failed to push repository metrics. Reason: %v", err.Error()),
			},
		},
	})
}

func SetRepositoryMetricsPushedConditionToTrue(stashClient cs.Interface, backupSession *v1beta1.BackupSession) (*v1beta1.BackupSession, error) {
	return invoker.UpdateBackupSessionStatus(stashClient, backupSession.ObjectMeta, &v1beta1.BackupSessionStatus{
		Conditions: []kmapi.Condition{
			{
				Type:    v1beta1.RepositoryMetricsPushed,
				Status:  core.ConditionTrue,
				Reason:  v1beta1.SuccessfullyPushedRepositoryMetrics,
				Message: "Successfully pushed repository metrics.",
			},
		},
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
