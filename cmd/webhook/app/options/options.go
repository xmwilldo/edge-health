package options

import "github.com/spf13/pflag"

type WebHookOptions struct {
	Port          int
	CertFile      string
	KeyFile       string
	SidecarConfig string
}

func NewWebHookOptions() WebHookOptions {
	return WebHookOptions{
		Port:          7777,
		CertFile:      "/etc/edge-service-autonomy-webhook/certs/cert.pem",
		KeyFile:       "/etc/edge-service-autonomy-webhook/certs/key.pem",
		SidecarConfig: "/etc/edge-service-autonomy-webhook/config/sidecarconfig.yaml",
	}
}

func (o *WebHookOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.IntVar(&o.Port, "port", o.Port, "The port of webhook server to listen.")
	fs.StringVar(&o.CertFile, "tlsCertPath", o.CertFile, "The path of tls cert")
	fs.StringVar(&o.KeyFile, "tlsKeyPath", o.KeyFile, "The path of tls key")
	fs.StringVar(&o.SidecarConfig, "sidecarConfig", o.SidecarConfig, "File containing the mutation configuration.")
}
