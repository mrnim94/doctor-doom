package doom

type DoomVictim struct {
	Path             string `json:"path"`
	Name             string `json:"name"`
	Size             int64  `json:"size"`
	LastModifiedUnix int64  `json:"last_modified_unix"`
	LiveIn           string `json:"live_in"`
}
