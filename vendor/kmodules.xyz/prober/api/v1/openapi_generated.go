//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"kmodules.xyz/prober/api/v1.FormEntry":      schema_kmodulesxyz_prober_api_v1_FormEntry(ref),
		"kmodules.xyz/prober/api/v1.HTTPPostAction": schema_kmodulesxyz_prober_api_v1_HTTPPostAction(ref),
		"kmodules.xyz/prober/api/v1.Handler":        schema_kmodulesxyz_prober_api_v1_Handler(ref),
	}
}

func schema_kmodulesxyz_prober_api_v1_FormEntry(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: []string{"object"},
				Properties: map[string]spec.Schema{
					"key": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"values": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func schema_kmodulesxyz_prober_api_v1_HTTPPostAction(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "HTTPPostAction describes an action based on HTTP Post requests.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"path": {
						SchemaProps: spec.SchemaProps{
							Description: "Path to access on the HTTP server.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"port": {
						SchemaProps: spec.SchemaProps{
							Description: "Name or number of the port to access on the container. Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME.",
							Default:     map[string]interface{}{},
							Ref:         ref("k8s.io/apimachinery/pkg/util/intstr.IntOrString"),
						},
					},
					"host": {
						SchemaProps: spec.SchemaProps{
							Description: "Host name to connect to, defaults to the pod IP. You probably want to set \"Host\" in httpHeaders instead.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"scheme": {
						SchemaProps: spec.SchemaProps{
							Description: "Scheme to use for connecting to the host. Defaults to HTTP.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"httpHeaders": {
						SchemaProps: spec.SchemaProps{
							Description: "Custom headers to set in the request. HTTP allows repeated headers.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("k8s.io/api/core/v1.HTTPHeader"),
									},
								},
							},
						},
					},
					"body": {
						SchemaProps: spec.SchemaProps{
							Description: "Body to set in the request.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"form": {
						SchemaProps: spec.SchemaProps{
							Description: "Form to set in the request body.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("kmodules.xyz/prober/api/v1.FormEntry"),
									},
								},
							},
						},
					},
				},
				Required: []string{"port"},
			},
		},
		Dependencies: []string{
			"k8s.io/api/core/v1.HTTPHeader", "k8s.io/apimachinery/pkg/util/intstr.IntOrString", "kmodules.xyz/prober/api/v1.FormEntry"},
	}
}

func schema_kmodulesxyz_prober_api_v1_Handler(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Handler defines a specific action that should be taken",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"exec": {
						SchemaProps: spec.SchemaProps{
							Description: "One and only one of the following should be specified. Exec specifies the action to take.",
							Ref:         ref("k8s.io/api/core/v1.ExecAction"),
						},
					},
					"httpGet": {
						SchemaProps: spec.SchemaProps{
							Description: "HTTPGet specifies the http Get request to perform.",
							Ref:         ref("k8s.io/api/core/v1.HTTPGetAction"),
						},
					},
					"httpPost": {
						SchemaProps: spec.SchemaProps{
							Description: "HTTPPost specifies the http Post request to perform.",
							Ref:         ref("kmodules.xyz/prober/api/v1.HTTPPostAction"),
						},
					},
					"tcpSocket": {
						SchemaProps: spec.SchemaProps{
							Description: "TCPSocket specifies an action involving a TCP port. TCP hooks not yet supported",
							Ref:         ref("k8s.io/api/core/v1.TCPSocketAction"),
						},
					},
					"containerName": {
						SchemaProps: spec.SchemaProps{
							Description: "ContainerName specifies the name of the container where to execute the commands for Exec probe or where to find the port for HTTP or TCP probe",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/api/core/v1.ExecAction", "k8s.io/api/core/v1.HTTPGetAction", "k8s.io/api/core/v1.TCPSocketAction", "kmodules.xyz/prober/api/v1.HTTPPostAction"},
	}
}
