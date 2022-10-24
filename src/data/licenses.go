package data

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
)

//go:embed spdx-licenses.json
var allLicenseBytes []byte

func GetSpdxLicenseIds() (licenseIds []string) {
	type LicenseFile struct {
		LicenseLastVersion string   `json:"licenseLastVersion"`
		ReleaseDate        string   `json:"releaseDate"`
		Licenses           []string `json:"licenses"`
	}

	// License data taken from https://github.com/spdx/license-list-data/blob/v3.18/json/licenses.json
	// Subsequently trimmed down to just the ids.
	var file LicenseFile
	var err = json.Unmarshal(allLicenseBytes, &file)
	if err != nil {
		fmt.Println("Unable to parse SPDX license data")
		os.Exit(1)
	}

	return file.Licenses
}
