package main

import (
	"fmt"

	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/streamingfast/cli"
	firecore "github.com/streamingfast/firehose-core"
	fhCmd "github.com/streamingfast/firehose-core/cmd"
	"github.com/streamingfast/firehose-core/substreams"
	"github.com/streamingfast/firehose-ethereum/codec"
	ethss "github.com/streamingfast/firehose-ethereum/substreams"
	"github.com/streamingfast/firehose-ethereum/transform"
	pbeth "github.com/streamingfast/firehose-ethereum/types/pb/sf/ethereum/type/v2"
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func main() {
	fhCmd.Main(Chain())
}

var chain *firecore.Chain[*pbeth.Block]

func Chain() *firecore.Chain[*pbeth.Block] {
	if chain != nil {
		return chain
	}

	chain = &firecore.Chain[*pbeth.Block]{
		ShortName:            "eth",
		LongName:             "Ethereum",
		ExecutableName:       "geth",
		FullyQualifiedModule: "github.com/streamingfast/firehose-ethereum",
		Version:              version,

		// Ensure that if you ever modify test, modify also `types/init.go#init` so that the `bstream.InitGeneric` there fits us
		BlockFactory: func() firecore.Block { return new(pbeth.Block) },

		BlockIndexerFactories: map[string]firecore.BlockIndexerFactory[*pbeth.Block]{
			transform.CombinedIndexerShortName: transform.NewEthCombinedIndexer,
		},

		BlockTransformerFactories: map[protoreflect.FullName]firecore.BlockTransformerFactory{
			transform.HeaderOnlyMessageName:     transform.NewHeaderOnlyTransformFactory,
			transform.CombinedFilterMessageName: transform.NewCombinedFilterTransformFactory,

			transform.MultiCallToFilterMessageName: transform.NewMultiCallToFilterTransformFactory,
			transform.MultiLogFilterMessageName:    transform.NewMultiLogFilterTransformFactory,
		},

		ConsoleReaderFactory: codec.NewConsoleReader,

		RegisterExtraStartFlags: func(flags *pflag.FlagSet) {
			// The "\n" is there on purpose to improve readability of the added elements
			flags.String("reader-node-bootstrap-data-url", "", firecore.DefaultReaderNodeBootstrapDataURLFlagDescription()+"\n"+cli.Dedent(`
				If the URL ends with json, it will be treated as a genesis.json file and the node will be bootstrapped with it. The bootstrapping
				is done by calling 'geth init' with the genesis file. If you need more custom logic, think about using 'bash://...' instead
				which execute any bash script and offers more flexibility.
			`)+"\n")

			flags.StringArray("substreams-rpc-endpoints", nil, "Remote endpoints to contact to satisfy Substreams 'eth_call's")
			flags.Uint64("substreams-rpc-gas-limit", 50_000_000, "Gas limit to set when calling RPC (set it to 0 for arbitrum chains, otherwise you should keep 50M)")
			flags.String("substreams-rpc-cache-store-url", "{data-dir}/rpc-cache", "where rpc cache will be store call responses")
			flags.Uint64("substreams-rpc-cache-chunk-size", uint64(1_000), "RPC cache chunk size in block")
		},

		RegisterSubstreamsExtensions: func(chain *firecore.Chain[*pbeth.Block]) ([]substreams.Extension, error) {
			dataDir := viper.GetString("global-data-dir")
			dataDirAbs, err := filepath.Abs(dataDir)
			if err != nil {
				return nil, fmt.Errorf("unable to setup directory structure: %w", err)
			}

			rpcGasLimit := viper.GetUint64("substreams-rpc-gas-limit")
			rpcEndpoints := viper.GetStringSlice("substreams-rpc-endpoints")
			rpcCacheStoreURL := firecore.MustReplaceDataDir(dataDirAbs, viper.GetString("substreams-rpc-cache-store-url"))
			rpcCacheChunkSize := viper.GetUint64("substreams-rpc-cache-chunk-size")
			rpcEngine, err := ethss.NewRPCEngine(
				rpcCacheStoreURL,
				rpcEndpoints,
				rpcCacheChunkSize,
				rpcGasLimit,
			)

			if err != nil {
				return nil, fmt.Errorf("creating rpc engine: %w", err)
			}

			return []substreams.Extension{
				{
					PipelineOptioner: rpcEngine,
					WASMExtensioner:  rpcEngine,
				},
			}, nil
		},

		ReaderNodeBootstrapperFactory: firecore.DefaultReaderNodeBootstrapper(newReaderNodeBootstrapper),

		Tools: &firecore.ToolsConfig[*pbeth.Block]{

			RegisterExtraCmd: func(chain *firecore.Chain[*pbeth.Block], parent *cobra.Command, zlog *zap.Logger, tracer logging.Tracer) error {
				parent.AddCommand(compareOneblockRPCCmd)
				parent.AddCommand(newCompareBlocksRPCCmd(zlog))
				parent.AddCommand(newFixOrdinalsCmd(zlog))
				parent.AddCommand(newFixAnyTypeCmd(zlog))
				parent.AddCommand(newPollRPCBlocksCmd(zlog))
				parent.AddCommand(newPollerCmd(zlog, tracer))
				parent.AddCommand(newOptimismPollerCmd(zlog, tracer))

				registerGethEnforcePeersCmd(parent, chain.BinaryName(), zlog, tracer)

				return nil
			},

			TransformFlags: &firecore.TransformFlags{
				Register: func(flags *pflag.FlagSet) {
					flags.Bool("header-only", false, "Apply the HeaderOnly transform sending back Block's header only (with few top-level fields), exclusive option")
					flags.String("call-filters", "", "call filters (format: '[address1[+address2[+...]]]:[eventsig1[+eventsig2[+...]]]")
					flags.String("log-filters", "", "log filters (format: '[address1[+address2[+...]]]:[eventsig1[+eventsig2[+...]]]")
					flags.Bool("send-all-block-headers", false, "ask for all the blocks to be sent (header-only if there is no match)")
				},

				Parse: parseTransformFlags,
			},
		},
	}
	return chain
}

// Version value, injected via go build `ldflags` at build time, **must** not be removed or inlined
var version = "dev"
