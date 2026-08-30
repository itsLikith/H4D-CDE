module github.com/itsLikith/h4d-cde/services/voxel-engine

go 1.25.0

replace github.com/itsLikith/h4d-cde/gen => ../../gen

require (
	github.com/itsLikith/h4d-cde/gen v0.0.0
	github.com/redis/go-redis/v9 v9.7.0
	github.com/segmentio/kafka-go v0.4.47
	github.com/stretchr/testify v1.11.1
	github.com/uber/h3-go/v4 v4.2.0
	google.golang.org/grpc v1.70.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/klauspost/compress v1.19.2 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.16.0 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
