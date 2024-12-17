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

package hooks

import (
	"fmt"
	"testing"

	"stash.appscode.dev/apimachinery/apis/stash/v1beta1"

	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	kmapi "kmodules.xyz/client-go/api/v1"
	prober "kmodules.xyz/prober/api/v1"
)

const (
	defaultErrorMsg             = "this session failed because destined to fail"
	errorMsgWithEscapeCharacter = "mysqldump: [ERROR] mysqldump: unknown option '-k'. {\"message_type\":\"error\",\"error\":{\"Op\":\"read\",\"Path\":\"/dumpfile.sql\",\"Err\":{}},\"during\":\"archival\",\"item\":\"/dumpfile.sql\"} Fatal: unable to save snapshot: snapshot is empty"
)

func TestHookExecutor_renderTemplate(t *testing.T) {
	type fields struct {
		Config      *rest.Config
		Hook        *prober.Handler
		ExecutorPod kmapi.ObjectReference
		Summary     *v1beta1.Summary
	}

	tests := []struct {
		name         string
		fields       fields
		wantErr      bool
		expectedBody string
	}{
		{
			name: "[HTTPPost] Successful session",
			fields: fields{
				Hook:    defaultHookTemplate(false),
				Summary: defaultSummary(),
			},
			wantErr:      false,
			expectedBody: "Name: test-session Namespace: test Phase: Succeeded",
		},
		{
			name: "[HTTPPost] Failed session",
			fields: fields{
				Hook:    defaultHookTemplate(false),
				Summary: failedSummary(),
			},
			wantErr:      false,
			expectedBody: "Name: test-session Namespace: test Phase: Failed",
		},
		{
			name: "[HTTPPost] Failed session with escape character in error message",
			fields: fields{
				Hook:    defaultHookTemplate(false),
				Summary: failedSummary(),
			},
			wantErr:      false,
			expectedBody: "Name: test-session Namespace: test Phase: Failed",
		},
		{
			name: "[HTTPPost] Conditional hook with Succeeded phase",
			fields: fields{
				Hook:    conditionalHookTemplate(defaultErrorMsg),
				Summary: defaultSummary(),
			},
			wantErr:      false,
			expectedBody: "Succeeded",
		},
		{
			name: "[HTTPPost] Conditional hook with Failed phase",
			fields: fields{
				Hook:    conditionalHookTemplate(defaultErrorMsg),
				Summary: failedSummary(),
			},
			wantErr:      false,
			expectedBody: fmt.Sprintf("Failed. Reason: %s", defaultErrorMsg),
		},
		{
			name: "[HTTPPost] Conditional hook with escape character in error message",
			fields: fields{
				Hook:    conditionalHookTemplate(errorMsgWithEscapeCharacter),
				Summary: failedSummary(),
			},
			wantErr:      false,
			expectedBody: fmt.Sprintf("Failed. Reason: %s", errorMsgWithEscapeCharacter),
		},
		{
			name: "[Exec] Successful session",
			fields: fields{
				Hook:    defaultHookTemplate(true),
				Summary: defaultSummary(),
			},
			wantErr:      false,
			expectedBody: "{\"text\":\":x: Name: test-session Namespace: test Phase: Succeeded\"}",
		},
		{
			name: "[Exec] Failed session",
			fields: fields{
				Hook:    defaultHookTemplate(true),
				Summary: failedSummary(),
			},
			wantErr:      false,
			expectedBody: "{\"text\":\":x: Name: test-session Namespace: test Phase: Failed\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &HookExecutor{
				Config:      tt.fields.Config,
				Hook:        tt.fields.Hook,
				ExecutorPod: tt.fields.ExecutorPod,
				Summary:     tt.fields.Summary,
			}
			if err := e.renderHookTemplate(); (err != nil) != tt.wantErr {
				t.Errorf("renderHookTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if e.Hook.HTTPPost != nil &&
				tt.expectedBody != e.Hook.HTTPPost.Body {
				t.Errorf("Expected: %v, found: %v", tt.expectedBody, e.Hook.HTTPPost.Body)
				return
			}
			if e.Hook.Exec != nil {
				if e.Hook.Exec.Command[6] != tt.expectedBody {
					t.Errorf("Expected: %v, found: %v", tt.expectedBody, e.Hook.Exec.Command[6])
					return
				}
			}
		})
	}
}

func defaultSummary(transformFuncs ...func(s *v1beta1.Summary)) *v1beta1.Summary {
	summary := &v1beta1.Summary{
		Name:      "test-session",
		Namespace: "test",

		Invoker: core.TypedLocalObjectReference{
			APIGroup: pointer.StringP("stash.appscode.com"),
			Kind:     "BackupConfiguration",
			Name:     "test",
		},

		Target: v1beta1.TargetRef{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
			Name:       "test-sts",
		},

		Status: v1beta1.TargetStatus{
			Phase:    "Succeeded",
			Duration: "2m",
		},
	}

	for _, f := range transformFuncs {
		f(summary)
	}
	return summary
}

func failedSummary() *v1beta1.Summary {
	return defaultSummary(func(s *v1beta1.Summary) {
		s.Status.Phase = "Failed"
		s.Status.Error = defaultErrorMsg
	})
}

func defaultHookTemplate(isExec bool) *prober.Handler {
	if isExec {
		return &prober.Handler{
			Exec: &core.ExecAction{
				Command: []string{
					"curl",
					"-X",
					"POST",
					"-H",
					"Content-Type: application/json",
					"-d",
					`{"text":":x: Name: {{ .Name }} Namespace: {{.Namespace}} Phase: {{.Status.Phase}}"}`,
					"https://slack-webhook.url/stash",
				},
			},
		}
	}
	return &prober.Handler{
		HTTPPost: &prober.HTTPPostAction{
			Body: "Name: {{ .Name }} Namespace: {{.Namespace}} Phase: {{.Status.Phase}}",
		},
	}
}

func conditionalHookTemplate(msg string) *prober.Handler {
	return &prober.Handler{
		HTTPPost: &prober.HTTPPostAction{
			Body: fmt.Sprintf("{{ if eq .Status.Phase `Succeeded`}}Succeeded{{ else }}Failed. Reason: %s{{ end}}", msg),
		},
	}
}
