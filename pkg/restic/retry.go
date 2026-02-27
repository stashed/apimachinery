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
	"context"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
)

const (
	maxRetries = 5
	delay      = 10 * time.Second
)

var retryablePatterns = []string{
	"Connection closed by foreign host",
}

type RetryConfig struct {
	MaxRetries  int
	Delay       time.Duration
	ShouldRetry func(error, string) bool
}

func NewRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries: maxRetries,
		Delay:      delay,
		ShouldRetry: func(err error, output string) bool {
			if err == nil {
				return false
			}
			combined := strings.ToLower(err.Error() + " " + output)
			for _, pattern := range retryablePatterns {
				if strings.Contains(combined, strings.ToLower(pattern)) {
					return true
				}
			}
			return false
		},
	}
}

func (rc *RetryConfig) RunWithRetry(ctx context.Context, execFunc func() ([]byte, error)) ([]byte, error) {
	var output []byte
	var lastErr error
	attempts := 1

	err := wait.PollUntilContextCancel(
		ctx,
		rc.Delay,
		true, // Run immediately on first call
		func(ctx context.Context) (bool, error) {
			// Stop if max retries reached
			if attempts > rc.MaxRetries {
				return false, fmt.Errorf("max retries reached")
			}
			output, lastErr = execFunc()
			klog.Infof("Attempt #%d: retrying in %v", attempts, rc.Delay)
			klog.Infof("Attempt #%d error: %v", attempts, lastErr)
			klog.Infof("Attempt #%d output: %s", attempts, string(output))

			if !rc.ShouldRetry(lastErr, string(output)) {
				return true, nil
			}
			attempts++
			return false, nil
		},
	)
	if err != nil {
		return output, fmt.Errorf("failed after %d attempts: %w", attempts, lastErr)
	}
	return output, lastErr
}
