package assets

import _ "embed"

//go:embed events.json
var EventsJSON []byte

//go:embed holidays.json
var HolidaysJSON []byte
