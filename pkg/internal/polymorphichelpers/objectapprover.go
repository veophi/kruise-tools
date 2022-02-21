package polymorphichelpers

import (
	"errors"
	"fmt"

	kruiserolloutsv1apha1 "github.com/openkruise/rollouts/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubectl/pkg/scheme"
)

func defaultObjectApprover (obj runtime.Object) ([]byte, error) {
	switch obj := obj.(type) {
	case *kruiserolloutsv1apha1.Rollout:
		if obj.Status.CanaryStatus == nil || obj.Status.CanaryStatus.CurrentStepState != kruiserolloutsv1apha1.CanaryStepStatePaused {
			return nil, errors.New("does not allow to approve, because current canary state is not 'StepInPaused'")
		}
		obj.Status.CanaryStatus.CurrentStepState = kruiserolloutsv1apha1.CanaryStepStateCompleted
		return runtime.Encode(scheme.Codecs.LegacyCodec(kruiserolloutsv1apha1.GroupVersion), obj)

	default:
		return nil, fmt.Errorf("approving is not supported")
	}
}