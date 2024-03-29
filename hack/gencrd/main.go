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

package main

import (
	"os"
	"path/filepath"

	repoinstall "stash.appscode.dev/apimachinery/apis/repositories/install"
	repov1alpha1 "stash.appscode.dev/apimachinery/apis/repositories/v1alpha1"
	stashinstall "stash.appscode.dev/apimachinery/apis/stash/install"
	stashv1alpha1 "stash.appscode.dev/apimachinery/apis/stash/v1alpha1"
	stashv1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"

	gort "gomodules.xyz/runtime"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/klog/v2"
	"k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/validation/spec"
	"kmodules.xyz/client-go/openapi"
)

func generateSwaggerJson() {
	var (
		Scheme = runtime.NewScheme()
		Codecs = serializer.NewCodecFactory(Scheme)
	)

	stashinstall.Install(Scheme)
	repoinstall.Install(Scheme)

	apispec, err := openapi.RenderOpenAPISpec(openapi.Config{
		Scheme: Scheme,
		Codecs: Codecs,
		Info: spec.InfoProps{
			Title:   "Stash",
			Version: "v0.9.0-rc.0",
			Contact: &spec.ContactInfo{
				Name:  "AppsCode Inc.",
				URL:   "https://appscode.com",
				Email: "hello@appscode.com",
			},
			License: &spec.License{
				Name: "Apache 2.0",
				URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
			},
		},
		OpenAPIDefinitions: []common.GetOpenAPIDefinitions{
			stashv1alpha1.GetOpenAPIDefinitions,
			stashv1beta1.GetOpenAPIDefinitions,
			repov1alpha1.GetOpenAPIDefinitions,
		},
		//nolint:govet
		Resources: []openapi.TypeInfo{
			// v1alpha1 resources
			{stashv1alpha1.SchemeGroupVersion, stashv1alpha1.ResourcePluralRepository, stashv1alpha1.ResourceKindRepository, true},

			// v1beta1 resources
			{stashv1beta1.SchemeGroupVersion, stashv1beta1.ResourcePluralBackupConfiguration, stashv1beta1.ResourceKindBackupConfiguration, true},
			{stashv1beta1.SchemeGroupVersion, stashv1beta1.ResourcePluralBackupSession, stashv1beta1.ResourceKindBackupSession, true},
			{stashv1beta1.SchemeGroupVersion, stashv1beta1.ResourcePluralBackupBatch, stashv1beta1.ResourceKindBackupBatch, false},
			{stashv1beta1.SchemeGroupVersion, stashv1beta1.ResourcePluralBackupBlueprint, stashv1beta1.ResourceKindBackupBlueprint, false},
			{stashv1beta1.SchemeGroupVersion, stashv1beta1.ResourcePluralRestoreSession, stashv1beta1.ResourceKindRestoreSession, true},
			{stashv1beta1.SchemeGroupVersion, stashv1beta1.ResourcePluralRestoreBatch, stashv1beta1.ResourceKindRestoreBatch, true},
			{stashv1beta1.SchemeGroupVersion, stashv1beta1.ResourcePluralFunction, stashv1beta1.ResourceKindFunction, false},
			{stashv1beta1.SchemeGroupVersion, stashv1beta1.ResourcePluralTask, stashv1beta1.ResourceKindTask, false},
		},
		//nolint:govet
		RDResources: []openapi.TypeInfo{
			{repov1alpha1.SchemeGroupVersion, repov1alpha1.ResourcePluralSnapshot, repov1alpha1.ResourceKindSnapshot, true},
		},
	})
	if err != nil {
		klog.Fatal(err)
	}

	filename := gort.GOPath() + "/src/stash.appscode.dev/apimachinery/openapi/swagger.json"
	err = os.MkdirAll(filepath.Dir(filename), 0o755)
	if err != nil {
		klog.Fatal(err)
	}
	err = os.WriteFile(filename, []byte(apispec), 0o644)
	if err != nil {
		klog.Fatal(err)
	}
}

func main() {
	generateSwaggerJson()
}
