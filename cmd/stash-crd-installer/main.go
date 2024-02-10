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
	"flag"
	"os"

	stashv1alpha1 "stash.appscode.dev/apimachinery/apis/stash/v1alpha1"
	stashv1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"

	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/apiextensions"
	appcatalog "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	metrics "kmodules.xyz/custom-resources/apis/metrics/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

var enableEnterpriseFeatures bool

func init() {
	flag.BoolVar(&enableEnterpriseFeatures, "enterprise", false, "Specify whether enterprise features enabled or not.")
}

func main() {
	flag.Parse()

	cfg, err := ctrl.GetConfig()
	if err != nil {
		klog.Errorf("Error building kubeconfig: %s", err.Error())
		os.Exit(1)
	}

	crdClient, err := crd_cs.NewForConfig(cfg)
	if err != nil {
		klog.Errorf("Error building CRD client: %s", err.Error())
		os.Exit(1)
	}

	if err := registerCRDs(crdClient); err != nil {
		klog.Errorf("Error building CRD client: %s", err.Error())
		os.Exit(1)
	}
	klog.Infoln("Successfully installed all CRDs.")
}

func registerCRDs(crdClient crd_cs.Interface) error {
	var resources []*apiextensions.CustomResourceDefinition

	resources = append(resources, getStashCRDs()...)
	resources = append(resources, getAppCatalogCRDs()...)

	if enableEnterpriseFeatures {
		resources = append(resources, getMetricCRDs()...)
	}
	return apiextensions.RegisterCRDs(crdClient, resources)
}

func getStashCRDs() []*apiextensions.CustomResourceDefinition {
	// community features CRDs
	var stashCRDs []*apiextensions.CustomResourceDefinition

	stashCRDs = append(stashCRDs,
		// v1alpha1 resources
		stashv1alpha1.Repository{}.CustomResourceDefinition(),

		// v1beta1 resources
		stashv1beta1.BackupConfiguration{}.CustomResourceDefinition(),
		stashv1beta1.BackupSession{}.CustomResourceDefinition(),
		stashv1beta1.RestoreSession{}.CustomResourceDefinition(),
		stashv1beta1.Function{}.CustomResourceDefinition(),
		stashv1beta1.Task{}.CustomResourceDefinition(),
	)

	// enterprise features CRDs
	if enableEnterpriseFeatures {
		stashCRDs = append(stashCRDs,
			// v1beta1 resources
			stashv1beta1.BackupBatch{}.CustomResourceDefinition(),
			stashv1beta1.BackupBlueprint{}.CustomResourceDefinition(),
			stashv1beta1.RestoreBatch{}.CustomResourceDefinition(),
		)
	}
	return stashCRDs
}

func getAppCatalogCRDs() []*apiextensions.CustomResourceDefinition {
	return []*apiextensions.CustomResourceDefinition{
		// v1alpha1 resources
		appcatalog.AppBinding{}.CustomResourceDefinition(),
	}
}

func getMetricCRDs() []*apiextensions.CustomResourceDefinition {
	return []*apiextensions.CustomResourceDefinition{
		// v1alpha1 resources
		metrics.MetricsConfiguration{}.CustomResourceDefinition(),
	}
}
