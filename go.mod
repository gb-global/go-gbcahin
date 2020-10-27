module gbchain-org/go-gbchain

go 1.13

replace (
	github.com/asdine/storm/v3 => github.com/simplechain-org/storm/v3 v3.2.1-0.20200521045524-c61eb1b00dec
	github.com/coreos/etcd => github.com/simplechain-org/quorum-etcd v0.1.0
)

require (
	github.com/Azure/azure-storage-blob-go v0.7.0
	github.com/Beyond-simplechain/foundation v0.0.0-20200417022121-620b0f2460ff
	github.com/OneOfOne/xxhash v1.2.5 // indirect
	github.com/VictoriaMetrics/fastcache v1.5.7
	github.com/aristanetworks/goarista v0.0.0-20200812190859-4cb0e71f3c0e
	github.com/asdine/storm/v3 v3.1.1
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/cespare/cp v0.1.0
	github.com/cloudflare/cloudflare-go v0.10.2-0.20190916151808-a80f83b9add9
	github.com/coreos/etcd v0.1.0
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v0.0.0-20180603214616-504e848d77ea
	github.com/docker/docker v1.4.2-0.20180625184442-8e610b2b55bf
	github.com/eapache/channels v1.1.0
	github.com/edsrzf/mmap-go v1.0.0
	github.com/elastic/gosigar v0.11.0
	github.com/ethereum/go-ethereum v1.9.21 // indirect
	github.com/fatih/color v1.9.0
	github.com/fjl/memsize v0.0.0-20180418122429-ca190fb6ffbc
	github.com/gballet/go-libpcsclite v0.0.0-20190607065134-2772fd86a8ff
	github.com/go-stack/stack v1.8.0
	github.com/golang/protobuf v1.4.2
	github.com/golang/snappy v0.0.2-0.20200707131729-196ae77b8a26
	github.com/gorilla/websocket v1.4.1-0.20190629185528-ae1634f6a989
	github.com/graph-gophers/graphql-go v0.0.0-20191115155744-f33e81362277
	github.com/hashicorp/golang-lru v0.5.4
	github.com/huin/goupnp v1.0.0
	github.com/influxdata/influxdb v1.2.3-0.20180221223340-01288bdb0883
	github.com/jackpal/go-nat-pmp v1.0.2-0.20160603034137-1fa385a6f458
	github.com/json-iterator/go v1.1.9
	github.com/julienschmidt/httprouter v1.2.0
	github.com/karalabe/usb v0.0.0-20190919080040-51dc0efba356
	github.com/karalabe/xgo v0.0.0-20191115072854-c5ccff8648a7 // indirect
	github.com/mattn/go-colorable v0.1.4
	github.com/mattn/go-isatty v0.0.11
	github.com/miguelmota/go-ethereum-hdwallet v0.0.0-20200123000308-a60dcd172b4c // indirect
	github.com/naoina/toml v0.1.2-0.20170918210437-9fafd6967416
	github.com/olekukonko/tablewriter v0.0.2-0.20190409134802-7e037d187b0c
	github.com/pborman/uuid v1.2.0
	github.com/peterh/liner v1.1.1-0.20190123174540-a2c9a5303de7
	github.com/prometheus/tsdb v0.6.2-0.20190402121629-4f204dcbc150
	github.com/rjeczalik/notify v0.9.2
	github.com/robertkrimen/otto v0.0.0-20170205013659-6a77b7cbc37d
	github.com/rs/cors v1.7.0
	github.com/satori/go.uuid v1.2.0
	github.com/spaolacci/murmur3 v1.0.1-0.20190317074736-539464a789e9 // indirect
	github.com/status-im/keycard-go v0.0.0-20190316090335-8537d3370df4
	github.com/steakknife/bloomfilter v0.0.0-20180922174646-6819c0d2a570
	github.com/stretchr/testify v1.4.0
	github.com/syndtr/goleveldb v1.0.1-0.20200815110645-5c35d600f0ca
	github.com/tyler-smith/go-bip39 v1.0.2
	github.com/wsddn/go-ecdh v0.0.0-20161211032359-48726bab9208
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20200917073148-efd3b9a0ff20
	golang.org/x/text v0.3.3
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20200619000410-60c24ae608a6
	gopkg.in/oleiade/lane.v1 v1.0.0
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/urfave/cli.v1 v1.20.0
)
