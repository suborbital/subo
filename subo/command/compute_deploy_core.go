package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/builder/context"
	"github.com/suborbital/subo/builder/template"
	"github.com/suborbital/subo/subo/input"
	"github.com/suborbital/subo/subo/localproxy"
	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/repl"
	"github.com/suborbital/subo/subo/util"
)

type deployData struct {
	SCCVersion       string
	EnvToken         string
	BuilderDomain    string
	StorageClassName string
}

const proxyDefaultPort int = 80

// ComputeDeployCoreCommand returns the compute deploy command
func ComputeDeployCoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "core",
		Short: "deploy the Suborbital Compute Core",
		Long:  `deploy the Suborbital Compute Core using Kubernetes or Docker Compose`,
		RunE: func(cmd *cobra.Command, args []string) error {
			localInstall := cmd.Flags().Changed(localFlag)

			if !localInstall {
				if err := introAcceptance(); err != nil {
					return err
				}
			}
			proxyPort, _ := cmd.Flags().GetInt(proxyPortFlag)
			if proxyPort < 1 || proxyPort > (2<<16)-1 {
				return errors.New("🚫 proxy-port must be between 1 and 65535")
			}

			cwd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to Getwd")
			}

			bctx, err := context.ForDirectory(cwd)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to get CurrentBuildContext")
			}

			util.LogStart("preparing deployment")

			// if there are any existing deployment manifests sitting around, let's replace them
			if err := removeExistingManifests(bctx); err != nil {
				return errors.Wrap(err, "failed to removeExistingManifests")
			}

			_, err = util.Mkdir(bctx.Cwd, ".suborbital")
			if err != nil {
				return errors.Wrap(err, "🚫 failed to Mkdir")
			}

			branch, _ := cmd.Flags().GetString(branchFlag)
			tag, _ := cmd.Flags().GetString(versionFlag)

			templatesPath, err := template.UpdateTemplates(defaultRepo, branch)
			if err != nil {
				return errors.Wrap(err, "🚫 failed to UpdateTemplates")
			}

			envToken, err := getEnvToken()
			if err != nil {
				return errors.Wrap(err, "🚫 failed to getEnvToken")
			}

			data := deployData{
				SCCVersion: tag,
				EnvToken:   envToken,
			}

			templateName := "scc-docker"

			if !localInstall {
				data.BuilderDomain, err = getBuilderDomain()
				if err != nil {
					return errors.Wrap(err, "🚫 failed to getBuilderDomain")
				}

				data.StorageClassName, err = getStorageClass()
				if err != nil {
					return errors.Wrap(err, "🚫 failed to getStorageClass")
				}

				templateName = "scc-k8s"
			}

			if err := template.ExecTmplDir(bctx.Cwd, "", templatesPath, templateName, data); err != nil {
				return errors.Wrap(err, "🚫 failed to ExecTmplDir")
			}

			util.LogDone("ready to start installation")

			dryRun, _ := cmd.Flags().GetBool(dryRunFlag)

			if dryRun {
				util.LogInfo("aborting due to dry-run, manifest files left in .suborbital")
				return nil
			}

			util.LogStart("installing...")

			if localInstall {
				if _, err := util.Run("docker-compose up -d"); err != nil {
					return errors.Wrap(err, "🚫 failed to docker-compose up")
				}

				util.LogInfo("use `docker ps` and `docker-compose logs` to check deployment status")

				proxyPortStr := strconv.Itoa(proxyPort)
				proxy := localproxy.New("editor.suborbital.network", proxyPortStr)

				go func() {
					if err := proxy.Start(); err != nil {
						log.Fatal(err)
					}
				}()

				// this is to give the proxy server some room to bind to the port and start up
				// it's not ideal, but the least gross way to ensure a good experience
				time.Sleep(time.Second * 2)

				repl := repl.New(proxyPortStr)
				repl.Run()

			} else {
				if _, err := util.Run("kubectl apply -f https://github.com/kedacore/keda/releases/download/v2.4.0/keda-2.4.0.yaml"); err != nil {
					return errors.Wrap(err, "🚫 failed to install KEDA")
				}

				// we don't care if this fails, so don't check error
				util.Run("kubectl create ns suborbital")

				if err := createConfigMap(cwd); err != nil {
					return errors.Wrap(err, "failed to createConfigMap")
				}

				if _, err := util.Run("kubectl apply -f .suborbital/"); err != nil {
					return errors.Wrap(err, "🚫 failed to kubectl apply")
				}

				util.LogInfo("use `kubectl get pods -n suborbital` and `kubectl get svc -n suborbital` to check deployment status")
			}

			util.LogDone("installation complete!")

			return nil
		},
	}

	cmd.Flags().String(branchFlag, "main", "git branch to download templates from")
	cmd.Flags().String(versionFlag, release.SCCTag, "Docker tag to use for control plane images")
	cmd.Flags().Int(proxyPortFlag, proxyDefaultPort, "port that the Editor proxy listens on")
	cmd.Flags().Bool(localFlag, false, "deploy locally using docker-compose")
	cmd.Flags().Bool(dryRunFlag, false, "prepare the installation in the .suborbital directory, but do not apply it")

	return cmd
}

