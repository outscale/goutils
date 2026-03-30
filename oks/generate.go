package oks

import (
	_ "github.com/outscale/goutils/oks/apis/oks.dev/v1beta2"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml apis/oks.dev/v1beta2/nodepool.json

//go:generate go run k8s.io/code-generator/cmd/deepcopy-gen github.com/outscale/goutils/oks/apis/oks.dev/v1beta2 --output-file zz_generated.deepcopy.go --go-header-file hack/boilerplate.go.txt

//go:generate go run k8s.io/code-generator/cmd/client-gen --input-base /home/outscale/go/src/github.com/outscale/goutils/oks/apis --output-dir . --output-pkg github.com/outscale/goutils/oks --clientset-name clientset --input oks.dev/v1beta2 --go-header-file hack/boilerplate.go.txt
