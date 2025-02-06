package main

import (
	"fmt"

	pbbstream "github.com/streamingfast/bstream/pb/sf/bstream/v1"
	firecore "github.com/streamingfast/firehose-core"
	pbeth "github.com/streamingfast/firehose-ethereum/types/pb/sf/ethereum/type/v2"
)

func SanitizeEthereumBlockForCompare(block *pbbstream.Block) *pbbstream.Block {
	untypedEthBlock, err := block.Payload.UnmarshalNew()
	if err != nil {
		panic(fmt.Errorf("unexpected block message %s, unable to unmarshal to proto.Message: %w", block.Payload.TypeUrl, err))
	}

	ethBlock, ok := untypedEthBlock.(*pbeth.Block)
	if !ok {
		panic(fmt.Errorf("unexpected block message %s, unable to cast to pbeth.Block", block.Payload.TypeUrl))
	}

	for _, tx := range ethBlock.TransactionTraces {
		for _, call := range tx.Calls {
			if call.FailureReason != "" {
				call.FailureReason = "<error replaced for comparison>"
			}
		}
	}

	out, err := firecore.EncodeBlock(firecore.BlockEnveloppe{
		Block:  ethBlock,
		LIBNum: block.LibNum,
	})
	if err != nil {
		panic(fmt.Errorf("unable to encode block: %w", err))
	}

	return out
}
