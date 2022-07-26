package publisher

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/deislabs/go-bindle/client"
	"github.com/deislabs/go-bindle/keyring"
	"github.com/deislabs/go-bindle/types"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/suborbital/velo/cli/util"
	"github.com/suborbital/velo/project"
)

const (
	BindlePublishJobType = "bindle"
	suboAuthor           = "Subo <subo@suborbital.dev>"
)

type BindlePublishJob struct{}

type parcelWrapper struct {
	parcel types.Parcel
	data   []byte
}

// NewBindlePublishJob returns a new PublishJob for Bindle.
func NewBindlePublishJob() PublishJob {
	b := &BindlePublishJob{}

	return b
}

// Type returns the publish job's type.
func (b *BindlePublishJob) Type() string {
	return BindlePublishJobType
}

// Publish publishes the application.
func (b *BindlePublishJob) Publish(log util.FriendlyLogger, ctx *project.Context) error {
	if ctx.Directive == nil {
		return errors.New("🚫 cannot push without Directive.yaml file")
	}

	log.LogStart(fmt.Sprintf("pushing %s@%s", ctx.Directive.Identifier, ctx.Directive.AppVersion))

	invoice := &types.Invoice{
		BindleVersion: "1.0.0",
		Bindle: types.BindleSpec{
			Name:    ctx.Directive.Identifier,
			Version: strings.TrimPrefix(ctx.Directive.AppVersion, "v"),
			Authors: []string{
				suboAuthor,
			},
		},
		Parcel: []types.Parcel{},
	}

	parcelsBySHA := map[string]parcelWrapper{}

	// add the Directive as a parcel.
	directiveBytes, err := yaml.Marshal(ctx.Directive)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal Directive")
	}

	directiveParcel := parcelForData("Directive.yaml", "application/yaml", directiveBytes)

	invoice.Parcel = append(invoice.Parcel, directiveParcel)

	parcelsBySHA[directiveParcel.Label.SHA256] = parcelWrapper{
		parcel: directiveParcel,
		data:   directiveBytes,
	}

	// add each Runnable as a parcel.
	for _, r := range ctx.Runnables {
		files, err := ioutil.ReadDir(r.Fullpath)
		if err != nil {
			return errors.Wrapf(err, "failed to ReadDir for %s", r.Fullpath)
		}

		for _, file := range files {
			if !strings.HasSuffix(file.Name(), ".wasm") {
				continue
			}

			fullPath := filepath.Join(r.Fullpath, file.Name())

			fileBytes, err := os.ReadFile(fullPath)
			if err != nil {
				return errors.Wrapf(err, "failed to Open %s", fullPath)
			}

			parcel := parcelForData(file.Name(), "application/wasm", fileBytes)

			invoice.Parcel = append(invoice.Parcel, parcel)

			parcelsBySHA[parcel.Label.SHA256] = parcelWrapper{
				parcel: parcel,
				data:   fileBytes,
			}
		}
	}

	sigKey, privKey, err := createOrReadKeypair(suboAuthor)
	if err != nil {
		return errors.Wrap(err, "failed to createOrReadKeypair")
	}

	if err := invoice.GenerateSignature(suboAuthor, types.RoleCreator, sigKey, privKey); err != nil {
		return errors.Wrap(err, "failed to GenerateCreatorSignaure")
	}

	publishClient, err := client.New("http://127.0.0.1:8080/v1", nil)
	if err != nil {
		return errors.Wrap(err, "failed to publishClient.New")
	}

	invResp, err := publishClient.CreateInvoice(*invoice)
	if err != nil {
		return errors.Wrap(err, "failed to CreateInvoice")
	}

	for _, p := range invResp.Missing {
		wrapper := parcelsBySHA[p.SHA256]

		if err := publishClient.CreateParcel(invoice.Name(), p.SHA256, wrapper.data); err != nil {
			return errors.Wrapf(err, "failed to CreateParcel for %s", wrapper.parcel.Label.Name)
		}
	}

	invoiceBytes, err := toml.Marshal(invoice)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal invoice")
	}

	invoiceBytes = append([]byte("# Autogenerated Bindle Invoice, do not edit\n\n"), invoiceBytes...)

	if err := os.WriteFile(filepath.Join(ctx.Cwd, "Invoice.toml"), invoiceBytes, util.PermFile); err != nil {
		return errors.Wrap(err, "failed to WriteFile for Invoice.toml")
	}

	util.LogDone("pushed")

	return nil
}

func parcelForData(name, mediaType string, data []byte) types.Parcel {
	sha := sha256.New()
	sha.Write(data)

	fileSHA := hex.EncodeToString(sha.Sum(nil))

	label := types.Label{
		SHA256:    fileSHA,
		MediaType: mediaType,
		Name:      name,
		Size:      uint64(len(data)),
	}

	parcel := types.Parcel{
		Label: label,
	}

	return parcel
}

func createOrReadKeypair(author string) (*types.SignatureKey, []byte, error) {
	var sigKey *types.SignatureKey
	var privKey []byte

	kr, err := keyring.LocalKeyring()
	if err != nil {
		sigKey, privKey, err = keyring.GenerateSignatureKey(author, "creator")
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to GenerateSignatureKey")
		}

		if err := keyring.AddLocalKey(sigKey); err != nil {
			return nil, nil, errors.Wrap(err, "failed to AddLocalKey")
		}

		if err := keyring.WritePrivKey(privKey, privKeyFilepath()); err != nil {
			return nil, nil, errors.Wrap(err, "failed to WritePrivateKey")
		}

		return sigKey, privKey, nil
	}

	// find the SignatureKey in the local Keyring.
	for i, k := range kr.Key {
		if k.Label == author {
			sigKey = &kr.Key[i]
			break
		}
	}

	// read the privkey from the '.ssh' location.
	privKey, err = keyring.ReadPrivKey(privKeyFilepath())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to ReadPrivKey")
	}

	return sigKey, privKey, nil
}

func privKeyFilepath() string {
	home := "$HOME"

	if usrHome, err := os.UserHomeDir(); err == nil {
		home = usrHome
	}

	return filepath.Join(home, ".ssh", "bindle_ed25519")
}
