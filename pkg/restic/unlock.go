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

package restic

import (
	"bytes"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	kutil "kmodules.xyz/client-go"
	"time"
)

func (w *ResticWrapper) UnlockRepository() error {
	_, err := w.unlock()
	return err
}

// getLockIDs lists every lock ID currently held in the repository.
func (w *ResticWrapper) getLockIDs() ([]string, error) {
	w.sh.ShowCMD = true
	out, err := w.listLocks()
	if err != nil {
		return nil, err
	}
	return extractLockIDs(bytes.NewReader(out))
}

// getLockStats returns the decoded JSON for a single lock.
func (w *ResticWrapper) getLockStats(lockID string) (*LockStats, error) {
	w.sh.ShowCMD = true
	out, err := w.lockStats(lockID)
	if err != nil {
		return nil, err
	}
	return extractLockStats(out)
}

// getPodNameIfAnyExclusiveLock scans every lock and returns the hostname aka (Pod name) of the first exclusive lock it finds, or "" if none exist.
func (w *ResticWrapper) getPodNameIfAnyExclusiveLock() (string, error) {
	ids, err := w.getLockIDs()
	if err != nil {
		return "", fmt.Errorf("failed to list locks: %w", err)
	}
	for _, id := range ids {
		st, err := w.getLockStats(id)
		if err != nil {
			return "", fmt.Errorf("failed to inspect lock %s: %w", id, err)
		}
		if st.Exclusive { // There's no chances to get multiple exclusive locks, so we can return the first one we find.
			return st.Hostname, nil
		}
	}
	return "", nil
}

// EnsureNoExclusiveLock blocks until any exclusive lock is released.
// If a lock is held by a Running Pod, it waits; otherwise it unlocks.
func (w *ResticWrapper) EnsureNoExclusiveLock(k8sClient kubernetes.Interface, namespace string) error {
	klog.Infoln("Ensuring no exclusive lock is held in the repository...")
	podName, err := w.getPodNameIfAnyExclusiveLock()
	if err != nil {
		return fmt.Errorf("failed to query exclusive lock: %w", err)
	}
	if podName == "" {
		return nil // nothing to do
	}

	return wait.PollUntilContextTimeout(
		context.Background(),
		5*time.Second,
		kutil.ReadinessTimeout,
		true,
		func(ctx context.Context) (bool, error) {
			klog.Infoln("Getting pod to check if it is still running:", podName)
			pod, err := k8sClient.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
			switch {
			case errors.IsNotFound(err): // Pod gone → unlock
				_, err := w.unlock()
				return true, err
			case err != nil: // API error → stop
				return false, err
			case pod.Status.Phase == corev1.PodRunning: // still running → keep waiting
				return false, nil
			default: // Pod finished → unlock
				_, err := w.unlock()
				return true, err
			}
		},
	)
}
