package pkg

import (
	"errors"
	"log"
	"sort"
	"strings"

	"github.com/sgirdhar/sbom/util"
)

type Package struct {
	Name         string
	Version      string
	Architecture string
	Type         string
	Licenses     []string
	URL          string
	PURL         string
}

func (pkg *Package) Empty() bool {
	return pkg.Name == "" || pkg.Version == ""
}

func AnalyzePkg(osRelease util.OsRelease, extractDir string) ([]Package, error) {

	if strings.Contains(osRelease.PRETTY_NAME, "Alpine") {
		log.Println("Getting packages for Alpine")
		pkgs, err := analyzeApk(extractDir)
		if err != nil {
			log.Println("error while getting alpine packages")
			return nil, err
		}
		sortPkgs(pkgs)
		return pkgs, nil
	}
	if strings.Contains(osRelease.PRETTY_NAME, "Ubuntu") {
		log.Println("Getting packages for Ubuntu")
		pkgs, err := analyzeDpkg(extractDir, "ubuntu")
		if err != nil {
			log.Printf("error while getting dpkg packages for %v\n", osRelease.PRETTY_NAME)
			return nil, err
		}
		sortPkgs(pkgs)
		return pkgs, nil
	}
	if strings.Contains(osRelease.PRETTY_NAME, "Debian") {
		log.Println("Getting packages for Debian")
		pkgs, err := analyzeDpkg(extractDir, "debian")
		if err != nil {
			log.Printf("error while getting dpkg packages for %v\n", osRelease.PRETTY_NAME)
			return nil, err
		}
		sortPkgs(pkgs)
		return pkgs, nil
	} else {
		log.Println("linux distribution not supported yet")
		err := errors.New("linux distribution not supported yet")
		return nil, err
	}
}

func sortPkgs(pkgs []Package) {
	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].Name < pkgs[j].Name
	})
}

func GetPkgMap(pkgs []Package) map[string]Package {
	var pkgMap = make(map[string]Package)
	for _, pkg := range pkgs {
		pkg.Type = "library"
		pkgMap[pkg.Name+"-"+pkg.Version] = pkg
	}
	return pkgMap
}
