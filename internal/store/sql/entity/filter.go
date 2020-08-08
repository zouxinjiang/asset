package entity

type (
	Paging struct {
		PageNumber int
		PageSize   int
	}

	Cond     string
	CondItem struct {
		OuterName string
		Cond      Cond
		Value     interface{}
	}

	Filter struct {
		Paging
		AndCond []CondItem
		OrCond  []CondItem
		Odr     []map[string]bool
	}
)

const (
	Eq        Cond = "="
	Ne        Cond = "!="
	Gt        Cond = ">"
	Ge        Cond = ">="
	Lt        Cond = "<"
	Le        Cond = "<="
	Like      Cond = "LIKE"
	NotLike   Cond = "NOT LIKE"
	IsNull    Cond = "IS NULL"
	IsNotNull Cond = "IS NOT NULL"
	In        Cond = "IN"
	NotIn     Cond = "NOT IN"
)

// Offset ...
func (p Paging) Offset() int {
	if p.PageNumber > 0 {
		return (p.PageNumber - 1) * p.PageSize
	}
	return 0
}

// Limit ...
func (p Paging) Limit() int {
	return p.PageSize
}

// And add a and cond
func (f Filter) And(outName string, cond Cond, value interface{}) Filter {
	ret := f.clone()
	ret.AndCond = append(f.AndCond, CondItem{OuterName: outName, Cond: cond, Value: value})
	return ret
}

// Or add a or cond
func (f Filter) Or(outName string, cond Cond, value interface{}) Filter {
	ret := f.clone()
	ret.OrCond = append(ret.OrCond, CondItem{OuterName: outName, Cond: cond, Value: value})
	return ret
}

// Order add a order cond
func (f Filter) Order(key string, sort bool) Filter {
	ret := f.clone()
	ret.Odr = append(ret.Odr, map[string]bool{key: sort})
	return ret
}

func (f Filter) clone() Filter {
	var ret = Filter{
		AndCond: make([]CondItem, len(f.AndCond)),
		OrCond:  make([]CondItem, len(f.OrCond)),
		Odr:     make([]map[string]bool, len(f.Odr)),
	}
	copy(ret.AndCond, f.AndCond)
	copy(ret.OrCond, f.OrCond)
	copy(ret.Odr, f.Odr)
	if f.Odr != nil {
		ret.Odr = f.Odr
	}
	ret.Paging = f.Paging

	return ret
}

func (f Filter) SetPageSize(pageSize int) Filter {
	ret := f.clone()
	ret.Paging.PageSize = pageSize
	return ret
}

func (f Filter) SetPage(pageNumber int) Filter {
	ret := f.clone()
	ret.Paging.PageNumber = pageNumber
	return ret
}
