package pagination

type PageOption struct {
	Size    int         `json:"size,omitempty"`
	Skip    int         `json:"skip,omitempty"`
	Total   int         `json:"total,omitempty"`
	PageNum int         `json:"page_num,omitempty"`
	LastId  interface{} `json:"last_id,omitempty"`
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
		LastId: LastId,
	}
}

func NewNumPage(pageNum, size int) *PageOption {
	p := PageOption{}
	p.Size = 20
	p.PageNum = pageNum
	if p.PageNum == 0 {
		p.PageNum = 1
	}
	if size > 0 {
		p.Size = size
	}
	p.Skip = p.Size * (p.PageNum - 1)
	return &p
}

func (page *PageOption) Offset() int {
	return page.Skip * page.Size
}

func (page *PageOption) IsOver() bool {
	return page.Offset()+page.Size > page.Total
}
