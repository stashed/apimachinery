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

	uiv1alpha1 "stash.appscode.dev/apimachinery/apis/ui/v1alpha1"

	stashv1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"

	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/apiextensions"
	stashv1alpha1 "stash.appscode.dev/apimachinery/apis/stash/v1alpha1"
	"stash.appscode.dev/apimachinery/crds"
)

var (
	masterURL  string
	kubeConfig string
)

func init() {
	flag.StringVar(&kubeConfig, "kube-config", "", "Path to a kube config file. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kube config file. Only required if out-of-cluster.")
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

	gvrs := []schema.GroupVersionResource{
		// v1alpha1 resources
		{Group: stashv1alpha1.SchemeGroupVersion.Group, Version: stashv1alpha1.SchemeGroupVersion.Version, Resource: stashv1alpha1.ResourcePluralRepository},

		// v1beta1 resources
		{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralBackupConfiguration},
		{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralBackupSession},
		{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralBackupBatch},
		{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralBackupBlueprint},
		{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralRestoreSession},
		{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralRestoreBatch},
		{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralFunction},
		{Group: stashv1beta1.SchemeGroupVersion.Group, Version: stashv1beta1.SchemeGroupVersion.Version, Resource: stashv1beta1.ResourcePluralTask},

		// UI resources
		{Group: uiv1alpha1.SchemeGroupVersion.Group, Version: uiv1alpha1.SchemeGroupVersion.Version, Resource: uiv1alpha1.ResourceBackupOverviews},
	}

	var resources []*apiextensions.CustomResourceDefinition
	for i := range gvrs {
		crd, err := crds.CustomResourceDefinition(gvrs[i])
		if err != nil {
			return err
		}
		resources = append(resources, crd)
	}
	return apiextensions.RegisterCRDs(crdClient, resources)
}
