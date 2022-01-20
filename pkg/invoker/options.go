package invoker

import (
	"fmt"

	"stash.appscode.dev/apimachinery/apis/stash/v1beta1"

	"k8s.io/apimachinery/pkg/runtime"
	kmapi "kmodules.xyz/client-go/api/v1"
)

type StashInvoker interface {
	FinalizerHandler
	ConditionHandler
	ExecutionOrderHandler
}

type FinalizerHandler interface {
	AddFinalizer() error
	RemoveFinalizer() error
}

type ConditionHandler interface {
	HasCondition(conditionType string) (bool, error)
	GetCondition(conditionType string) (int, *kmapi.Condition, error)
	SetCondition(newCondition kmapi.Condition) error
	IsConditionTrue(conditionType string) (bool, error)
}

type ExecutionOrderHandler interface {
	NextInOrder(curTarget v1beta1.TargetRef) bool
}

func New(obj runtime.Object) (StashInvoker, error) {
	switch obj.(type) {
	case *v1beta1.BackupConfiguration:
		return &BackupConfigurationInvoker{}, nil
	default:
		return nil, fmt.Errorf("unknown invoker type: %s", obj.GetObjectKind().GroupVersionKind().String())
	}
}
