package util

import (
	"context"
	"encoding/json"
	"strings"

	"stash.appscode.dev/apimachinery/apis/stash/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	appcatalog_cs "kmodules.xyz/custom-resources/client/clientset/versioned"
)

func ExtractAddonInfo(appClient appcatalog_cs.Interface, task v1beta1.TaskRef, targetRef v1beta1.TargetRef, namespace string) (*appcat.StashAddonSpec, error) {
	addonInfo := appcat.StashAddonSpec{}

	// If the target is AppBinding and it has addon information set in the parameters section, then extract the addon info.
	if targetOfGroupKind(targetRef, appcat.SchemeGroupVersion.Group, appcat.ResourceKindApp) {
		// get the AppBinding
		appBinding, err := appClient.AppcatalogV1alpha1().AppBindings(namespace).Get(context.TODO(), targetRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		// extract the parameters
		if appBinding.Spec.Parameters != nil {
			err = json.Unmarshal(appBinding.Spec.Parameters.Raw, &addonInfo)
			if err != nil {
				return nil, err
			}
		}
	}

	// If the user provides Task information in the backup/restore invoker spec, it should have higher precedence.
	// We don't know whether this function was called from BackupSession controller or RestoreSession controller.
	// Hence, we are going to overwrite the task name & parameters in both backupTask & restoreTask section.
	// It does not have any adverse effect because when it is called from the BackupSession controller, we will overwrite with backup task info
	// and when it is called from the RestoreSession controller, we will overwrite with restore task info.
	if task.Name != "" {
		addonInfo.Addon.BackupTask.Name = task.Name
		addonInfo.Addon.RestoreTask.Name = task.Name
	}
	if len(task.Params) != 0 {
		addonInfo.Addon.BackupTask.Params = getTaskParams(task)
		addonInfo.Addon.RestoreTask.Params = getTaskParams(task)
	}

	return &addonInfo, nil
}

func targetOfGroupKind(targetRef v1beta1.TargetRef, group, kind string) bool {
	gv := strings.Split(targetRef.APIVersion, "/")
	if gv[0] == group && targetRef.Kind == kind {
		return true
	}
	return false
}

func getTaskParams(task v1beta1.TaskRef) []appcat.Param {
	params := make([]appcat.Param, len(task.Params))
	for i := range task.Params {
		params[i].Name = task.Params[i].Name
		params[i].Value = task.Params[i].Value
	}
	return params
}
