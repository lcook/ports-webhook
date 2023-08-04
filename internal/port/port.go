package port

type Port struct {
	Category string
	Name     string
}

func (p *Port) Fullname() string {
	return p.Category + "/" + p.Name
}
