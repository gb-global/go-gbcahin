package discv5

import "gbchain-org/go-gbchain/metrics"

var (
	ingressTrafficMeter = metrics.NewRegisteredMeter("discv5/InboundTraffic", nil)
	egressTrafficMeter  = metrics.NewRegisteredMeter("discv5/OutboundTraffic", nil)
)
