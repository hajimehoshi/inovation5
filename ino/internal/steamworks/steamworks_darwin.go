package steamworks

import (
	_ "embed"
)

//go:embed libsteam_api.dylib
var libSteamAPI []byte
