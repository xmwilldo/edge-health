package webhook

import (
	"net/http"

	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

type Server interface {
	mutating(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse
	validating(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse
	Start()
	Stop()
}

type webHookServer struct {
	server        *http.Server
	sidecarConfig *Config
}

type Config struct {
	Containers []corev1.Container `yaml:"containers"`
	Volumes    []corev1.Volume    `yaml:"volumes"`
}

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}
