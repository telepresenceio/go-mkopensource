package dependency

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/telepresenceio/go-mkopensource/pkg/dependencies"
	"github.com/telepresenceio/go-mkopensource/pkg/detectlicense"
	"github.com/telepresenceio/go-mkopensource/pkg/scanningerrors"
)

type NodeDependencies map[string]nodeDependency

type nodeDependency struct {
	Licenses       interface{} `json:"licenses"`
	Repository     string      `json:"repository"`
	DependencyPath string      `json:"dependencyPath"`
	Name           string      `json:"name"`
	Version        string      `json:"version"`
	Path           string      `json:"path"`
	URL            string      `json:"url"`
	LicenseFile    string      `json:"licenseFile"`
	LicenseText    string      `json:"licenseText"`
}

func (n *nodeDependency) licenses() (string, error) {
	if licenses, ok := n.Licenses.(string); ok {
		return licenses, nil
	}

	if licenseArray, ok := n.Licenses.([]interface{}); ok {
		var licenses []string
		for _, v := range licenseArray {
			license, ok := v.(string)
			if !ok {
				return "", fmt.Errorf("dependency '%s@%s' has an invalid license field: %#v", n.Name, n.Version, n.Licenses)
			}
			licenses = append(licenses, license)
		}

		return strings.Join(licenses, " AND "), nil
	}

	return "", fmt.Errorf("dependency '%s@%s' has an invalid license field: %v", n.Name, n.Version, n.Licenses)
}

func GetDependencyInformation(r io.Reader, licenseRestriction detectlicense.LicenseRestriction) (dependencyInfo dependencies.DependencyInfo, err error) {
	nodeDependencies := &NodeDependencies{}
	data, err := io.ReadAll(r)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, nodeDependencies)
	if err != nil {
		return
	}

	sortedDependencies := getSortedDependencies(nodeDependencies)

	dependencyInfo = dependencies.NewDependencyInfo()
	licErrs := []error{}

loop:
	for _, dependencyId := range sortedDependencies {
		nodeDependency := (*nodeDependencies)[dependencyId]

		dependency, dependencyErr := getDependencyDetails(nodeDependency, dependencyId)
		if dependencyErr != nil {
			licErrs = append(licErrs, dependencyErr)
			continue
		}

		for _, license := range dependency.Licenses {
			if licenseErr := dependencies.CheckLicenseRestrictions(*dependency, license, licenseRestriction); licenseErr != nil {
				licErrs = append(licErrs, licenseErr)
				continue loop
			}
		}

		dependencyInfo.Dependencies = append(dependencyInfo.Dependencies, *dependency)
	}

	if len(licErrs) > 0 {
		return dependencyInfo, scanningerrors.ExplainErrors(licErrs)
	}

	if err := dependencyInfo.UpdateLicenseList(); err != nil {
		return dependencyInfo, fmt.Errorf("could not generate list of license URLs for JavaScript: %w", err)
	}

	return dependencyInfo, err
}

func getDependencyDetails(nodeDependency nodeDependency, dependencyId string) (*dependencies.Dependency, error) {
	name, version := splitDependencyIdentifier(dependencyId)

	dependency := &dependencies.Dependency{
		Name:     name,
		Version:  version,
		Licenses: []string{},
	}

	allLicenses, err := getDependencyLicenses(dependencyId, nodeDependency)
	if err != nil {
		return nil, err
	}
	dependency.Licenses = allLicenses

	return dependency, nil
}

func getDependencyLicenses(dependencyId string, nodeDependency nodeDependency) ([]string, error) {
	licenseString, err := nodeDependency.licenses()
	if err != nil {
		return nil, err
	}

	if licenseString == "" {
		return nil, fmt.Errorf("dependency '%s@%s' is missing a license identifier", nodeDependency.Name, nodeDependency.Version)
	}

	parenthesisRe, err := regexp.Compile(`^\(|\)$`)
	if err != nil {
		return nil, err
	}
	licenseString = parenthesisRe.ReplaceAllString(licenseString, "")

	separatorRe, err := regexp.Compile(` OR | AND `)
	if err != nil {
		return nil, err
	}
	licenses := separatorRe.Split(licenseString, -1)

	allLicenses := []string{}
	for _, spdxId := range licenses {
		license, ok := detectlicense.SpdxIdentifiers[spdxId]
		if ok {
			allLicenses = append(allLicenses, license.Name)
			continue
		}

		licenses, ok := hardcodedJsDependencies[dependencyId]
		if ok {
			allLicenses = licenses
			break
		}

		return nil, fmt.Errorf("dependency '%s@%s' has an unknown SPDX Identifier '%s'",
			nodeDependency.Name, nodeDependency.Version, spdxId)
	}

	sort.Strings(allLicenses)
	return allLicenses, nil
}

func getSortedDependencies(nodeDependencies *NodeDependencies) []string {
	sortedDependencies := make([]string, 0, len(*nodeDependencies))
	for k := range *nodeDependencies {
		sortedDependencies = append(sortedDependencies, k)
	}
	sort.Strings(sortedDependencies)
	return sortedDependencies
}

func splitDependencyIdentifier(identifier string) (name string, version string) {
	parts := strings.Split(identifier, "@")

	numberOfParts := len(parts)
	return strings.Join(parts[:numberOfParts-1], "@"), parts[numberOfParts-1]
}
