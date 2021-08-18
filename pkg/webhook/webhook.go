package webhook

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/xmwilldo/edge-service-autonomy/cmd/webhook/app/options"
	"github.com/xmwilldo/edge-service-autonomy/pkg/util"

	"github.com/ghodss/yaml"
	admissionv1 "k8s.io/api/admission/v1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	ignoredNamespaces = []string{
		metav1.NamespaceSystem,
		metav1.NamespacePublic,
	}
)

const (
	admissionWebhookAnnotationMutateKey = "edge-service-autonomy/enabled"
	admissionWebhookAnnotationStatusKey = "edge-service-autonomy/status"
)

var (
	runtimeScheme = runtime.NewScheme()
)

func init() {
	_ = corev1.AddToScheme(runtimeScheme)
	_ = admissionregistrationv1beta1.AddToScheme(runtimeScheme)
	_ = v1.AddToScheme(runtimeScheme)
}

var _ admission.Handler = &EdgeServiceAutonomy{}

func NewEdgeServiceAutonomy(options options.WebHookOptions) (*EdgeServiceAutonomy, error) {
	sidecarConfig, err := loadConfig(options.SidecarConfig)
	if err != nil {
		return nil, err
	}
	return &EdgeServiceAutonomy{
		sidecarConfig: sidecarConfig,
	}, nil
}

func (e *EdgeServiceAutonomy) Handle(ctx context.Context, req admission.Request) admission.Response {
	response := admission.Response{}
	var (
		availableLabels, availableAnnotations, templateLabels map[string]string
		objectMeta                                            *metav1.ObjectMeta
		resourceNamespace, resourceName                       string
		deployment                                            appsv1.Deployment
	)

	klog.Infof("AdmissionReview for Kind=%v, Namespace=%v Name=%v (%v) UID=%v patchOperation=%v UserInfo=%v",
		req.Kind, req.Namespace, req.Name, resourceName, req.UID, req.Operation, req.UserInfo)

	switch req.Kind.Kind {
	case "Deployment":
		klog.V(4).Infof("Mutating admission Deployment AdmissionRequest = %s", util.AdmissionRequestDebugString(req))

		if err := json.Unmarshal(req.Object.Raw, &deployment); err != nil {
			klog.Errorf("Could not unmarshal raw object: %v", err)
			return admission.Response{
				AdmissionResponse: admissionv1.AdmissionResponse{
					Allowed: false,
					Result: &metav1.Status{
						Status: metav1.StatusFailure, Code: http.StatusBadRequest, Reason: metav1.StatusReasonBadRequest,
						Message: err.Error(),
					},
				},
			}
		}

		klog.V(4).Infof("Admitting deployment = %+v", deployment)

		resourceName, resourceNamespace, objectMeta = deployment.Name, deployment.Namespace, &deployment.ObjectMeta
		availableLabels = deployment.Labels
		templateLabels = deployment.Spec.Template.ObjectMeta.Labels
	}

	if !mutationRequired(ignoredNamespaces, objectMeta) {
		klog.Infof("Skipping mutation for %s/%s due to policy check", resourceNamespace, resourceName)
		return admission.Response{
			AdmissionResponse: admissionv1.AdmissionResponse{
				Allowed: true,
			},
		}
	}

	annotations := map[string]string{admissionWebhookAnnotationStatusKey: "enabled"}
	patchBytes, err := createPatch(&deployment, e.sidecarConfig, availableAnnotations, annotations, availableLabels, templateLabels, nil)
	if err != nil {
		return admission.Response{
			AdmissionResponse: admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Status: metav1.StatusFailure, Code: http.StatusBadRequest, Reason: metav1.StatusReasonBadRequest,
					Message: err.Error(),
				},
			},
		}
	}

	klog.Infof("AdmissionResponse: patch=%v\n", string(patchBytes))
	response.PatchType = new(admissionv1.PatchType)
	*response.PatchType = admissionv1.PatchTypeJSONPatch
	response.Patch = patchBytes
	response.Allowed = true
	return response
}

func admissionRequired(ignoredList []string, admissionAnnotationKey string, metadata *metav1.ObjectMeta) bool {
	// skip special kubernetes system namespaces
	for _, namespace := range ignoredList {
		if metadata.Namespace == namespace {
			klog.Infof("Skip validation for %v for it's in special namespace:%v", metadata.Name, metadata.Namespace)
			return false
		}
	}

	annotations := metadata.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	var required bool
	switch strings.ToLower(annotations[admissionAnnotationKey]) {
	default:
		required = true
	case "n", "no", "false", "off":
		required = false
	}
	return required
}

func mutationRequired(ignoredList []string, metadata *metav1.ObjectMeta) bool {
	required := admissionRequired(ignoredList, admissionWebhookAnnotationMutateKey, metadata)
	annotations := metadata.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	status := annotations[admissionWebhookAnnotationStatusKey]

	if strings.ToLower(status) == "enabled" {
		required = false
	}

	klog.Infof("Mutation policy for %v/%v: required:%v", metadata.Namespace, metadata.Name, required)
	return required
}

func updateAnnotation(target map[string]string, added map[string]string) (patch []patchOperation) {
	for key, value := range added {
		if target == nil || target[key] == "" {
			target = map[string]string{}
			patch = append(patch, patchOperation{
				Op:   "add",
				Path: "/metadata/annotations",
				Value: map[string]string{
					key: value,
				},
			})
		} else {
			patch = append(patch, patchOperation{
				Op:    "replace",
				Path:  "/metadata/annotations/" + key,
				Value: value,
			})
		}
	}
	return patch
}

func updateLabels(target map[string]string, added map[string]string) (patch []patchOperation) {
	values := make(map[string]string)
	for key, value := range added {
		if target == nil || target[key] == "" {
			values[key] = value
		}
	}
	patch = append(patch, patchOperation{
		Op:    "add",
		Path:  "/metadata/labels",
		Value: values,
	})
	return patch
}

func updateTemplateLabels(target map[string]string, added map[string]string) (patch []patchOperation) {
	values := target
	for key, value := range added {
		v, ok := values[key]
		if !ok {
			values[key] = value
		} else {
			if v != value {
				values[key] = value
			}
		}
	}
	patch = append(patch, patchOperation{
		Op:    "add",
		Path:  "/spec/template/metadata/labels",
		Value: values,
	})
	return patch
}

func createPatch(deployment *appsv1.Deployment, sidecarConfig *Config, availableAnnotations map[string]string, annotations map[string]string, availableLabels map[string]string, templateLabels map[string]string, labels map[string]string) ([]byte, error) {
	var patch []patchOperation

	if !reflect.DeepEqual(deployment, &appsv1.Deployment{}) {
		patch = append(patch, addContainer(deployment.Spec.Template.Spec.Containers, sidecarConfig.Containers, "/spec/template/spec/containers")...)
		patch = append(patch, updateTemplateLabels(templateLabels, labels)...)
	}
	patch = append(patch, updateAnnotation(availableAnnotations, annotations)...)
	patch = append(patch, updateLabels(availableLabels, labels)...)

	return json.Marshal(patch)
}

func addContainer(target, added []corev1.Container, basePath string) (patch []patchOperation) {
	first := len(target) == 0
	var value interface{}
	for _, add := range added {
		value = add
		path := basePath
		if first {
			first = false
			value = []corev1.Container{add}
		} else {
			path = path + "/-"
		}
		patch = append(patch, patchOperation{
			Op:    "add",
			Path:  path,
			Value: value,
		})
	}
	return patch
}

func loadConfig(configFile string) (*Config, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	klog.Infof("New configuration: sha256sum %x", sha256.Sum256(data))

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
