package webhook

import (
	corev1 "k8s.io/api/core/v1"
)

type EdgeServiceAutonomy struct {
	sidecarConfig *Config
}

type Config struct {
	Containers []corev1.Container `yaml:"containers"`
}

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}
