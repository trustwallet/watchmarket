package coingecko

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/watchmarket/services/assets"
	"net/http"
	"testing"
)

func TestInitProvider(t *testing.T) {
	provider := InitProvider("web.api", "USD", assets.Init("assets.api"))
	assert.NotNil(t, provider)
	assert.Equal(t, "web.api", provider.client.baseURL)
	assert.Equal(t, "USD", provider.currency)
	assert.Equal(t, "coingecko", provider.id)
}

func TestProvider_GetProvider(t *testing.T) {
	provider := InitProvider("web.api", "USD", assets.Init("assets.api"))
	assert.Equal(t, "coingecko", provider.GetProvider())
}

func createMockedAPI() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/v3/coins/list", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		coin1 := Coin{
			Id:        "01coin",
			Symbol:    "zoc",
			Name:      "01coin",
			Platforms: nil,
		}
		coin2 := Coin{
			Id:        "02-token",
			Symbol:    "o2t",
			Name:      "O2 Token",
			Platforms: map[string]string{"ethereum": "0xb1bafca3737268a96673a250173b6ed8f1b5b65f"},
		}
		coin3 := Coin{
			Id:        "lovehearts",
			Symbol:    "lvh",
			Name:      "LoveHearts",
			Platforms: map[string]string{"tron": "1000451"},
		}
		coin4 := Coin{
			Id:        "xrp-bep2",
			Symbol:    "xrp-bf2",
			Name:      "XRP BEP2",
			Platforms: map[string]string{"binancecoin": "XRP-BF2"},
		}
		coin5 := Coin{
			Id:        "bitcoin",
			Symbol:    "btc",
			Name:      "Bitcoin",
			Platforms: nil,
		}
		coin6 := Coin{
			Id:        "binancecoin",
			Symbol:    "bnb",
			Name:      "Binance Coin",
			Platforms: map[string]string{"binancecoin": "BNB"},
		}
		coin7 := Coin{
			Id:        "tron",
			Symbol:    "trx",
			Name:      "TRON",
			Platforms: nil,
		}
		coin8 := Coin{
			Id:        "ethereum",
			Symbol:    "eth",
			Name:      "ethereum",
			Platforms: nil,
		}

		rawBytes, err := json.Marshal(Coins([]Coin{coin1, coin2, coin3, coin4, coin5, coin6, coin7, coin8}))
		if err != nil {
			panic(err)
		}
		if _, err := w.Write(rawBytes); err != nil {
			panic(err)
		}
	})

	r.HandleFunc("/v3/coins/markets", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, mockedMarketsResponse); err != nil {
			panic(err)
		}
	})

	r.HandleFunc("/v3/coins/ethereum/market_chart/range", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, mockedChartResponse); err != nil {
			panic(err)
		}
	})
	r.HandleFunc("/ethereum/info/info.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, mockedInfoResponse); err != nil {
			panic(err)
		}
	})

	return r
}

