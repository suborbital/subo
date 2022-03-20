package deployer

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/pkg/errors"

	"github.com/suborbital/atmo/atmo/appsource"
	"github.com/suborbital/atmo/atmo/options"
	"github.com/suborbital/atmo/directive"
	subotemplate "github.com/suborbital/subo/builder/template"
	"github.com/suborbital/subo/project"
	"github.com/suborbital/subo/subo/util"
)

const (
	cloudflareDeployJobType = "cloudflare"
	routerScriptName        = "subo-router"

	bundlePath = "./runnables.wasm.zip"
)

//go:embed cloudflare/shared.js
var shared string

//go:embed cloudflare/itty-router.js
var ittyRouter string

type CloudflareDeployJob struct {
	domain          string
	repo            string
	branch          string
	updateTemplates bool
}

func NewCloudflareDeployer(repo, branch, domain string, updateTemplates bool) DeployJob {
	c := &CloudflareDeployJob{
		repo:            repo,
		branch:          branch,
		domain:          domain,
		updateTemplates: updateTemplates,
	}
	return c
}

func (k *CloudflareDeployJob) Type() string {
	return cloudflareDeployJobType
}

type cfWorker struct {
	Path           string
	Method         string
	ServiceBinding string
}

func (job *CloudflareDeployJob) Deploy(log util.FriendlyLogger, ctx *project.Context) error {
	atmoOpts := options.NewWithModifiers()
	appSource := appsource.NewBundleSource(bundlePath)

	if err := appSource.Start(*atmoOpts); err != nil {
		return errors.Wrap(err, "failed to start")
	}

	api, err := getAPIClient()
	if err != nil {
		return errors.Wrap(err, "failed to initalize Cloudflare api client")
	}

	if err := os.RemoveAll(filepath.Join(ctx.Cwd, ".deployment")); err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "failed to RemoveAll deployment files")
		}
	}

	if err := os.MkdirAll(filepath.Join(ctx.Cwd, ".deployment"), util.PermDirectory); err != nil {
		return errors.Wrap(err, "failed to MkdirAll .deployment")
	}

	templatesPath, err := subotemplate.FullPath(job.repo, job.branch)
	if err != nil {
		return errors.Wrap(err, "failed to template.FullPath")
	}

	if job.updateTemplates {
		templatesPath, err = subotemplate.UpdateTemplates(job.repo, job.branch)
		if err != nil {
			return errors.Wrap(err, "🚫 failed to UpdateTemplates")
		}
	}

	for _, meta := range appSource.Applications() {
		log.LogStart(fmt.Sprintf("deploying %s (%s) runnables as workers", meta.Identifier, job.domain))

		runnablesByName := getRunnablesByName(appSource, meta)

		workersToRoute := make(map[string]cfWorker)

		for _, handler := range appSource.Handlers(meta.Identifier, meta.AppVersion) {
			if handler.Input.Type != "request" {
				log.LogWarn(fmt.Sprintf("unsupported handler type: %s; ignoring.", handler.Input.Type))
				continue
			}

			if len(handler.Steps) != 1 {
				log.LogWarn(fmt.Sprintf("only one step is currently supported per handler"))
				continue
			}

			for _, step := range handler.Steps {
				runnable, ok := runnablesByName[step.CallableFn.Fn]
				if !ok {
					return errors.Errorf("unexpected runnable %s", step.CallableFn.Fn)
				}

				script, err := renderRequestWorker(templatesPath, handler.Input.Method, handler.Input.Resource)
				if err != nil {
					return errors.Wrap(err, "failed to render request worker")
				}
				scriptName := fmt.Sprintf("%s-request-%s", strings.ReplaceAll(meta.Identifier, ".", "-"), runnable.Name)
				if err := deployRunnable(api, &runnable, scriptName, script, job.domain); err != nil {
					return errors.Wrap(err, "failed to deploy runnable")
				}

				path := handler.Input.Resource
				workersToRoute[scriptName] = cfWorker{
					Path:           path,
					Method:         handler.Input.Method,
					ServiceBinding: toWorkerBinding(scriptName),
				}

				log.LogInfo(fmt.Sprintf("deployed runnable %s", runnable.Name))
			}
		}

		if err := deployRouter(templatesPath, api, workersToRoute, job.domain); err != nil {
			return errors.Wrap(err, "failed to create router")
		}
		routerRoute := fmt.Sprintf("%s/*", job.domain)
		if err := createRoute(api, routerRoute, routerScriptName, job.domain); err != nil {
			log.LogWarn(fmt.Sprintf("error while creating route: %s", err))
		}

		log.LogDone(fmt.Sprintf("deployed to Cloudflare at %s", routerRoute))
	}

	return nil
}

