package cmd

import (
	"fmt"
	"path"
	"strings"

	"github.com/bcrusu/kcm/util"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func newKubectlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "kubectl ARGS",
		Aliases:            []string{"ctl", "kctl"},
		Short:              "Runs kubectl for the current cluster",
		SilenceUsage:       true,
		DisableFlagParsing: true,
	}

	cmd.RunE = kubectlCmdRunE
	return cmd
}

func kubectlCmdRunE(cmd *cobra.Command, args []string) error {
	if len(args) == 1 && args[0] == "--help" {
		return cmd.Help()
	}

	for _, arg := range args {
		if strings.HasPrefix(arg, "--kubeconfig") {
			fmt.Println("Do not set 'kubeconfig' option. It will be set automatically by kcm")
			return nil
		}
	}

	clusterRepository, err := newClusterRepository()
	if err != nil {
		return err
	}

	cluster, err := getWorkingCluster(clusterRepository, "")
	if err != nil {
		return err
	}

	kubectlPath := path.Join(kubernetesBinDir(cluster.KubernetesVersion), "kubectl")
	kubeconfigPath := path.Join(*dataDir, "config", cluster.Name, "kubeconfig", "kubectl")
	kubectlArgs := append(args, fmt.Sprintf(`--kubeconfig=%s`, kubeconfigPath))

	if err := util.ExecCommandAndWait(kubectlPath, kubectlArgs...); err != nil {
		glog.Warningf("failed to execute kubectl. Error: %v", err)
	}

	return nil
}
