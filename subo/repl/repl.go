package repl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/suborbital/atmo/fqfn"
	"github.com/suborbital/subo/subo/input"
)

// Repl is a 'local proxy repl' that allows the user to perform simple actions against their local install of Compute
type Repl struct {
	proxyPort string
}

type tokenResp struct {
	Token string `json:"token"`
}

// New creates a new "local proxy repl"
func New(proxyPort string) *Repl {
	return &Repl{proxyPort: proxyPort}
}

func (r *Repl) Run() error {
	fmt.Print("\n\nPress enter to launch the local Compute REPL...")

	if _, err := input.ReadStdinString(); err != nil {
		return errors.Wrap(err, "failed to ReadStdinString")
	}

	for {
		fmt.Println("\n\n1. Create or edit a function")
		fmt.Print("\nChoose an option: ")

		opt, err := input.ReadStdinString()
		if err != nil {
			return errors.Wrap(err, "failed to ReadStdinString")
		}

		var replErr error

		switch opt {
		case "1":
			replErr = r.editFunction()
		default:
			fmt.Println("invalid, choose again.")
		}

		if replErr != nil {
			return errors.Wrap(err, "error produced by option "+opt)
		}
	}
}

func (r *Repl) editFunction() error {
	fmt.Print("\n\nTo create or edit a function, enter its name (or FQFN): ")
	name, err := input.ReadStdinString()
	if err != nil {
		return errors.Wrap(err, "failed to ReadStdinString")
	}

	ident := "com.suborbital.acmeco"
	namespace := "default"

	FQFN := fqfn.Parse(name)
	if FQFN.Identifier != "" {
		ident = FQFN.Identifier
	}

	if FQFN.Namespace != "" {
		namespace = FQFN.Namespace
	}

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://local.suborbital.network:8081/api/v1/token/%s/%s/%s", ident, namespace, FQFN.Fn), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to Do request")
	}

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	token := tokenResp{}
	json.Unmarshal(body, &token)

	editorHost := "local.suborbital.network"
	if r.proxyPort != "80" {
		editorHost += ":" + r.proxyPort
	}

	editorURL := fmt.Sprintf("http://%s/?builder=http://local.suborbital.network:8082&token=%s&ident=%s&namespace=%s&fn=%s", editorHost, token.Token, ident, namespace, FQFN.Fn)

	fmt.Println("\n✅ visit", editorURL, "to access the editor")

	return nil
}
