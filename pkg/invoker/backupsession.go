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

package invoker

import (
	"context"
	"time"

	"stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	cs "stash.appscode.dev/apimachinery/client/clientset/versioned"
	stash_util "stash.appscode.dev/apimachinery/client/clientset/versioned/typed/stash/v1beta1/util"

	"gomodules.xyz/x/arrays"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kmapi "kmodules.xyz/client-go/api/v1"
)

func UpdateBackupSessionStatus(stashClient cs.Interface, meta metav1.ObjectMeta, status *v1beta1.BackupSessionStatus) (*v1beta1.BackupSession, error) {
	return stash_util.UpdateBackupSessionStatus(
		context.TODO(),
		stashClient.StashV1beta1(),
		meta,
		func(in *v1beta1.BackupSessionStatus) (types.UID, *v1beta1.BackupSessionStatus) {
			in.Conditions = upsertConditions(in.Conditions, status.Conditions)

			if len(status.Targets) > 0 {
				for i := range status.Targets {
					in.Targets = upsertBackupMembersStatus(in.Targets, status.Targets[i])
				}
			}

			in.Phase = calculateBackupSessionPhase(in)
			if IsBackupCompleted(in.Phase) {
				in.SessionDuration = time.Since(meta.CreationTimestamp.Time).Round(time.Second).String()
			}
			return meta.UID, in
		},
		metav1.UpdateOptions{},
	)
}

func IsBackupCompleted(phase v1beta1.BackupSessionPhase) bool {
	return phase == v1beta1.BackupSessionSucceeded ||
		phase == v1beta1.BackupSessionFailed ||
		phase == v1beta1.BackupSessionSkipped ||
		phase == v1beta1.BackupSessionUnknown
}

func upsertBackupMembersStatus(cur []v1beta1.BackupTargetStatus, new v1beta1.BackupTargetStatus) []v1beta1.BackupTargetStatus {
	// if the member status already exist, then update it
	for i := range cur {
		if TargetMatched(cur[i].Ref, new.Ref) {
			cur[i] = upsertBackupTargetStatus(cur[i], new)
			return cur
		}
	}
	// the member status does not exist. so, add new entry.
	cur = append(cur, new)
	return cur
}

func upsertBackupTargetStatus(cur, new v1beta1.BackupTargetStatus) v1beta1.BackupTargetStatus {
	if len(new.Conditions) > 0 {
		cur.Conditions = upsertConditions(cur.Conditions, new.Conditions)
	}

	if new.TotalHosts != nil {
		cur.TotalHosts = new.TotalHosts
	}

	if len(new.Stats) > 0 {
		cur.Stats = upsertBackupHostStatus(cur.Stats, new.Stats)
	}

	if len(new.PreBackupActions) > 0 {
		cur.PreBackupActions = upsertArray(cur.PreBackupActions, new.PreBackupActions)
	}

	if len(new.PostBackupActions) > 0 {
		cur.PostBackupActions = upsertArray(cur.PostBackupActions, new.PostBackupActions)
	}

	cur.Phase = calculateBackupTargetPhase(cur)
	return cur
}

func upsertBackupHostStatus(cur, new []v1beta1.HostBackupStats) []v1beta1.HostBackupStats {
	for i := range new {
		index, hostEntryExist := backupHostEntryIndex(cur, new[i])
		if hostEntryExist {
			cur[index] = new[i]
		} else {
			cur = append(cur, new[i])
		}
	}
	return cur
}

func calculateBackupTargetPhase(status v1beta1.BackupTargetStatus) v1beta1.TargetPhase {
	if status.TotalHosts == nil {
		return v1beta1.TargetBackupPending
	}

	if kmapi.IsConditionFalse(status.Conditions, v1beta1.BackupExecutorEnsured) {
		return v1beta1.TargetBackupFailed
	}

	failedHostCount := int32(0)
	successfulHostCount := int32(0)
	for _, hostStats := range status.Stats {
		switch hostStats.Phase {
		case v1beta1.HostBackupFailed:
			failedHostCount++
		case v1beta1.HostBackupSucceeded:
			successfulHostCount++
		}
	}
	completedHosts := successfulHostCount + failedHostCount

	if completedHosts == *status.TotalHosts {
		if failedHostCount > 0 ||
			kmapi.IsConditionFalse(status.Conditions, v1beta1.PreBackupHookExecutionSucceeded) ||
			kmapi.IsConditionFalse(status.Conditions, v1beta1.PostBackupHookExecutionSucceeded) {
			return v1beta1.TargetBackupFailed
		}
		return v1beta1.TargetBackupSucceeded
	}
	return v1beta1.TargetBackupRunning
}

func calculateBackupSessionPhase(status *v1beta1.BackupSessionStatus) v1beta1.BackupSessionPhase {
	if len(status.Conditions) == 0 || len(status.Targets) == 0 {
		return v1beta1.BackupSessionPending
	}

	if kmapi.IsConditionTrue(status.Conditions, v1beta1.BackupSkipped) {
		return v1beta1.BackupSessionSkipped
	}

	failedTargetCount := 0
	successfulTargetCount := 0

	for _, t := range status.Targets {
		switch t.Phase {
		case v1beta1.TargetBackupFailed:
			failedTargetCount++
		case v1beta1.TargetBackupSucceeded:
			successfulTargetCount++
		}
	}
	completedTargets := successfulTargetCount + failedTargetCount

	if completedTargets == len(status.Targets) {
		if failedTargetCount > 0 ||
			kmapi.IsConditionFalse(status.Conditions, v1beta1.RetentionPolicyApplied) ||
			kmapi.IsConditionFalse(status.Conditions, v1beta1.RepositoryMetricsPushed) ||
			kmapi.IsConditionFalse(status.Conditions, v1beta1.MetricsPushed) ||
			kmapi.IsConditionFalse(status.Conditions, v1beta1.BackupHistoryCleaned) ||
			kmapi.IsConditionFalse(status.Conditions, v1beta1.RepositoryIntegrityVerified) ||
			kmapi.IsConditionFalse(status.Conditions, v1beta1.GlobalPreBackupHookSucceeded) ||
			kmapi.IsConditionFalse(status.Conditions, v1beta1.GlobalPostBackupHookSucceeded) {
			return v1beta1.BackupSessionFailed
		}

		if kmapi.IsConditionTrue(status.Conditions, v1beta1.MetricsPushed) &&
			kmapi.IsConditionTrue(status.Conditions, v1beta1.BackupHistoryCleaned) {
			return v1beta1.BackupSessionSucceeded
		}
	}

	return v1beta1.BackupSessionRunning
}

func backupHostEntryIndex(entries []v1beta1.HostBackupStats, target v1beta1.HostBackupStats) (int, bool) {
	for i := range entries {
		if entries[i].Hostname == target.Hostname {
			return i, true
		}
	}
	return -1, false
}

func upsertArray(cur, new []string) []string {
	for i := range new {
		if exist, idx := arrays.Contains(cur, new[i]); exist {
			cur[idx] = new[i]
			continue
		}
		cur = append(cur, new[i])
	}
	return cur
}
