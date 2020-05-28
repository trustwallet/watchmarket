package fixer

type Provider struct {
	id, currency string
	client       Client
}

func InitProvider(api, key, currency string) Provider {
	return Provider{id: id, currency: currency, client: NewClient(api, key, currency)}
}

func (p Provider) GetProvider() string {
	return p.id
}
