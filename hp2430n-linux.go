//go:build !tinygo

package hp2430n

import "embed"

//go:embed css go.mod html images js template
var fs embed.FS
