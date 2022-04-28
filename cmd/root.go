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
for container images. It takes image string or image tar as input.`,
	Example: `  Input image string
	sbom scan alpine:latest
	sbom scan alpine@sha256:51103b3f2993cbc1b45ff9d941b5d461484002792e02aa29580ec5282de719d4

  Input tarfile
	sbom scan --tar /path/to/tarfile/alpine.tar`,
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
