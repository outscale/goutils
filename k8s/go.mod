module github.com/outscale/goutils/k8s

go 1.25.3

require (
	dario.cat/mergo v1.0.2
	github.com/outscale/goutils/sdk v0.0.0-20260127143749-d95db5597c97
	github.com/outscale/osc-sdk-go/v3 v3.0.0-beta.3
	github.com/spf13/pflag v1.0.10
	github.com/stretchr/testify v1.11.1
	go.uber.org/mock v0.6.0
	k8s.io/klog/v2 v2.130.1
)

// replace github.com/outscale/goutils/sdk => ../sdk

require (
	github.com/aws/smithy-go/aws-http-auth v1.0.0 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	github.com/oapi-codegen/runtime v1.1.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.uber.org/ratelimit v0.3.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
