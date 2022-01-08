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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindBackupOverview = "BackupOverview"
	ResourceBackupOverview     = "backupoverview"
	ResourceBackupOverviews    = "backupoverviews"
)

// BackupOverviewRequest defines the request fields of BackupOverview
type BackupOverviewRequest struct {
	Group    string `json:"group" protobuf:"bytes,1,opt,name=group"`
	Version  string `json:"version" protobuf:"version,2,opt,name=version"`
	Resource string `json:"resource" protobuf:"resource,3,opt,name=resource"`
}

type BackupStatus string

const (
	BackupStatusActive = "Active"
	BackupStatusPaused = "Paused"
)

// BackupOverviewResponse defines the response fields of BackupOverview
type BackupOverviewResponse struct {
	Schedule           string       `json:"schedule,omitempty" protobuf:"bytes,1,opt,name=schedule"`
	Status             BackupStatus `json:"status,omitempty" protobuf:"bytes,2,opt,name=status"`
	LastBackupTime     *metav1.Time `json:"lastBackupTime,omitempty" protobuf:"bytes,3,opt,name=lastBackupTime"`
	UpcomingBackupTime *metav1.Time `json:"upcomingBackupTime,omitempty" protobuf:"bytes,4,opt,name=upcomingBackupTime"`
	Repository         string       `json:"repository,omitempty" protobuf:"bytes,5,opt,name=repository"`
	DataSize           string       `json:"dataSize,omitempty" protobuf:"bytes,6,opt,name=dataSize"`
	NumberOfSnapshots  int64        `json:"numberOfSnapshots,omitempty" protobuf:"bytes,7,opt,name=numberOfSnapshots"`
	DataIntegrity      bool         `json:"dataIntegrity,omitempty" protobuf:"bytes,8,opt,name=dataIntegrity"`
}

// BackupOverview is the Schema for the BackupOverviews API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BackupOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Request  BackupOverviewRequest  `json:"request,omitempty" protobuf:"bytes,2,opt,name=request"`
	Response BackupOverviewResponse `json:"response,omitempty" protobuf:"bytes,3,opt,name=response"`
}

// BackupOverviewList contains a list of BackupOverview

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BackupOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []BackupOverview `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func init() {
	SchemeBuilder.Register(&BackupOverview{}, &BackupOverviewList{})
}
