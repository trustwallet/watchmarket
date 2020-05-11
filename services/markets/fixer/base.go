package fixer

type Provider struct {
	ID, currency string
	client       Client
}

func InitProvider(api, key, currency string) Provider {
	return Provider{
		ID:       id,
		currency: currency,
		client:   NewClient(api, key, currency),
	}
}
