module github.com/paysuper/paysuper-billing-server

require (
	github.com/InVisionApp/go-health v2.1.0+incompatible
	github.com/ProtocolONE/geoip-service v0.0.0-20190903084234-1d5ae6b96679
	github.com/ProtocolONE/go-micro-plugins v0.3.0
	github.com/alicebob/miniredis v2.5.0+incompatible
	github.com/centrifugal/gocent v2.0.2+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/elliotchance/redismock v1.5.1
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/gogo/protobuf v1.3.0
	github.com/golang-migrate/migrate/v4 v4.3.1
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/jinzhu/copier v0.0.0-20190625015134-976e0346caa8
	github.com/jinzhu/now v1.0.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.8.0
	github.com/micro/go-plugins v1.2.0
	github.com/morikuni/aec v0.0.0-20170113033406-39771216ff4c // indirect
	github.com/paysuper/document-signer v0.0.0-20190923080905-cb9cc2665d8b
	github.com/paysuper/paysuper-currencies v0.0.0-20190903083641-668b8b2b997d
	github.com/paysuper/paysuper-database-mongo v0.1.1
	github.com/paysuper/paysuper-i18n v0.0.0-20190926113224-7eaca4563c7b
	github.com/paysuper/paysuper-recurring-repository v1.0.123
	github.com/paysuper/paysuper-reporter v0.0.0-20190926160409-3dfa3c2d811f
	github.com/paysuper/paysuper-tax-service v0.0.0-20190903084038-7849f394f122
	github.com/paysuper/postmark-sender v0.0.0-20190926173213-1b6244122c49
	github.com/prometheus/client_golang v1.0.0
	github.com/stoewer/go-strcase v1.0.2
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
	github.com/stretchr/testify v1.4.0
	github.com/ttacon/libphonenumber v1.0.1
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20190926114937-fa1a29108794
	gopkg.in/ProtocolONE/rabbitmq.v1 v1.0.0-20190719062839-9858d727f3ef
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
)

replace (
	github.com/gogo/protobuf v0.0.0-20190410021324-65acae22fc9 => github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1
)

go 1.13
