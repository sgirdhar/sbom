package cmd

import (
	"github.com/spf13/cobra"
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generate SBOM using image name and tag",
	Long: `image command is used to generate sbom from a container image.
Expected input sbom image image_name:tag`,
	Run: func(cmd *cobra.Command, args []string) {
		// scanImageInitial(args)
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// func scanImageInitial(image []string) {
// 	log.Printf("image called with args: %v\n", image)

// 	v1Image, err := util.PullImage(image)
// 	if err != nil {
// 		log.Fatalln("error while pulling image: ", err.Error())
// 	}

// 	tempDir, err := util.SaveAndUntarImage(v1Image, image)
// 	if err != nil {
// 		log.Fatalln("error while saving or untarring image: ", err.Error())
// 	}

// 	manifest, err := util.ReadImageManifest(tempDir)
// 	if err != nil {
// 		log.Fatalln("error while reading image manifest: ", err.Error())
// 	}

// 	extractLayer, err := util.ExtractLayer(tempDir, manifest)
// 	if err != nil {
// 		log.Fatalln("error while extracting layer: ", err.Error())
// 	}

// 	osRelease, err := util.IdentifyOsRelease(extractLayer)
// 	if err != nil {
// 		log.Fatalln("error while identifying OS release: ", err.Error())
// 	}

// 	pkgs, err := pkg.AnalyzePkg(osRelease, extractLayer)
// 	if err != nil {
// 		log.Fatalln("error while fetching package information: ", err.Error())
// 	}

// 	util.RemoveTempDir(tempDir)

// 	err = report.GenerateTableReport(pkgs)
// 	if err != nil {
// 		log.Fatalln("error while generating report: ", err.Error())
// 	}

// }
