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

package v1alpha1

import (
	"testing"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRepository_UsageAllowed(t *testing.T) {
	repoNamespace := "test-repo-namespace"
	notRepoNamespace := "not-repo-namespace"

	sameNamespace := NamespacesFromSame
	allNamespaces := NamespacesFromAll
	namespacesFromSelector := NamespacesFromSelector

	tests := []struct {
		name        string
		usagePolicy *UsagePolicy
		namespace   string
		want        bool
	}{
		{
			name:        "Allow from same namespace if UsagePolicy is nil",
			usagePolicy: nil,
			namespace:   repoNamespace,
			want:        true,
		},
		{
			name:        "Don't allow from different namespace if UsagePolicy is nil",
			usagePolicy: nil,
			namespace:   notRepoNamespace,
			want:        false,
		},
		{
			name: "Allow from same namespace if allowedNamespaces.From = Same",
			usagePolicy: &UsagePolicy{
				AllowedNamespaces: AllowedNamespaces{
					From: &sameNamespace,
				},
			},
			namespace: repoNamespace,
			want:      true,
		},
		{
			name: "Don't allow from different namespace if allowedNamespaces.From = Same",
			usagePolicy: &UsagePolicy{
				AllowedNamespaces: AllowedNamespaces{
					From: &sameNamespace,
				},
			},
			namespace: notRepoNamespace,
			want:      false,
		},
		{
			name: "Allow from same namespace if allowedNamespaces.From = All",
			usagePolicy: &UsagePolicy{
				AllowedNamespaces: AllowedNamespaces{
					From: &allNamespaces,
				},
			},
			namespace: repoNamespace,
			want:      true,
		},
		{
			name: "Allow from different namespace if allowedNamespaces.From = All",
			usagePolicy: &UsagePolicy{
				AllowedNamespaces: AllowedNamespaces{
					From: &allNamespaces,
				},
			},
			namespace: notRepoNamespace,
			want:      true,
		},
		{
			name: "Allow namespace that matches MatchLabels selector",
			usagePolicy: &UsagePolicy{
				AllowedNamespaces: AllowedNamespaces{
					From: &namespacesFromSelector,
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"kubernetes.io/metadata.name": notRepoNamespace,
						},
					},
				},
			},
			namespace: notRepoNamespace,
			want:      true,
		},
		{
			name: "Don't allow namespace that does not match MatchLabels selector",
			usagePolicy: &UsagePolicy{
				AllowedNamespaces: AllowedNamespaces{
					From: &namespacesFromSelector,
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"kubernetes.io/metadata.name": repoNamespace,
						},
					},
				},
			},
			namespace: notRepoNamespace,
			want:      false,
		},
		{
			name: "Allow namespace that matches MatchExpression selector",
			usagePolicy: &UsagePolicy{
				AllowedNamespaces: AllowedNamespaces{
					From: &namespacesFromSelector,
					Selector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "kubernetes.io/metadata.name",
								Operator: metav1.LabelSelectorOpIn,
								Values:   []string{repoNamespace, notRepoNamespace},
							},
						},
					},
				},
			},
			namespace: notRepoNamespace,
			want:      true,
		},
		{
			name: "Don't allow namespace that does not match MatchExpression selector",
			usagePolicy: &UsagePolicy{
				AllowedNamespaces: AllowedNamespaces{
					From: &namespacesFromSelector,
					Selector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "kubernetes.io/metadata.name",
								Operator: metav1.LabelSelectorOpIn,
								Values:   []string{repoNamespace, notRepoNamespace},
							},
						},
					},
				},
			},
			namespace: "some-other-namespace",
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newTestRepository(repoNamespace, tt.usagePolicy)
			ns := newTestNamespace(tt.namespace)
			if got := r.UsageAllowed(ns); got != tt.want {
				t.Errorf("UsageAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newTestRepository(namespace string, policy *UsagePolicy) *Repository {
	return &Repository{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-repo",
			Namespace: namespace,
		},
		Spec: RepositorySpec{
			UsagePolicy: policy,
		},
	}
}

func newTestNamespace(name string) *core.Namespace {
	return &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"kubernetes.io/metadata.name": name,
			},
		},
	}
}
