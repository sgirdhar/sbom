package util

import (
	"errors"
)

func ListComp(identifiedPkgs, readPkgs []string) error {
	// to do : package builder logic

	if len(identifiedPkgs) != len(readPkgs) {
		return errors.New("length of packages list unequal")
	}
	return nil
}
