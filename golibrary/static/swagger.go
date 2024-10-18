package swagger

import "embed"

//go:embed swagger.json
var Swaggerfile embed.FS
