package report

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/sgirdhar/sbom/pkg"
)

func GenerateTableReport(pkgs []pkg.Package) error {
	// second cell of each line, belong to different columns.
	log.Println("Generating tabular report")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	// fmt.Fprintln(w, "NAME\tVERSION\tARCHITECTURE\tPURL")
	fmt.Fprintln(w, "NAME\tVERSION\tARCHITECTURE\tTYPE")
	for _, pkg := range pkgs {
		fmt.Fprintln(w, pkg.Name+"\t"+pkg.Version+"\t"+pkg.Architecture+"\t"+pkg.Type)
	}
	w.Flush()

	return nil
}
