/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/platform9/pf9ctl/pkg/qbert"
	"github.com/spf13/cobra"
)

var (
	ccName                  string
	ccContainerCIDR         string
	ccServiceCIDR           string
	ccMasterVirtualIP       string
	ccMasterVirtualIPIface  string
	ccAllowWorkloadOnMaster bool
	ccPrivileged            bool
	ccDnsName               string
	ccNetworkPlugin         string
	ccMetalLBAddressPool    string
	ccEnableMetalLb         bool
	ccMasterless            bool
	ccNodePoolUUID          string
)

// clusterCmd represents the cluster command
var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Create a new managed PMK cluster",
	Run:   createCluster,
}

func init() {
	createCmd.AddCommand(createClusterCmd)

	createClusterCmd.Flags().StringVarP(&ccName, "name", "n", ccName, "Cluster name")
	createClusterCmd.Flags().StringVarP(&ccContainerCIDR, "container-cidr", "", ccContainerCIDR, "Container CIDR range")
	createClusterCmd.Flags().StringVarP(&ccServiceCIDR, "service-cidr", "", ccServiceCIDR, "Service CIDR range")
	createClusterCmd.Flags().StringVarP(&ccMasterVirtualIP, "master-vip", "", ccMasterVirtualIP, "Master virtual IP address")
	createClusterCmd.Flags().StringVarP(&ccMasterVirtualIPIface, "master-veth", "", ccMasterVirtualIPIface, "Master virtual IP interface")
	createClusterCmd.Flags().BoolVarP(&ccAllowWorkloadOnMaster, "workload-on-master", "", ccAllowWorkloadOnMaster, "Allow workloads on master node(s)")
	createClusterCmd.Flags().BoolVarP(&ccPrivileged, "privileged", "", privileged, "Allow privileged workloads")
	createClusterCmd.Flags().StringVarP(&ccDnsName, "dns-name", "", ccDnsName, "DNS hostname of cluster")
	createClusterCmd.Flags().StringVarP(&ccNetworkPlugin, "network-plugin", "", ccNetworkPlugin, "CNI to use. Can be either \"flannel\" or \"calico\".")
	createClusterCmd.Flags().StringVarP(&ccMetalLBAddressPool, "metallb-pool", "", ccMetalLBAddressPool, "Address pool for MetalLB")
	createClusterCmd.Flags().BoolVarP(&ccEnableMetalLb, "metallb", "", ccEnableMetalLb, "Use MetalLB")
	createClusterCmd.Flags().BoolVarP(&ccMasterless, "masterless", "", ccMasterless, "Run in masterless configuration")
	createClusterCmd.Flags().StringVarP(&ccNodePoolUUID, "node-pool-uuid", "", ccNodePoolUUID, "UUID of node pool")
}

func createCluster(cmd *cobra.Command, args []string) {
	var (
		fqdn string
		cni  qbert.CNIBackend
	)

	cni = qbert.CNIBackend(ccNetworkPlugin)
	qb := qbert.NewQbert(fqdn)

	req := qbert.ClusterCreateRequest{
		Name:                  ccName,
		ContainerCIDR:         ccContainerCIDR,
		ServiceCIDR:           ccServiceCIDR,
		MasterVirtualIP:       ccMasterVirtualIP,
		MasterVirtualIPIface:  ccMasterVirtualIPIface,
		AllowWorkloadOnMaster: ccAllowWorkloadOnMaster,
		Privileged:            ccPrivileged,
		ExternalDNSName:       ccDnsName,
		NetworkPlugin:         cni,
		MetalLBAddressPool:    ccMetalLBAddressPool,
		NodePoolUUID:          ccNodePoolUUID,
		EnableMetalLb:         ccEnableMetalLb,
		Masterless:            ccMasterless,
	}

	// Wrapper functions to help load source tokens
	_, client := loadCredentials()
	ks, _ := getAuthConfig(client)

	uuid, _ := qb.CreateCluster(req, ks.ProjectID, ks.Token)
	fmt.Println(uuid)
}
