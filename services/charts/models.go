package charts

type (
	Charts struct {
		Prices     []ChartVolume `json:"prices"`
		MarketCaps []ChartVolume `json:"market_caps"`
		Volumes    []ChartVolume `json:"total_volumes"`
	}

	ChartVolume []float64

	CoinType string

	Data struct {
		Prices []Price `json:"prices,omitempty"`
		Error  string  `json:"error,omitempty"`
	}

	Price struct {
		Price float64 `json:"price"`
		Date  int64   `json:"date"`
	}

	CoinDetails struct {
		Vol24             float64 `json:"volume_24"`
		MarketCap         float64 `json:"market_cap"`
		CirculatingSupply float64 `json:"circulating_supply"`
		TotalSupply       float64 `json:"total_supply"`
		Info              Info    `json:"info,omitempty"`
	}

	Info struct {
		Name             string       `json:"name,omitempty"`
		Website          string       `json:"website,omitempty"`
		SourceCode       string       `json:"source_code,omitempty"`
		WhitePaper       string       `json:"white_paper,omitempty"`
		Description      string       `json:"description,omitempty"`
		ShortDescription string       `json:"short_description,omitempty"`
		Explorer         string       `json:"explorer,omitempty"`
		Socials          []SocialLink `json:"socials,omitempty"`
	}

	SocialLink struct {
		Name   string `json:"name"`
		Url    string `json:"url"`
		Handle string `json:"handle"`
	}
)

func (d Data) IsEmpty() bool {
	return len(d.Prices) == 0
}

func (i CoinDetails) IsEmpty() bool {
	return i.Info.Name == ""
}

const (
	Token CoinType = "token"
	Coin  CoinType = "coin"
)
