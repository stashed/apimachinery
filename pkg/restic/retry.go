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
			klog.Info("Combined output: " + combined)
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
	attempts := 0

	err := wait.PollUntilContextCancel(
		ctx,
		rc.Delay,
		true, // Run immediately on first call
		func(ctx context.Context) (bool, error) {
			// Stop if max retries reached
			if attempts >= rc.MaxRetries {
				return false, fmt.Errorf("max retries reached")
			}
			output, lastErr = execFunc()
			if !rc.ShouldRetry(lastErr, string(output)) {
				return true, nil
			}
			klog.Info("Retrying command after error",
				"attempt", attempts,
				"maxRetries", rc.MaxRetries,
				"error", fmt.Sprintf("%s %s", lastErr, string(output)))
			attempts++
			return false, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed after %d attempts: %w", attempts, lastErr)
	}

	return output, lastErr
}
