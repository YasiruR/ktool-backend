module github.com/YasiruR/ktool-backend

go 1.14

require (
	github.com/Shopify/sarama v1.26.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.3
	github.com/pickme-go/log v1.2.3
	github.com/pickme-go/traceable-context v1.0.0
	github.com/reiver/go-oi v1.0.0 // indirect
	github.com/reiver/go-telnet v0.0.0-20180421082511-9ff0b2ab096e
	github.com/rs/cors v1.7.0
	github.com/sparrc/go-ping v0.0.0-20190613174326-4e5b6552494c
	golang.org/x/crypto v0.0.0-20200128174031-69ecbb4d6d5d
	gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v3 v3.0.0-20200506231410-2ff61e1afc86
//github.com/pickme-go/k-stream v0.0.0
)

//replace github.com/pickme-go/k-stream => github.com/pickme-go/k-stream.git v0.0.0
