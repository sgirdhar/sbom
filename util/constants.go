package util

import "log"

const (
	ApplicationName       = "sbom"
	ApplicationVersion    = "0.5"
	AssociationThresholod = 3
)

func PrintSbom() {
	log.Println("|----|  |-----|  |-----|  |----|----|")
	log.Println("|       |     |  |     |  |    |    |")
	log.Println("|       |     |  |     |  |    |    |")
	log.Println("|----|  |-----|  |     |  |         |")
	log.Println("     |  |     |  |     |  |         |")
	log.Println("     |  |     |  |     |  |         |")
	log.Println("|----|  |-----|  |-----|  |         |")
}
