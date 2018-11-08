package da

type InfoList struct {
	Pagination string   `json:"paginationToken"`
	Data       []string `json:"data"`
}
