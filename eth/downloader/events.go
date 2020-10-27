package downloader

import "gbchain-org/go-gbchain/core/types"

type DoneEvent struct {
	Latest *types.Header
}
type StartEvent struct{}
type FailedEvent struct{ Err error }
