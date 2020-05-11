package binancedex

type Provider struct {
	ID     string
	client Client
}

func InitProvider(api string) Provider {
	m := Provider{
		ID:     id,
		client: NewClient(api),
	}
	return m
}
