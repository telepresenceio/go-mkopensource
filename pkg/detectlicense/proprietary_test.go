package detectlicense

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAmbassadorProprietarySoftware(t *testing.T) {
	proprietarySoftware := GetProprietarySoftware("github.com/datawire/secretprogram")
	require.Len(t, proprietarySoftware, 1)
	require.Contains(t, proprietarySoftware, "github.com/datawire/secretprogram")
}

func TestReadProprietarySoftwareFile(t *testing.T) {
	proprietarySoftware := GetProprietarySoftware()

	err := proprietarySoftware.ReadProprietarySoftwareFile("./testdata/proprietary_software.yaml")

	require.NoError(t, err)
	require.Len(t, proprietarySoftware, 2)
	require.Contains(t, proprietarySoftware, "github.com/datawire/secretprogram")
	require.Contains(t, proprietarySoftware, "github.com/datawire/othersecretprogram")
}
