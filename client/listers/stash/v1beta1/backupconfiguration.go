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

// Code generated by lister-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
)

// BackupConfigurationLister helps list BackupConfigurations.
// All objects returned here must be treated as read-only.
type BackupConfigurationLister interface {
	// List lists all BackupConfigurations in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.BackupConfiguration, err error)
	// BackupConfigurations returns an object that can list and get BackupConfigurations.
	BackupConfigurations(namespace string) BackupConfigurationNamespaceLister
	BackupConfigurationListerExpansion
}

// backupConfigurationLister implements the BackupConfigurationLister interface.
type backupConfigurationLister struct {
	indexer cache.Indexer
}

// NewBackupConfigurationLister returns a new BackupConfigurationLister.
func NewBackupConfigurationLister(indexer cache.Indexer) BackupConfigurationLister {
	return &backupConfigurationLister{indexer: indexer}
}

// List lists all BackupConfigurations in the indexer.
func (s *backupConfigurationLister) List(selector labels.Selector) (ret []*v1beta1.BackupConfiguration, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.BackupConfiguration))
	})
	return ret, err
}

// BackupConfigurations returns an object that can list and get BackupConfigurations.
func (s *backupConfigurationLister) BackupConfigurations(namespace string) BackupConfigurationNamespaceLister {
	return backupConfigurationNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// BackupConfigurationNamespaceLister helps list and get BackupConfigurations.
// All objects returned here must be treated as read-only.
type BackupConfigurationNamespaceLister interface {
	// List lists all BackupConfigurations in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.BackupConfiguration, err error)
	// Get retrieves the BackupConfiguration from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta1.BackupConfiguration, error)
	BackupConfigurationNamespaceListerExpansion
}

// backupConfigurationNamespaceLister implements the BackupConfigurationNamespaceLister
// interface.
type backupConfigurationNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all BackupConfigurations in the indexer for a given namespace.
func (s backupConfigurationNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.BackupConfiguration, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.BackupConfiguration))
	})
	return ret, err
}

// Get retrieves the BackupConfiguration from the indexer for a given namespace and name.
func (s backupConfigurationNamespaceLister) Get(name string) (*v1beta1.BackupConfiguration, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("backupconfiguration"), name)
	}
	return obj.(*v1beta1.BackupConfiguration), nil
}
