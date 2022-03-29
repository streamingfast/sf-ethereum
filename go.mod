module github.com/streamingfast/sf-ethereum

go 1.16

require (
	github.com/RoaringBitmap/roaring v0.9.4
	github.com/ShinyTrinkets/overseer v0.3.0
	github.com/golang/protobuf v1.5.2
	github.com/google/cel-go v0.4.1
	github.com/google/go-cmp v0.5.7
	github.com/lithammer/dedent v1.1.0
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/manifoldco/promptui v0.8.0
	github.com/mitchellh/go-testing-interface v1.14.1
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	github.com/streamingfast/bstream v0.0.2-0.20220328194113-e36a9f64f958
	github.com/streamingfast/cli v0.0.4-0.20220113202443-f7bcefa38f7e
	github.com/streamingfast/client-go v0.2.1-0.20220328132410-afd23e7857ca
	github.com/streamingfast/dauth v0.0.0-20220307162109-cca1810ae757
	github.com/streamingfast/dbin v0.9.0
	github.com/streamingfast/derr v0.0.0-20220307162255-f277e08753fa
	github.com/streamingfast/dgrpc v0.0.0-20220307180102-b2d417ac8da7
	github.com/streamingfast/dlauncher v0.0.0-20220307153121-5674e1b64d40
	github.com/streamingfast/dmetering v0.0.0-20220307162406-37261b4b3de9
	github.com/streamingfast/dmetrics v0.0.0-20220307162521-2389094ab4a1
	github.com/streamingfast/dstore v0.1.1-0.20220315134935-980696943a79
	github.com/streamingfast/eth-go v0.0.0-20220312041930-62a1ff104ff6
	github.com/streamingfast/firehose v0.1.1-0.20220328200521-97fd1f6b0e5c
	github.com/streamingfast/jsonpb v0.0.0-20210811021341-3670f0aa02d0
	github.com/streamingfast/logging v0.0.0-20220304214715-bc750a74b424
	github.com/streamingfast/merger v0.0.3-0.20220318152213-9ab5185b44e8
	github.com/streamingfast/node-manager v0.0.2-0.20220319133814-75361e421e41
	github.com/streamingfast/pbgo v0.0.6-0.20220304191603-f73822f471ff
	github.com/streamingfast/relayer v0.0.2-0.20220307182103-5f4178c54fde
	github.com/streamingfast/sf-tools v0.0.0-20220307162924-1a39f7035cd5
	github.com/streamingfast/shutter v1.5.0
	github.com/streamingfast/substreams v0.0.0-20220328171236-39ac0f2ee03d
	github.com/stretchr/testify v1.7.1-0.20210427113832-6241f9ab9942
	github.com/tidwall/gjson v1.12.1
	go.uber.org/goleak v1.1.12 // indirect
	go.uber.org/zap v1.21.0
	golang.org/x/oauth2 v0.0.0-20220223155221-ee480838109b
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
)

replace (
	github.com/ShinyTrinkets/overseer => github.com/streamingfast/overseer v0.2.1-0.20210326144022-ee491780e3ef
	github.com/gorilla/rpc => github.com/streamingfast/rpc v1.2.1-0.20201124195002-f9fc01524e38
	github.com/graph-gophers/graphql-go => github.com/streamingfast/graphql-go v0.0.0-20210204202750-0e485a040a3c
)
