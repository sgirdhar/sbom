package pkg

import (
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
	License      string
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
			log.Println("error while getting dpkg packages")
			return nil, err
		}
		sortPkgs(pkgs)
		return pkgs, nil
	}
	if strings.Contains(osRelease.PRETTY_NAME, "Debian") {
		log.Println("Getting packages for Debian")
		pkgs, err := analyzeDpkg(extractDir, "debian")
		if err != nil {
			log.Println("error while getting dpkg packages")
			return nil, err
		}
		sortPkgs(pkgs)
		return pkgs, nil
	} else {
		log.Fatalln("Linux distribution not supported yet!!!!")
		return nil, nil
	}
}

func sortPkgs(pkgs []Package) {
	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].Name < pkgs[j].Name
	})
}
