package util

import (
	"log"

	"github.com/acobaugh/osrelease"
)

type OsRelease struct {
	NAME           string
	ID             string
	VERSION_ID     string
	PRETTY_NAME    string
	HOME_URL       string
	BUG_REPORT_URL string
}

func IdentifyOsRelease(extractDir string) (OsRelease, error) {
	const etcOsRelease = "/etc/os-release"
	const usrLibOsRelease = "/usr/lib/os-release"
	var osRelease OsRelease

	osReleaseMap, err := osrelease.ReadFile(extractDir + etcOsRelease)
	if err != nil || len(osReleaseMap) == 0 {
		osReleaseMap, err = osrelease.ReadFile(extractDir + usrLibOsRelease)
	}
	if err != nil {
		log.Println("error while reading os-release file: ", err.Error())
		return osRelease, err
	}

	osRelease = OsRelease{
		NAME:           osReleaseMap["NAME"],
		ID:             osReleaseMap["ID"],
		VERSION_ID:     osReleaseMap["VERSION_ID"],
		PRETTY_NAME:    osReleaseMap["PRETTY_NAME"],
		HOME_URL:       osReleaseMap["HOME_URL"],
		BUG_REPORT_URL: osReleaseMap["BUG_REPORT_URL"],
	}

	log.Println("Identified linux distro: ", osRelease.PRETTY_NAME)
	return osRelease, err
}
