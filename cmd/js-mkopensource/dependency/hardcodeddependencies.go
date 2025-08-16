package dependency

import license "github.com/telepresenceio/go-mkopensource/pkg/detectlicense"

//nolint:gochecknoglobals // Would be 'const'.
var hardcodedJsDependencies = map[string][]string{
	"cyclist@0.2.2":                {license.MIT.Name},
	"doctrine@1.5.0":               {license.BSD2.Name, license.Apache2.Name},
	"emitter-component@1.1.1":      {license.MIT.Name},
	"flexboxgrid@6.3.1":            {license.Apache2.Name},
	"indexof@0.0.1":                {license.MIT.Name},
	"intro.js@4.1.0":               {license.AGPL3OrLater.Name},
	"json-schema@0.2.3":            {license.AFL21.Name},
	"node-forge@0.10.0":            {license.BSD3.Name},
	"pako@1.0.10":                  {license.MIT.Name},
	"regenerator-transform@0.10.1": {license.BSD2.Name},
	"regjsparser@0.1.5":            {license.BSD2.Name},
	"node-forge@1.3.1":             {license.BSD3.Name},
}
