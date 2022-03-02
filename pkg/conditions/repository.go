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

func SetRepositoryFoundConditionToUnknown(i interface{}, err error) error {
	switch in := i.(type) {
	case invoker.BackupInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:   v1beta1.RepositoryFound,
			Status: core.ConditionUnknown,
			Reason: v1beta1.UnableToCheckRepositoryAvailability,
			Message: fmt.Sprintf("Failed to check whether the Repository %s/%s exist or not. Reason: %v",
				in.GetRepoRef().Namespace,
				in.GetRepoRef().Name,
				err.Error(),
			),
		})
	case invoker.RestoreInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:   v1beta1.RepositoryFound,
			Status: core.ConditionUnknown,
			Reason: v1beta1.UnableToCheckRepositoryAvailability,
			Message: fmt.Sprintf("Failed to check whether the Repository %s/%s exist or not. Reason: %v",
				in.GetRepoRef().Namespace,
				in.GetRepoRef().Name,
				err.Error(),
			),
		})
	default:
		return fmt.Errorf("unable to set %s condition. Reason: invoker type unknown", v1beta1.RepositoryFound)
	}
}

func SetRepositoryFoundConditionToFalse(i interface{}) error {
	switch in := i.(type) {
	case invoker.BackupInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:   v1beta1.RepositoryFound,
			Status: core.ConditionFalse,
			Reason: v1beta1.RepositoryNotAvailable,
			Message: fmt.Sprintf("Repository %s/%s does not exist.",
				in.GetRepoRef().Namespace,
				in.GetRepoRef().Name,
			),
		})
	case invoker.RestoreInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:   v1beta1.RepositoryFound,
			Status: core.ConditionFalse,
			Reason: v1beta1.RepositoryNotAvailable,
			Message: fmt.Sprintf("Repository %s/%s does not exist.",
				in.GetRepoRef().Namespace,
				in.GetRepoRef().Name,
			),
		})
	default:
		return fmt.Errorf("unable to set %s condition. Reason: invoker type unknown", v1beta1.RepositoryFound)
	}
}

func SetRepositoryFoundConditionToTrue(i interface{}) error {
	switch in := i.(type) {
	case invoker.BackupInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:   v1beta1.RepositoryFound,
			Status: core.ConditionTrue,
			Reason: v1beta1.RepositoryAvailable,
			Message: fmt.Sprintf("Repository %s/%s exist.",
				in.GetRepoRef().Namespace,
				in.GetRepoRef().Name,
			),
		})
	case invoker.RestoreInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:   v1beta1.RepositoryFound,
			Status: core.ConditionTrue,
			Reason: v1beta1.RepositoryAvailable,
			Message: fmt.Sprintf("Repository %s/%s exist.",
				in.GetRepoRef().Namespace,
				in.GetRepoRef().Name,
			),
		})
	default:
		return fmt.Errorf("unable to set %s condition. Reason: invoker type unknown", v1beta1.RepositoryFound)
	}
}

func SetValidationPassedToTrue(i interface{}) error {
	switch in := i.(type) {
	case invoker.BackupInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:    v1beta1.ValidationPassed,
			Status:  core.ConditionTrue,
			Reason:  v1beta1.ResourceValidationPassed,
			Message: "Successfully validated.",
		})
	case invoker.RestoreInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:    v1beta1.ValidationPassed,
			Status:  core.ConditionTrue,
			Reason:  v1beta1.ResourceValidationPassed,
			Message: "Successfully validated.",
		})
	default:
		return fmt.Errorf("unable to set %s condition. Reason: invoker type unknown", v1beta1.ValidationPassed)
	}
}

func SetValidationPassedToFalse(i interface{}, err error) error {
	switch in := i.(type) {
	case invoker.BackupInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:    v1beta1.ValidationPassed,
			Status:  core.ConditionFalse,
			Reason:  v1beta1.ResourceValidationFailed,
			Message: err.Error(),
		})
	case invoker.RestoreInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:    v1beta1.ValidationPassed,
			Status:  core.ConditionFalse,
			Reason:  v1beta1.ResourceValidationFailed,
			Message: err.Error(),
		})
	default:
		return fmt.Errorf("unable to set %s condition. Reason: invoker type unknown", v1beta1.ValidationPassed)
	}
}

func SetBackendSecretFoundConditionToUnknown(i interface{}, secretName string, err error) error {
	switch in := i.(type) {
	case invoker.BackupInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:   v1beta1.BackendSecretFound,
			Status: core.ConditionUnknown,
			Reason: v1beta1.UnableToCheckBackendSecretAvailability,
			Message: fmt.Sprintf("Failed to check whether the backend Secret %s/%s exist or not. Reason: %v",
				in.GetRepoRef().Namespace,
				secretName,
				err.Error(),
			),
		})
	case invoker.RestoreInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:   v1beta1.BackendSecretFound,
			Status: core.ConditionUnknown,
			Reason: v1beta1.UnableToCheckBackendSecretAvailability,
			Message: fmt.Sprintf("Failed to check whether the backend Secret %s/%s exist or not. Reason: %v",
				in.GetRepoRef().Namespace,
				secretName,
				err.Error(),
			),
		})
	default:
		return fmt.Errorf("unable to set %s condition. Reason: invoker type unknown", v1beta1.BackendSecretFound)
	}
}

func SetBackendSecretFoundConditionToFalse(i interface{}, secretName string) error {
	switch in := i.(type) {
	case invoker.BackupInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:   v1beta1.BackendSecretFound,
			Status: core.ConditionFalse,
			Reason: v1beta1.BackendSecretNotAvailable,
			Message: fmt.Sprintf("Backend Secret %s/%s does not exist.",
				in.GetRepoRef().Namespace,
				secretName,
			),
		})
	case invoker.RestoreInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:   v1beta1.BackendSecretFound,
			Status: core.ConditionFalse,
			Reason: v1beta1.BackendSecretNotAvailable,
			Message: fmt.Sprintf("Backend Secret %s/%s does not exist.",
				in.GetRepoRef().Namespace,
				secretName,
			),
		})
	default:
		return fmt.Errorf("unable to set %s condition. Reason: invoker type unknown", v1beta1.BackendSecretFound)
	}
}

func SetBackendSecretFoundConditionToTrue(i interface{}, secretName string) error {
	switch in := i.(type) {
	case invoker.BackupInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:   v1beta1.BackendSecretFound,
			Status: core.ConditionTrue,
			Reason: v1beta1.BackendSecretAvailable,
			Message: fmt.Sprintf("Backend Secret %s/%s exist.",
				in.GetRepoRef().Namespace,
				secretName,
			),
		})
	case invoker.RestoreInvoker:
		return in.SetCondition(nil, kmapi.Condition{
			Type:   v1beta1.BackendSecretFound,
			Status: core.ConditionTrue,
			Reason: v1beta1.BackendSecretAvailable,
			Message: fmt.Sprintf("Backend Secret %s/%s exist.",
				in.GetRepoRef().Namespace,
				secretName,
			),
		})
	default:
		return fmt.Errorf("unable to set %s condition. Reason: invoker type unknown", v1beta1.BackendSecretFound)
	}
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