var (
	wantedRates           = `[{"currency":"BTC","provider":"coingecko","rate":9696.96,"timestamp":1588871554},{"currency":"ETH","provider":"coingecko","rate":206.55,"timestamp":1588871558},{"currency":"BNB","provider":"coingecko","rate":16.76,"timestamp":1588871427},{"currency":"TRX","provider":"coingecko","rate":0.01594768,"timestamp":1588871427},{"currency":"ZOC","provider":"coingecko","rate":0.00135115,"timestamp":1588870632},{"currency":"O2T","provider":"coingecko","rate":0.00083971,"timestamp":1577332821},{"currency":"XRP-BF2","provider":"coingecko","rate":0.21726,"timestamp":1588871653},{"currency":"LVH","provider":"coingecko","rate":0.00000808,"timestamp":1588871413}]`
	mockedMarketsResponse = `[ { "id": "bitcoin", "symbol": "btc", "name": "Bitcoin", "image": "https://assets.coingecko.com/coins/images/1/large/bitcoin.png?1547033579", "current_price": 9696.96, "market_cap": 177446468003, "market_cap_rank": 1, "total_volume": 51778003346, "high_24h": 9661.04, "low_24h": 9099.5, "price_change_24h": 459.99, "price_change_percentage_24h": 4.97984, "market_cap_change_24h": 7519669329, "market_cap_change_percentage_24h": 4.42524, "circulating_supply": 18367225.0, "total_supply": 21000000.0, "ath": 19665.39, "ath_change_percentage": -51.20717, "ath_date": "2017-12-16T00:00:00.000Z", "atl": 67.81, "atl_change_percentage": 14050.4866, "atl_date": "2013-07-06T00:00:00.000Z", "roi": null, "last_updated": "2020-05-07T17:12:34.220Z" }, { "id": "ethereum", "symbol": "eth", "name": "Ethereum", "image": "https://assets.coingecko.com/coins/images/279/large/ethereum.png?1547034048", "current_price": 206.55, "market_cap": 22851909019, "market_cap_rank": 2, "total_volume": 17356592769, "high_24h": 207.93, "low_24h": 199.21, "price_change_24h": -1.32300849, "price_change_percentage_24h": -0.63646, "market_cap_change_24h": -218200222.339359, "market_cap_change_percentage_24h": -0.94581, "circulating_supply": 110830501.624, "total_supply": null, "ath": 1448.18, "ath_change_percentage": -85.76061, "ath_date": "2018-01-13T00:00:00.000Z", "atl": 0.432979, "atl_change_percentage": 47526.36618, "atl_date": "2015-10-20T00:00:00.000Z", "roi": { "times": 27.445825905402025, "currency": "btc", "percentage": 2744.5825905402025 }, "last_updated": "2020-05-07T17:12:38.629Z" }, { "id": "binancecoin", "symbol": "bnb", "name": "Binance Coin", "image": "https://assets.coingecko.com/coins/images/825/large/binance-coin-logo.png?1547034615", "current_price": 16.76, "market_cap": 2474368090, "market_cap_rank": 9, "total_volume": 386293713, "high_24h": 16.96, "low_24h": 16.34, "price_change_24h": -0.15954337, "price_change_percentage_24h": -0.9431, "market_cap_change_24h": -28933861.2562523, "market_cap_change_percentage_24h": -1.15583, "circulating_supply": 147883948.0, "total_supply": 179883948.0, "ath": 39.68, "ath_change_percentage": -57.95222, "ath_date": "2019-06-22T12:20:21.894Z", "atl": 0.0398177, "atl_change_percentage": 41801.07312, "atl_date": "2017-10-19T00:00:00.000Z", "roi": null, "last_updated": "2020-05-07T17:10:27.413Z" }, { "id": "tron", "symbol": "trx", "name": "TRON", "image": "https://assets.coingecko.com/coins/images/1094/large/tron-logo.png?1547035066", "current_price": 0.01594768, "market_cap": 1057389936, "market_cap_rank": 17, "total_volume": 1697721558, "high_24h": 0.01620348, "low_24h": 0.0155473, "price_change_24h": -8.271e-05, "price_change_percentage_24h": -0.51593, "market_cap_change_24h": -4768379.66360569, "market_cap_change_percentage_24h": -0.44893, "circulating_supply": 66140232427.0, "total_supply": 99281283754.0, "ath": 0.231673, "ath_change_percentage": -93.12752, "ath_date": "2018-01-05T00:00:00.000Z", "atl": 0.00180434, "atl_change_percentage": 782.40789, "atl_date": "2017-11-12T00:00:00.000Z", "roi": { "times": 7.39351342638589, "currency": "usd", "percentage": 739.3513426385889 }, "last_updated": "2020-05-07T17:10:27.759Z" }, { "id": "01coin", "symbol": "zoc", "name": "01coin", "image": "https://assets.coingecko.com/coins/images/5720/large/F1nTlw9I_400x400.jpg?1547041588", "current_price": 0.00135115, "market_cap": 14384.83, "market_cap_rank": 1587, "total_volume": 888.82, "high_24h": 0.0013974, "low_24h": 0.00120559, "price_change_24h": 3.876e-05, "price_change_percentage_24h": 2.95316, "market_cap_change_24h": 480.77, "market_cap_change_percentage_24h": 3.45776, "circulating_supply": 10646360.834599, "total_supply": 65658824.0, "ath": 0.03418169, "ath_change_percentage": -96.04715, "ath_date": "2018-10-10T17:27:38.034Z", "atl": 0.00070641, "atl_change_percentage": 91.26875, "atl_date": "2020-03-16T10:22:30.944Z", "roi": null, "last_updated": "2020-05-07T16:57:12.616Z" }, { "id": "02-token", "symbol": "o2t", "name": "O2 Token", "image": "https://assets.coingecko.com/coins/images/6925/large/44429612.jpeg?1547043298", "current_price": 0.00083971, "market_cap": 0.0, "market_cap_rank": 7111, "total_volume": 69.52, "high_24h": null, "low_24h": null, "price_change_24h": null, "price_change_percentage_24h": null, "market_cap_change_24h": null, "market_cap_change_percentage_24h": null, "circulating_supply": 0.0, "total_supply": 28520100.0, "ath": 0.00439107, "ath_change_percentage": -80.87694, "ath_date": "2018-11-20T05:12:22.611Z", "atl": 0.00057411, "atl_change_percentage": 46.26319, "atl_date": "2018-11-26T00:00:00.000Z", "roi": null, "last_updated": "2019-12-26T04:00:21.046Z" }, { "id": "xrp-bep2", "symbol": "xrp-bf2", "name": "XRP BEP2", "image": "https://assets.coingecko.com/coins/images/9686/large/12-122739_xrp-logo-png-clipart.png?1570790408", "current_price": 0.21726, "market_cap": 0.0, "market_cap_rank": 5069, "total_volume": 301.34, "high_24h": 0.219438, "low_24h": 0.212267, "price_change_24h": -0.00202035, "price_change_percentage_24h": -0.92136, "market_cap_change_24h": 0.0, "market_cap_change_percentage_24h": 0.0, "circulating_supply": 0.0, "total_supply": 10000000.0, "ath": 0.360995, "ath_change_percentage": -40.44337, "ath_date": "2019-10-21T13:44:21.822Z", "atl": 0.115982, "atl_change_percentage": 85.36984, "atl_date": "2020-03-13T02:02:33.103Z", "roi": null, "last_updated": "2020-05-07T17:14:13.364Z" }, { "id": "lovehearts", "symbol": "lvh", "name": "LoveHearts", "image": "https://assets.coingecko.com/coins/images/9360/large/1_d3hJ7JQeQ84goeTVWLI9Qw.png?1566528108", "current_price": 8.08e-06, "market_cap": 0.0, "market_cap_rank": 5528, "total_volume": 7.87, "high_24h": 8.73e-06, "low_24h": 7.85e-06, "price_change_24h": 1.9e-07, "price_change_percentage_24h": 2.39645, "market_cap_change_24h": 0.0, "market_cap_change_percentage_24h": 0.0, "circulating_supply": 0.0, "total_supply": 100000000000.0, "ath": 8.596e-05, "ath_change_percentage": -90.70838, "ath_date": "2019-08-23T03:49:38.791Z", "atl": 3.13e-06, "atl_change_percentage": 155.38143, "atl_date": "2020-02-21T21:25:35.813Z", "roi": null, "last_updated": "2020-05-07T17:10:13.895Z" } ]`
	wantedInfo            = `{"provider":"coingecko","provider_url":"https://www.coingecko.com/en/coins/bitcoin","volume_24":51778003346,"market_cap":177446468003,"circulating_supply":18367225,"total_supply":21000000,"info":{"name":"Ethereum","website":"https://ethereum.org/","source_code":"https://github.com/ethereum","white_paper":"https://github.com/ethereum/wiki/wiki/White-Paper","description":"Open source platform to write and distribute decentralized applications.","short_description":"Open source platform to write and distribute decentralized applications.","explorer":"https://etherscan.io/","socials":[{"name":"Twitter","url":"https://twitter.com/ethereum","handle":"ethereum"},{"name":"Reddit","url":"https://www.reddit.com/r/ethereum","handle":"ethereum"}]}}`
	mockedInfoResponse    = `{ "name": "Ethereum", "website": "https://ethereum.org/", "source_code": "https://github.com/ethereum", "white_paper": "https://github.com/ethereum/wiki/wiki/White-Paper", "short_description": "Open source platform to write and distribute decentralized applications.", "description": "Open source platform to write and distribute decentralized applications.", "socials": [ { "name": "Twitter", "url": "https://twitter.com/ethereum", "handle": "ethereum" }, { "name": "Reddit", "url": "https://www.reddit.com/r/ethereum", "handle": "ethereum" } ], "explorer": "https://etherscan.io/" }`
	wantedCharts          = `{"provider":"coingecko","prices":[{"price":130.4846850311141,"date":1577923200},{"price":127.04525801179804,"date":1578009600},{"price":133.70264861844188,"date":1578096000},{"price":134.1368826978575,"date":1578182400},{"price":135.00571364233375,"date":1578268800},{"price":143.80639795566535,"date":1578355200},{"price":143.019432264986,"date":1578441600},{"price":140.27393527388705,"date":1578528000},{"price":137.86055928182662,"date":1578614400},{"price":144.6047743464792,"date":1578700800},{"price":142.18233327497873,"date":1578787200},{"price":145.42287633530452,"date":1578873600},{"price":143.58796905036286,"date":1578960000},{"price":165.99143854583525,"date":1579046400},{"price":166.2530969048888,"date":1579132800},{"price":163.80590111423027,"date":1579219200},{"price":171.15621787671842,"date":1579305600},{"price":174.23837318258896,"date":1579392000},{"price":166.6271059505621,"date":1579478400},{"price":166.93852805776982,"date":1579564800},{"price":169.27993585839258,"date":1579651200},{"price":167.8265296405669,"date":1579737600},{"price":162.51937110634034,"date":1579824000},{"price":162.40947975163584,"date":1579910400},{"price":160.6736108617405,"date":1579996800},{"price":167.6474095857562,"date":1580083200},{"price":169.7364357354222,"date":1580169600},{"price":175.19000305345378,"date":1580256000},{"price":173.70688823620796,"date":1580342400},{"price":184.72621843141044,"date":1580428800},{"price":179.22910385240368,"date":1580515200},{"price":183.33693755210257,"date":1580601600},{"price":188.55064179431113,"date":1580688000},{"price":189.86176346298166,"date":1580774400},{"price":188.84220736948276,"date":1580860800},{"price":203.8593183747095,"date":1580947200},{"price":212.73341421352782,"date":1581033600},{"price":223.27631974791998,"date":1581120000},{"price":223.30077697426225,"date":1581206400},{"price":228.29226083217347,"date":1581292800},{"price":224.14699652028605,"date":1581379200},{"price":236.78534852972567,"date":1581465600},{"price":264.03276824299627,"date":1581552000},{"price":267.6704445042514,"date":1581638400},{"price":284.2318914741183,"date":1581724800},{"price":263.8999463915962,"date":1581811200},{"price":262.1562443168135,"date":1581897600},{"price":267.9356656159883,"date":1581984000},{"price":281.9457403055278,"date":1582070400},{"price":259.1838037480287,"date":1582156800},{"price":257.9883953048601,"date":1582243200},{"price":265.1644591338591,"date":1582329600},{"price":260.99654013152536,"date":1582416000},{"price":274.6312157785517,"date":1582502400},{"price":265.2421912203732,"date":1582588800},{"price":248.31738139395102,"date":1582675200},{"price":224.81050829744677,"date":1582761600},{"price":225.42971511683754,"date":1582848000},{"price":227.70528448376498,"date":1582934400},{"price":218.34835692985286,"date":1583020800},{"price":218.93979697447656,"date":1583107200},{"price":231.67944604592213,"date":1583193600},{"price":223.96011601810258,"date":1583280000},{"price":224.13487737180202,"date":1583366400},{"price":228.0809351016884,"date":1583452800},{"price":244.23846951978683,"date":1583539200},{"price":237.38079020191347,"date":1583625600},{"price":198.81964228768578,"date":1583712000},{"price":200.84687148188746,"date":1583798400},{"price":200.74903854851766,"date":1583884800},{"price":194.21794124025206,"date":1583971200},{"price":110.5978978308351,"date":1584057600},{"price":132.57285770180832,"date":1584144000},{"price":123.030844108903,"date":1584230400},{"price":124.60342553062202,"date":1584316800},{"price":110.99159845334059,"date":1584403200},{"price":117.21704138518466,"date":1584489600},{"price":117.7004246469663,"date":1584576000},{"price":136.81187578973578,"date":1584662400},{"price":131.96469397290025,"date":1584748800},{"price":131.94321845316688,"date":1584835200},{"price":122.5205987665882,"date":1584921600},{"price":135.58865798037067,"date":1585008000},{"price":138.77657127268222,"date":1585094400},{"price":136.23392142176775,"date":1585180800},{"price":138.76693759902926,"date":1585267200},{"price":130.2865618875773,"date":1585353600},{"price":131.42405214081145,"date":1585440000},{"price":125.31887451568916,"date":1585526400},{"price":132.36373895414985,"date":1585612800},{"price":133.2364468438201,"date":1585699200},{"price":136.21638557312852,"date":1585785600},{"price":141.4537053523866,"date":1585872000},{"price":141.26154326644877,"date":1585958400},{"price":144.200802681986,"date":1586044800},{"price":142.8508834648482,"date":1586131200},{"price":169.85504979336824,"date":1586217600},{"price":164.5160197107307,"date":1586304000},{"price":172.80356826747519,"date":1586390400},{"price":170.09510352387568,"date":1586476800},{"price":157.74015800780134,"date":1586563200},{"price":158.32787794181905,"date":1586649600},{"price":158.8638257195104,"date":1586736000},{"price":156.70135879140688,"date":1586822400},{"price":158.26715118298955,"date":1586908800},{"price":153.2228636800284,"date":1586995200},{"price":171.7759907508837,"date":1587081600},{"price":170.445890477055,"date":1587168000},{"price":187.14354411560592,"date":1587254400},{"price":180.05859654051866,"date":1587340800},{"price":170.70405868653958,"date":1587427200},{"price":170.45154288047368,"date":1587513600},{"price":182.2767576620704,"date":1587600000},{"price":184.59370248021838,"date":1587686400},{"price":187.34336861419638,"date":1587772800},{"price":194.11423682570447,"date":1587859200},{"price":197.22990444623989,"date":1587945600},{"price":196.46262535084685,"date":1588032000},{"price":197.15473156634306,"date":1588118400},{"price":215.54817266567386,"date":1588204800},{"price":205.55600541063674,"date":1588291200},{"price":211.96829534904975,"date":1588377600},{"price":213.94301121407548,"date":1588464000},{"price":210.02470947464775,"date":1588550400},{"price":206.8325342346069,"date":1588636800},{"price":205.22517973985205,"date":1588723200},{"price":200.25479450009442,"date":1588809600},{"price":212.29247620086525,"date":1588896000},{"price":211.68904835113835,"date":1588982400}]}`
	mockedChartResponse   = `{"prices":[[1577923200000,130.4846850311141],[1578009600000,127.04525801179804],[1578096000000,133.70264861844188],[1578182400000,134.1368826978575],[1578268800000,135.00571364233375],[1578355200000,143.80639795566535],[1578441600000,143.019432264986],[1578528000000,140.27393527388705],[1578614400000,137.86055928182662],[1578700800000,144.6047743464792],[1578787200000,142.18233327497873],[1578873600000,145.42287633530452],[1578960000000,143.58796905036286],[1579046400000,165.99143854583525],[1579132800000,166.2530969048888],[1579219200000,163.80590111423027],[1579305600000,171.15621787671842],[1579392000000,174.23837318258896],[1579478400000,166.6271059505621],[1579564800000,166.93852805776982],[1579651200000,169.27993585839258],[1579737600000,167.8265296405669],[1579824000000,162.51937110634034],[1579910400000,162.40947975163584],[1579996800000,160.6736108617405],[1580083200000,167.6474095857562],[1580169600000,169.7364357354222],[1580256000000,175.19000305345378],[1580342400000,173.70688823620796],[1580428800000,184.72621843141044],[1580515200000,179.22910385240368],[1580601600000,183.33693755210257],[1580688000000,188.55064179431113],[1580774400000,189.86176346298166],[1580860800000,188.84220736948276],[1580947200000,203.8593183747095],[1581033600000,212.73341421352782],[1581120000000,223.27631974791998],[1581206400000,223.30077697426225],[1581292800000,228.29226083217347],[1581379200000,224.14699652028605],[1581465600000,236.78534852972567],[1581552000000,264.03276824299627],[1581638400000,267.6704445042514],[1581724800000,284.2318914741183],[1581811200000,263.8999463915962],[1581897600000,262.1562443168135],[1581984000000,267.9356656159883],[1582070400000,281.9457403055278],[1582156800000,259.1838037480287],[1582243200000,257.9883953048601],[1582329600000,265.1644591338591],[1582416000000,260.99654013152536],[1582502400000,274.6312157785517],[1582588800000,265.2421912203732],[1582675200000,248.31738139395102],[1582761600000,224.81050829744677],[1582848000000,225.42971511683754],[1582934400000,227.70528448376498],[1583020800000,218.34835692985286],[1583107200000,218.93979697447656],[1583193600000,231.67944604592213],[1583280000000,223.96011601810258],[1583366400000,224.13487737180202],[1583452800000,228.0809351016884],[1583539200000,244.23846951978683],[1583625600000,237.38079020191347],[1583712000000,198.81964228768578],[1583798400000,200.84687148188746],[1583884800000,200.74903854851766],[1583971200000,194.21794124025206],[1584057600000,110.5978978308351],[1584144000000,132.57285770180832],[1584230400000,123.030844108903],[1584316800000,124.60342553062202],[1584403200000,110.99159845334059],[1584489600000,117.21704138518466],[1584576000000,117.7004246469663],[1584662400000,136.81187578973578],[1584748800000,131.96469397290025],[1584835200000,131.94321845316688],[1584921600000,122.5205987665882],[1585008000000,135.58865798037067],[1585094400000,138.77657127268222],[1585180800000,136.23392142176775],[1585267200000,138.76693759902926],[1585353600000,130.2865618875773],[1585440000000,131.42405214081145],[1585526400000,125.31887451568916],[1585612800000,132.36373895414985],[1585699200000,133.2364468438201],[1585785600000,136.21638557312852],[1585872000000,141.4537053523866],[1585958400000,141.26154326644877],[1586044800000,144.200802681986],[1586131200000,142.8508834648482],[1586217600000,169.85504979336824],[1586304000000,164.5160197107307],[1586390400000,172.80356826747519],[1586476800000,170.09510352387568],[1586563200000,157.74015800780134],[1586649600000,158.32787794181905],[1586736000000,158.8638257195104],[1586822400000,156.70135879140688],[1586908800000,158.26715118298955],[1586995200000,153.2228636800284],[1587081600000,171.7759907508837],[1587168000000,170.445890477055],[1587254400000,187.14354411560592],[1587340800000,180.05859654051866],[1587427200000,170.70405868653958],[1587513600000,170.45154288047368],[1587600000000,182.2767576620704],[1587686400000,184.59370248021838],[1587772800000,187.34336861419638],[1587859200000,194.11423682570447],[1587945600000,197.22990444623989],[1588032000000,196.46262535084685],[1588118400000,197.15473156634306],[1588204800000,215.54817266567386],[1588291200000,205.55600541063674],[1588377600000,211.96829534904975],[1588464000000,213.94301121407548],[1588550400000,210.02470947464775],[1588636800000,206.8325342346069],[1588723200000,205.22517973985205],[1588809600000,200.25479450009442],[1588896000000,212.29247620086525],[1588982400000,211.68904835113835]],"market_caps":[[1577923200000,14230588104.281237],[1578009600000,13859793329.814827],[1578096000000,14613194775.854177],[1578182400000,14636213160.069315],[1578268800000,14738834254.266087],[1578355200000,15692578635.119993],[1578441600000,15656491188.064734],[1578528000000,15364547561.67077],[1578614400000,14998595868.023438],[1578700800000,15647071448.436348],[1578787200000,15564837467.520994],[1578873600000,15859031094.44091],[1578960000000,15715287657.358463],[1579046400000,18075462340.568443],[1579132800000,18104742352.015366],[1579219200000,17869543294.566414],[1579305600000,18727969032.79025],[1579392000000,19138755594.52358],[1579478400000,18220366702.244606],[1579564800000,18243871990.6739],[1579651200000,18460475135.804626],[1579737600000,18358377435.95901],[1579824000000,17722969842.693172],[1579910400000,17777723072.283943],[1579996800000,17587666193.195988],[1580083200000,18279197138.315884],[1580169600000,18576520021.386745],[1580256000000,19212409346.36024],[1580342400000,19140249467.28951],[1580428800000,20090242968.224346],[1580515200000,19751392940.997242],[1580601600000,20079997942.89275],[1580688000000,20619310835.276108],[1580774400000,20791014743.301926],[1580860800000,20665245888.10715],[1580947200000,22414850796.683693],[1581033600000,23314707733.86607],[1581120000000,24407934017.312393],[1581206400000,24478412143.01876],[1581292800000,25011126326.73704],[1581379200000,24577280465.17777],[1581465600000,25966283262.662064],[1581552000000,29001954205.996136],[1581638400000,29312880235.867893],[1581724800000,31059453816.69378],[1581811200000,28909877146.05339],[1581897600000,28746450516.039288],[1581984000000,29075052114.63163],[1582070400000,31034929114.66128],[1582156800000,28662558273.135616],[1582243200000,28258309183.25315],[1582329600000,29114570356.366127],[1582416000000,28657670954.134422],[1582502400000,30135977087.940575],[1582588800000,28974292556.854683],[1582675200000,27274600634.783813],[1582761600000,24823014146.598347],[1582848000000,24879553838.25944],[1582934400000,25137364865.97101],[1583020800000,24105658618.96741],[1583107200000,24049874527.90718],[1583193600000,25392608994.495888],[1583280000000,24709361638.74632],[1583366400000,24635535324.34428],[1583452800000,25072941968.371387],[1583539200000,26861442576.626507],[1583625600000,26199142579.509506],[1583712000000,22145350491.45419],[1583798400000,22001125399.41239],[1583884800000,22178131746.591427],[1583971200000,21375073649.951973],[1584057600000,11956631650.303354],[1584144000000,14567054564.19822],[1584230400000,13505998666.239887],[1584316800000,13756072458.523424],[1584403200000,12178760413.015514],[1584489600000,12991193734.51619],[1584576000000,12900856057.87349],[1584662400000,15068113354.89011],[1584748800000,14451395504.502293],[1584835200000,14575207631.098148],[1584921600000,13499834314.773163],[1585008000000,14944699607.313162],[1585094400000,15385255695.77982],[1585180800000,15020580578.564753],[1585267200000,15277363230.972204],[1585353600000,14958737456.74538],[1585440000000,14468066254.835968],[1585526400000,13819843925.360764],[1585612800000,14615131244.892355],[1585699200000,14699970260.177477],[1585785600000,14896572551.520117],[1585872000000,15522802609.007648],[1585958400000,15597757554.270603],[1586044800000,15898519073.894499],[1586131200000,15770466064.128725],[1586217600000,18679968361.111393],[1586304000000,18166748814.85624],[1586390400000,19089126835.76443],[1586476800000,18787453383.948494],[1586563200000,17375118634.991734],[1586649600000,17489565908.771362],[1586736000000,17790409219.818943],[1586822400000,17316641314.145325],[1586908800000,17559917949.258163],[1586995200000,16933905371.155535],[1587081600000,18989526634.23506],[1587168000000,18844786940.26037],[1587254400000,20686869606.075146],[1587340800000,19894212517.443478],[1587427200000,18894679854.84961],[1587513600000,18883297426.60788],[1587600000000,20165266729.19022],[1587686400000,20371607313.859978],[1587772800000,20756708748.58927],[1587859200000,21449928671.504494],[1587945600000,21798344583.464844],[1588032000000,21715118585.744934],[1588118400000,21824818237.698734],[1588204800000,23796586177.48911],[1588291200000,22894189506.373463],[1588377600000,23540629181.114624],[1588464000000,23700169401.049606],[1588550400000,23254687127.813503],[1588636800000,22819858077.111965],[1588723200000,22740460359.466698],[1588809600000,22529347281.111214],[1588896000000,23678773321.927563],[1588982400000,23476076684.66088]],"total_volumes":[[1577923200000,6623732040.964472],[1578009600000,6497297884.22486],[1578096000000,9655245099.720835],[1578182400000,6958393690.018681],[1578268800000,7260645128.406949],[1578355200000,9514201882.396584],[1578441600000,9440938787.519835],[1578528000000,10095847759.294762],[1578614400000,7280640218.292043],[1578700800000,10273370788.312328],[1578787200000,9737605323.282524],[1578873600000,8488793477.5950365],[1578960000000,8019794567.486078],[1579046400000,18970204899.197784],[1579132800000,17862385531.73745],[1579219200000,14605493158.53319],[1579305600000,15915322019.78247],[1579392000000,15808136411.02366],[1579478400000,14003433869.000366],[1579564800000,10983216126.222286],[1579651200000,8267420452.9121065],[1579737600000,8809189828.884594],[1579824000000,9017537604.34573],[1579910400000,11043498594.604164],[1579996800000,8673816085.484137],[1580083200000,10518901915.39258],[1580169600000,12435457692.695446],[1580256000000,14300512538.139254],[1580342400000,15565369994.796808],[1580428800000,18170945326.39071],[1580515200000,18506953476.934044],[1580601600000,19083817630.29935],[1580688000000,21543158277.49522],[1580774400000,23519056716.26084],[1580860800000,21734656893.64477],[1580947200000,27708721578.171207],[1581033600000,30730370695.67516],[1581120000000,34073387340.7219],[1581206400000,35852075948.97941],[1581292800000,37065973224.819984],[1581379200000,40180881296.506096],[1581465600000,17260236147.946793],[1581552000000,25083848405.774593],[1581638400000,31919298278.23887],[1581724800000,32486783383.59592],[1581811200000,31021139510.840595],[1581897600000,26307305450.095947],[1581984000000,28621805406.267815],[1582070400000,33038477053.26572],[1582156800000,24634102799.07547],[1582243200000,34355959852.746918],[1582329600000,18672669187.760765],[1582416000000,14860160863.37393],[1582502400000,15734092874.96997],[1582588800000,21019871769.54907],[1582675200000,15216321566.877127],[1582761600000,20470479538.054836],[1582848000000,23994578699.806026],[1582934400000,19715536659.914295],[1583020800000,13967572974.821035],[1583107200000,16557363028.554333],[1583193600000,18301098060.91034],[1583280000000,19133763158.937866],[1583366400000,14844945103.395332],[1583452800000,16056678004.768778],[1583539200000,16938887970.316448],[1583625600000,15033855133.023136],[1583712000000,18839637103.648037],[1583798400000,17900321700.644154],[1583884800000,10622282127.512653],[1583971200000,14699687559.439268],[1584057600000,17357843633.111378],[1584144000000,24412190319.85676],[1584230400000,8678897518.712492],[1584316800000,10971643911.795376],[1584403200000,11161656313.89003],[1584489600000,10514054745.732649],[1584576000000,9829287702.00789],[1584662400000,13910341816.437582],[1584748800000,16267979993.668526],[1584835200000,12633534610.010733],[1584921600000,10177508398.873247],[1585008000000,11596650460.730272],[1585094400000,12119249137.978954],[1585180800000,10616935307.838089],[1585267200000,7749490793.658535],[1585353600000,10183491264.54011],[1585440000000,11067855577.856102],[1585526400000,7718875104.828544],[1585612800000,8469234427.790874],[1585699200000,7869741296.803818],[1585785600000,7848601147.336895],[1585872000000,13645052823.766476],[1585958400000,11621521038.709623],[1586044800000,9863619646.24454],[1586131200000,8271697077.94318],[1586217600000,15337221519.354977],[1586304000000,16711178690.848326],[1586390400000,12116309211.63937],[1586476800000,11971791025.372421],[1586563200000,14135915530.589357],[1586649600000,11128989594.360977],[1586736000000,13002801329.50846],[1586822400000,13286483792.048328],[1586908800000,11468431211.774742],[1586995200000,8578826590.313573],[1587081600000,17430629608.3789],[1587168000000,15220772945.895311],[1587254400000,14284812698.801844],[1587340800000,13763176713.456295],[1587427200000,15901061387.704433],[1587513600000,12178214643.078756],[1587600000000,14853382137.991045],[1587686400000,6903118239.477871],[1587772800000,13419651870.125168],[1587859200000,14452583265.899439],[1587945600000,13530756223.951258],[1588032000000,17659706480.60195],[1588118400000,14644384537.090769],[1588204800000,20966892392.04868],[1588291200000,23630952297.4432],[1588377600000,17137878522.355497],[1588464000000,14287655444.044804],[1588550400000,16158461349.767452],[1588636800000,17901489283.101368],[1588723200000,14345915703.711891],[1588809600000,18128961143.16994],[1588896000000,15142506242.294548],[1588982400000,18337241719.88773]]}`
	wantedTickers         = `[{"coin":0,"coin_name":"BTC","type":"coin","price":{"change_24h":4.97984,"currency":"USD","provider":"coingecko","value":9696.96},"last_update":"2020-05-07T17:12:34.22Z","volume":51778003346,"market_cap":177446468003},{"coin":60,"coin_name":"ETH","type":"coin","price":{"change_24h":-0.63646,"currency":"USD","provider":"coingecko","value":206.55},"last_update":"2020-05-07T17:12:38.629Z","volume":17356592769,"market_cap":22851909019},{"coin":20000714,"coin_name":"BNB","type":"coin","price":{"change_24h":-0.9431,"currency":"USD","provider":"coingecko","value":16.76},"last_update":"2020-05-07T17:10:27.413Z","volume":386293713,"market_cap":2474368090},{"coin":195,"coin_name":"TRX","type":"coin","price":{"change_24h":-0.51593,"currency":"USD","provider":"coingecko","value":0.01594768},"last_update":"2020-05-07T17:10:27.759Z","volume":1697721558,"market_cap":1057389936},{"coin":60,"coin_name":"ETH","token_id":"0xb1bafca3737268a96673a250173b6ed8f1b5b65f","type":"token","price":{"change_24h":0,"currency":"USD","provider":"coingecko","value":0.00083971},"last_update":"2019-12-26T04:00:21.046Z","volume":69.52,"market_cap":0},{"coin":714,"coin_name":"BNB","token_id":"xrp-bf2","type":"token","price":{"change_24h":-0.92136,"currency":"USD","provider":"coingecko","value":0.21726},"last_update":"2020-05-07T17:14:13.364Z","volume":301.34,"market_cap":0},{"coin":195,"coin_name":"TRX","token_id":"1000451","type":"token","price":{"change_24h":2.39645,"currency":"USD","provider":"coingecko","value":0.00000808},"last_update":"2020-05-07T17:10:13.895Z","volume":7.87,"market_cap":0}]`
	wantedTickers2        = `[{"coin":0,"coin_name":"BTC","type":"coin","price":{"change_24h":4.97984,"currency":"USD","provider":"coingecko","value":9696.96},"last_update":"2020-05-07T17:12:34.22Z","volume":51778003346,"market_cap":177446468003},{"coin":60,"coin_name":"ETH","type":"coin","price":{"change_24h":-0.63646,"currency":"USD","provider":"coingecko","value":206.55},"last_update":"2020-05-07T17:12:38.629Z","volume":17356592769,"market_cap":22851909019},{"coin":714,"coin_name":"BNB","type":"coin","price":{"change_24h":-0.9431,"currency":"USD","provider":"coingecko","value":16.76},"last_update":"2020-05-07T17:10:27.413Z","volume":386293713,"market_cap":2474368090},{"coin":195,"coin_name":"TRX","type":"coin","price":{"change_24h":-0.51593,"currency":"USD","provider":"coingecko","value":0.01594768},"last_update":"2020-05-07T17:10:27.759Z","volume":1697721558,"market_cap":1057389936},{"coin":60,"coin_name":"ETH","token_id":"0xb1bafca3737268a96673a250173b6ed8f1b5b65f","type":"token","price":{"change_24h":0,"currency":"USD","provider":"coingecko","value":0.00083971},"last_update":"2019-12-26T04:00:21.046Z","volume":69.52,"market_cap":0},{"coin":714,"coin_name":"BNB","token_id":"xrp-bf2","type":"token","price":{"change_24h":-0.92136,"currency":"USD","provider":"coingecko","value":0.21726},"last_update":"2020-05-07T17:14:13.364Z","volume":301.34,"market_cap":0},{"coin":195,"coin_name":"TRX","token_id":"1000451","type":"token","price":{"change_24h":2.39645,"currency":"USD","provider":"coingecko","value":0.00000808},"last_update":"2020-05-07T17:10:13.895Z","volume":7.87,"market_cap":0}]`
	wantedTickers3        = `[{"coin":0,"coin_name":"BTC","type":"coin","price":{"change_24h":4.97984,"currency":"USD","provider":"coingecko","value":9696.96},"last_update":"2020-05-07T17:12:34.22Z","volume":51778003346,"market_cap":177446468003},{"coin":60,"coin_name":"ETH","type":"coin","price":{"change_24h":-0.63646,"currency":"USD","provider":"coingecko","value":206.55},"last_update":"2020-05-07T17:12:38.629Z","volume":17356592769,"market_cap":22851909019},{"coin":10000714,"coin_name":"BNB","type":"coin","price":{"change_24h":-0.9431,"currency":"USD","provider":"coingecko","value":16.76},"last_update":"2020-05-07T17:10:27.413Z","volume":386293713,"market_cap":2474368090},{"coin":195,"coin_name":"TRX","type":"coin","price":{"change_24h":-0.51593,"currency":"USD","provider":"coingecko","value":0.01594768},"last_update":"2020-05-07T17:10:27.759Z","volume":1697721558,"market_cap":1057389936},{"coin":60,"coin_name":"ETH","token_id":"0xb1bafca3737268a96673a250173b6ed8f1b5b65f","type":"token","price":{"change_24h":0,"currency":"USD","provider":"coingecko","value":0.00083971},"last_update":"2019-12-26T04:00:21.046Z","volume":69.52,"market_cap":0},{"coin":714,"coin_name":"BNB","token_id":"xrp-bf2","type":"token","price":{"change_24h":-0.92136,"currency":"USD","provider":"coingecko","value":0.21726},"last_update":"2020-05-07T17:14:13.364Z","volume":301.34,"market_cap":0},{"coin":195,"coin_name":"TRX","token_id":"1000451","type":"token","price":{"change_24h":2.39645,"currency":"USD","provider":"coingecko","value":0.00000808},"last_update":"2020-05-07T17:10:13.895Z","volume":7.87,"market_cap":0}]`
)
