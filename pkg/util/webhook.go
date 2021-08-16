package util

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func AdmissionRequestDebugString(a admission.Request) string {
	return fmt.Sprintf("UID=%v Kind={%v} Resource=%+v SubResource=%v Name=%v Namespace=%v Operation=%v UserInfo=%+v DryRun=%v",
		a.UID, a.Kind, a.Resource, a.SubResource, a.Name, a.Namespace, a.Operation, a.UserInfo, *a.DryRun)
}
