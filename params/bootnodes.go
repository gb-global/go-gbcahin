package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main GBChain network.
var MainnetBootnodes = []string{
	//GBChain Foundation Go Bootnodes
	//TODO MainnetBootnodes
	"enode://781b3a4cbd3474d2695cd0d842dca092f713a23cac3cdcf16a5c45c1a4e8c7216aff14b617b3ddbeed54b627a83b9e03f161a30ae166082e4e3e6e8923660c0d@81.209.211.154:30312", //JPN3


}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// test network.
var TestnetBootnodes = []string{
	"enode://bc6858a4de55d8715834a203def74162474e6ff8062c30093def22d577bcc96cd4755e9738be6b0c7e6f3ee7fcec5cc84a7c94b509692737e0744ada8bbde507@87.110.48.207:30312", // CN
	"enode://c72b5cb21086dac58bb9235bc68b217475e050e9c8c2a827867242193deb68a9c6abe13fd8da7cb64d3c2eb1d7ce6e4cdf5f48cd174b772934ef2446a21136a8@97.74.52.42:30312",   // JPN
	"enode://2e1162b335c72cfd767d2dffe617df942b9f71817557fffb28b24bff2aff5f2a18881ec7b58578498985400816e3fd62dcceed8cf842b9fd7dfa2fcbb464dea0@107.88.58.252:30312",  // US
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	"enode://b9f34d999d0a719967f2b3e55f34b3938a9ff4c0c87e8064a3cd4102ad54ea89834f881177ffa0759e298c3e7e561426d366183836d8c81b0c7fb520fedf73db@47.115.87.110:30312",  // US

}
