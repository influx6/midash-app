// +build ignore

// Generates the manifests.go package file.

//go:generate go run generate.go

package main

import (
	"encoding/json"
	"fmt"
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/gu-io/gu/shell"
	"github.com/gu-io/gu/shell/parse"
)

var pkg = "// Package {{PKG}} defines a package which embeds all external files which are used\n// within the project.\n// This is automatically genereated and do not edit by hand.\n\n//go:generate go run generate.go\n\npackage {{PKG}}\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"errors\"\n\n\t\"github.com/gu-io/gu/shell\"\n)\n\n// Manifests defines the slice of manifests files loaded from the generated\n// data.\nvar Manifests []shell.AppManifest\n\n// Get returns the a shell.AppManifest if it exists with the given name.\nfunc Get(name string) (shell.AppManifest, error) {\n\tfor _, manifest := range Manifests {\n\t\tif manifest.Name == name {\n\t\t\treturn manifest, nil\n\t\t}\n\t}\n\n\treturn shell.AppManifest{}, errors.New(\"Not Found\")\n}\n\nfunc init (){\n  if err := json.Unmarshal([]byte({{MANIFEST}}),&Manifests); err != nil {\n  \tpanic(fmt.Sprintf(\"Failed to unmarshal manifests json: %+q\\n\", err))\n  }\n}\n"

func main() {
	cdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	manifestDir := filepath.Join(cdir, "../manifests")
	componentDir := filepath.Join(cdir, "../components")

	res, err := parse.ShellResources(componentDir)
	if err != nil {
		panic(err)
	}

	var manifests []shell.AppManifest

	for _, rs := range res {
		ms, terr := rs.GenManifests()
		if terr != nil {
			panic(terr)
		}

		manifests = append(manifests, *ms)
	}

	file, err := os.Create(filepath.Join(manifestDir, "manifests.go"))
	if err != nil {
		panic("Failed to create manifest pkg file: " + err.Error())
	}

	defer file.Close()

	maniJSON, terr := json.MarshalIndent(manifests, "", "\t")
	if terr != nil {
		panic("Failed to create manifest json: " + terr.Error())
	}

	if len(maniJSON) == 0 || bytes.Equal(maniJSON,[]byte("null")){
		maniJSON = []byte("[]")
	}

	maniJSONQuoted := fmt.Sprintf("%+q", maniJSON)
	pkg = strings.Replace(pkg, "{{PKG}}", "manifests", -1)
	pkg = strings.Replace(pkg, "{{MANIFEST}}", maniJSONQuoted, -1)

	file.Write([]byte(pkg))
}
