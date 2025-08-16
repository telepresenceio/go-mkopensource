package dependencies

import (
	"encoding/json"
	"fmt"

	license "github.com/telepresenceio/go-mkopensource/pkg/detectlicense"
)

//nolint:gochecknoglobals // Can't be a constant
var licensesByName = map[string]license.License{
	license.ZeroBSD.Name:       license.ZeroBSD,
	license.Apache2.Name:       license.Apache2,
	license.AFL21.Name:         license.AFL21,
	license.AGPL1Only.Name:     license.AGPL1Only,
	license.AGPL1OrLater.Name:  license.AGPL1OrLater,
	license.AGPL3Only.Name:     license.AGPL3Only,
	license.AGPL3OrLater.Name:  license.AGPL3OrLater,
	license.BSD1.Name:          license.BSD1,
	license.BSD2.Name:          license.BSD2,
	license.BSD3.Name:          license.BSD3,
	license.Cc010.Name:         license.Cc010,
	license.CcBy30.Name:        license.CcBy30,
	license.CcBy40.Name:        license.CcBy40,
	license.CcBySa40.Name:      license.CcBySa40,
	license.EPL10.Name:         license.EPL10,
	license.GPL1Only.Name:      license.GPL1Only,
	license.GPL1OrLater.Name:   license.GPL1OrLater,
	license.GPL2Only.Name:      license.GPL2Only,
	license.GPL2OrLater.Name:   license.GPL2OrLater,
	license.GPL3Only.Name:      license.GPL3Only,
	license.GPL3OrLater.Name:   license.GPL3OrLater,
	license.ISC.Name:           license.ISC,
	license.LGPL2Only.Name:     license.LGPL2Only,
	license.LGPL2OrLater.Name:  license.LGPL2OrLater,
	license.LGPL21Only.Name:    license.LGPL21Only,
	license.LGPL21OrLater.Name: license.LGPL21OrLater,
	license.LGPL3Only.Name:     license.LGPL3Only,
	license.LGPL3OrLater.Name:  license.LGPL3OrLater,
	license.MIT.Name:           license.MIT,
	license.MPL11.Name:         license.MPL11,
	license.MPL2.Name:          license.MPL2,
	license.ODCBy10.Name:       license.ODCBy10,
	license.OFL11.Name:         license.OFL11,
	license.PSF.Name:           license.PSF,
	license.Python20.Name:      license.Python20,
	license.PublicDomain.Name:  license.PublicDomain,
	license.Unicode2015.Name:   license.Unicode2015,
	license.Unlicense.Name:     license.Unlicense,
	license.WTFPL.Name:         license.WTFPL,
}

type DependencyInfo struct {
	Dependencies []Dependency      `json:"dependencies"`
	Licenses     map[string]string `json:"licenseInfo"`
}

type Dependency struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Licenses []string `json:"licenses"`
}

func NewDependencyInfo() DependencyInfo {
	return DependencyInfo{
		Dependencies: []Dependency{},
		Licenses:     map[string]string{},
	}
}

func (d *DependencyInfo) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, d); err != nil {
		return err
	}

	return nil
}

func (d *DependencyInfo) UpdateLicenseList() error {
	usedLicenses := map[string]license.License{}

	for _, dependency := range d.Dependencies {
		for _, licenseName := range dependency.Licenses {
			license, err := getLicenseFromName(licenseName)
			if err != nil {
				return err
			}
			usedLicenses[license.Name] = license
		}
	}

	for k, v := range usedLicenses {
		d.Licenses[k] = v.URL
	}

	return nil
}

func getLicenseFromName(licenseName string) (license.License, error) {
	lic, ok := licensesByName[licenseName]
	if !ok {
		return license.License{}, fmt.Errorf("license details for '%s' are not known", licenseName)
	}
	return lic, nil
}

func CheckLicenseRestrictions(dependency Dependency, licenseName string, licenseRestriction license.LicenseRestriction) error {
	lic, err := getLicenseFromName(licenseName)
	if err != nil {
		return err
	}

	if lic.Restriction == license.Forbidden {
		return fmt.Errorf("dependency '%s@%s' uses license '%s' which is forbidden", dependency.Name,
			dependency.Version, lic.Name)
	}

	if lic.Restriction < licenseRestriction {
		return fmt.Errorf("dependency '%s@%s' uses license '%s' which is not allowed on applications that run on customer machines",
			dependency.Name, dependency.Version, lic.Name)
	}
	return nil
}
