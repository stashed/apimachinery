/*
Copyright The KubeVault Authors.

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

package install

import (
	"testing"

	"stash.appscode.dev/apimachinery/apis/stash/fuzzer"
	"stash.appscode.dev/apimachinery/apis/stash/v1alpha1"
	"stash.appscode.dev/apimachinery/apis/stash/v1beta1"

	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	crdfuzz "kmodules.xyz/crd-schema-fuzz"
)

func TestPruneTypes(t *testing.T) {
	Install(clientsetscheme.Scheme)

	// v1alpha1
	if crd := (v1alpha1.Restic{}).CustomResourceDefinition(); crd.V1 != nil {
		crdfuzz.SchemaFuzzTestForV1CRD(t, clientsetscheme.Scheme, crd.V1, fuzzer.Funcs)
	}
	if crd := (v1alpha1.Recovery{}.CustomResourceDefinition()); crd.V1 != nil {
		crdfuzz.SchemaFuzzTestForV1CRD(t, clientsetscheme.Scheme, crd.V1, fuzzer.Funcs)
	}
	if crd := (v1alpha1.Repository{}).CustomResourceDefinition(); crd.V1 != nil {
		crdfuzz.SchemaFuzzTestForV1CRD(t, clientsetscheme.Scheme, crd.V1, fuzzer.Funcs)
	}

	// v1beta1
	if crd := (v1beta1.BackupBatch{}).CustomResourceDefinition(); crd.V1 != nil {
		crdfuzz.SchemaFuzzTestForV1CRD(t, clientsetscheme.Scheme, crd.V1, fuzzer.Funcs)
	}
	if crd := (v1beta1.BackupBlueprint{}).CustomResourceDefinition(); crd.V1 != nil {
		crdfuzz.SchemaFuzzTestForV1CRD(t, clientsetscheme.Scheme, crd.V1, fuzzer.Funcs)
	}
	if crd := (v1beta1.BackupConfiguration{}).CustomResourceDefinition(); crd.V1 != nil {
		crdfuzz.SchemaFuzzTestForV1CRD(t, clientsetscheme.Scheme, crd.V1, fuzzer.Funcs)
	}
	if crd := (v1beta1.BackupSession{}).CustomResourceDefinition(); crd.V1 != nil {
		crdfuzz.SchemaFuzzTestForV1CRD(t, clientsetscheme.Scheme, crd.V1, fuzzer.Funcs)
	}
	if crd := (v1beta1.Function{}).CustomResourceDefinition(); crd.V1 != nil {
		crdfuzz.SchemaFuzzTestForV1CRD(t, clientsetscheme.Scheme, crd.V1, fuzzer.Funcs)
	}
	if crd := (v1beta1.RestoreSession{}).CustomResourceDefinition(); crd.V1 != nil {
		crdfuzz.SchemaFuzzTestForV1CRD(t, clientsetscheme.Scheme, crd.V1, fuzzer.Funcs)
	}
	if crd := (v1beta1.Task{}).CustomResourceDefinition(); crd.V1 != nil {
		crdfuzz.SchemaFuzzTestForV1CRD(t, clientsetscheme.Scheme, crd.V1, fuzzer.Funcs)
	}

	// v1alpha1
	if crd := (v1alpha1.Restic{}).CustomResourceDefinition(); crd.V1beta1 != nil {
		crdfuzz.SchemaFuzzTestForV1beta1CRD(t, clientsetscheme.Scheme, crd.V1beta1, fuzzer.Funcs)
	}
	if crd := (v1alpha1.Recovery{}.CustomResourceDefinition()); crd.V1beta1 != nil {
		crdfuzz.SchemaFuzzTestForV1beta1CRD(t, clientsetscheme.Scheme, crd.V1beta1, fuzzer.Funcs)
	}
	if crd := (v1alpha1.Repository{}).CustomResourceDefinition(); crd.V1beta1 != nil {
		crdfuzz.SchemaFuzzTestForV1beta1CRD(t, clientsetscheme.Scheme, crd.V1beta1, fuzzer.Funcs)
	}

	// v1beta1
	if crd := (v1beta1.BackupBatch{}).CustomResourceDefinition(); crd.V1beta1 != nil {
		crdfuzz.SchemaFuzzTestForV1beta1CRD(t, clientsetscheme.Scheme, crd.V1beta1, fuzzer.Funcs)
	}
	if crd := (v1beta1.BackupBlueprint{}).CustomResourceDefinition(); crd.V1beta1 != nil {
		crdfuzz.SchemaFuzzTestForV1beta1CRD(t, clientsetscheme.Scheme, crd.V1beta1, fuzzer.Funcs)
	}
	if crd := (v1beta1.BackupConfiguration{}).CustomResourceDefinition(); crd.V1beta1 != nil {
		crdfuzz.SchemaFuzzTestForV1beta1CRD(t, clientsetscheme.Scheme, crd.V1beta1, fuzzer.Funcs)
	}
	if crd := (v1beta1.BackupSession{}).CustomResourceDefinition(); crd.V1beta1 != nil {
		crdfuzz.SchemaFuzzTestForV1beta1CRD(t, clientsetscheme.Scheme, crd.V1beta1, fuzzer.Funcs)
	}
	if crd := (v1beta1.Function{}).CustomResourceDefinition(); crd.V1beta1 != nil {
		crdfuzz.SchemaFuzzTestForV1beta1CRD(t, clientsetscheme.Scheme, crd.V1beta1, fuzzer.Funcs)
	}
	if crd := (v1beta1.RestoreSession{}).CustomResourceDefinition(); crd.V1beta1 != nil {
		crdfuzz.SchemaFuzzTestForV1beta1CRD(t, clientsetscheme.Scheme, crd.V1beta1, fuzzer.Funcs)
	}
	if crd := (v1beta1.Task{}).CustomResourceDefinition(); crd.V1beta1 != nil {
		crdfuzz.SchemaFuzzTestForV1beta1CRD(t, clientsetscheme.Scheme, crd.V1beta1, fuzzer.Funcs)
	}
}
