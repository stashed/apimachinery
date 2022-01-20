package invoker

import (
	"context"

	"stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	cs "stash.appscode.dev/apimachinery/client/clientset/versioned"
	v1beta1_util "stash.appscode.dev/apimachinery/client/clientset/versioned/typed/stash/v1beta1/util"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kmapi "kmodules.xyz/client-go/api/v1"
	core_util "kmodules.xyz/client-go/core/v1"
)

type BackupConfigurationInvoker struct {
	backupConfig *v1beta1.BackupConfiguration
	stashClient  cs.Interface
}

func (inv *BackupConfigurationInvoker) AddFinalizer() error {
	_, _, err := v1beta1_util.PatchBackupConfiguration(context.TODO(), inv.stashClient.StashV1beta1(), inv.backupConfig, func(in *v1beta1.BackupConfiguration) *v1beta1.BackupConfiguration {
		in.ObjectMeta = core_util.AddFinalizer(in.ObjectMeta, v1beta1.StashKey)
		return in
	}, metav1.PatchOptions{})
	return err
}

func (inv *BackupConfigurationInvoker) RemoveFinalizer() error {
	_, _, err := v1beta1_util.PatchBackupConfiguration(context.TODO(), inv.stashClient.StashV1beta1(), inv.backupConfig, func(in *v1beta1.BackupConfiguration) *v1beta1.BackupConfiguration {
		in.ObjectMeta = core_util.RemoveFinalizer(in.ObjectMeta, v1beta1.StashKey)
		return in
	}, metav1.PatchOptions{})
	return err
}

func (inv *BackupConfigurationInvoker) HasCondition(conditionType string) (bool, error) {
	backupConfig, err := inv.stashClient.StashV1beta1().BackupConfigurations(inv.backupConfig.Namespace).Get(context.TODO(), inv.backupConfig.Name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	return kmapi.HasCondition(backupConfig.Status.Conditions, conditionType), nil
}

func (inv *BackupConfigurationInvoker) GetCondition(conditionType string) (int, *kmapi.Condition, error) {
	backupConfig, err := inv.stashClient.StashV1beta1().BackupConfigurations(inv.backupConfig.Namespace).Get(context.TODO(), inv.backupConfig.Name, metav1.GetOptions{})
	if err != nil {
		return -1, nil, err
	}
	idx, cond := kmapi.GetCondition(backupConfig.Status.Conditions, conditionType)
	return idx, cond, nil
}

func (inv *BackupConfigurationInvoker) SetCondition(newCondition kmapi.Condition) error {
	_, err := v1beta1_util.UpdateBackupConfigurationStatus(context.TODO(), inv.stashClient.StashV1beta1(), inv.backupConfig.ObjectMeta, func(in *v1beta1.BackupConfigurationStatus) (types.UID, *v1beta1.BackupConfigurationStatus) {
		in.Conditions = kmapi.SetCondition(in.Conditions, newCondition)
		return inv.backupConfig.UID, in
	}, metav1.UpdateOptions{})
	return err
}

func (inv *BackupConfigurationInvoker) IsConditionTrue(conditionType string) (bool, error) {
	backupConfig, err := inv.stashClient.StashV1beta1().BackupConfigurations(inv.backupConfig.Namespace).Get(context.TODO(), inv.backupConfig.Name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	return kmapi.IsConditionTrue(backupConfig.Status.Conditions, conditionType), nil
}

func (inv *BackupConfigurationInvoker) NextInOrder(curTarget v1beta1.TargetRef) bool {
	// BackupConfiguration has only one target. So, it will be always in front of execution order.
	return true
}
