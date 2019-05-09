/*
Copyright 2019 The Stash Authors.

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
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
	v1beta1 "stash.appscode.dev/stash/client/clientset/versioned/typed/stash/v1beta1"
)

type FakeStashV1beta1 struct {
	*testing.Fake
}

func (c *FakeStashV1beta1) BackupConfigurations(namespace string) v1beta1.BackupConfigurationInterface {
	return &FakeBackupConfigurations{c, namespace}
}

func (c *FakeStashV1beta1) BackupConfigurationTemplates() v1beta1.BackupConfigurationTemplateInterface {
	return &FakeBackupConfigurationTemplates{c}
}

func (c *FakeStashV1beta1) BackupSessions(namespace string) v1beta1.BackupSessionInterface {
	return &FakeBackupSessions{c, namespace}
}

func (c *FakeStashV1beta1) Functions() v1beta1.FunctionInterface {
	return &FakeFunctions{c}
}

func (c *FakeStashV1beta1) RestoreSessions(namespace string) v1beta1.RestoreSessionInterface {
	return &FakeRestoreSessions{c, namespace}
}

func (c *FakeStashV1beta1) Tasks() v1beta1.TaskInterface {
	return &FakeTasks{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeStashV1beta1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
