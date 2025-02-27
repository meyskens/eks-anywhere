//go:build e2e
// +build e2e

package e2e

import (
	"github.com/aws/eks-anywhere/internal/pkg/api"
	"github.com/aws/eks-anywhere/pkg/api/v1alpha1"
	"github.com/aws/eks-anywhere/test/framework"
)

func runSimpleUpgradeFlow(test *framework.ClusterE2ETest, updateVersion v1alpha1.KubernetesVersion, clusterOpts ...framework.ClusterE2ETestOpt) {
	test.GenerateClusterConfig()
	test.CreateCluster()
	test.UpgradeClusterWithNewConfig(clusterOpts)
	test.ValidateCluster(updateVersion)
	test.StopIfFailed()
	test.DeleteCluster()
}

// runSimpleUpgradeFlowWithoutClusterConfigGeneration runs the Create, Upgrade and Delete cluster flows
// that use the cluster config generated by the WithClusterConfig method when the test object is created,
// and avoids regenerating a cluster config with defaults.
func runSimpleUpgradeFlowWithoutClusterConfigGeneration(test *framework.ClusterE2ETest, updateVersion v1alpha1.KubernetesVersion, clusterOpts ...framework.ClusterE2ETestOpt) {
	test.CreateCluster()
	test.UpgradeClusterWithNewConfig(clusterOpts)
	test.ValidateCluster(updateVersion)
	test.StopIfFailed()
	test.DeleteCluster()
}

func runUpgradeFlowWithCheckpoint(test *framework.ClusterE2ETest, updateVersion v1alpha1.KubernetesVersion, clusterOpts, clusterOpts2 []framework.ClusterE2ETestOpt, commandOpts []framework.CommandOpt) {
	test.GenerateClusterConfig()
	test.CreateCluster()
	test.UpgradeClusterWithNewConfig(clusterOpts, commandOpts...)
	test.UpgradeClusterWithNewConfig(clusterOpts2)
	test.ValidateCluster(updateVersion)
	test.StopIfFailed()
	test.DeleteCluster()
}

func runSimpleUpgradeFlowForBareMetal(test *framework.ClusterE2ETest, updateVersion v1alpha1.KubernetesVersion, clusterOpts ...framework.ClusterE2ETestOpt) {
	test.GenerateClusterConfig()
	test.GenerateHardwareConfig()
	test.PowerOffHardware()
	test.CreateCluster(framework.WithControlPlaneWaitTimeout("20m"))
	test.UpgradeClusterWithNewConfig(clusterOpts)
	test.ValidateCluster(updateVersion)
	test.StopIfFailed()
	test.DeleteCluster()
	test.ValidateHardwareDecommissioned()
}

// runSimpleUpgradeFlowForBaremetalWithoutClusterConfigGeneration runs the Create, Upgrade and Delete cluster flows
// for Baremetal that use the cluster config generated by the WithClusterConfig method when the test object is created,
// and avoids regenerating a cluster config with defaults.
func runSimpleUpgradeFlowForBaremetalWithoutClusterConfigGeneration(test *framework.ClusterE2ETest, updateVersion v1alpha1.KubernetesVersion, clusterOpts ...framework.ClusterE2ETestOpt) {
	test.GenerateHardwareConfig()
	test.PowerOffHardware()
	test.CreateCluster(framework.WithControlPlaneWaitTimeout("20m"))
	test.UpgradeClusterWithNewConfig(clusterOpts)
	test.ValidateCluster(updateVersion)
	test.StopIfFailed()
	test.DeleteCluster()
	test.ValidateHardwareDecommissioned()
}

func runUpgradeFlowWithAPI(test *framework.ClusterE2ETest, fillers ...api.ClusterConfigFiller) {
	test.CreateCluster()
	test.LoadClusterConfigGeneratedByCLI()
	test.UpdateClusterConfig(fillers...)
	test.ApplyClusterManifest()
	test.ValidateClusterState()
	test.StopIfFailed()
	test.DeleteCluster()
}

func runUpgradeFlowForBareMetalWithAPI(test *framework.ClusterE2ETest, fillers ...api.ClusterConfigFiller) {
	test.GenerateClusterConfig()
	test.GenerateHardwareConfig()
	test.PowerOffHardware()
	test.CreateCluster(framework.WithControlPlaneWaitTimeout("20m"))
	test.LoadClusterConfigGeneratedByCLI()
	test.UpdateClusterConfig(fillers...)
	test.ApplyClusterManifest()
	test.ValidateClusterState()
	test.StopIfFailed()
	test.DeleteCluster()
	test.ValidateHardwareDecommissioned()
}

func runWorkloadClusterUpgradeFlowAPI(test *framework.MulticlusterE2ETest, filler ...api.ClusterConfigFiller) {
	test.CreateManagementCluster()
	test.RunConcurrentlyInWorkloadClusters(func(wc *framework.WorkloadCluster) {
		wc.ApplyClusterManifest()
		wc.WaitForKubeconfig()
		wc.ValidateClusterState()
		wc.UpdateClusterConfig(filler...)
		wc.ApplyClusterManifest()
		wc.ValidateClusterState()
		wc.DeleteClusterWithKubectl()
		wc.ValidateClusterDelete()
	})
	test.ManagementCluster.StopIfFailed()
	test.DeleteManagementCluster()
}

func runWorkloadClusterUpgradeFlowAPIWithFlux(test *framework.MulticlusterE2ETest, filler ...api.ClusterConfigFiller) {
	test.CreateManagementCluster()
	test.RunConcurrentlyInWorkloadClusters(func(wc *framework.WorkloadCluster) {
		test.PushWorkloadClusterToGit(wc)
		wc.WaitForKubeconfig()
		wc.ValidateClusterState()
		test.PushWorkloadClusterToGit(wc, filler...)
		wc.ValidateClusterState()
		test.DeleteWorkloadClusterFromGit(wc)
		wc.ValidateClusterDelete()
	})
	test.ManagementCluster.StopIfFailed()
	test.DeleteManagementCluster()
}
