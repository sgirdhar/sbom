package pkg

import (
	"fmt"
	"strings"

	"github.com/sgirdhar/sbom/util"
)

func analyzeDpkg(extractDir, osDistro string) ([]Package, error) {
	const pkgPath string = "/var/lib/dpkg/status"
	pkgLines, err := util.ReadFile(extractDir + pkgPath)
	if err != nil {
		fmt.Printf(util.Red+"error while reading %v file", pkgPath)
		return nil, err
	}

	pkgs := parseDpkg(pkgLines, osDistro)
	return pkgs, nil
}

func parseDpkg(pkgLines []string, osDistro string) []Package {
	var (
		pkgs []Package
		pkg  Package
	)
	for _, line := range pkgLines {
		// check package if paragraph end
		if len(line) < 1 {
			if !pkg.Empty() {
				pkg.Type = "deb"
				pkg.PURL = GeneratePurl(pkg, osDistro)
				pkgs = append(pkgs, pkg)
			}
			pkg = Package{}
			continue
		}

		if strings.HasPrefix(line, "Package: ") {
			pkg.Name = strings.TrimSpace(strings.TrimPrefix(line, "Package: "))
		} else if strings.HasPrefix(line, "Architecture: ") {
			pkg.Architecture = strings.TrimSpace(strings.TrimPrefix(line, "Architecture: "))
		} else if strings.HasPrefix(line, "Version: ") {
			pkg.Version = strings.TrimSpace(strings.TrimPrefix(line, "Version: "))
		} else if strings.HasPrefix(line, "Homepage: ") {
			pkg.URL = strings.TrimSpace(strings.TrimPrefix(line, "Homepage: "))
		}

	}
	// in case of last paragraph
	if !pkg.Empty() {
		pkgs = append(pkgs, pkg)
	}
	return pkgs
}
