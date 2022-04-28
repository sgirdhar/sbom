package cmd

import (
	"errors"
	"log"
	"strings"

	"github.com/sgirdhar/sbom/pkg"
	"github.com/sgirdhar/sbom/report"
	"github.com/sgirdhar/sbom/util"

	"github.com/spf13/cobra"
)

var tarFile, outputFormat string

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Generate SBOM using image name or image tarfile",
	Long: `scan command is used to generate sbom from a container image.
Image name (with either tag or digest) or image tar can be passed as input`,
	Example: `  Input image string
	sbom scan my_image:my_tag
	sbom scan my_image@my_digest

  Input tarfile
	sbom scan --tar /path/to/tarfile/my_image.tar`,
	// Args: cobra.ExactArgs(1),
	Args: func(cmd *cobra.Command, args []string) error {
		if tarFile == "" && len(args) == 0 {
			return errors.New("accepts 1 arg, received 0")
		}
		if outputFormat == "" && len(args) == 0 {
			return errors.New("accepts 1 arg, received 0")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tarFile != "" {
			scanTarImage(tarFile)
		} else {
			scanStringImage(args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	scanCmd.PersistentFlags().StringVarP(&tarFile, "tarball", "t", "", "tar file path")

	scanCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "output format")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func scanStringImage(image string) {

	log.Printf("scan called with args: %v\n", image)

	v1Image, err := util.PullImage(image)
	if err != nil {
		log.Fatalln("error while pulling image: ", err.Error())
	}

	tempDir, err := util.SaveAndUntarImage(v1Image, image)
	if err != nil {
		log.Fatalln("error while saving or untarring image: ", err.Error())
	}
	commonProcessing(tempDir, image)
}

func scanTarImage(tarFile string) {
	log.Printf("scan called with args: %v\n", tarFile)

	tempDir, err := util.UntarImage(tarFile)
	if err != nil {
		log.Fatalln("error while untarring tar image: ", err.Error())
	}
	// log.Println("tempDir: ", tempDir)
	commonProcessing(tempDir, tarFile)
}

func commonProcessing(tempDir, image string) {
	manifest, err := util.ReadImageManifest(tempDir)
	if err != nil {
		log.Fatalln("error while reading image manifest: ", err.Error())
	}

	configFile, err := util.ReadImageConfig(tempDir, manifest)
	if err != nil {
		log.Fatalln("error while reading image config: ", err.Error())
	}

	extractLayer, err := util.ExtractLayer(tempDir, manifest)
	if err != nil {
		log.Fatalln("error while extracting layer: ", err.Error())
	}

	osRelease, err := util.IdentifyOsRelease(extractLayer)
	if err != nil {
		log.Fatalln("error while identifying OS release: ", err.Error())
	}

	pkgs, err := pkg.AnalyzePkg(osRelease, extractLayer)
	if err != nil {
		log.Fatalln("error while fetching package information: ", err.Error())
	}

	util.RemoveDir(tempDir)

	log.Println("No. of packages found: ", len(pkgs))

	if strings.Contains(outputFormat, "cyclonedx") {
		err = report.GenerateCycloneDxReport(image, outputFormat, configFile, pkgs, osRelease)
	} else if outputFormat == "" || strings.Contains(outputFormat, "table") {
		err = report.GenerateTableReport(pkgs)
	} else {
		err = errors.New("invalid output format")
	}
	if err != nil {
		log.Fatalln("error while generating report: ", err.Error())
	}
}
