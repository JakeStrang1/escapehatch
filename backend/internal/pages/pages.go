package pages

type PageResult struct {
	Page       int
	PerPage    int
	TotalPages int
}

func (p *PageResult) HasMore() bool {
	return p.Page < p.TotalPages
}
