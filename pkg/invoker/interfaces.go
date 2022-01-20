package invoker

import (
	"stash.appscode.dev/apimachinery/apis/stash/v1beta1"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

type MetadataHandler interface {
	GetObjectMeta() metav1.ObjectMeta
	GetTypeMeta() metav1.TypeMeta
	GetObjectRef() (*core.ObjectReference, error)
	GetOwnerRef() *metav1.OwnerReference
	GetLabels() map[string]string
	AddFinalizer() error
	RemoveFinalizer() error
}

type ConditionHandler interface {
	HasCondition(target *v1beta1.TargetRef, conditionType string) (bool, error)
	GetCondition(target *v1beta1.TargetRef, conditionType string) (int, *kmapi.Condition, error)
	SetCondition(target *v1beta1.TargetRef, newCondition kmapi.Condition) error
	IsConditionTrue(target *v1beta1.TargetRef, conditionType string) (bool, error)
}
