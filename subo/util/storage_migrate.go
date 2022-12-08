package util

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	art "github.com/plar/go-adaptive-radix-tree"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"
	"gopkg.in/yaml.v3"

	"github.com/suborbital/systemspec/fqmn"
	"github.com/suborbital/systemspec/system"
	"github.com/suborbital/systemspec/tenant"
)

var versionPattern = regexp.MustCompile(`v[0-9]*\.0\.0`)

type Tenant struct {
	Overview *system.TenantOverview `json:"overview,omitempty"`
	Stale    bool                   `json:"stale"`
}

// Runnable stores the relevant parts of a .runnable.yaml file
type Runnable struct {
	Identifier   string `json:"identifier"`
	Name         string `yaml:"name" json:"name"`
	Namespace    string `yaml:"namespace" json:"namespace"`
	Lang         string `yaml:"lang" json:"lang"`
	Version      string `yaml:"version" json:"version"`
	DraftVersion string `yaml:"draftVersion,omitempty" json:"draftVersion,omitempty"`
	APIVersion   string `yaml:"apiVersion,omitempty" json:"apiVersion,omitempty"`
	FQFN         string `yaml:"fqfn" json:"fqfn"`
}

func main() {
	ctx := context.Background()

	root := "/Users/ryanpridgeon/Workspaces/scn/tests/migration/old/local"

	if err := FileStoreMigration(root); err != nil {
		panic(err)
	}

	os.Setenv("AWS_ACCESS_KEY_ID", "minioadmin")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "minioadmin")
	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")

	root = "s3://mybucket?endpoint=127.0.0.1:9000&disableSSL=true&s3ForcePathStyle=true"
	if err := BlobStoreMigration(ctx, root); err != nil {
		panic(err)
	}
}

func FileStoreMigration(storageRoot string) error {
	recorder := art.New()

	if err := MapFileStorageV0(recorder, storageRoot); err != nil {
		return err
	}

	return MigrateFileStorageV0(recorder, storageRoot)
}

func MapFileStorageV0(recorder art.Tree, storageRoot string) error {
	var runnable Runnable
	err := filepath.WalkDir(storageRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// map plugins
		if d.Name() == ".runnable.yaml" {
			err = loadYaml(path, &runnable)
			if err != nil {
				return err
			}

			// track tenant for constructing tenant.json
			recorder.Insert(art.Key("tenant"), identifier(&runnable))
			recorder.Insert(canonicalize(&runnable, "runnable"), runnable)

			return err
		}

		// .runnable.yaml precedes all resources of interest
		if runnable.Name == "" {
			return nil
		}

		// map modules
		if !d.IsDir() && d.Name() == runnable.Name+".wasm" {
			recorder.Insert(canonicalize(&runnable, "bin", version(path)), path)
			return nil
		}

		// map recorder
		if !d.IsDir() && d.Name() == libFile(runnable.Lang) {
			recorder.Insert(canonicalize(&runnable, "src", version(path)), path)
			return nil
		}

		return nil
	})

	return err
}

