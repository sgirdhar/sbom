package util

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func CreateTempDir(imageName string) (string, error) {
	dir, err := os.MkdirTemp("", imageName)
	log.Printf("***** Creating the temp directory: %v *****\n", dir)

	if err != nil {
		log.Println("error while creating directory")
		return "", err
	}
	return dir, nil
}

func RemoveDir(dir string) {
	log.Printf("***** Removing directory: %v *****\n", dir)
	os.RemoveAll(dir) // clean up
}

func CreateDir(dir string) error {
	log.Printf("***** Creating directory: %v *****\n", dir)
	err := os.MkdirAll(dir, 0750)
	if err != nil && !os.IsExist(err) {
		log.Println("error while creating directory")
		return err
	}
	return nil
}

func GetCurrentWorkDir() (string, error) {
	cmd := exec.Command("pwd")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Println("running command: ", cmd)
	cmd.Run()

	if cmd.ProcessState.ExitCode() != 0 {
		// log.Fatalf("could not execute command:%v, error:%v", cmd, stderr.String())
		return "", stderr.UnreadByte()
	} else {
		log.Printf("Returning current working directory: %v", stdout.String())
		return stdout.String(), nil
	}
}

func Untar(tarball, target string) error {
	log.Printf("Untarring %v at\n %v\n", tarball, target)
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

func UnGzip(source, target string) error {
	log.Printf("Unzipping %v", source)
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}

func ReadFile(filePath string) ([]string, error) {
	var txtlines []string
	file, err := os.Open(filePath)

	if err != nil {
		log.Printf("error while opening %v\n", filePath)
		return txtlines, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	return txtlines, nil
}

func CheckFile(name string) error {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		log.Printf("%v does not exist.\n", name)
	}
	if err != nil {
		return err
	} else {
		return nil
	}
}
