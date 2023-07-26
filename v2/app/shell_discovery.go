package app

import (
	"golang.org/x/exp/slices"

	"encr.dev/pkg/paths"
	"encr.dev/v2/internals/parsectx"
	"encr.dev/v2/internals/pkginfo"
	"encr.dev/v2/parser"
	"encr.dev/v2/parser/shell/cmd"
)

type ShellPackages struct {
	paths    []paths.FS
	Packages []*pkginfo.Package
}

// Contains returns true if the given package is contained or is
// a package in which a shell command has been directly created in
func (p *ShellPackages) Contains(pkg *pkginfo.Package) bool {
	for _, path := range p.paths {
		if pkg.FSPath.HasPrefix(path) {
			return true
		}
	}
	return false
}

// discoverShellPackages discovers any packages which
// create shell commands, and then marks the highest level
// package as a shell package.
func discoverShellPackages(pc *parsectx.Context, result *parser.Result) *ShellPackages {
	foundPaths := make(map[paths.FS]struct{})
	allPkgs := make(map[*pkginfo.Package]struct{})

	registerPath := func(path paths.FS) {
		for existing := range foundPaths {
			if path.HasPrefix(existing) {
				return
			}
			if existing.HasPrefix(path) {
				delete(foundPaths, existing)
				break
			}
		}

		foundPaths[path] = struct{}{}
	}

	// Find all shell things
	for _, r := range result.Resources() {
		switch r := r.(type) {
		case *cmd.Command:
			allPkgs[r.Package()] = struct{}{}
			registerPath(r.Package().FSPath)
		}
	}

	// Now sort them
	asSlice := make([]paths.FS, 0, len(foundPaths))
	for path := range foundPaths {
		asSlice = append(asSlice, path)
	}
	slices.SortStableFunc(asSlice, func(a, b paths.FS) bool {
		return a.ToIO() < b.ToIO()
	})

	pkgsAsSlic := make([]*pkginfo.Package, 0, len(allPkgs))
	for pkg := range allPkgs {
		pkgsAsSlic = append(pkgsAsSlic, pkg)
	}
	slices.SortStableFunc(pkgsAsSlic, func(a, b *pkginfo.Package) bool {
		return a.FSPath.ToIO() < b.FSPath.ToIO()
	})

	return &ShellPackages{asSlice, pkgsAsSlic}
}
