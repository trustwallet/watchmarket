package binancedex

type Provider struct {
	id     string
	client Client
}

func InitProvider(api string) Provider {
	m := Provider{
		id:     id,
		client: NewClient(api),
	}
	return m
}

func (p Provider) GetProvider() string {
	return p.id
}
