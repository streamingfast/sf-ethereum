package main

import (
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/streamingfast/cli/sflags"
	"github.com/streamingfast/eth-go/rpc"
	firecore "github.com/streamingfast/firehose-core"
	"github.com/streamingfast/firehose-core/blockpoller"
	firecorerpc "github.com/streamingfast/firehose-core/rpc"
	"github.com/streamingfast/firehose-ethereum/blockfetcher"
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
)

func newPollerCmd(logger *zap.Logger, tracer logging.Tracer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "poller",
		Short: "poll blocks from different sources",
	}

	cmd.AddCommand(newOptimismPollerCmd(logger, tracer))
	cmd.AddCommand(newArbOnePollerCmd(logger, tracer))
	cmd.AddCommand(newGenericEVMPollerCmd(logger, tracer))
	return cmd
}

func newOptimismPollerCmd(logger *zap.Logger, tracer logging.Tracer) *cobra.Command {
	// identical as generic-evm for now
	cmd := &cobra.Command{
		Use:   "optimism <rpc-endpoint> <first-streamable-block>",
		Short: "poll blocks from optimism rpc",
		Args:  cobra.ExactArgs(2),
		RunE:  pollerRunE(logger, tracer),
	}
	cmd.Flags().Duration("interval-between-fetch", 0, "interval between fetch")

	return cmd
}
func newArbOnePollerCmd(logger *zap.Logger, tracer logging.Tracer) *cobra.Command {
	// identical as generic-evm for now
	cmd := &cobra.Command{
		Use:   "arb-one <rpc-endpoint> <first-streamable-block>",
		Short: "poll blocks from arb-one rpc",
		Args:  cobra.ExactArgs(2),
		RunE:  pollerRunE(logger, tracer),
	}
	cmd.Flags().Duration("interval-between-fetch", 0, "interval between fetch")

	return cmd
}
func newGenericEVMPollerCmd(logger *zap.Logger, tracer logging.Tracer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generic-evm <rpc-endpoint> <first-streamable-block>",
		Short: "poll blocks from a generic EVM RPC endpoint",
		Args:  cobra.ExactArgs(2),
		RunE:  pollerRunE(logger, tracer),
	}
	cmd.Flags().Duration("interval-between-fetch", 0, "interval between fetch")
	cmd.Flags().Duration("max-block-fetch-duration", 5*time.Second, "maximum delay before retrying a block fetch")

	return cmd
}

func pollerRunE(logger *zap.Logger, tracer logging.Tracer) firecore.CommandExecutor {
	return func(cmd *cobra.Command, args []string) (err error) {
		rpcEndpoint := args[0]
		//dataDir := cmd.Flag("data-dir").Value.String()
		fetchInterval := sflags.MustGetDuration(cmd, "interval-between-fetch")
		maxBlockFetchDuration := sflags.MustGetDuration(cmd, "max-block-fetch-duration")

		dataDir := sflags.MustGetString(cmd, "data-dir")
		stateDir := path.Join(dataDir, "poller-state")

		firstStreamableBlock, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("unable to parse first streamable block %d: %w", firstStreamableBlock, err)
		}

		logger.Info("launching firehose-ethereum poller",
			zap.String("rpc_endpoint", rpcEndpoint),
			zap.String("data_dir", dataDir),
			zap.String("state_dir", stateDir),
			zap.Duration("fetch_interval", fetchInterval),
			zap.Duration("max_block_fetch_duration", maxBlockFetchDuration),
			zap.Uint64("first_streamable_block", firstStreamableBlock),
		)
		rpcClients := firecorerpc.NewClients[*rpc.Client](maxBlockFetchDuration, firecorerpc.NewStickyRollingStrategy[*rpc.Client](), logger)
		rpcClients.Add(rpc.NewClient(rpcEndpoint))

		fetcher := blockfetcher.NewOptimismBlockFetcher(fetchInterval, 1*time.Second, logger)
		handler := blockpoller.NewFireBlockHandler("type.googleapis.com/sf.ethereum.type.v2.Block")
		poller := blockpoller.New[*rpc.Client](fetcher, handler, rpcClients, blockpoller.WithStoringState[*rpc.Client](stateDir), blockpoller.WithLogger[*rpc.Client](logger))

		err = poller.Run(firstStreamableBlock, nil, 1)
		if err != nil {
			return fmt.Errorf("running poller: %w", err)
		}

		return nil
	}
}
