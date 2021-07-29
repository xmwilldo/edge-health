package options

import (
	"fmt"
	"strconv"
	"strings"
)

type ServiceAutonomyEnhancementOptions struct {
	Enabled        bool
	AppStatusSvc   string
	UpdateInterval int
	Delay          int
}

func (p ServiceAutonomyEnhancementOptions) Name() string {
	return "ServiceAutonomyEnhancement"
}

func (p *ServiceAutonomyEnhancementOptions) Set(s string) error {
	var err error

	for _, para := range strings.Split(s, ",") {
		if len(para) == 0 {
			continue
		}
		arr := strings.Split(para, "=")
		trimkey := strings.TrimSpace(arr[0])
		switch trimkey {
		case "address":
			(*p).AppStatusSvc = strings.TrimSpace(arr[1])
		case "interval":
			interval, _ := strconv.Atoi(strings.TrimSpace(arr[1]))
			(*p).UpdateInterval = interval
		case "enabled":
			enabled, _ := strconv.ParseBool(strings.TrimSpace(arr[1]))
			(*p).Enabled = enabled
		}
	}
	return err
}

func (p *ServiceAutonomyEnhancementOptions) String() string {
	return fmt.Sprintf("%#v", *p)
}

func (*ServiceAutonomyEnhancementOptions) Type() string {
	return "ServiceAutonomyEnhancement"
}
