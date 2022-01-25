/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Free Trial License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Free-Trial-1.0.0.md

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
	uiv1alpha1 "stash.appscode.dev/apimachinery/apis/ui/v1alpha1"
	"stash.appscode.dev/apimachinery/crds"

	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/apiextensions"
	appcatalog "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	metrics "kmodules.xyz/custom-resources/apis/metrics/v1alpha1"
	kmodules_crds "kmodules.xyz/custom-resources/crds"
)

var (
	masterURL                string
	kubeConfig               string
	enableEnterpriseFeatures bool
)

func init() {
	flag.StringVar(&kubeConfig, "kube-config", "", "Path to a kube config file. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kube config file. Only required if out-of-cluster.")
	flag.BoolVar(&enableEnterpriseFeatures, "enterprise", false, "Specify whether enterprise features enabled or not.")
}

func main() {
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeConfig)
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

	stashCRDs, err := getStashCRDs()
	if err != nil {
		return err
	}
	resources = append(resources, stashCRDs...)

	appCatalogCRDs, err := getAppCatalogCRDs()
	if err != nil {
		return err
	}
	resources = append(resources, appCatalogCRDs...)

	if enableEnterpriseFeatures {
		metricCRDs, err := getMetricCRDs()
		if err != nil {
			return err
		}
		resources = append(resources, metricCRDs...)
	}
	return apiextensions.RegisterCRDs(crdClient, resources)
}

func getStashCRDs() ([]*apiextensions.CustomResourceDefinition, error) {
	// community features CRDs
	gvrs := []schema.GroupVersionResource{
		// v1alpha1 resources
		{Group: stashv1alpha1.SchemeGroupVersion.Group, Version: stashv1alpha1.SchemeGroupVersion.Version, Resource: stashv1alpha1.ResourcePluralRepository},

		// v1beta1 resources
		{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralBackupConfiguration},
		{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralBackupSession},
		{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralRestoreSession},
		{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralFunction},
		{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralTask},
	}

	// enterprise features CRDs
	if enableEnterpriseFeatures {
		gvrs = append(gvrs, []schema.GroupVersionResource{
			// v1beta1 resources
			{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralBackupBatch},
			{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralBackupBlueprint},
			{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralRestoreBatch},

			// UI resources
			{Group: uiv1alpha1.SchemeGroupVersion.Group, Version: uiv1alpha1.SchemeGroupVersion.Version, Resource: uiv1alpha1.ResourceBackupOverviews},
		}...)
	}

	var stashCRDs []*apiextensions.CustomResourceDefinition
	for i := range gvrs {
		crd, err := crds.CustomResourceDefinition(gvrs[i])
		if err != nil {
			return nil, err
		}
		stashCRDs = append(stashCRDs, crd)
	}
	return stashCRDs, nil
}

func getAppCatalogCRDs() ([]*apiextensions.CustomResourceDefinition, error) {
	gvrs := []schema.GroupVersionResource{
		// v1alpha1 resources
		{Group: appcatalog.SchemeGroupVersion.Group, Version: appcatalog.SchemeGroupVersion.Version, Resource: appcatalog.ResourceApps},
	}
	var appCatalogCRDs []*apiextensions.CustomResourceDefinition
	for i := range gvrs {
		crd, err := kmodules_crds.CustomResourceDefinition(gvrs[i])
		if err != nil {
			return nil, err
		}
		appCatalogCRDs = append(appCatalogCRDs, crd)
	}
	return appCatalogCRDs, nil
}

func getMetricCRDs() ([]*apiextensions.CustomResourceDefinition, error) {
	gvrs := []schema.GroupVersionResource{
		// v1alpha1 resources
		{Group: metrics.SchemeGroupVersion.Group, Version: metrics.SchemeGroupVersion.Version, Resource: metrics.ResourceMetricsConfigurations},
	}

	var metricCRDs []*apiextensions.CustomResourceDefinition
	for i := range gvrs {
		crd, err := kmodules_crds.CustomResourceDefinition(gvrs[i])
		if err != nil {
			return nil, err
		}
		metricCRDs = append(metricCRDs, crd)
	}
	return metricCRDs, nil
}
