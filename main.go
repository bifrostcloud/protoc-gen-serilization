package main

import (
	gomodule "github.com/bifrostcloud/protoc-gen-serialization/modules/go"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func main() {
	mod := pgs.Init(pgs.DebugEnv("DEBUG"))

	mod.RegisterModule(gomodule.New())
	// TODO : THis may cause issues once other target languages are added
	mod.RegisterPostProcessor(pgsgo.GoFmt())
	mod.Render()

}
