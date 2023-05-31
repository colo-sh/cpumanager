module github.com/colo-sh/cpumanager

go 1.16

require (
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/prometheus/client_golang v1.13.0 // indirect
	github.com/stretchr/testify v1.8.1
	golang.org/x/net v0.4.0 // indirect
	google.golang.org/grpc v1.40.0 // indirect
	k8s.io/api v0.22.6
	k8s.io/apimachinery v0.22.6
	k8s.io/apiserver v0.22.6
	k8s.io/cri-api v0.22.6
	k8s.io/klog/v2 v2.10.0
	k8s.io/kubernetes v1.22.6
)

require (
	github.com/Microsoft/go-winio v0.4.17 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/cadvisor v0.39.3
	github.com/google/go-cmp v0.5.9
	github.com/google/uuid v1.3.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/onsi/gomega v1.15.0 // indirect
	github.com/opencontainers/runc v1.0.2
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/zap v1.19.0 // indirect
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	k8s.io/component-base v0.27.2
	k8s.io/kube-scheduler v0.22.6 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect

)

replace (
	github.com/alecthomas/units => github.com/alecthomas/units v0.0.0-20210912230133-d1bdfacee922
	github.com/cespare/xxhash/v2 => github.com/cespare/xxhash/v2 v2.1.2
	github.com/colo-sh/cpumanager => ./
	github.com/docker/docker => github.com/docker/engine v0.0.0-20190822180741-9552f2b2fdde
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.2
	github.com/iovisor/gobpf => github.com/iovisor/gobpf v0.2.0
	github.com/prometheus/common => github.com/prometheus/common v0.30.0
	github.com/prometheus/procfs => github.com/prometheus/procfs v0.7.3
	golang.org/x/sys => golang.org/x/sys v0.0.0-20210923061019-b8560ed6a9b7
	google.golang.org/protobuf => google.golang.org/protobuf v1.27.1
	k8s.io/api => k8s.io/api v0.22.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.22.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.6
	k8s.io/apiserver => k8s.io/apiserver v0.22.6
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.22.6
	k8s.io/client-go => k8s.io/client-go v0.22.6
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.22.6
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.22.6
	k8s.io/code-generator => k8s.io/code-generator v0.22.6
	k8s.io/component-base => k8s.io/component-base v0.22.6
	k8s.io/component-helpers => k8s.io/component-helpers v0.22.6
	k8s.io/controller-manager => k8s.io/controller-manager v0.22.6
	k8s.io/cri-api => k8s.io/cri-api v0.22.6
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.22.6
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.22.6
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.22.6
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.22.6
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.22.6
	k8s.io/kubectl => k8s.io/kubectl v0.22.6
	k8s.io/kubelet => k8s.io/kubelet v0.22.6
	k8s.io/kubernetes => k8s.io/kubernetes v1.22.6
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.22.6
	k8s.io/metrics => k8s.io/metrics v0.22.6
	k8s.io/mount-utils => k8s.io/mount-utils v0.22.6
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.22.6
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.22.6
)
