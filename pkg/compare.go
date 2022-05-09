package pkg

import (
	"fmt"
	"log"

	"github.com/sgirdhar/sbom/util"

	levenshtein "github.com/ka-weihe/fast-levenshtein"
)

func ListComp(identifiedMap, readMap map[string]Package, toolName string) error {
	// package builder logic

	if len(identifiedMap) != len(readMap) {

		fmt.Println("Unequal number of components identified by tools")
	}

	identifiedDiff, readDiff := getDiffLists(identifiedMap, readMap)

	handleDiff(identifiedDiff, readDiff, identifiedMap, readMap, toolName)

	return nil
}

func getDiffLists(identifiedMap, readMap map[string]Package) ([]string, []string) {
	var identifiedDiff, readDiff []string

	for key := range identifiedMap {
		if _, exists := readMap[key]; !exists {
			// key does not exist in readMap
			identifiedDiff = append(identifiedDiff, key)
		}
	}

	for key := range readMap {
		if _, exists := identifiedMap[key]; !exists {
			// key does not exist in identifiedMap
			readDiff = append(readDiff, key)
		}
	}

	return identifiedDiff, readDiff
}

func handleDiff(identifiedDiff, readDiff []string, identifiedMap, readMap map[string]Package, toolName string) {

	// matching SBOM
	if len(identifiedDiff) == 0 && len(readDiff) == 0 {
		fmt.Println("Matching SBOM")
		fmt.Printf("Components identified by %v: %v\n", util.ApplicationName, len(identifiedMap))
		fmt.Printf("Components identified by %v: %v\n", toolName, len(readMap))
		return
	}

	// read list has extra components
	if len(identifiedDiff) == 0 && len(readDiff) != 0 {
		fmt.Printf("Components identified by %v: %v\n", util.ApplicationName, len(identifiedMap))
		fmt.Printf("Components identified by %v: %v\n", toolName, len(readMap))
		fmt.Println("-----------------------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Extra component(s) identified by %v\n", toolName)
		printComponents(readDiff, readMap)
		return
	}

	// identified list has extra components
	if len(identifiedDiff) != 0 && len(readDiff) == 0 {
		fmt.Printf("Components identified by %v: %v\n", util.ApplicationName, len(identifiedMap))
		fmt.Printf("Components identified by %v: %v\n", toolName, len(readMap))
		fmt.Println("-----------------------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Extra component(s) identified by %v\n", util.ApplicationName)
		printComponents(identifiedDiff, identifiedMap)
		return
	}

	// both lists have extra components - further investigation needed
	if len(identifiedDiff) != 0 && len(readDiff) != 0 {
		fmt.Printf("Components identified by %v: %v\n", util.ApplicationName, len(identifiedMap))
		fmt.Printf("Components identified by %v: %v\n", toolName, len(readMap))
		fmt.Println("-----------------------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Extra component(s) identified by %v\n", util.ApplicationName)
		printComponents(identifiedDiff, identifiedMap)
		fmt.Printf("Extra component(s) identified by %v\n", toolName)
		printComponents(readDiff, readMap)

		associationMap := guessAssociation(identifiedDiff, readDiff)
		log.Println(associationMap)

		if len(associationMap) != 0 && associationMap != nil {
			fmt.Println("Disputed component(s):")
			for key1, key2 := range associationMap {
				printDisputedComponents(identifiedMap[key1], readMap[key2], toolName)
			}
		}
		return
	}
	// default case
	fmt.Println("Comparison unsuccessful !")
}

func printComponents(readDiff []string, pkgMap map[string]Package) {
	for _, key := range readDiff {
		fmt.Println("-----------------------------------------------------------------------------------------------------------------------------")
		fmt.Println("Name: ", pkgMap[key].Name)
		fmt.Println("Version: ", pkgMap[key].Version)
		fmt.Println("Type: ", pkgMap[key].Type)
		fmt.Println("Purl: ", pkgMap[key].PURL)
		fmt.Println("-----------------------------------------------------------------------------------------------------------------------------")
	}
}

func guessAssociation(identifiedDiff, readDiff []string) map[string]string {
	var associationMap = make(map[string]string)
	for _, id := range identifiedDiff {
		for _, rd := range readDiff {
			distance := levenshtein.Distance(id, rd)
			log.Printf("The distance between %s and %s is %d.\n", id, rd, distance)
			if distance < util.AssociationThresholod {
				associationMap[id] = rd
			}
		}
	}
	return associationMap
}

func printDisputedComponents(identifiedPkg, readPkg Package, toolName string) {
	fmt.Println("-----------------------------------------------------------------------------------------------------------------------------")
	if identifiedPkg.Name == readPkg.Name {
		fmt.Println("Name: ", identifiedPkg.Name)
	} else {
		fmt.Printf("Name (%v): %v | Name (%v): %v\n", util.ApplicationName, identifiedPkg.Name, toolName, readPkg.Name)
	}
	if identifiedPkg.Version == readPkg.Version {
		fmt.Println("Version: ", identifiedPkg.Version)
	} else {
		fmt.Printf("Version (%v): %v | Version (%v): %v\n", util.ApplicationName, identifiedPkg.Version, toolName, readPkg.Version)
	}
	if identifiedPkg.Type == readPkg.Type {
		fmt.Println("Type: ", identifiedPkg.Type)
	} else {
		fmt.Printf("Type (%v): %v | Type (%v): %v\n", util.ApplicationName, identifiedPkg.Type, toolName, readPkg.Type)
	}
	if identifiedPkg.PURL == readPkg.PURL {
		fmt.Println("Purl: ", identifiedPkg.Name)
	} else {
		fmt.Printf("Purl (%v): %v | Purl (%v): %v\n", util.ApplicationName, identifiedPkg.PURL, toolName, readPkg.PURL)
	}
	fmt.Println("-----------------------------------------------------------------------------------------------------------------------------")
}
