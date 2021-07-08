module github.com/xmwilldo/edge-service-autonomy

go 1.16

require (
	github.com/hashicorp/serf v0.9.5
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	go.uber.org/zap v1.14.1 // indirect
	golang.org/x/mod v0.3.0 // indirect
	golang.org/x/sys v0.0.0-20201112073958-5cba982894dd
	golang.org/x/tools v0.0.0-20200616133436-c1934b75d054 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.20.5
	k8s.io/apimachinery v0.20.5
	k8s.io/client-go v0.20.5
	k8s.io/component-base v0.20.5
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.8.0
)

replace (
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.1
	github.com/moby/term => github.com/moby/term v0.0.0-20200312100748-672ec06f55cd
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
	k8s.io/api => k8s.io/api v0.20.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.5
	k8s.io/client-go => k8s.io/client-go v0.20.5
	k8s.io/component-base => k8s.io/component-base v0.20.5
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20201113171705-d219536bb9fd
)
