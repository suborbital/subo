package packager

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/suborbital/atmo/bundle"
	"github.com/suborbital/atmo/directive"
	"github.com/suborbital/subo/project"
	"github.com/suborbital/subo/subo/release"
	"github.com/suborbital/subo/subo/util"
	"golang.org/x/mod/semver"
)

const bundleJobType = "bundle"

type BundlePackageJob struct{}

func NewBundlePackageJob() PackageJob {
	b := &BundlePackageJob{}

	return b
}

// Type returns the job type
func (b *BundlePackageJob) Type() string {
	return bundleJobType
}

// Package packages the application
func (b *BundlePackageJob) Package(log util.FriendlyLogger, ctx *project.Context) error {
	for _, r := range ctx.Runnables {
		if err := r.HasModule(); err != nil {
			return errors.Wrap(err, "missing built module")
		}
	}

	if ctx.Directive == nil {
		ctx.Directive = &directive.Directive{
			Identifier: "com.suborbital.app",
			// TODO: insert some git smarts here?
			AppVersion:  "v0.0.1",
			AtmoVersion: fmt.Sprintf("v%s", release.AtmoVersion),
		}
	} else if ctx.Directive.Headless {
		log.LogInfo("updating Directive")

		// Bump the appVersion since we're in headless mode.
		majorStr := strings.TrimPrefix(semver.Major(ctx.Directive.AppVersion), "v")
		major, _ := strconv.Atoi(majorStr)
		new := fmt.Sprintf("v%d.0.0", major+1)

		ctx.Directive.AppVersion = new

		if err := project.WriteDirectiveFile(ctx.Cwd, ctx.Directive); err != nil {
			return errors.Wrap(err, "failed to WriteDirectiveFile")
		}
	}

	if err := project.AugmentAndValidateDirectiveFns(ctx.Directive, ctx.Runnables); err != nil {
		return errors.Wrap(err, "🚫 failed to AugmentAndValidateDirectiveFns")
	}

	if err := ctx.Directive.Validate(); err != nil {
		return errors.Wrap(err, "🚫 failed to Validate Directive")
	}

	static, err := CollectStaticFiles(ctx.Cwd)
	if err != nil {
		return errors.Wrap(err, "failed to CollectStaticFiles")
	}

	if static != nil {
		log.LogInfo("adding static files to bundle")
	}

	directiveBytes, err := ctx.Directive.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to Directive.Marshal")
	}

	modules, err := ctx.Modules()
	if err != nil {
		return errors.Wrap(err, "failed to Modules for build")
	}

	if err := bundle.Write(directiveBytes, modules, static, ctx.Bundle.Fullpath); err != nil {
		return errors.Wrap(err, "🚫 failed to WriteBundle")
	}

	bundleRef := project.BundleRef{
		Exists:   true,
		Fullpath: filepath.Join(ctx.Cwd, "runnables.wasm.zip"),
	}

	ctx.Bundle = bundleRef

	log.LogDone(fmt.Sprintf("bundle was created -> %s @ %s", ctx.Bundle.Fullpath, ctx.Directive.AppVersion))

	return nil
}
