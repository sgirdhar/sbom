package pkg

import (
	"github.com/package-url/packageurl-go"
)

func GeneratePurl(pkg Package, osDistro string) string {

	var qualifier1 packageurl.Qualifier
	qualifier1.Key = "arch"
	qualifier1.Value = pkg.Architecture

	var qualifiers []packageurl.Qualifier
	qualifiers = append(qualifiers, qualifier1)

	purl := packageurl.NewPackageURL(pkg.Type, osDistro, pkg.Name, pkg.Version, qualifiers, "")
	return purl.ToString()
}
