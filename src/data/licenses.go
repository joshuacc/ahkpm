package data

import (
	_ "embed"
	"encoding/json"
)

//go:embed spdx-licenses.json
var allLicenseBytes []byte

func GetSpdxLicenseIds() (licenseIds []string) {
	type License struct {
		Reference             string   `json:"reference"`
		IsDeprecatedLicenseID bool     `json:"isDeprecatedLicenseId"`
		DetailsURL            string   `json:"detailsUrl"`
		ReferenceNumber       int      `json:"referenceNumber"`
		Name                  string   `json:"name"`
		LicenseId             string   `json:"licenseId"`
		SeeAlso               []string `json:"seeAlso"`
		IsOsiApproved         bool     `json:"isOsiApproved"`
	}

	type LicenseFile struct {
		LicenseLastVersion string `json:"licenseLastVersion"`
		ReleaseDate        string `json:"releaseDate"`
		Licenses           []License
	}
	var allLicenses []string

	// License data taken from https://github.com/spdx/license-list-data/blob/v3.18/json/licenses.json
	var file LicenseFile
	var err = json.Unmarshal(allLicenseBytes, &file)
	if err != nil {
		panic(err)
	}

	for _, license := range file.Licenses {
		allLicenses = append(allLicenses, license.LicenseId)
	}

	return allLicenses
}
