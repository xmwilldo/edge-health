package options

import "github.com/spf13/pflag"

type AppHealthOptions struct {
	JoinIp    string
	Namespace string
	SvcName   string
}

func NewAppHealthOptions() *AppHealthOptions {
	return &AppHealthOptions{
		JoinIp:    "",
		Namespace: "",
		SvcName:   "",
	}
}

func (o *AppHealthOptions) Validate() []error {
	return nil
}

func (o *AppHealthOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.JoinIp, "join", o.JoinIp, "Join IP")
	fs.StringVar(&o.Namespace, "namespace", o.Namespace, "")
	fs.StringVar(&o.SvcName, "svc_name", o.SvcName, "")
}
