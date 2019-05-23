package module

import (
	templatex "text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

// template -
var template *templatex.Template

type plugin struct {
	pgs.ModuleBase
	pgsgo.Context
}

// New -
func New() *plugin {
	return &plugin{
		ModuleBase: pgs.ModuleBase{},
	}
}

// Name -
func (*plugin) Name() string {
	return "serialization"
}

func (p *plugin) InitContext(c pgs.BuildContext) {
	p.ModuleBase.InitContext(c)
	p.Context = pgsgo.InitContext(c.Parameters())
}
