package blockfetcher

import (
	"context"
	"time"

	"go.uber.org/zap"

	pbbstream "github.com/streamingfast/bstream/pb/sf/bstream/v1"
	"github.com/streamingfast/eth-go/rpc"
	"github.com/streamingfast/firehose-ethereum/block"
)

type ArbOneBlockFetcher struct {
	fetcher *BlockFetcher
}

func NewArbOneBlockFetcher(intervalBetweenFetch time.Duration, latestBlockRetryInterval time.Duration, logger *zap.Logger) *OptimismBlockFetcher {
	fetcher := NewBlockFetcher(intervalBetweenFetch, latestBlockRetryInterval, block.RpcToEthBlock, logger)
	return &OptimismBlockFetcher{
		fetcher: fetcher,
	}
}

func (f *ArbOneBlockFetcher) PollingInterval() time.Duration { return 1 * time.Second }

func (f *ArbOneBlockFetcher) Fetch(ctx context.Context, rpcClient *rpc.Client, blockNum uint64) (*pbbstream.Block, error) {
	return f.fetcher.Fetch(ctx, rpcClient, blockNum)
}
