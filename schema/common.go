package schema

type SortType string

const (
	ASC  SortType = "ASC"
	DESC SortType = "DESC"
)

type Page struct {
	Offset  int
	Limit   int
	OrderBy string
	SortBy  SortType
}

func (p *Page) LoadDefault() {
	if p.Offset == 0 {
		p.Offset = 0
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.OrderBy == "" {
		p.OrderBy = "id"
	}
	if p.SortBy != ASC && p.SortBy != DESC {
		p.SortBy = ASC
	}
}
