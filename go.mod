module github.com/trustwallet/watchmarket

go 1.15

// +heroku goVersion go1.15
// +heroku install ./cmd/...

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/alicebob/miniredis/v2 v2.14.3
	github.com/chenjiandongx/ginprom v0.0.0-20200410120253-7cfb22707fa6
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.8.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	github.com/trustwallet/golibs v0.1.5
	github.com/trustwallet/golibs/network v0.0.0-20210124080535-8638b407c4ab
	golang.org/x/tools v0.0.0-20200513175351-0951661448da // indirect
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.20.12
)
