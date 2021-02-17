// Copyright © 2020 The pf9ctl authors

package cmd

import (
	"github.com/platform9/pf9ctl/pkg/pmk"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var checkNodeCmd = &cobra.Command{
	Use:   "check-node",
	Short: "Check prerequisites for k8s",
	Long: `Check if a node satisfies prerequisites to be ready to be added to a Kubernetes cluster. Read more
	at https://platform9.com/blog/support/managed-container-cloud-requirements-checklist/`,
	Run: checkNodeRun,
}

func init() {
	checkNodeCmd.Flags().StringVarP(&user, "user", "u", "", "ssh username for the nodes")
	checkNodeCmd.Flags().StringVarP(&password, "password", "p", "", "ssh password for the nodes")
	checkNodeCmd.Flags().StringVarP(&sshKey, "ssh-key", "s", "", "ssh key file for connecting to the nodes")
	checkNodeCmd.Flags().StringSliceVarP(&ips, "ips", "i", []string{}, "ips of host to be prepared")
	checkNodeCmd.Flags().BoolVarP(&floatingIP, "floating-ip", "f", false, "")

	rootCmd.AddCommand(checkNodeCmd)
}

func checkNodeRun(cmd *cobra.Command, args []string) {
	zap.S().Debug("==========Running check-node==========")
	ctx, err := pmk.LoadConfig(Pf9DBLoc)
	if err != nil {
		zap.S().Fatalf("Unable to load the context: %s\n", err.Error())
	}

	executor, err := getExecutor()
	if err != nil {
		zap.S().Fatalf("Error connecting to host %s", err.Error())
	}
	c, err := pmk.NewClient(ctx.Fqdn, executor, ctx.AllowInsecure, false)
	if err != nil {
		zap.S().Fatalf("Unable to load clients needed for the Cmd. Error: %s", err.Error())
	}

	result := pmk.CheckNode(c)

	if !result {
		zap.S().Errorf("Node not ready. See %s or use --verbose for logs", Pf9Log)
	}
	zap.S().Debug("==========Finished running check-node==========")
}