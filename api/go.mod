module github.com/fluxcd/notification-controller/api

go 1.18

require (
	github.com/fluxcd/pkg/apis/meta v0.18.0
	k8s.io/apimachinery v0.25.4
	sigs.k8s.io/controller-runtime v0.13.1
)

// Fix CVE-2022-32149
replace golang.org/x/text => golang.org/x/text v0.4.0

// Fix CVE-2022-28948
replace gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	k8s.io/utils v0.0.0-20221108210102-8e77b1f39fe2 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)
