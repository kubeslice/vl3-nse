// Copyright 2019 Cisco Systems, Inc.
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/networkservice"
	"github.com/networkservicemesh/networkservicemesh/pkg/tools"
	"github.com/networkservicemesh/networkservicemesh/sdk/common"
	"github.com/networkservicemesh/networkservicemesh/sdk/endpoint"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connectioncontext"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/cisco-app-networking/nsm-nse/pkg/metrics"
	"github.com/cisco-app-networking/nsm-nse/pkg/nseconfig"
	"github.com/cisco-app-networking/nsm-nse/pkg/universal-cnf/ucnf"
	"github.com/cisco-app-networking/nsm-nse/pkg/universal-cnf/vppagent"
)

const (
	metricsPortEnv     = "METRICS_PORT"
	metricsPath        = "/metrics"
	metricsPortDefault = "2112"
)

const (
	defaultConfigPath   = "/etc/universal-cnf/config.yaml"
	defaultPluginModule = ""
)

// Flags holds the command line flags as supplied with the binary invocation
type Flags struct {
	ConfigPath string
	Verify     bool
}

type fnGetNseName func() string

// Process will parse the command line flags and init the structure members
func (mf *Flags) Process() {
	flag.StringVar(&mf.ConfigPath, "file", defaultConfigPath, " full path to the configuration file")
	flag.BoolVar(&mf.Verify, "verify", false, "only verify the configuration, don't run")
	flag.Parse()
}

type vL3CompositeEndpoint struct {
}

func (e vL3CompositeEndpoint) AddCompositeEndpoints(nsConfig *common.NSConfiguration, ucnfEndpoint *nseconfig.Endpoint) *[]networkservice.NetworkServiceServer {

	logrus.WithFields(logrus.Fields{
		"prefixPool":         nsConfig.IPAddress,
		"nsConfig.IPAddress": nsConfig.IPAddress,
	}).Infof("Creating vL3 IPAM endpoint")

	var nsRemoteIpList []string
	nsRemoteIpListStr, ok := os.LookupEnv("NSM_REMOTE_NS_IP_LIST")
	if ok {
		nsRemoteIpList = strings.Split(nsRemoteIpListStr, ",")
	}
	compositeEndpoints := []networkservice.NetworkServiceServer{
		newVL3ConnectComposite(nsConfig, nsConfig.IPAddress,
			&vppagent.UniversalCNFVPPAgentBackend{}, nsRemoteIpList, func() string {
				return ucnfEndpoint.NseName
			}, ucnfEndpoint.VL3.IPAM.DefaultPrefixPool, ucnfEndpoint.VL3.IPAM.ServerAddress, ucnfEndpoint.NseControl.ConnectivityDomain),
	}

	return &compositeEndpoints
}

func InitializeMetrics() {
	metricsPort := os.Getenv(metricsPortEnv)
	if metricsPort == "" {
		metricsPort = metricsPortDefault
	}
	addr := fmt.Sprintf("0.0.0.0:%v", metricsPort)
	logrus.WithField("path", metricsPath).Infof("Serving metrics on: %v", addr)
	metrics.ServeMetrics(addr, metricsPath)
}

var (
	nsmEndpoint endpoint.NsmEndpoint
)

func getStringListFromString(str string) []string {
	return strings.Split(str, ",")
}

func main() {

	// Capture signals to cleanup before exiting
	c := tools.NewOSSignalChannel()

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)

	dataplane := os.Getenv("DATAPLANE")

	if dataplane == "vpp" {

		mainFlags := &Flags{}
		mainFlags.Process()

		InitializeMetrics()

		// Capture signals to cleanup before exiting
		prometheus.NewBuildInfoCollector()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		vl3 := vL3CompositeEndpoint{}
		ucnfNse := ucnf.NewUcnfNse(mainFlags.ConfigPath, mainFlags.Verify, &vppagent.UniversalCNFVPPAgentBackend{}, vl3, ctx)
		logrus.Info("endpoint started")

		defer ucnfNse.Cleanup()
		<-c

	} else {

		configuration := common.FromEnv()

		endpoints := []networkservice.NetworkServiceServer{
			endpoint.NewMonitorEndpoint(configuration),
			endpoint.NewConnectionEndpoint(configuration),
			endpoint.NewIpamEndpoint(configuration),
		}

		dstRoutes := os.Getenv("DST_ROUTES")
		if dstRoutes != "" {
			dstRouteMutator := endpoint.CreateRouteMutator(getStringListFromString(dstRoutes))
			endpoints = append(endpoints, endpoint.NewCustomFuncEndpoint("route", dstRouteMutator))
		}

		dnsNameServers := os.Getenv("DNS_NAMESERVERS")
		dnsDomains := os.Getenv("DNS_DOMAINS")
		if dnsNameServers != "" {
			dnsMutator := endpoint.NewAddDNSConfigs(&connectioncontext.DNSConfig{
				DnsServerIps: getStringListFromString(dnsNameServers),
				SearchDomains: getStringListFromString(dnsDomains),
			})
			endpoints = append(endpoints, dnsMutator)
		}

		composite := endpoint.NewCompositeEndpoint(endpoints...)

		nsme, err := endpoint.NewNSMEndpoint(context.Background(), configuration, composite)
		if err != nil {
			logrus.Fatalf("%v", err)
		}

		nsmEndpoint = nsme

		err = nsme.Start()
		if err != nil {
			logrus.Fatalf("Unable to start the endpoint: %v", err)
		}

		logrus.Infof("Started NSE --got name %s", nsme.GetName())
		defer func() { _ = nsme.Delete() }()

		// Capture signals to cleanup before exiting
		<-c
	}
}

func GetMyNseName() string {
	return nsmEndpoint.GetName()
}
