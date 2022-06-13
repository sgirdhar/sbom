package cmd

import (
	"os"

	"github.com/sgirdhar/sbom/util"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sbom",
	Short: "Generate SBOM for container images",
	Long: `Sbom is a CLI tool to generate the Software Bill of Material (SBOM) 
for container images. It takes image string or image tarball as input.`,
	Example: `  Possible Inputs
	sbom scan my_image:my_tag			Get the image from dockerHub. If no tag specified, default is latest.
	sbom scan my_image@my_digest			Get the image from dockerHub. No default value for digest.
	sbom scan --tar /path/to/tarfile/my_image.tar	Docker tar or OCI tar.

  Output Formats
	sbom scan my_image:my_tag				Default output is human readable summary table
	sbom scan my_image:my_tag --output cyclonedx		Generates CycloneDX xml output
	sbom scan my_image:my_tag --output cyclonedx-json	Generates CycloneDX json output
	
  Compare SBOM
	sbom scan my_image:my_tag --compare /path/to/cyclonedx-json/my_sbom.json	Generates SBOM and compares with CycloneDX json file provided`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Version: util.ApplicationVersion,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("argument: %v\n", args)
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sbom.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
