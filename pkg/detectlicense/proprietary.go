package detectlicense

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ProprietarySoftware map[string]struct{}

func GetProprietarySoftware(proprietarySoftware ...string) ProprietarySoftware {
	ambProprietarySoftware := make(ProprietarySoftware, len(proprietarySoftware))
	for _, v := range proprietarySoftware {
		ambProprietarySoftware[v] = struct{}{}
	}
	return ambProprietarySoftware
}

func (a ProprietarySoftware) IsProprietarySoftware(packageName string) bool {
	_, ok := a[packageName]
	return ok
}

func (a ProprietarySoftware) ReadProprietarySoftwareFile(name string) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	var proprietarySoftware []string
	if err = yaml.Unmarshal(data, &proprietarySoftware); err != nil {
		return err
	}

	for _, v := range proprietarySoftware {
		a[v] = struct{}{}
	}
	return nil
}