func MigrateFileStorageV0(recorder art.Tree, storageRoot string) error {
	// process modules
	recorder.ForEachPrefix(art.Key("bin"), func(node art.Node) (cont bool) {
		if node.Kind() != art.Leaf {
			return true
		}

		bin := node.Value().(string)

		revBytes, revDigest, err := readAndHashFile(bin)
		if err != nil {
			panic(errors.Wrap(err, "compute module digest"))
		}

		revPath := versionPattern.ReplaceAllString(bin, revDigest)
		err = writeTo(revPath, revBytes)
		if err != nil {
			panic(errors.Wrapf(err, "copy binary from %s to %s", bin, revPath))
		}

		key := art.Key(bytes.Replace(node.Key(), []byte("bin"), []byte("digest"), 1))
		recorder.Insert(key, revDigest)

		return true
	})

	// process drafts
	recorder.ForEachPrefix(art.Key("src"), func(node art.Node) (cont bool) {
		if node.Kind() != art.Leaf {
			return true
		}

		src := node.Value().(string)

		key := art.Key(bytes.Replace(node.Key(), []byte("src"), []byte("digest"), 1))
		revDigest, ok := recorder.Search(key)
		if !ok {
			panic(errors.Wrapf(os.ErrNotExist, "look up module digest associated with %s", src))
		}

		revPath := versionPattern.ReplaceAllString(src, revDigest.(string))
		err := copyfile(revPath, src)
		if err != nil {
			panic(errors.Wrapf(err, "copy %s to %s", src, revPath))
		}

		return true
	})

	// process runnables(plugins)
	recorder.ForEachPrefix(art.Key("runnable"), func(node art.Node) (cont bool) {
		if node.Kind() != art.Leaf {
			return true
		}

		runnable := node.Value().(Runnable)

		ident := identifier(&runnable)

		ref, ok := recorder.Search(canonicalize(&runnable, "digest", runnable.Version))
		if !ok {
			panic(errors.Wrapf(os.ErrNotExist, "find ref for plugin %s", runnable.Name))
		}

		modRef := ref.(string)

		draftRef, ok := recorder.Search(canonicalize(&runnable, "digest", runnable.Version))
		if !ok && runnable.DraftVersion != "" {
			panic(errors.Wrapf(os.ErrNotExist, "find ref for plugin %s", runnable.Name))
		}

		fqmnString, err := fqmn.FromParts(ident, runnable.Namespace, runnable.Name, modRef)
		if err != nil {
			panic(errors.Wrapf(err, "construct fqmn for %s/%s", ident, runnable.Name))
		}

		module := tenant.Module{
			Name:       runnable.Name,
			Namespace:  runnable.Namespace,
			Lang:       runnable.Lang,
			APIVersion: runnable.APIVersion,
			Ref:        modRef,
			DraftRef:   optionalCoercion[string](draftRef),
			FQMN:       fqmnString,
			URI:        fmt.Sprintf("/%s/%s/%s/%s", ident, ref.(string), runnable.Namespace, runnable.Name),
			Revisions:  []tenant.ModuleRevision{},
		}

		recorder.Insert(canonicalize(&runnable, "module"), module)
		modBytes, err := yaml.Marshal(module)
		if err != nil {
			panic(errors.Wrapf(err, "serialize .module.yaml for %s/%s", ident, runnable.Name))
		}

		err = writeTo(string(canonicalize(&runnable, storageRoot, ".module.yaml")), modBytes)
		if err != nil {
			panic(errors.Wrapf(err, "write .module.yaml for %s/%s", ident, runnable.Name))
		}

		return true
	})

	var sysVersion int64 = 0
	recorder.ForEachPrefix(art.Key("tenant"), func(node art.Node) (cont bool) {
		if node.Kind() != art.Leaf {
			return true
		}

		overview := &system.TenantOverview{
			Identifier: node.Value().(string),
			Config:     new(tenant.Config),
		}

		recorder.ForEachPrefix(art.Key("module/"+overview.Identifier), func(node art.Node) (cont bool) {
			if node.Kind() != art.Leaf {
				return true
			}

			// calculate tenant version
			module := node.Value().(tenant.Module)
			binKey := art.Key(filepath.Join("bin", overview.Identifier, module.Namespace, module.Name))
			recorder.ForEachPrefix(binKey, func(node art.Node) (cont bool) {
				if node.Kind() != art.Leaf {
					return true
				}

				sysVersion += 1
				overview.Version += 1
				overview.Config.TenantVersion += 1

				return true
			})

			overview.Config.Modules = append(overview.Config.Modules, module)

			return true
		})

		tntBytes, err := json.Marshal(&Tenant{Overview: overview})
		if err != nil {
			panic(errors.Wrapf(err, "serialize tenant.json for %s", overview.Identifier))
		}

		err = writeTo(filepath.Join(storageRoot, overview.Identifier, "tenant.json"), tntBytes)
		if err != nil {
			panic(errors.Wrapf(err, "write tenant.json for %s", overview.Identifier))
		}

		return true
	})

	return writeTo(filepath.Join(storageRoot, "system", "version"), []byte(strconv.FormatInt(sysVersion, 10)))
}

func BlobStoreMigration(ctx context.Context, root string) error {
	recorder := art.New()

	var bucket *blob.Bucket
	bucket, err := blob.OpenBucket(ctx, root)
	if err != nil {
		return err
	}

	err = MapBlobStorageV0(ctx, recorder, bucket)
	if err != nil && err != io.EOF {
		return err
	}

	err = MigrateBlobStorageV0(ctx, recorder, bucket)
	if err != nil {
		return err
	}

	return nil
}

