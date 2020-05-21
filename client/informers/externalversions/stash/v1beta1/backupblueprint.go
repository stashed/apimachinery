/*
Copyright The Stash Authors.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	time "time"

	stashv1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	versioned "stash.appscode.dev/apimachinery/client/clientset/versioned"
	internalinterfaces "stash.appscode.dev/apimachinery/client/informers/externalversions/internalinterfaces"
	v1beta1 "stash.appscode.dev/apimachinery/client/listers/stash/v1beta1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// BackupBlueprintInformer provides access to a shared informer and lister for
// BackupBlueprints.
type BackupBlueprintInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.BackupBlueprintLister
}

type backupBlueprintInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewBackupBlueprintInformer constructs a new informer for BackupBlueprint type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewBackupBlueprintInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredBackupBlueprintInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredBackupBlueprintInformer constructs a new informer for BackupBlueprint type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredBackupBlueprintInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.StashV1beta1().BackupBlueprints().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.StashV1beta1().BackupBlueprints().Watch(context.TODO(), options)
			},
		},
		&stashv1beta1.BackupBlueprint{},
		resyncPeriod,
		indexers,
	)
}

func (f *backupBlueprintInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredBackupBlueprintInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *backupBlueprintInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&stashv1beta1.BackupBlueprint{}, f.defaultInformer)
}

func (f *backupBlueprintInformer) Lister() v1beta1.BackupBlueprintLister {
	return v1beta1.NewBackupBlueprintLister(f.Informer().GetIndexer())
}
