package doom

type DoomVictim struct {
	Path             string `json:"path"`
	Name             string `json:"name"`
	Ext              string `json:"ext"`
	LastModifiedUnix int64  `json:"last_modified_unix"`
}
