package data_test

import (
	"ahkpm/src/data"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSpdxLicenseIds(t *testing.T) {
	licenseIds := data.GetSpdxLicenseIds()
	assert.IsType(t, []string{}, licenseIds)
	assert.Contains(t, licenseIds, "MIT")
}
