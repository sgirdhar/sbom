package pkg

import (
	"log"
	"strings"

	"github.com/sgirdhar/sbom/util"
)

type void struct{}

func analyzeDpkg(extractDir, osDistro string) ([]Package, error) {
	const pkgPath string = "/var/lib/dpkg/status"
	pkgLines, err := util.ReadFile(extractDir + pkgPath)
	if err != nil {
		log.Printf("error while reading %v file\n", pkgPath)
		return nil, err
	}

	pkgs, err := parseDpkg(pkgLines, osDistro, extractDir)
	if err != nil {
		log.Println("error while parsing dpkg packages")
		return nil, err
	}
	return pkgs, nil
}

func parseDpkg(pkgLines []string, osDistro, extractDir string) ([]Package, error) {
	log.Println("Parsing dpkg")
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
				licenses, err := getLicence(extractDir, pkg.Name)
				if err != nil {
					log.Printf("error while reading getting licenses for %v\n", pkg.Name)
					return nil, err
				}
				pkg.Licenses = licenses
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
	return pkgs, nil
}

func getLicence(extractDir, name string) ([]string, error) {
	const docPath = "/usr/share/doc/"
	copyrightFile := extractDir + docPath + name + "/copyright"

	// check for alias folders - skipping these at the moment
	// TODO: Find a way to access alias folders
	err := util.CheckFile(copyrightFile)
	if err != nil {
		log.Printf("Problamatic file: skipping license info for %v \n", copyrightFile)
		return nil, nil
	}

	copyrightLines, err := util.ReadFile(copyrightFile)
	if err != nil {
		log.Printf("error while reading %v file: %v\n", copyrightFile, err)
		return nil, err
	}
	return parseCopyright(copyrightLines), nil
}

func parseCopyright(copyrightLines []string) []string {
	var member void
	set := make(map[string]void)
	var licenses []string
	for _, line := range copyrightLines {
		if strings.HasPrefix(line, "License: ") {
			set[strings.TrimSpace(strings.TrimPrefix(line, "License: "))] = member
		}
	}
	for key := range set {
		licenses = append(licenses, key)
	}
	return licenses
}
