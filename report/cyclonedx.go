package report

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	cdx "github.com/CycloneDX/cyclonedx-go"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/uuid"
	"github.com/sgirdhar/sbom/pkg"
	"github.com/sgirdhar/sbom/util"
)

func GenerateCycloneDxReport(image, outputFormat string, configFile v1.ConfigFile, pkgs []pkg.Package, osRelease util.OsRelease) error {

	metadata := getMetadata(image, configFile, osRelease)
	components := getComponents(pkgs)
	// dependencies := getDependencies()

	// Assemble the BOM
	bom := cdx.NewBOM()
	bom.SerialNumber = uuid.New().URN()
	bom.Metadata = &metadata
	bom.Components = &components
	// bom.Dependencies = &dependencies

	// Encode the BOM
	var encoder cdx.BOMEncoder
	if strings.Contains(outputFormat, "json") {
		encoder = cdx.NewBOMEncoder(os.Stdout, cdx.BOMFileFormatJSON)
	} else {
		encoder = cdx.NewBOMEncoder(os.Stdout, cdx.BOMFileFormatXML)
	}

	encoder.SetPretty(true)
	if err := encoder.Encode(bom); err != nil {
		fmt.Println(util.Red + "Error while encoding BOM")
		return err
	}
	return nil
}

func getMetadata(image string, configFile v1.ConfigFile, osRelease util.OsRelease) cdx.Metadata {
	return cdx.Metadata{
		// Define metadata about the main component
		// (the component which the BOM will describe)
		Timestamp: time.Now().Format(time.RFC3339),
		Tools: &[]cdx.Tool{
			{
				Vendor:  "Open Source Software",
				Name:    util.ApplicationName,
				Version: util.ApplicationVersion,
			},
		},
		Component: &cdx.Component{
			// BOMRef:  "pkg:golang/acme-inc/acme-app@v1.0.0",
			Type:    cdx.ComponentTypeContainer,
			Name:    image,
			Version: configFile.Config.Image,
		},
		// Use properties to include an internal identifier for this BOM
		// https://cyclonedx.org/use-cases/#properties--name-value-store
		Properties: &[]cdx.Property{
			{
				Name:  "Created",
				Value: configFile.Created.String(),
			},
			{
				Name:  "Architecture",
				Value: configFile.Architecture,
			},
			{
				Name:  "Identified OS",
				Value: osRelease.PRETTY_NAME,
			},
		},
	}
}

func getComponents(pkgs []pkg.Package) []cdx.Component {
	// Define the components that image ships with
	// https://cyclonedx.org/use-cases/#inventory
	var componentList []cdx.Component
	for _, pkg := range pkgs {
		bom := cdx.Component{
			Type:       cdx.ComponentTypeLibrary,
			Name:       pkg.Name,
			Version:    pkg.Version,
			PackageURL: pkg.PURL,
		}
		if pkg.License != "" {
			bom.Licenses = &cdx.Licenses{
				// cdx.LicenseChoice{Expression: pkg.License},
				cdx.LicenseChoice{License: &cdx.License{
					Name: pkg.License,
				}},
			}
		}
		componentList = append(componentList, bom)
	}

	return componentList
}

// func getDependencies() []cdx.Dependency {
// 	// Define the dependency graph
// 	// https://cyclonedx.org/use-cases/#dependency-graph
// 	return []cdx.Dependency{
// 		{
// 			Ref: "pkg:golang/acme-inc/acme-app@v1.0.0",
// 			Dependencies: &[]cdx.Dependency{
// 				{Ref: "pkg:golang/github.com/CycloneDX/cyclonedx-go@v0.3.0"},
// 			},
// 		},
// 		{
// 			Ref: "pkg:golang/github.com/CycloneDX/cyclonedx-go@v0.3.0",
// 		},
// 	}
// }

func ReadCycloneDxReport(sbomFile string) (*cdx.BOM, error) {

	// Acquire BOM
	file, err := os.Open(sbomFile)
	if err != nil {
		fmt.Printf(util.Red+"error opening %v", sbomFile)
		return nil, err
	}
	defer file.Close()

	// Decode BOM
	bom := new(cdx.BOM)
	decoder := cdx.NewBOMDecoder(file, cdx.BOMFileFormatJSON)
	if err := decoder.Decode(bom); err != nil {
		fmt.Println(util.Red + "error while decoding BOM")
		return nil, err
	}

	log.Printf("Successfully decoded %s\n", sbomFile)
	// log.Printf("Generated: %s\n", bom.Metadata.Timestamp)
	log.Printf("Components identified by %v: %v\n", (*bom.Metadata.Tools)[0].Name, len(*bom.Components))

	return bom, nil
}

func GetPkgMap(sbomFile string) (map[string]pkg.Package, string, error) {
	bom, err := ReadCycloneDxReport(sbomFile)
	if err != nil {
		fmt.Println(util.Red + "error while reading cyclonedx sbom")
		return nil, "", err
	}
	var componentMap = make(map[string]pkg.Package)
	for _, component := range *bom.Components {
		componentMap[component.Name+"-"+component.Version] = pkg.Package{
			Name:    component.Name,
			Version: component.Version,
			PURL:    component.PackageURL,
			Type:    string(component.Type),
		}
	}
	return componentMap, (*bom.Metadata.Tools)[0].Name, nil
}
