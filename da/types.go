package da

type InfoList struct {
	Pagination string   `json:"paginationToken"`
	Data       []string `json:"data"`
}

type AliasesList struct {
	Pagination string  `json:"paginationToken"`
	Data       []Alias `json:"data"`
}

type VersionList struct {
	Pagination string `json:"paginationToken"`
	Data       []uint `json:"data"`
}

type Alias struct {
	ID      string `json:"id"`
	Version uint   `json:"version"`
}
