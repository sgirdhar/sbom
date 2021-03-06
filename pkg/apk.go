package pkg

import (
	"log"

	"github.com/sgirdhar/sbom/util"
)

func analyzeApk(extractDir string) ([]Package, error) {
	const pkgPath string = "/lib/apk/db/installed"
	pkgLines, err := util.ReadFile(extractDir + pkgPath)
	if err != nil {
		log.Printf("error while reading %v file", pkgPath)
		return nil, err
	}
	pkgs := parseApk(pkgLines)
	return pkgs, nil
}

func parseApk(pkgLines []string) []Package {
	log.Println("Parsing apk")
	var (
		pkgs []Package
		pkg  Package
	)
	for _, line := range pkgLines {
		// check package if paragraph end
		if len(line) < 2 {
			if !pkg.Empty() {
				pkg.Type = "apk"
				pkg.PURL = GeneratePurl(pkg, "alpine")
				pkgs = append(pkgs, pkg)
			}
			pkg = Package{}
			continue
		}

		switch line[:2] {
		case "P:":
			pkg.Name = line[2:]
		case "V:":
			pkg.Version = line[2:]
		case "A:":
			pkg.Architecture = line[2:]
		case "L:":
			pkg.Licenses = []string{line[2:]}
		case "U:":
			pkg.URL = line[2:]
		}
	}
	// in case of last paragraph
	if !pkg.Empty() {
		pkgs = append(pkgs, pkg)
	}
	return pkgs
}
