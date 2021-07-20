package command

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/builder/context"
	"github.com/suborbital/subo/builder/template"
	"github.com/suborbital/subo/subo/input"
	"github.com/suborbital/subo/subo/util"
)

type deployData struct {
	SCCVersion       string
	EnvToken         string
	BuilderDomain    string
	StorageClassName string
}

// ComputeDeployCoreCommand returns the compute deploy command
func ComputeDeployCoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "core",
		Short: "deploy the Suborbital Compute Core",
		Long:  `deploy the Suborbital Compute core using Kubernetes or Docker Compose`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to Getwd")
			}

			bctx, err := context.ForDirectory(cwd)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to get CurrentBuildContext")
			}

			util.LogStart("preparing deployment")

			_, err = util.Mkdir(bctx.Cwd, ".suborbital")
			if err != nil {
				return errors.Wrap(err, "🚫 failed to Mkdir")
			}

			branch, _ := cmd.Flags().GetString(branchFlag)

			templatesPath, err := template.UpdateTemplates(defaultRepo, branch)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to UpdateTemplates")
			}

			envToken, err := getEnvToken()
			if err != nil {
				return errors.Wrap(err, "🚫 failed to getEnvToken")
			}

			builderDomain, err := getBuilderDomain()
			if err != nil {
				return errors.Wrap(err, "🚫 failed to getBuilderDomain")
			}

			storageClass, err := getStorageClass()
			if err != nil {
				return errors.Wrap(err, "🚫 failed to getStorageClass")
			}

			data := deployData{
				SCCVersion:       "dev",
				EnvToken:         envToken,
				BuilderDomain:    builderDomain,
				StorageClassName: storageClass,
			}

			if err := template.ExecTmplDir(bctx.Cwd, ".suborbital", templatesPath, "scc-k8s", data); err != nil {
				return errors.Wrap(err, "🚫 failed to ExecTmplDir")
			}

			util.LogDone("ready to start installation")

			dryRun, _ := cmd.Flags().GetBool(dryRunFlag)

			if !dryRun {
				util.LogStart("installing...")

				if _, err := util.Run("kubectl apply -f .suborbital/"); err != nil {
					return errors.Wrap(err, "🚫 failed to kubectl apply")
				}

				util.LogDone("installation complete!")
			} else {
				util.LogInfo("aborting due to dry-run")
			}

			return nil
		},
	}

	cmd.Flags().String(branchFlag, "main", "git branch to download templates from")
	cmd.Flags().Bool(dryRunFlag, false, "prepare the installation in the .suborbital directory, but do not apply it")

	return cmd
}

// getEnvToken gets the environment token from stdin
func getEnvToken() (string, error) {
	fmt.Print("Enter your environment token:")
	token, err := input.ReadStdinString()
	if err != nil {
		return "", errors.Wrap(err, "failed to ReadStdinString")
	}

	if len(token) != 32 {
		return "", errors.New("token must be 32 characters in length")
	}

	return token, nil
}

// getBuilderDomain gets the environment token from stdin
func getBuilderDomain() (string, error) {
	fmt.Print("Enter the domain name that will be used for the builder service:")
	domain, err := input.ReadStdinString()
	if err != nil {
		return "", errors.Wrap(err, "failed to ReadStdinString")
	}

	if len(domain) == 0 {
		return "", errors.New("domain must not be empty")
	}

	return domain, nil
}

// getStorageClass gets the storage class to use
func getStorageClass() (string, error) {
	fmt.Print("Enter the Kubernetes storage class to use:")
	storageClass, err := input.ReadStdinString()
	if err != nil {
		return "", errors.Wrap(err, "failed to ReadStdinString")
	}

	if len(storageClass) == 0 {
		return "", errors.New("storage class must not be empty")
	}

	return storageClass, nil
}
