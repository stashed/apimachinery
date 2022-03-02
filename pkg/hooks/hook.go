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
	"bytes"
	"encoding/json"
	"strings"
	"sync"
	"text/template"

	"stash.appscode.dev/apimachinery/apis"
	"stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	cs "stash.appscode.dev/apimachinery/client/clientset/versioned"
	"stash.appscode.dev/apimachinery/pkg/conditions"
	"stash.appscode.dev/apimachinery/pkg/invoker"

	"k8s.io/client-go/rest"
	kmapi "kmodules.xyz/client-go/api/v1"
	prober "kmodules.xyz/prober/api/v1"
	"kmodules.xyz/prober/probe"
)

type HookExecutor struct {
	Config      *rest.Config
	Hook        *prober.Handler
	ExecutorPod kmapi.ObjectReference
	Summary     *v1beta1.Summary
}

func (e *HookExecutor) Execute() error {
	if strings.Contains(e.Hook.String(), "{{") {
		if err := e.renderTemplate(); err != nil {
			return err
		}
	}
	return probe.RunProbe(e.Config, e.Hook, e.ExecutorPod.Name, e.ExecutorPod.Namespace)
}

var pool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func (e *HookExecutor) renderTemplate() error {
	hookContent, err := json.Marshal(e.Hook)
	if err != nil {
		return err
	}

	tpl, err := template.New("hook-template").Parse(string(hookContent))
	if err != nil {
		return err
	}
	tpl.Option("missingkey=default")

	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()
	defer pool.Put(buf)

	err = tpl.Execute(buf, e.Summary)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf.Bytes(), &e.Hook)
}

type BackupHookExecutor struct {
	Config        *rest.Config
	StashClient   cs.Interface
	BackupSession *v1beta1.BackupSession
	Invoker       invoker.BackupInvoker
	Target        v1beta1.TargetRef
	ExecutorPod   kmapi.ObjectReference
	Hook          *prober.Handler
	HookType      string
}

func (e *BackupHookExecutor) Execute() (*v1beta1.BackupSession, error) {
	hookExecutor := HookExecutor{
		Config:      e.Config,
		Hook:        e.Hook,
		ExecutorPod: e.ExecutorPod,
		Summary: e.Invoker.GetSummary(e.Target, kmapi.ObjectReference{
			Namespace: e.BackupSession.Namespace,
			Name:      e.BackupSession.Name,
		}),
	}

	if e.alreadyExecuted() {
		return e.BackupSession, nil
	}

	if err := hookExecutor.Execute(); err != nil {
		return e.setBackupHookExecutionSucceededToFalse(err)
	}

	return e.setBackupHookExecutionSucceededToTrue()
}

func (e *BackupHookExecutor) alreadyExecuted() bool {
	if e.HookType == apis.PreBackupHook {
		return kmapi.HasCondition(e.BackupSession.Status.Conditions, v1beta1.PreBackupHookExecutionSucceeded)
	}
	return kmapi.HasCondition(e.BackupSession.Status.Conditions, v1beta1.PreBackupHookExecutionSucceeded)
}

func (e *BackupHookExecutor) setBackupHookExecutionSucceededToFalse(err error) (*v1beta1.BackupSession, error) {
	if e.HookType == apis.PreBackupHook {
		return conditions.SetPreBackupHookExecutionSucceededToFalse(e.StashClient, e.BackupSession, err)
	} else {
		return conditions.SetPostBackupHookExecutionSucceededToFalse(e.StashClient, e.BackupSession, err)
	}
}

func (e *BackupHookExecutor) setBackupHookExecutionSucceededToTrue() (*v1beta1.BackupSession, error) {
	if e.HookType == apis.PreBackupHook {
		return conditions.SetPreBackupHookExecutionSucceededToTrue(e.StashClient, e.BackupSession)
	} else {
		return conditions.SetPostBackupHookExecutionSucceededToTrue(e.StashClient, e.BackupSession)
	}
}

type RestoreHookExecutor struct {
	Config      *rest.Config
	Invoker     invoker.RestoreInvoker
	Target      v1beta1.TargetRef
	ExecutorPod kmapi.ObjectReference
	Hook        *prober.Handler
	HookType    string
}

func (e *RestoreHookExecutor) Execute() error {
	hookExecutor := HookExecutor{
		Config:      e.Config,
		Hook:        e.Hook,
		ExecutorPod: e.ExecutorPod,
		Summary: e.Invoker.GetSummary(e.Target, kmapi.ObjectReference{
			Namespace: e.Invoker.GetObjectMeta().Namespace,
			Name:      e.Invoker.GetObjectMeta().Name,
		}),
	}

	if yes, err := e.alreadyExecuted(); yes || err != nil {
		return err
	}

	if err := hookExecutor.Execute(); err != nil {
		return e.setRestoreHookExecutionSucceededToFalse(err)
	}

	return e.setRestoreHookExecutionSucceededToTrue()
}

func (e *RestoreHookExecutor) alreadyExecuted() (bool, error) {
	if e.HookType == apis.PreBackupHook {
		return e.Invoker.HasCondition(&e.Target, v1beta1.PreRestoreHookExecutionSucceeded)
	}
	return e.Invoker.HasCondition(&e.Target, v1beta1.PostRestoreHookExecutionSucceeded)
}

func (e *RestoreHookExecutor) setRestoreHookExecutionSucceededToFalse(err error) error {
	if e.HookType == apis.PreRestoreHook {
		return conditions.SetPreRestoreHookExecutionSucceededToFalse(e.Invoker, err)
	} else {
		return conditions.SetPostRestoreHookExecutionSucceededToFalse(e.Invoker, err)
	}
}

func (e *RestoreHookExecutor) setRestoreHookExecutionSucceededToTrue() error {
	if e.HookType == apis.PreRestoreHook {
		return conditions.SetPreRestoreHookExecutionSucceededToTrue(e.Invoker)
	} else {
		return conditions.SetPostRestoreHookExecutionSucceededToTrue(e.Invoker)
	}
}
