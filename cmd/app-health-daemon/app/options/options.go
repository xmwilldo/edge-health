package options

import "github.com/spf13/pflag"

type AppHealthOptions struct {
	JoinIP string
}

func NewAppHealthOptions() *AppHealthOptions {
	return &AppHealthOptions{
		JoinIP: "",
	}
}

func (o *AppHealthOptions) Validate() []error {
	return nil
}

func (o *AppHealthOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.JoinIP, "join", o.JoinIP, "Join IP")
}
