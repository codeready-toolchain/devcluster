// +build ignore

package main

import (
	"github.com/codeready-toolchain/devcluster/pkg/log"
	"github.com/codeready-toolchain/devcluster/pkg/static"

	"github.com/shurcooL/vfsgen"
)

func main() {
	log.Init("assert-generate")
	err := vfsgen.Generate(static.Assets, vfsgen.Options{
		PackageName:  "static",
		BuildTags:    "!dev",
		VariableName: "Assets",
		Filename:     "pkg/static/generated_assets.go",
	})
	if err != nil {
		log.Error(nil, err, err.Error())
	}
}
