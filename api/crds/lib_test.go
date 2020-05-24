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

package crds

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestCustomResourceDefinition(t *testing.T) {
	type args struct {
		gvr schema.GroupVersionResource
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "tasks",
			args: args{
				gvr: schema.GroupVersionResource{
					Group:    "stash.appscode.com",
					Version:  "v1beta1",
					Resource: "tasks",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CustomResourceDefinition(tt.args.gvr)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomResourceDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.V1 == nil {
				t.Errorf("Missing V1 CustomResourceDefinition for gvr = %v", tt.args.gvr)
			}
			if got.V1beta1 == nil {
				t.Errorf("Missing V1beta1 CustomResourceDefinition for gvr = %v", tt.args.gvr)
			}
		})
	}
}