func MapBlobStorageV0(ctx context.Context, source art.Tree, bucket *blob.Bucket) error {
	var runnable Runnable
	var obj *blob.ListObject
	var err error

	ittr := bucket.List(nil)
	for obj, err = ittr.Next(ctx); err != io.EOF; obj, err = ittr.Next(ctx) {
		if err != nil {
			return errors.Wrap(err, "list bucket objects")
		}

		if strings.HasSuffix(obj.Key, ".runnable.yaml") {
			var dotRunnableBytes []byte
			dotRunnableBytes, err = bucket.ReadAll(ctx, obj.Key)
			if err != nil {
				return errors.Wrapf(err, "read %s", obj.Key)
			}

			if err = yaml.Unmarshal(dotRunnableBytes, &runnable); err != nil {
				return errors.Wrapf(err, "deserialized %s", obj.Key)
			}

			// track tenant for constructing tenant.json
			source.Insert(art.Key("tenant"), identifier(&runnable))
			source.Insert(canonicalize(&runnable, "runnable"), runnable)

			continue
		}

		// .runnable.yaml precedes all resources of interest
		if runnable.Name == "" {
			continue
		}

		// map modules
		if strings.HasSuffix(obj.Key, runnable.Name+".wasm") {
			source.Insert(canonicalize(&runnable, "bin", version(obj.Key)), obj.Key)
			continue
		}

		// map source
		if strings.HasSuffix(obj.Key, libFile(runnable.Lang)) {
			source.Insert(canonicalize(&runnable, "src", version(obj.Key)), obj.Key)
			continue
		}

		continue
	}

	return err
}

func MigrateBlobStorageV0(ctx context.Context, recorder art.Tree, bucket *blob.Bucket) error {
	// process modules
	recorder.ForEachPrefix(art.Key("bin"), func(node art.Node) (cont bool) {
		if node.Kind() != art.Leaf {
			return true
		}

		bin := node.Value().(string)

		revBytes, revDigest, err := readAndHashBucket(ctx, bucket, bin)
		if err != nil {
			panic(errors.Wrapf(err, "compute module digest"))
		}

		revPath := versionPattern.ReplaceAllString(bin, revDigest)
		err = writeToBucket(ctx, bucket, revPath, revBytes)
		if err != nil {
			panic(errors.Wrapf(err, "copy binary from %s to %s", bin, revPath))
		}

		key := art.Key(bytes.Replace(node.Key(), []byte("bin"), []byte("digest"), 1))
		recorder.Insert(key, revDigest)

		return true
	})

	// process drafts
	recorder.ForEachPrefix(art.Key("src"), func(node art.Node) (cont bool) {
		if node.Kind() != art.Leaf {
			return true
		}

		src := node.Value().(string)

		key := art.Key(bytes.Replace(node.Key(), []byte("src"), []byte("digest"), 1))
		revDigest, ok := recorder.Search(key)
		if !ok {
			panic(errors.Wrapf(os.ErrNotExist, "look up module digest associated with %s", src))
		}

		revPath := versionPattern.ReplaceAllString(src, revDigest.(string))
		err := bucket.Copy(ctx, revPath, src, nil)
		if err != nil {
			panic(errors.Wrapf(err, "copy %s to %s", src, revPath))
		}

		return true
	})

	// process runnables(plugins)
	recorder.ForEachPrefix(art.Key("runnable"), func(node art.Node) (cont bool) {
		if node.Kind() != art.Leaf {
			return true
		}

		runnable := node.Value().(Runnable)

		ident := identifier(&runnable)

		ref, ok := recorder.Search(canonicalize(&runnable, "digest", runnable.Version))
		if !ok {
			panic(errors.Wrapf(os.ErrNotExist, "find ref for plugin %s", runnable.Name))
		}

		modRef := ref.(string)

		draftRef, ok := recorder.Search(canonicalize(&runnable, "digest", runnable.Version))
		if !ok && runnable.DraftVersion != "" {
			panic(errors.Wrapf(os.ErrNotExist, "find ref for plugin %s", runnable.Name))
		}

		fqmnString, err := fqmn.FromParts(ident, runnable.Namespace, runnable.Name, modRef)
		if err != nil {
			panic(errors.Wrapf(os.ErrNotExist, "construct fqmn for %s/%s", ident, runnable.Name))
		}

		module := tenant.Module{
			Name:       runnable.Name,
			Namespace:  runnable.Namespace,
			Lang:       runnable.Lang,
			APIVersion: runnable.APIVersion,
			Ref:        modRef,
			DraftRef:   optionalCoercion[string](draftRef),
			FQMN:       fqmnString,
			URI:        fmt.Sprintf("/%s/%s/%s/%s", ident, ref.(string), runnable.Namespace, runnable.Name),
			Revisions:  []tenant.ModuleRevision{},
		}

		recorder.Insert(canonicalize(&runnable, "module"), module)
		modBytes, err := yaml.Marshal(module)
		if err != nil {
			panic(errors.Wrapf(err, "serialize .module.yaml for %s/%s", ident, runnable.Name))
		}

		err = writeToBucket(ctx, bucket, filepath.Join(ident, runnable.Namespace, runnable.Name, ".module.yaml"), modBytes)
		if err != nil {
			panic(errors.Wrapf(err, "write .module.yaml for %s/%s", ident, runnable.Name))
		}

		return true
	})

	// process tenant
	var sysVersion int64 = 0
	recorder.ForEachPrefix(art.Key("tenant"), func(node art.Node) (cont bool) {
		if node.Kind() != art.Leaf {
			return true
		}

		overview := &system.TenantOverview{
			Identifier: node.Value().(string),
			Config:     new(tenant.Config),
		}

		recorder.ForEachPrefix(art.Key("module/"+overview.Identifier), func(node art.Node) (cont bool) {
			if node.Kind() != art.Leaf {
				return true
			}

			// calculate tenant version
			module := node.Value().(tenant.Module)
			binKey := art.Key(filepath.Join("bin", overview.Identifier, module.Namespace, module.Name))
			recorder.ForEachPrefix(binKey, func(node art.Node) (cont bool) {
				if node.Kind() != art.Leaf {
					return true
				}

				sysVersion += 1
				overview.Version += 1
				overview.Config.TenantVersion += 1

				return true
			})

			overview.Config.Modules = append(overview.Config.Modules, module)

			return true
		})

		tntBytes, err := json.Marshal(&Tenant{Overview: overview})
		if err != nil {
			panic(errors.Wrapf(err, "serialize tenant.json for %s", overview.Identifier))
		}

		err = writeToBucket(ctx, bucket, filepath.Join(overview.Identifier, "tenant.json"), tntBytes)
		if err != nil {
			panic(errors.Wrapf(err, "write tenant.json for %s", overview.Identifier))
		}

		return true
	})

	return bucket.WriteAll(ctx, "/system/version", []byte(strconv.FormatInt(sysVersion, 10)), nil)
}

