package pagination

type PageOption struct {
	Size    int
	Skip    int
	Total   int
	PageNum int
	Order   string
	LastId  interface{}
}

func NewPageOption(size, skip int, order string, LastId interface{}) *PageOption {
	if size == 0 {
		size = 20
	}
	if skip < 0 {
		skip = 0
	}
	return &PageOption{
		Size:   size,
		Skip:   skip,
		Total:  0,
		Order:  order,
		LastId: LastId,
	}
}

func NewNumPage(pageNum, size int) *PageOption {
	p := PageOption{}
	p.Size = 10
	p.PageNum = pageNum
	if p.PageNum == 0 {
		p.PageNum = 1
	}
	if size > 0 {
		p.Size = size
	}
	p.Order = "ctime"
	p.Skip = p.Size * (p.PageNum - 1)
	return &p
}

func (page *PageOption) Offset() int {
	return page.Skip * page.Size
}

func (page *PageOption) IsOver() bool {
	return page.Offset()+page.Size > page.Total
}
