module github.com/cisco-app-networking/nsm-nse

go 1.14

require (
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/golang/protobuf v1.5.3
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/networkservicemesh/networkservicemesh/controlplane/api v0.3.0
	github.com/networkservicemesh/networkservicemesh/pkg v0.3.0
	github.com/networkservicemesh/networkservicemesh/sdk v0.3.0
	github.com/onsi/gomega v1.10.3
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.1.0
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.8.3
	go.ligato.io/cn-infra/v2 v2.5.0-alpha.0.20200313154441-b0d4c1b11c73
	go.ligato.io/vpp-agent/v3 v3.3.0-alpha.0.20210111123645-a04d009c61c5
	golang.org/x/sync v0.4.0
	google.golang.org/grpc v1.58.2
	google.golang.org/grpc/examples v0.0.0-20231024221354-6e14274d008d // indirect
	gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 v3.0.1
	gotest.tools v2.2.0+incompatible
)

replace (
	github.com/census-instrumentation/opencensus-proto v0.1.0-0.20181214143942-ba49f56771b8 => github.com/census-instrumentation/opencensus-proto v0.0.3-0.20181214143942-ba49f56771b8
	github.com/networkservicemesh/networkservicemesh/controlplane => github.com/cisco-app-networking/networkservicemesh/controlplane v0.2.3-vl3
	github.com/networkservicemesh/networkservicemesh/controlplane/api => github.com/cisco-app-networking/networkservicemesh/controlplane/api v0.2.3-vl3
	github.com/networkservicemesh/networkservicemesh/forwarder/api => github.com/cisco-app-networking/networkservicemesh/forwarder/api v0.2.3-vl3
	github.com/networkservicemesh/networkservicemesh/k8s => github.com/cisco-app-networking/networkservicemesh/k8s v0.2.3-vl3
	github.com/networkservicemesh/networkservicemesh/k8s/pkg/apis => github.com/cisco-app-networking/networkservicemesh/k8s/pkg/apis v0.2.3-vl3
	github.com/networkservicemesh/networkservicemesh/pkg => github.com/cisco-app-networking/networkservicemesh/pkg v0.2.3-vl3
	github.com/networkservicemesh/networkservicemesh/sdk => github.com/cisco-app-networking/networkservicemesh/sdk v0.2.3-vl3
	github.com/networkservicemesh/networkservicemesh/utils => github.com/cisco-app-networking/networkservicemesh/utils v0.2.3-vl3
)
