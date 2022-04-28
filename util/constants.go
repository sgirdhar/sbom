package util

import "log"

const (
	ApplicationName    = "sbom"
	ApplicationVersion = "0.3"
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

// termBanner := `ICAgICAgICAgICAgICAgXyAgICAgIAogX19fIF8gX18gICBfX3wgfF8gIF9fCi8gX198ICdfIFwg
// LyBfYCBcIFwvIC8KXF9fIFwgfF8pIHwgKF98IHw+ICA8IAp8X19fLyAuX18vIFxfXyxfL18vXF9c
// CiAgICB8X3wgICAgICAgICAgICAgICAK`
// 		d, err := base64.StdEncoding.DecodeString(termBanner)
// 		if err != nil {
// 			log.Println("error while testing: ", err.Error())
// 		}
// 		log.Println(d)
