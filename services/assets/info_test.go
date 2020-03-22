package assets

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/magiconair/properties/assert"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"testing"
)

func Test_GetCoinInfo(t *testing.T) {
	type args struct {
		c     coin.Coin
		token string
	}

	httpClient := resty.New()
	httpmock.ActivateNonDefault(httpClient.GetClient())
	client := HttpAssetClient{HttpClient: httpClient}

	mockEthCoinInfo := watchmarket.CoinInfo{
		Name:             "Ethereum",
		Website:          "https://ethereum.org/",
		SourceCode:       "https://github.com/ethereum",
		WhitePaper:       "https://github.com/ethereum/wiki/wiki/White-Paper",
		Description:      "Open source platform to write and distribute decentralized applications.",
		ShortDescription: "Open source platform to write and distribute decentralized applications.",
		Explorer:         "https://etherscan.io/",
		Socials:          []watchmarket.SocialLink{{Name: "Twitter", Url: "https://twitter.com/ethereum", Handle: "ethereum"}},
	}

	ethResponder, err := httpmock.NewJsonResponder(200, mockEthCoinInfo)
	if err != nil {
		t.Fatal(err)
	}
	httpmock.RegisterResponder("GET", "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/bitcoin/info/info.json", httpmock.NewStringResponder(200, "Bad data"))
	httpmock.RegisterResponder("GET", "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/info/info.json", ethResponder)
	httpmock.RegisterResponder("GET", "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/binance/info/info.json", httpmock.NewStringResponder(404, "Not Found"))
	httpmock.RegisterResponder("GET", "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/vechain/info/info.json", httpmock.NewStringResponder(400, "Boom!"))

	tests := []struct {
		name     string
		args     args
		wantErr  error
		wantInfo watchmarket.CoinInfo
	}{
		{"test bad response from cmc", args{coin.Bitcoin(), ""}, errors.New("Failed to unmarshal invalid character 'B' looking for beginning of value"), watchmarket.CoinInfo{}},
		{"test nominal", args{coin.Ethereum(), ""}, nil, mockEthCoinInfo},
		{"test not found", args{coin.Binance(), ""}, watchmarket.ErrNotFound, watchmarket.CoinInfo{}},
		{"test not found", args{coin.Vechain(), ""}, errors.New("Request to https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/vechain/info/info.json failed with HTTP 400: Boom!"), watchmarket.CoinInfo{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := client.GetCoinInfo(int(tt.args.c.ID), tt.args.token)
			if err != nil {
				assert.Equal(t, err, tt.wantErr)
			} else {
				assert.Equal(t, *info, tt.wantInfo)
			}
		})
	}
}

func Test_getCoinInfoUrl(t *testing.T) {
	type args struct {
		c     coin.Coin
		token string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test Ethereum coin", args{coin.Ethereum(), ""}, AssetsURL + coin.Ethereum().Handle + "/info"},
		{"test Ethereum token", args{coin.Ethereum(), "0x0000000000b3F879cb30FE243b4Dfee438691c04"}, AssetsURL + coin.Ethereum().Handle + "/assets/" + "0x0000000000b3F879cb30FE243b4Dfee438691c04"},
		{"test Binance coin", args{coin.Binance(), ""}, AssetsURL + coin.Binance().Handle + "/info"},
		{"test Binance token", args{coin.Binance(), "BUSD-BD1"}, AssetsURL + coin.Binance().Handle + "/assets/" + "BUSD-BD1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCoinInfoUrl(tt.args.c, tt.args.token); got != tt.want {
				t.Errorf("getCoinInfoUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
