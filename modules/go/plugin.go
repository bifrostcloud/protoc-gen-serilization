package module

import (

	// "github.com/bifrostcloud/protoc-gen-serialization/pkg/tags"

	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"regexp"
	"strings"

	"github.com/fatih/astrewrite"

	model "github.com/bifrostcloud/protoc-gen-serialization/modules/proto"
	"github.com/fatih/camelcase"
	pgs "github.com/lyft/protoc-gen-star"
)

var (
	rTags = regexp.MustCompile(`[\w_]+:"[^"]+"`)
)

func (p *plugin) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	target := p.Parameters().Str("target")
	if len(target) == 0 || strings.ToLower(target) == "go" {
		messages := map[string]bool{}

		for _, file := range targets {
			b := model.Base{Package: p.Context.PackageName(file).String()}
			imports := map[string]model.Package{
				"json": model.Package{
					PackagePath: "encoding/json",
				},
				"xml": model.Package{
					PackagePath: "encoding/xml",
				},
				"reflect": model.Package{
					PackagePath: "reflect",
				},
				"strings": model.Package{
					PackagePath: "strings",
				},
				"stacktrace": model.Package{
					PackageName: "stacktrace",
					PackagePath: "github.com/palantir/stacktrace",
				},
				"mapstructure": model.Package{
					PackageName: "mapstructure",
					PackagePath: "github.com/mitchellh/mapstructure",
				},
			}
			for _, srv := range file.Services() {
				s := model.Service{}
				s.UpperCamelCaseServiceName = srv.Name().UpperCamelCase().String()
				s.LowerCamelCaseServiceName = srv.Name().LowerCamelCase().String()
				for _, method := range srv.Methods() {
					upperCamelCaseMethodName := p.Context.Name(method).UpperCamelCase().String()
					m := model.Method{}
					m.UpperCamelCaseServiceName = srv.Name().UpperCamelCase().String()
					m.UpperCamelCaseMethodName = upperCamelCaseMethodName
					inputKey := p.Context.Name(method.Input()).UpperCamelCase().String()
					outputKey := p.Context.Name(method.Output()).UpperCamelCase().String()

					if !messages[inputKey] {
						messages[inputKey] = true
						m.InputType = p.Context.Name(method.Input()).String()
						if !method.Input().BuildTarget() {
							path := p.Context.ImportPath(method.Input()).String()
							imports[path] = model.Package{
								PackageName: p.Context.PackageName(method.Input()).String(),
								PackagePath: path,
							}
							m.InputType = p.Context.PackageName(method.Input()).String() + "." + p.Context.Name(method.Input()).String()
						}
						m.InputFields.TypeName = p.Context.Name(method.Input()).UpperCamelCase().String()
						m.InputFields.Type = m.InputType
						for _, field := range method.Input().Fields() {
							m.InputFields.FieldImport = append(m.InputFields.FieldImport, model.FieldImport{
								Name: field.Name().UpperCamelCase().String(),
								Tag:  field.Name().String(),
							})
							m.InputFields.Base = append(m.InputFields.Base, strings.ToLower(field.Name().String()))
							m.InputFields.Lowercase = append(m.InputFields.Lowercase, strings.ToLower(field.Name().LowerCamelCase().String()))
							m.InputFields.DotNotation = append(m.InputFields.DotNotation, strings.ToLower(field.Name().LowerDotNotation().String()))
							spl := camelcase.Split(field.Name().UpperCamelCase().String())
							m.InputFields.ParamCase = append(m.InputFields.ParamCase, strings.ToLower(strings.Join(spl, "-")))
						}

					}
					if !messages[outputKey] {
						messages[outputKey] = true
						m.OutputType = p.Context.Name(method.Output()).String()
						if !method.Output().BuildTarget() {
							path := p.Context.ImportPath(method.Output()).String()
							imports[path] = model.Package{
								PackagePath: path,
								PackageName: p.Context.PackageName(method.Output()).String(),
							}
							m.OutputType = p.Context.PackageName(method.Output()).String() + "." + p.Context.Name(method.Output()).String()
						}
						m.OutputFields.TypeName = p.Context.Name(method.Output()).UpperCamelCase().String()
						m.OutputFields.Type = m.OutputType
						for _, field := range method.Output().Fields() {
							m.OutputFields.FieldImport = append(m.OutputFields.FieldImport, model.FieldImport{
								Name: field.Name().UpperCamelCase().String(),
								Tag:  field.Name().String(),
							})
							m.OutputFields.Base = append(m.OutputFields.Base, strings.ToLower(field.Name().String()))
							m.OutputFields.Lowercase = append(m.OutputFields.Lowercase, strings.ToLower(field.Name().LowerCamelCase().String()))
							m.OutputFields.DotNotation = append(m.OutputFields.DotNotation, strings.ToLower(field.Name().LowerDotNotation().String()))
							spl := camelcase.Split(field.Name().UpperCamelCase().String())
							m.OutputFields.ParamCase = append(m.OutputFields.ParamCase, strings.ToLower(strings.Join(spl, "-")))
						}
					}

					s.Methods = append(s.Methods, m)
				}
				b.Services = append(b.Services, s)
			}
			if len(b.Services) == 0 {
				continue
			}
			for _, pkg := range imports {
				b.Imports = append(b.Imports, pkg)
			}
			pname := p.Context.OutputPath(file).SetExt(".serialization.go").String()
			p.OverwriteGeneratorTemplateFile(
				pname,
				template.Lookup("Base"),
				&b,
			)
			gfname := p.Context.OutputPath(file).SetExt(".go").String()

			fs := token.NewFileSet()
			fn, err := parser.ParseFile(fs, gfname, nil, parser.ParseComments)
			p.CheckErr(err)
			// Adding tags ...
			visited := map[string]bool{}

			rewriteFunc := func(n ast.Node) (ast.Node, bool) {
				typeSpec, ok := n.(*ast.TypeSpec)
				if !ok {
					return n, true
				}
				structDecl, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					return n, true
				}
				for _, field := range structDecl.Fields.List {
					if len(field.Names) > 0 {
						if !visited[field.Names[0].Name] {
							visited[field.Names[0].Name] = true
							if !strings.HasPrefix(field.Names[0].Name, "XXX") {
								fieldSnakeCase := pgs.Name(field.Names[0].Name).LowerSnakeCase().String()
								splitted := rTags.FindAllString(field.Tag.Value, -1)
								splitted = append(splitted, `xml:"`+fieldSnakeCase+`,omitempty"`)
								splitted = append(splitted, `mapstructure:"`+fieldSnakeCase+`"`)
								field.Tag.Value = "`" + strings.Join(splitted, " ") + "`"
							}
						}

					}
				}
				return typeSpec, true
			}
			astrewrite.Walk(fn, rewriteFunc)
			rewritten := astrewrite.Walk(fn, rewriteFunc)
			var buf bytes.Buffer
			printer.Fprint(&buf, fs, rewritten)
			p.OverwriteGeneratorFile(gfname, buf.String())
		}
	}

	return p.Artifacts()
}
