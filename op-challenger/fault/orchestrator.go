package fault

import (
	"context"
	"fmt"
)

type Orchestrator struct {
	agents    []Agent
	outputChs []chan Claim
	responses chan Claim
}

func NewOrchestrator(maxDepth int, traces []TraceProvider, root, counter Claim) Orchestrator {
	o := Orchestrator{
		responses: make(chan Claim, 100),
		outputChs: make([]chan Claim, len(traces)),
		agents:    make([]Agent, len(traces)),
	}
	for i, trace := range traces {
		game := NewGameState()
		game.Put(root)
		game.Put(counter)
		o.agents[i] = NewAgent(game, maxDepth, trace, &o)
		o.outputChs[i] = make(chan Claim)
	}
	return o
}

func (o *Orchestrator) Respond(_ context.Context, response Claim) error {
	o.responses <- response
	return nil
}

func (o *Orchestrator) Start() {
	// TODO handle shutdown
	for i := 0; i < len(o.agents); i++ {
		go runAgent(&o.agents[i], o.outputChs[i])
	}
	o.reponderThread()
}

func runAgent(agent *Agent, claimCh <-chan Claim) {
	for {
		agent.PerformActions()
		// TODO: Multiple claims / how to balance performing actions with
		// accepting new claims
		claim := <-claimCh
		agent.AddClaim(claim)

	}
}

func (o *Orchestrator) reponderThread() {
	for {
		resp := <-o.responses
		PrettyPrintAlphabetClaim("Got response", resp)
		for _, ch := range o.outputChs {
			// Copy it. Should be immutable, but be sure.
			resp := resp
			ch <- resp
		}
	}
}

func PrettyPrintAlphabetClaim(name string, claim Claim) {
	value := claim.Value
	idx := value[30]
	letter := value[31]
	if claim.IsRoot() {
		fmt.Printf("%s\ttrace %v letter %c\n", name, idx, letter)
	} else {
		fmt.Printf("%s\ttrace %v letter %c is attack %v\n", name, idx, letter, !claim.DefendsParent())
	}

}
