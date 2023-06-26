package main

import (
	"github.com/ethereum-optimism/optimism/op-challenger/fault"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	// examples.SolverExampleOne()

	canonical := "abcdefgh"
	disputed := "abcdexyz"
	maxDepth := 3
	canonicalProvider := fault.NewAlphabetProvider(canonical, uint64(maxDepth))
	disputedProvider := fault.NewAlphabetProvider(disputed, uint64(maxDepth))

	root := fault.Claim{
		ClaimData: fault.ClaimData{
			Value:    common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000077a"),
			Position: fault.NewPosition(0, 0),
		},
	}
	counter := fault.Claim{
		ClaimData: fault.ClaimData{
			Value:    common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000364"),
			Position: fault.NewPosition(1, 0),
		},
		Parent: root.ClaimData,
	}

	o := fault.NewOrchestrator(maxDepth, []fault.TraceProvider{canonicalProvider, disputedProvider}, root, counter)
	o.Start()

	// examples.PositionExampleOne()
	// examples.PositionExampleTwo()
}
