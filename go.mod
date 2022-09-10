module apart-deal-api

go 1.18

require (
	github.com/Netflix/go-env v0.0.0-20220526054621-78278af1949d
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/go-redis/redis/v8 v8.11.5
	github.com/labstack/echo/v4 v4.7.2
	github.com/onsi/ginkgo/v2 v2.1.4
	github.com/onsi/gomega v1.20.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	gitlab.com/apart-deals/openapi/go/api v1.0.0
	go.mongodb.org/mongo-driver v1.10.1
	go.uber.org/fx v1.17.1
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
)

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/labstack/gommon v0.3.1 // indirect
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.1.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/dig v1.14.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220422013727-9388b58f7150 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace gitlab.com/apart-deals/openapi/go/api => ./openapi/go
