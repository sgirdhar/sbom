package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/sgirdhar/sbom/pkg"
	"github.com/sgirdhar/sbom/report"
	"github.com/sgirdhar/sbom/util"

	"github.com/spf13/cobra"
)

var tarFile, outputFormat, compareFile string
var verbose bool

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Generate SBOM using image name or image tarfile",
	Long: `scan command can be used to generate sbom from a container image.
It can also compare this sbom with sbom generated using some other tool.
Image name (with either tag or digest) or image tarball can be passed as input.`,
	Example: `  Possible Inputs
	sbom scan alpine:latest
	sbom scan alpine@sha256:51103b3f2993cbc1b45ff9d941b5d461484002792e02aa29580ec5282de719d4
	sbom scan --tar /path/to/tarfile/my_image.tar
	
  Output Formats
	sbom scan alpine:latest
	sbom scan alpine:latest --output cyclonedx
	sbom scan alpine:latest --output cyclonedx-json
		
  Compare SBOM
	sbom scan alpine:latest --compare /path/to/cyclonedx-json/my_sbom.json`,
	// Args: cobra.ExactArgs(1),
	Args: func(cmd *cobra.Command, args []string) error {
		if tarFile == "" && len(args) == 0 {
			return errors.New("accepts 1 arg, received 0")
		}
		// if outputFormat == "" && len(args) == 0 {
		// 	return errors.New("accepts 1 arg, received 0")
		// }
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetOutput(os.Stderr)
		} else {
			log.SetOutput(ioutil.Discard)
		}
		if tarFile != "" {
			scanTarImage(tarFile)
		} else {
			scanStringImage(args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	scanCmd.PersistentFlags().StringVarP(&tarFile, "tar", "t", "", "tarball file path")

	scanCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "output format")

	scanCmd.PersistentFlags().StringVarP(&compareFile, "compare", "c", "", "compare generated sbom with the result of another cyclonedx output")

	scanCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output including logs")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func scanStringImage(image string) {

	log.Printf("scan called with args: %v\n", image)

	v1Image, err := util.PullImage(image)
	if err != nil {
		fmt.Println(util.Red+"error while pulling image: ", err)
		log.Fatalln()
	}

	tempDir, err := util.SaveAndUntarImage(v1Image, image)
	if err != nil {
		fmt.Println(util.Red+"error while saving or untarring image: ", err)
		log.Fatalln()
	}
	commonProcessing(tempDir, image)
}

func scanTarImage(tarFile string) {
	log.Printf("scan called with args: %v\n", tarFile)

	tempDir, err := util.UntarImage(tarFile)
	if err != nil {
		fmt.Println(util.Red+"error while untarring tar image: ", err)
		log.Fatalln()
	}
	// log.Println("tempDir: ", tempDir)
	commonProcessing(tempDir, tarFile)
}

func commonProcessing(tempDir, image string) {
	manifest, err := util.ReadImageManifest(tempDir)
	if err != nil {
		fmt.Println(util.Red+"error while reading image manifest: ", err)
		log.Fatalln()
	}

	configFile, err := util.ReadImageConfig(tempDir, manifest)
	if err != nil {
		fmt.Println(util.Red+"error while reading image config: ", err)
		log.Fatalln()
	}

	extractLayer, err := util.ExtractLayer(tempDir, manifest)
	if err != nil {
		fmt.Println(util.Red+"error while extracting layer: ", err)
		log.Fatalln()
	}

	osRelease, err := util.IdentifyOsRelease(extractLayer)
	if err != nil {
		fmt.Println(util.Red+"error while identifying OS release: ", err)
		log.Fatalln()
	}

	pkgs, err := pkg.AnalyzePkg(osRelease, extractLayer)
	if err != nil {
		fmt.Println(util.Red+"error while fetching package information: ", err)
		log.Fatalln()
	}

	// cleanup
	defer util.RemoveDir(tempDir)

	log.Printf("Components identified by %v: %v\n", util.ApplicationName, len(pkgs))

	// TO DO: Refactor - create new functions for output generation and comparison
	if compareFile == "" && len(compareFile) == 0 {
		generateSbom(image, configFile, pkgs, osRelease)
	} else {
		compareSbom(pkgs)
	}
}

func generateSbom(image string, configFile v1.ConfigFile, pkgs []pkg.Package, osRelease util.OsRelease) {
	log.Println("Generating sbom...")
	var err error
	if strings.Contains(outputFormat, "cyclonedx") {
		err = report.GenerateCycloneDxReport(image, outputFormat, configFile, pkgs, osRelease)
	} else if outputFormat == "" || strings.Contains(outputFormat, "table") {
		err = report.GenerateTableReport(pkgs)
	} else {
		err = errors.New("invalid output format")
	}
	if err != nil {
		fmt.Println(util.Red+"error while generating report:", err)
		log.Fatalln()
	}
}

func compareSbom(pkgs []pkg.Package) {
	log.Println("Comparing...")
	identifiedMap := pkg.GetPkgMap(pkgs)

	readMap, toolName, err := report.GetPkgMap(compareFile)
	if err != nil {
		fmt.Println(util.Red+"error while reading json report", err)
		log.Fatalln()
	}

	err = pkg.ListComp(identifiedMap, readMap, toolName)
	if err != nil {
		fmt.Println(util.Red+"error while comparing sbom: ", err)
		log.Fatalln()
	}
}