func introAcceptance() error {
	fmt.Print(`
Suborbital Compute Core Installer

BEFORE YOU CONTINUE:
	- You must first run "subo compute create token <email>" to get an environment token

	- You must have kubectl installed in PATH, and it must be connected to the cluster you'd like to use

	- You must be able to set up DNS records for the builder service after this installation completes
			- Choose the DNS name you'd like to use before continuing, e.g. builder.acmeco.com

	- Subo will attempt to determine the default storage class for your Kubernetes cluster, 
	  but if is unable to do so you will need to provide one
			- See the Compute documentation for more details

	- Subo will install the KEDA autoscaler into your cluster. It will not affect any existing deployments.

Are you ready to continue? (y/N): `)

	answer, err := input.ReadStdinString()
	if err != nil {
		return errors.Wrap(err, "failed to ReadStdinString")
	}

	if !strings.EqualFold(answer, "y") {
		return errors.New("aborting")
	}

	return nil
}

// getEnvToken gets the environment token from stdin
func getEnvToken() (string, error) {
	buf, err := util.ReadEnvironmentToken()
	if err == nil {
		return buf, nil
	}

	fmt.Print("Enter your environment token: ")
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
	fmt.Print("Enter the domain name that will be used for the builder service: ")
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
	defaultClass, err := detectStorageClass()
	if err != nil {
		// that's fine, continue
		fmt.Println("failed to automatically detect Kubernetes storage class:", err.Error())
	} else if defaultClass != "" {
		fmt.Println("using default storage class: ", defaultClass)
		return defaultClass, nil
	}

	fmt.Print("Enter the Kubernetes storage class to use: ")
	storageClass, err := input.ReadStdinString()
	if err != nil {
		return "", errors.Wrap(err, "failed to ReadStdinString")
	}

	if len(storageClass) == 0 {
		return "", errors.New("storage class must not be empty")
	}

	return storageClass, nil
}

func detectStorageClass() (string, error) {
	output, err := util.Run("kubectl get storageclass --output=name")
	if err != nil {
		return "", errors.Wrap(err, "failed to get default storageclass")
	}

	// output will look like: storageclass.storage.k8s.io/do-block-storage
	// so split on the / and return the last part

	outputParts := strings.Split(output, "/")
	if len(outputParts) != 2 {
		return "", errors.New("could not automatically determine storage class")
	}

	return outputParts[1], nil
}

func createConfigMap(cwd string) error {
	configFilepath := filepath.Join(cwd, "config", "scc-config.yaml")

	_, err := os.Stat(configFilepath)
	if err != nil {
		return errors.Wrap(err, "failed to Stat scc-config.yaml")
	}

	if _, err := util.Run(fmt.Sprintf("kubectl create configmap scc-config --from-file=scc-config.yaml=%s -n suborbital", configFilepath)); err != nil {
		return errors.Wrap(err, "failed to create configmap (you may need to run `kubectl delete configmap scc-config -n suborbital`)")
	}

	return nil
}

func removeExistingManifests(bctx *context.BuildContext) error {
	// start with a clean slate
	if _, err := os.Stat(filepath.Join(bctx.Cwd, ".suborbital")); err == nil {
		if err := os.RemoveAll(filepath.Join(bctx.Cwd, ".suborbital")); err != nil {
			return errors.Wrap(err, "failed to RemoveAll .suborbital")
		}
	}

	if _, err := os.Stat(filepath.Join(bctx.Cwd, "docker-compose.yml")); err == nil {
		if err := os.Remove(filepath.Join(bctx.Cwd, "docker-compose.yml")); err != nil {
			return errors.Wrap(err, "failed to Remove docker-compose.yml")
		}
	}

	return nil
}
