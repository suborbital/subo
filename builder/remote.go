package builder

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/suborbital/subo/subo/util"
)

// BuildStatus represents the status of a build
type BuildStatus struct {
	UUID      string `json:"uuid"`
	Status    string `json:"status"`
	HasBundle bool   `json:"hasBundle"`

	Results []BuildResult `json:"results"`
}

// buildStartedResponse is a response to a build started request
type buildStartedResponse struct {
	UUID string `json:"uuid"`
}

func (b *Builder) doRemoteBuild() error {
	b.log.LogInfo("preparing remote build")

	archive, err := archiveForCwd(b.Context.Cwd)
	if err != nil {
		return errors.Wrap(err, "failed to archiveForCwd")
	}

	b.log.LogInfo("starting remote build")

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:8082/api/v1/build/source/%s", b.Context.Directive.Identifier), archive)
	if err != nil {
		return errors.Wrap(err, "failed to NewRequest")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to Do")
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to start build: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	build := &buildStartedResponse{}

	if err := json.NewDecoder(resp.Body).Decode(build); err != nil {
		return errors.Wrap(err, "failed to Decode")
	}

	status := BuildStatus{}

	for {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8082/api/v1/build/source/%s/status", build.UUID), nil)
		if err != nil {
			return errors.Wrap(err, "failed to NewRequest")
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return errors.Wrap(err, "failed to Do")
		}

		if resp.StatusCode != http.StatusOK {
			b.log.LogWarn(fmt.Sprintf("status check failed: %d... will retry", resp.StatusCode))
		}

		defer resp.Body.Close()

		newStatus := &BuildStatus{}

		if err := json.NewDecoder(resp.Body).Decode(newStatus); err != nil {
			return errors.Wrap(err, "failed to Decode")
		}

		if newStatus.Status != status.Status {
			b.log.LogInfo("build has " + newStatus.Status)
			status = *newStatus
		}

		if status.Status == "completed" || status.Status == "failed" {
			break
		}

		time.Sleep(time.Second * 3)
	}

	resultsReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8082/api/v1/build/source/%s/results", build.UUID), nil)
	if err != nil {
		return errors.Wrap(err, "failed to NewRequest")
	}

	resultsResp, err := http.DefaultClient.Do(resultsReq)
	if err != nil {
		return errors.Wrap(err, "failed to Do")
	}

	if resultsResp.StatusCode != http.StatusOK {
		return fmt.Errorf("results request failed: %d", resultsResp.StatusCode)
	}

	defer resultsResp.Body.Close()

	finalStatus := &BuildStatus{}

	if err := json.NewDecoder(resultsResp.Body).Decode(finalStatus); err != nil {
		return errors.Wrap(err, "failed to Decode")
	}

	for _, result := range finalStatus.Results {
		if result.Succeeded {
			b.log.LogInfo("built:\n" + result.OutputLog)
		} else {
			b.log.LogFail("failed:\n" + result.OutputLog)
		}
	}

	if !finalStatus.HasBundle {
		return errors.New("build did not include a bundle")
	}

	b.log.LogInfo("downloading build result")

	bundleReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8082/api/v1/build/source/%s/bundle?shouldDelete=true", build.UUID), nil)
	if err != nil {
		return errors.Wrap(err, "failed to NewRequest")
	}

	bundleResp, err := http.DefaultClient.Do(bundleReq)
	if err != nil {
		return errors.Wrap(err, "failed to Do")
	}

	if bundleResp.StatusCode != http.StatusOK {
		b.log.LogWarn(fmt.Sprintf("status check failed: %d... will retry", bundleResp.StatusCode))
	}

	defer bundleResp.Body.Close()

	bundleFile, err := os.OpenFile("runnables.wasm.zip", os.O_CREATE|os.O_RDWR, util.PermFile)
	if err != nil {
		return errors.Wrap(err, "failed to OpenFile")
	}

	if _, err := io.Copy(bundleFile, bundleResp.Body); err != nil {
		return errors.Wrap(err, "failed to Copy")
	}

	b.log.LogDone("remote build complete")

	return nil
}

func archiveForCwd(cwd string) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})

	writer := zip.NewWriter(buf)

	if walkErr := filepath.WalkDir(cwd, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == cwd {
			return nil
		}

		relPath := strings.TrimPrefix(path, cwd)

		compressedFile, err := writer.Create(relPath)
		if err != nil {
			return errors.Wrap(err, "failed to Create")
		}

		// just creating the file is enough for directories
		if d.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return errors.Wrap(err, "failed to Open")
		}

		defer file.Close()

		if _, err := io.Copy(compressedFile, file); err != nil {
			return errors.Wrap(err, "failed to Copy file to compressedFile")
		}

		return nil
	}); walkErr != nil {
		return nil, errors.Wrap(walkErr, "failed to WalkDir")
	}

	if err := writer.Close(); err != nil {
		return nil, errors.Wrap(err, "failed to writer.Close")
	}

	return buf, nil
}
