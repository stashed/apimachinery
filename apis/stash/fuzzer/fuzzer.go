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

package fuzzer

import (
	"stash.appscode.dev/apimachinery/apis/stash/v1alpha1"
	"stash.appscode.dev/apimachinery/apis/stash/v1beta1"

	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/randfill"
)

// Funcs returns the fuzzer functions for this api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []any {
	return []any{
		// v1alpha1
		func(s *v1alpha1.Repository, c randfill.Continue) {
			c.Fill(s) // fuzz self without calling this function again
		},
		// v1beta1
		func(s *v1beta1.BackupBatch, c randfill.Continue) {
			c.Fill(s) // fuzz self without calling this function again
		},
		func(s *v1beta1.BackupBlueprint, c randfill.Continue) {
			c.Fill(s) // fuzz self without calling this function again
		},
		func(s *v1beta1.BackupConfiguration, c randfill.Continue) {
			c.Fill(s) // fuzz self without calling this function again
		},
		func(s *v1beta1.BackupSession, c randfill.Continue) {
			c.Fill(s) // fuzz self without calling this function again
		},
		func(s *v1beta1.Function, c randfill.Continue) {
			c.Fill(s) // fuzz self without calling this function again
		},
		func(s *v1beta1.RestoreSession, c randfill.Continue) {
			c.Fill(s) // fuzz self without calling this function again
		},
		func(s *v1beta1.Task, c randfill.Continue) {
			c.Fill(s) // fuzz self without calling this function again
		},
	}
}
