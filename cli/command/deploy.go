package command

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/suborbital/velo/cli/util"
	"github.com/suborbital/velo/deployer"
	"github.com/suborbital/velo/project"
)

var validDeployTypes = map[string]bool{
	"kubernetes": true,
	"k8s":        true,
	"fly":        true,
	"fly.io":     true,
}

//DeployCmd deploys the current project.
func DeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy an application",
		Long:  "deploy the current project to a remote environment (Kubernetes, etc.)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			deployType := args[0]
			if _, valid := validDeployTypes[deployType]; !valid {
				return fmt.Errorf("invalid deployment type %s", deployType)
			}

			cwd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to Getwd")
			}

			ctx, err := project.ForDirectory(cwd)
			if err != nil {
				return errors.Wrap(err, "failed to project.ForDirectory")
			}

			dplyr := deployer.New(&util.PrintLogger{})
			var deployJob deployer.DeployJob

			repo, _ := cmd.Flags().GetString(repoFlag)
			branch, _ := cmd.Flags().GetString(branchFlag)
			domain, _ := cmd.Flags().GetString(domainFlag)
			updateTemplates := cmd.Flags().Changed(updateTemplatesFlag)
			org, _ := cmd.Flags().GetString(orgFlag)
			region, _ := cmd.Flags().GetString(regionFlag)
			local := cmd.Flags().Changed(localFlag)

			switch deployType {
			case "kubernetes", "k8s":
				deployJob = deployer.NewK8sDeployJob(repo, branch, domain, updateTemplates)
			case "fly", "fly.io":
				deployJob = deployer.NewFlyDeployJob(org, region, local)
			}

			if err := dplyr.Deploy(ctx, deployJob); err != nil {
				return errors.Wrap(err, "failed to Deploy")
			}

			return nil
		},
	}

	cmd.Flags().String(domainFlag, "", "domain name to configure TLS for (DNS must be configured post-deploy)")
	cmd.Flags().String(repoFlag, "suborbital/runnable-templates", "git repo to download templates from")
	cmd.Flags().String(branchFlag, "vmain", "git branch to download templates from")
	cmd.Flags().String(orgFlag, "personal", "the organisation in which to deploy (for fly.io)")
	cmd.Flags().String(regionFlag, "lax", "the region in which to deploy (for fly.io)")
	cmd.Flags().Bool(updateTemplatesFlag, false, "update with the newest runnable templates")
	cmd.Flags().Bool(localFlag, false, "use the local Docker installation for container builds (for fly.io)")

	return cmd
}