func canonicalize(runnable *Runnable, kind string, opts ...string) art.Key {
	key := append([]string{kind, identifier(runnable), runnable.Namespace, runnable.Name}, opts...)
	return art.Key(filepath.Join(key...))
}

func identifier(runnable *Runnable) string {
	return runnable.FQFN[:strings.Index(runnable.FQFN, "#")]
}

func version(binPath string) string {
	return versionPattern.FindString(binPath)
}

func optionalCoercion[T any](v any) T {
	var empty T

	c, ok := v.(T)
	if ok {
		return c
	}

	return empty
}

const (
	AssemblyScript = "assemblyscript"
	JavaScript     = "javascript"
	Rust           = "rust"
	TinyGo         = "tinygo"
	TypeScript     = "typescript"
)

func libFile(lang string) string {
	switch lang {
	case AssemblyScript:
		return "lib.ts"
	case JavaScript:
		return "lib.js"
	case Rust:
		return "lib.rs"
	case TinyGo:
		return "main.go"
	case TypeScript:
		return "lib.ts"
	default:
		return ""
	}
}

func writeTo(dstPath string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
		return err
	}

	return os.WriteFile(dstPath, data, os.ModePerm)
}

func copyfile(dstPath string, srcPath string) error {
	if err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
		return err
	}

	var dst *os.File
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	var src *os.File
	src, err = os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = io.Copy(dst, src)

	return err
}

func loadYaml(path string, obj any) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	defer reader.Close()

	return yaml.NewDecoder(reader).Decode(obj)
}

func readAndHashFile(path string) ([]byte, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}

	digest := sha256.Sum256(data)
	return data, hex.EncodeToString(digest[:]), nil
}

func writeToBucket(ctx context.Context, bucket *blob.Bucket, key string, data []byte) error {
	return bucket.WriteAll(ctx, key, data, nil)
}

func readAndHashBucket(ctx context.Context, bucket *blob.Bucket, key string) ([]byte, string, error) {
	data, err := bucket.ReadAll(ctx, key)
	if err != nil {
		return nil, "", err
	}

	digest := sha256.Sum256(data)
	return data, hex.EncodeToString(digest[:]), nil
}
