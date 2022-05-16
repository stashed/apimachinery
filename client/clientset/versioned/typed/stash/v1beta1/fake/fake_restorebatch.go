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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
)

// FakeRestoreBatches implements RestoreBatchInterface
type FakeRestoreBatches struct {
	Fake *FakeStashV1beta1
	ns   string
}

var restorebatchesResource = schema.GroupVersionResource{Group: "stash.appscode.com", Version: "v1beta1", Resource: "restorebatches"}

var restorebatchesKind = schema.GroupVersionKind{Group: "stash.appscode.com", Version: "v1beta1", Kind: "RestoreBatch"}

// Get takes name of the restoreBatch, and returns the corresponding restoreBatch object, and an error if there is any.
func (c *FakeRestoreBatches) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.RestoreBatch, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(restorebatchesResource, c.ns, name), &v1beta1.RestoreBatch{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.RestoreBatch), err
}

// List takes label and field selectors, and returns the list of RestoreBatches that match those selectors.
func (c *FakeRestoreBatches) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.RestoreBatchList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(restorebatchesResource, restorebatchesKind, c.ns, opts), &v1beta1.RestoreBatchList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.RestoreBatchList{ListMeta: obj.(*v1beta1.RestoreBatchList).ListMeta}
	for _, item := range obj.(*v1beta1.RestoreBatchList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested restoreBatches.
func (c *FakeRestoreBatches) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(restorebatchesResource, c.ns, opts))

}

// Create takes the representation of a restoreBatch and creates it.  Returns the server's representation of the restoreBatch, and an error, if there is any.
func (c *FakeRestoreBatches) Create(ctx context.Context, restoreBatch *v1beta1.RestoreBatch, opts v1.CreateOptions) (result *v1beta1.RestoreBatch, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(restorebatchesResource, c.ns, restoreBatch), &v1beta1.RestoreBatch{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.RestoreBatch), err
}

// Update takes the representation of a restoreBatch and updates it. Returns the server's representation of the restoreBatch, and an error, if there is any.
func (c *FakeRestoreBatches) Update(ctx context.Context, restoreBatch *v1beta1.RestoreBatch, opts v1.UpdateOptions) (result *v1beta1.RestoreBatch, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(restorebatchesResource, c.ns, restoreBatch), &v1beta1.RestoreBatch{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.RestoreBatch), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRestoreBatches) UpdateStatus(ctx context.Context, restoreBatch *v1beta1.RestoreBatch, opts v1.UpdateOptions) (*v1beta1.RestoreBatch, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(restorebatchesResource, "status", c.ns, restoreBatch), &v1beta1.RestoreBatch{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.RestoreBatch), err
}

// Delete takes name of the restoreBatch and deletes it. Returns an error if one occurs.
func (c *FakeRestoreBatches) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(restorebatchesResource, c.ns, name), &v1beta1.RestoreBatch{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRestoreBatches) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(restorebatchesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.RestoreBatchList{})
	return err
}

// Patch applies the patch and returns the patched restoreBatch.
func (c *FakeRestoreBatches) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.RestoreBatch, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(restorebatchesResource, c.ns, name, pt, data, subresources...), &v1beta1.RestoreBatch{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.RestoreBatch), err
}
