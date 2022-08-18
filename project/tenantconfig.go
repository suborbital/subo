package project

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/suborbital/appspec/tenant"
	"github.com/suborbital/subo/subo/util"
)

// WriteTenantConfig writes a tenant config to disk.
func WriteTenantConfig(cwd string, cfg *tenant.Config) error {
	filePath := filepath.Join(cwd, "tenant.json")

	configBytes, err := cfg.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to Marshal")
	}

	if err := ioutil.WriteFile(filePath, configBytes, util.PermFilePrivate); err != nil {
		return errors.Wrap(err, "failed to WriteFile")
	}

	return nil
}

// readTenantConfig finds a tenant.json from disk but does not validate it.
func readTenantConfig(cwd string) (*tenant.Config, error) {
	filePath := filepath.Join(cwd, "tenant.json")

	tenantBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ReadFile for Directive")
	}

	tenant := &tenant.Config{}
	if err := tenant.Unmarshal(tenantBytes); err != nil {
		return nil, errors.Wrap(err, "failed to Unmarshal Directive")
	}

	return tenant, nil
}

// readQueriesFile finds a queries.yaml from disk.
func readQueriesFile(cwd string) ([]tenant.DBQuery, error) {
	filePath := filepath.Join(cwd, "Queries.yaml")

	configBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to ReadFile for Queries.yaml")
	}

	tenant := &tenant.Config{}
	if err := tenant.Unmarshal(configBytes); err != nil {
		return nil, errors.Wrap(err, "failed to Unmarshal Directive")
	}

	return tenant.DefaultNamespace.Queries, nil
}

// AugmentAndValidateModules ensures that all modules referenced in a workflow exist
// in the project and then adds the module list to the provided config.
func AugmentAndValidateModules(cfg *tenant.Config, mods []ModuleDir) error {
	modMap := map[string]bool{}
	for _, fn := range mods {
		modMap[fn.Name] = true
	}

	workflowMods := getWorkflowModList(cfg)

	for _, fn := range workflowMods {
		if good, exists := modMap[fn]; !good || !exists {
			return fmt.Errorf("project does not contain function %s listed in Directive", fn)
		}
	}

	dirModules := make([]tenant.Module, len(mods))

	// for each module, calculate its ref (a.k.a. its hash), and then add it to the context
	for i := range mods {
		mod := mods[i]
		modFile, err := mod.WasmFile()
		if err != nil {
			return errors.Wrap(err, "failed to WasmFile")
		}

		defer modFile.Close()

		hash, err := calculateModuleRef(modFile)
		if err != nil {
			return errors.Wrap(err, "failed to calculateModuleRef")
		}

		mod.Module.Ref = hash
		rev := tenant.ModuleRevision{
			Ref: hash,
		}

		if mod.Module.Revisions == nil {
			mod.Module.Revisions = []tenant.ModuleRevision{rev}
		} else {
			mod.Module.Revisions = append(mod.Module.Revisions, rev)
		}

		dirModules[i] = *mod.Module
	}

	cfg.Modules = dirModules

	return cfg.Validate()
}

// calculateModuleRef calculates the base64url-encoded sha256 hash of a module file
func calculateModuleRef(mod *os.File) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, mod); err != nil {
		return "", errors.Wrap(err, "failed to Copy module contents")
	}

	hashBytes := hasher.Sum(nil)

	hashString := base64.URLEncoding.EncodeToString(hashBytes)

	return hashString, nil
}

// getWorkflowModList gets a full list of all functions used in the config's workflows.
func getWorkflowModList(cfg *tenant.Config) []string {
	modMap := map[string]bool{}

	// collect all the workflows in all of the namespaces
	workflows := []tenant.Workflow{}
	workflows = append(workflows, cfg.DefaultNamespace.Workflows...)
	for _, ns := range cfg.Namespaces {
		workflows = append(workflows, ns.Workflows...)
	}

	for _, h := range workflows {
		for _, step := range h.Steps {
			if step.IsFn() {
				modMap[step.ExecutableMod.FQMN] = true
			} else if step.IsGroup() {
				for _, mod := range step.Group {
					modMap[mod.FQMN] = true
				}
			}
		}
	}

	mods := []string{}
	for fn := range modMap {
		mods = append(mods, fn)
	}

	return mods
}

func DockerNameFromConfig(cfg *tenant.Config) (string, error) {
	identParts := strings.Split(cfg.Identifier, ".")
	if len(identParts) != 3 {
		return "", errors.New("ident has incorrect number of parts")
	}

	org := identParts[1]
	repo := identParts[2]

	name := fmt.Sprintf("%s/%s:%d", org, repo, cfg.TenantVersion)

	return name, nil
}
