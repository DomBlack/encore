package cmd

import (
	"go/ast"
	"go/token"

	"encr.dev/pkg/paths"
	"encr.dev/v2/internals/pkginfo"
	"encr.dev/v2/parser/internal/parseutil"
	"encr.dev/v2/parser/resource"
	"encr.dev/v2/parser/resource/resourceparser"
)

type Command struct {
	AST  *ast.CallExpr
	File *pkginfo.File
	Name string // The unique name of the command
}

func (c *Command) Kind() resource.Kind       { return resource.ShellCommand }
func (c *Command) Package() *pkginfo.Package { return c.File.Pkg }
func (c *Command) ASTExpr() ast.Expr         { return c.AST }
func (c *Command) ResourceName() string      { return c.Name }
func (c *Command) Pos() token.Pos            { return c.AST.Pos() }
func (c *Command) End() token.Pos            { return c.AST.End() }
func (c *Command) SortKey() string           { return c.Name }

var CommandParser = &resourceparser.Parser{
	Name: "Shell Command",

	InterestingImports: []paths.Pkg{"encore.dev/shell"},
	Run: func(p *resourceparser.Pass) {
		name := pkginfo.QualifiedName{PkgPath: "encore.dev/shell", Name: "NewCommand"}

		spec := &parseutil.ReferenceSpec{
			MinTypeArgs: 0,
			MaxTypeArgs: 0,
			Parse:       parseCommand,
		}

		parseutil.FindPkgNameRefs(p.Pkg, []pkginfo.QualifiedName{name}, func(file *pkginfo.File, name pkginfo.QualifiedName, stack []ast.Node) {
			parseutil.ParseReference(p, spec, parseutil.ReferenceData{
				File:         file,
				Stack:        stack,
				ResourceFunc: name,
			})
		})
	},
}

func parseCommand(d parseutil.ReferenceInfo) {
	name := d.ResourceFunc.NaiveDisplayName()

	command := &Command{
		AST:  d.Call,
		File: d.File,
		Name: name,
	}

	d.Pass.RegisterResource(command)
	d.Pass.AddBind(d.File, d.Ident, command)
}
