package core

import "gbchain-org/go-gbchain/common"

func (c *core) handleFinalCommitted() error {
	logger := c.logger.New("state", c.state)
	logger.Trace("Received a final committed proposal")
	c.startNewRound(common.Big0)
	return nil
}