func getAPIClient() (*cloudflare.API, error) {
	apiToken := os.Getenv("CF_API_TOKEN")
	accountId := os.Getenv("CF_ACCOUNT_ID")

	if apiToken == "" || accountId == "" {
		return nil, errors.New("missing CF_API_TOKEN or CF_ACCOUNT_ID")
	}

	api, err := cloudflare.NewWithAPIToken(apiToken, cloudflare.UsingAccount(accountId))
	if err != nil {
		return nil, err
	}

	return api, nil
}

func getRunnablesByName(appSource appsource.AppSource, meta appsource.Meta) map[string]directive.Runnable {
	runnablesByName := make(map[string]directive.Runnable)
	runnables := appSource.Runnables(meta.Identifier, meta.AppVersion)

	for _, runnable := range runnables {
		runnablesByName[runnable.Name] = runnable
	}

	return runnablesByName
}

func renderRouterWorker(templatesPath string, workersToRoute map[string]cfWorker) (string, error) {
	templatePath := path.Join(templatesPath, "cf/router_worker.js.tmpl")
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse template")
	}

	type routerWorkerData struct {
		IttyRouter     string
		WorkersToRoute map[string]cfWorker
	}

	data := routerWorkerData{
		IttyRouter:     ittyRouter,
		WorkersToRoute: workersToRoute,
	}

	var res bytes.Buffer

	if err := t.Execute(&res, data); err != nil {
		return "", errors.Wrap(err, "failed to render template")
	}

	return res.String(), nil
}

func renderRequestWorker(templatesPath, method, urlPath string) (string, error) {
	templatePath := path.Join(templatesPath, "cf/request_worker.js.tmpl")
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse template")
	}

	type requestWorkerData struct {
		Method     string
		Path       string
		IttyRouter string
		Shared     string
	}

	data := requestWorkerData{
		Method:     strings.ToLower(method),
		Path:       urlPath,
		IttyRouter: ittyRouter,
		Shared:     shared,
	}

	var res bytes.Buffer

	if err := t.Execute(&res, data); err != nil {
		return "", errors.Wrap(err, "failed to render template")
	}

	return res.String(), nil
}

func deployRunnable(api *cloudflare.API, runnable *directive.Runnable, scriptName string, script string, zoneName string) error {
	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return errors.Wrap(err, "failed to get zone id")
	}

	data, err := runnable.ModuleRef.Bytes()
	if err != nil {
		return errors.Wrap(err, "failed to get module ref")
	}
	wasmReader := bytes.NewReader(data)

	scriptParams := cloudflare.WorkerScriptParams{
		Script: script,
		Bindings: map[string]cloudflare.WorkerBinding{
			"WASM_MODULE": cloudflare.WorkerWebAssemblyBinding{Module: wasmReader},
		},
	}

	_, err = api.UploadWorkerWithBindings(context.Background(), &cloudflare.WorkerRequestParams{ZoneID: zoneID, ScriptName: scriptName}, &scriptParams)
	if err != nil {
		return errors.Wrap(err, "failed to upload runnable worker")
	}

	return nil
}

func toWorkerBinding(name string) string {
	return strings.ReplaceAll(strings.ToUpper(name), "-", "_")
}

func deployRouter(templatesPath string, api *cloudflare.API, workersToRoute map[string]cfWorker, zoneName string) error {
	script, err := renderRouterWorker(templatesPath, workersToRoute)
	if err != nil {
		return errors.Wrap(err, "failed to render script")
	}
	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return errors.Wrap(err, "failed to get zone id")
	}

	bindings := make(map[string]cloudflare.WorkerBinding, len(workersToRoute))
	for name, route := range workersToRoute {
		bindings[route.ServiceBinding] = cloudflare.WorkerServiceBinding{Service: name, Environment: "production"}
	}

	scriptParams := cloudflare.WorkerScriptParams{
		Script:   script,
		Bindings: bindings,
	}

	_, err = api.UploadWorkerWithBindings(context.Background(), &cloudflare.WorkerRequestParams{ZoneID: zoneID, ScriptName: routerScriptName}, &scriptParams)
	if err != nil {
		return errors.Wrap(err, "failed to upload router worker")
	}

	return nil
}

func createRoute(api *cloudflare.API, pattern string, scriptName string, zoneName string) error {
	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return errors.Wrap(err, "failed to get domain's zone tag")
	}

	_, err = api.CreateWorkerRoute(context.Background(), zoneID, cloudflare.WorkerRoute{Pattern: pattern, Script: scriptName})
	if err != nil {
		return errors.Wrap(err, "failed to create worker route")
	}

	return nil
}
