package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

type Manifest struct {
	Config   string
	RepoTags []string
	Layers   []string
}

func PullImage(image string) (v1.Image, error) {

	log.Println("pulling image: ", image)

	v1Image, err := crane.Pull(image)
	if err != nil {
		log.Println("error while pulling image: ", err.Error())
		return nil, err
	}
	return v1Image, nil
}

func SaveAndUntarImage(v1Image v1.Image, image string) (string, error) {

	imageName, _ := GetImageNameAndVersion(image)
	tarFile := imageName + ".tar"

	tempDir, err := CreateTempDir(imageName)
	if err != nil {
		log.Println("error while creating temp directory: ", err.Error())
		return "", err
	}

	imageLocation := tempDir + "/" + tarFile

	err = crane.Save(v1Image, image, imageLocation)
	if err != nil {
		log.Println("error while saving image: ", err.Error())
		return "", err
	}

	err = Untar(imageLocation, tempDir)
	if err != nil {
		log.Println("error while untarring image: ", err.Error())
		return "", err
	}
	return tempDir, nil
}

func GetImageNameAndVersion(image string) (string, string) {
	if strings.ContainsAny(image, ":") && strings.ContainsAny(image, "@") {
		return strings.Split(image, "@")[0], strings.Split(image, "@")[1]
	}
	if strings.ContainsAny(image, ":") {
		return strings.Split(image, ":")[0], strings.Split(image, ":")[1]
	}
	return image, "latest"
}

func UntarImage(tarPath string) (string, error) {
	Slice := strings.Split(tarPath, "/")
	image := Slice[len(Slice)-1]
	tarDir := strings.Join(Slice[:len(Slice)-1], "/")
	log.Printf("tarFile:%v, tarDir:%v", image, tarDir)

	tempDir := tarDir + "/" + strings.Split(image, ".")[0]

	err := CreateDir(tempDir)
	if err != nil {
		log.Println("error while creating directory: ", err.Error())
		return "", err
	}
	err = Untar(tarPath, tempDir)
	if err != nil {
		log.Println("error while untarring image: ", err.Error())
		return "", err
	}
	return tempDir, nil
}

func ReadImageManifest(tempDir string) (Manifest, error) {

	// Let's first read the `manifest.json` file
	manifestLocation := tempDir + "/manifest.json"
	log.Println("Reading manifest: ", manifestLocation)

	content, err := ioutil.ReadFile(manifestLocation)
	if err != nil {
		log.Fatal("Error when opening file: ", err.Error())
	}

	// Now let's unmarshall the data into `payload`
	var payload []Manifest
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err.Error())
	}

	manifest := payload[0]
	return manifest, nil
}

func ReadImageConfig(tempDir string, manifest Manifest) (v1.ConfigFile, error) {
	configLocation := tempDir + "/" + manifest.Config
	log.Println("Reading image config: ", configLocation)

	content, err := ioutil.ReadFile(configLocation)
	if err != nil {
		log.Fatal("Error when opening file: ", err.Error())
	}

	// Now let's unmarshall the data into `payload`
	var configFile v1.ConfigFile
	err = json.Unmarshal(content, &configFile)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err.Error())
	}

	return configFile, nil
}

func ExtractLayer(tempDir string, manifest Manifest) (string, error) {
	source := tempDir + "/" + manifest.Layers[0]
	target := tempDir + "/" + strings.Split(manifest.Layers[0], ".")[0] + "." + strings.Split(manifest.Layers[0], ".")[1]

	if strings.Contains(source, ".gz") {
		err := UnGzip(source, target)
		if err != nil {
			log.Println("error while unzipping image: ", err.Error())
			return "", err
		}
	}

	newDir := tempDir + "/" + strings.Split(manifest.Layers[0], ".")[0]
	err := CreateDir(newDir)
	if err != nil {
		log.Println("error while creating directory: ", err.Error())
		return "", err
	}

	err = Untar(target, newDir)
	if err != nil {
		log.Println("error while untarring image: ", err.Error())
		return "", err
	}
	return newDir, nil
}
