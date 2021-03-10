package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suborbital/subo/subo/util"
	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v2"
)

const (
	preReleaseFlag = "prerelease"
)

// DotSuboFile describes a .subo file for controlling releases
type DotSuboFile struct {
	DotVersionFiles []string `yaml:"dotVersionFiles"`
	PreMakeTargets  []string `yaml:"preMakeTargets"`
	PostMakeTargets []string `yaml:"postMakeTargets"`
}

// CreateReleaseCmd returns the create release command
// this is only available for development builds
func CreateReleaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release <version> <title>",
		Short: "create a new release",
		Long:  `tag a new version and create a new GitHub release, configured using the .subo.yml file.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			logStart("checking release conditions")
			cwd, _ := cmd.Flags().GetString("dir")

			newVersion := args[0]
			releaseName := args[1]

			if err := validateVersion(newVersion); err != nil {
				return errors.Wrap(err, "failed to validateVersion")
			}

			if err := checkGitCleanliness(); err != nil {
				return errors.Wrap(err, "failed to checkGitCleanliness")
			}

			dotSubo, err := findDotSubo(cwd)
			if err != nil {
				return errors.Wrap(err, "failed to findDotSubo")
			} else if dotSubo == nil {
				return errors.New(".subo.yml file is missing")
			}

			changelogFilePath := filepath.Join(cwd, "changelogs", fmt.Sprintf("%s.md", newVersion))

			if err := checkChangelogFileExists(changelogFilePath); err != nil {
				return errors.Wrap(err, "failed to checkChangelogFileExists")
			}

			for _, f := range dotSubo.DotVersionFiles {
				filePath := filepath.Join(cwd, f)

				if err := util.CheckFileForVersionString(filePath, newVersion); err != nil {
					if errors.Is(err, util.ErrVersionNotPresent) {
						return fmt.Errorf("required dotVersionFile %s does not contain the release version number %s", filePath, newVersion)
					}

					return errors.Wrap(err, "failed to CheckFileForVersionString")
				}
			}

			logDone("release is ready to go")
			logStart("running pre-make targets")

			for _, target := range dotSubo.PreMakeTargets {
				if _, _, err := util.Run(fmt.Sprintf("make %s", target)); err != nil {
					return errors.Wrapf(err, "failed to run preMakeTarget %s", target)
				}
			}

			logDone("pre-make targets complete")
			logStart("creating release")

			if _, _, err := util.Run("git push"); err != nil {
				return errors.Wrap(err, "failed to Run git push")
			}

			ghCommand := fmt.Sprintf("gh release create %s --title=%s --notes-file=%s", newVersion, releaseName, changelogFilePath)
			if preRelease, _ := cmd.Flags().GetBool(preReleaseFlag); preRelease {
				ghCommand += " --prerelease"
			}

			if _, _, err := util.Run(ghCommand); err != nil {
				return errors.Wrap(err, "failed to Run gh command")
			}

			if _, _, err := util.Run("git pull --tags"); err != nil {
				return errors.Wrap(err, "failed to Run git pull command")
			}

			logDone("release created!")
			logStart("running post-make targets")

			for _, target := range dotSubo.PostMakeTargets {
				if _, _, err := util.Run(fmt.Sprintf("make %s", target)); err != nil {
					return errors.Wrapf(err, "failed to run postMakeTarget %s", target)
				}
			}

			logDone("post-make targets complete")

			return nil
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "$HOME"
	}

	cmd.Flags().String(dirFlag, cwd, "the directory to create the release for")
	cmd.Flags().Bool(preReleaseFlag, false, "pass --prelease to mark the release as such")

	return cmd
}

func findDotSubo(cwd string) (*DotSuboFile, error) {
	dotSuboPath := filepath.Join(cwd, ".subo.yml")

	dotSuboBytes, err := ioutil.ReadFile(dotSuboPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to ReadFile")
	}

	dotSubo := &DotSuboFile{}
	if err := yaml.Unmarshal(dotSuboBytes, dotSubo); err != nil {
		return nil, errors.Wrap(err, "failed to Unmarshal dotSubo file")
	}

	return dotSubo, nil
}

func checkChangelogFileExists(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		return errors.Wrap(err, "failed to Stat changelog file")
	}

	return nil
}

func checkGitCleanliness() error {
	if out, _, err := util.Run("git diff-index --name-only HEAD"); err != nil {
		return errors.Wrap(err, "failed to git diff-index")
	} else if out != "" {
		return errors.New("project has modified files")
	}

	if out, _, err := util.Run("git ls-files --exclude-standard --others"); err != nil {
		return errors.Wrap(err, "failed to git ls-files")
	} else if out != "" {
		return errors.New("project has untracked files")
	}

	return nil
}

func validateVersion(version string) error {
	if !strings.HasPrefix(version, "v") {
		return errors.New("version does not start with v")
	}

	if !semver.IsValid(version) {
		return errors.New("version is not valid semver")
	}

	return nil
}
