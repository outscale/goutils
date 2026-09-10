module github.com/outscale/goutils/k8s

go 1.27.0

toolchain go1.27.1

require (
	dario.cat/mergo v1.0.2
	github.com/outscale/goutils/sdk v0.0.9
	github.com/outscale/osc-sdk-go/v3 v3.0.0-rc.5
	github.com/spf13/pflag v1.0.10
	github.com/stretchr/testify v1.12.1
	go.uber.org/mock v0.6.0
	k8s.io/klog/v2 v2.140.0
)

// replace github.com/outscale/goutils/sdk => ../sdk

require (
	github.com/aws/smithy-go/aws-http-auth v1.2.1 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/samber/lo v1.53.0 // indirect
	go.uber.org/ratelimit v0.3.1 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
