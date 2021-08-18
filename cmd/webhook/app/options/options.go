package options

import (
	"flag"

	"github.com/spf13/pflag"
)

type WebHookOptions struct {
	Port          int
	KubeConfig    string
	MasterURL     string
	CertDir       string
	SidecarConfig string
	VerFlag       bool
}

func NewDefaultWebHookOptions() WebHookOptions {
	return WebHookOptions{
		Port:          443,
		KubeConfig:    "/root/.kube/config",
		MasterURL:     "",
		CertDir:       "/etc/edge-service-autonomy-webhook/certs",
		SidecarConfig: "/etc/edge-service-autonomy-webhook/config/sidecarconfig.yaml",
		VerFlag:       false,
	}
}

func (o *WebHookOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	// Add the command line flags from other dependencies(klog, kubebuilder, etc.)
	fs.AddGoFlagSet(flag.CommandLine)

	fs.IntVar(&o.Port, "secure_port", o.Port, "The port on which to serve HTTPS.")
	fs.StringVar(&o.KubeConfig, "kubeconfig", o.KubeConfig, "Path to a kubeconfig. Only required if out-of-cluster.")
	fs.StringVar(&o.MasterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	fs.StringVar(&o.CertDir, "cert_dir", "", "The directory where the TLS certs are located.")
	fs.StringVar(&o.SidecarConfig, "sidecar_config", o.SidecarConfig, "File containing the mutation configuration.")
	fs.BoolVar(&o.VerFlag, "version", o.VerFlag, "Prints the Version info of webhook.")
}
