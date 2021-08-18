module github.com/xmwilldo/edge-health

go 1.16

require (
	github.com/ghodss/yaml v1.0.0
	github.com/hashicorp/serf v0.9.5
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v0.21.3
	k8s.io/component-base v0.21.3
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.8.0
	sigs.k8s.io/controller-runtime v0.9.1
)

replace (
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.1
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)
