package report

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/sgirdhar/sbom/pkg"
)

// func GenerateTableReport(pkgs []pkg.Package) error {
// 	log.Println("Generating tabular report")
// 	var (
// 		rows   [][]string
// 		output io.Writer
// 	)

// 	columns := []string{"Name", "Version", "Type"}
// 	for _, p := range pkgs {
// 		row := []string{
// 			p.Name,
// 			p.Version,
// 			p.Type,
// 		}
// 		rows = append(rows, row)
// 	}

// 	if len(rows) == 0 {
// 		_, err := fmt.Fprintln(output, "No packages discovered")
// 		return err
// 	}

// 	// sort by name, version, then type
// 	sort.SliceStable(rows, func(i, j int) bool {
// 		for col := 0; col < len(columns); col++ {
// 			if rows[i][col] != rows[j][col] {
// 				return rows[i][col] < rows[j][col]
// 			}
// 		}
// 		return false
// 	})

// 	table := tablewriter.NewWriter(output)

// 	table.SetHeader(columns)
// 	table.SetHeaderLine(false)
// 	table.SetBorder(false)
// 	table.SetAutoWrapText(false)
// 	table.SetAutoFormatHeaders(true)
// 	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
// 	table.SetAlignment(tablewriter.ALIGN_LEFT)
// 	table.SetCenterSeparator("")
// 	table.SetColumnSeparator("")
// 	table.SetRowSeparator("")
// 	table.SetTablePadding("  ")
// 	table.SetNoWhiteSpace(true)

// 	table.AppendBulk(rows)
// 	table.Render()

// 	return nil
// }

func GenerateTableReport(pkgs []pkg.Package) error {
	// Observe how the b's and the d's, despite appearing in the
	// second cell of each line, belong to different columns.
	log.Println("Generating tabular report")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	// fmt.Fprintln(w, "------------------\t\t------------------\t\t------------------\t\t------------------\t\t------------------\t\t------------------")
	fmt.Fprintln(w, "NAME\tVERSION\tARCHITECTURE\tPURL")
	for _, pkg := range pkgs {
		fmt.Fprintln(w, pkg.Name+"\t"+pkg.Version+"\t"+pkg.Architecture+"\t"+pkg.PURL)
	}
	// fmt.Fprintln(w, "------------------\t\t------------------\t\t------------------\t\t------------------\t\t------------------\t\t------------------")
	w.Flush()

	return nil
}
