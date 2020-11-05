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

package util

import (
	"context"
	"fmt"
	"time"

	api "stash.appscode.dev/apimachinery/apis/stash/v1alpha1"
	cs "stash.appscode.dev/apimachinery/client/clientset/versioned/typed/stash/v1alpha1"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/golang/glog"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	kutil "kmodules.xyz/client-go"
)

func CreateOrPatchRecovery(ctx context.Context, c cs.StashV1alpha1Interface, meta metav1.ObjectMeta, transform func(in *api.Recovery) *api.Recovery, opts metav1.PatchOptions) (*api.Recovery, kutil.VerbType, error) {
	cur, err := c.Recoveries(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating Recovery %s/%s.", meta.Namespace, meta.Name)
		out, err := c.Recoveries(meta.Namespace).Create(ctx, transform(&api.Recovery{
			TypeMeta: metav1.TypeMeta{
				Kind:       api.ResourceKindRecovery,
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}), metav1.CreateOptions{
			DryRun:       opts.DryRun,
			FieldManager: opts.FieldManager,
		})
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return PatchRecovery(ctx, c, cur, transform, opts)
}

func PatchRecovery(ctx context.Context, c cs.StashV1alpha1Interface, cur *api.Recovery, transform func(*api.Recovery) *api.Recovery, opts metav1.PatchOptions) (*api.Recovery, kutil.VerbType, error) {
	return PatchRecoveryObject(ctx, c, cur, transform(cur.DeepCopy()), opts)
}

func PatchRecoveryObject(ctx context.Context, c cs.StashV1alpha1Interface, cur, mod *api.Recovery, opts metav1.PatchOptions) (*api.Recovery, kutil.VerbType, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(mod)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := jsonpatch.CreateMergePatch(curJson, modJson)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}
	glog.V(3).Infof("Patching Recovery %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.Recoveries(cur.Namespace).Patch(ctx, cur.Name, types.MergePatchType, patch, opts)
	return out, kutil.VerbPatched, err
}

func TryUpdateRecovery(ctx context.Context, c cs.StashV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.Recovery) *api.Recovery, opts metav1.UpdateOptions) (result *api.Recovery, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Recoveries(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.Recoveries(cur.Namespace).Update(ctx, transform(cur.DeepCopy()), opts)
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update Recovery %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update Recovery %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func UpdateRecoveryStatus(
	ctx context.Context,
	c cs.StashV1alpha1Interface,
	meta metav1.ObjectMeta,
	transform func(*api.RecoveryStatus) (types.UID, *api.RecoveryStatus),
	opts metav1.UpdateOptions,
) (result *api.Recovery, err error) {
	apply := func(x *api.Recovery) *api.Recovery {
		uid, updatedStatus := transform(x.Status.DeepCopy())
		// Ignore status update when uid does not match
		if uid != "" && uid != x.UID {
			return x
		}
		out := &api.Recovery{
			TypeMeta:   x.TypeMeta,
			ObjectMeta: x.ObjectMeta,
			Spec:       x.Spec,
			Status:     *updatedStatus,
		}
		return out
	}

	attempt := 0
	cur, err := c.Recoveries(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		var e2 error
		result, e2 = c.Recoveries(meta.Namespace).UpdateStatus(ctx, apply(cur), opts)
		if kerr.IsConflict(e2) {
			latest, e3 := c.Recoveries(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
			switch {
			case e3 == nil:
				cur = latest
				return false, nil
			case kutil.IsRequestRetryable(e3):
				return false, nil
			default:
				return false, e3
			}
		} else if err != nil && !kutil.IsRequestRetryable(e2) {
			return false, e2
		}
		return e2 == nil, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update status of Recovery %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func SetRecoveryStats(ctx context.Context, c cs.StashV1alpha1Interface, recovery metav1.ObjectMeta, path string, d time.Duration, phase api.RecoveryPhase, opts metav1.UpdateOptions) (*api.Recovery, error) {
	return UpdateRecoveryStatus(ctx, c, recovery, func(in *api.RecoveryStatus) (types.UID, *api.RecoveryStatus) {
		found := false
		for _, stats := range in.Stats {
			if stats.Path == path {
				found = true
				stats.Duration = d.String()
				stats.Phase = phase
			}
		}
		if !found {
			in.Stats = append(in.Stats, api.RestoreStats{
				Path:     path,
				Duration: d.String(),
				Phase:    phase,
			})
		}
		return "", in
	}, opts)
}
