package okto

var (
	// TODO: get this mapping dynalically at app start and cache it
	// okto supports only a sub set of networks returned  by /supported-networks. ex, solana is not supported
	NETWORK_NAME_TO_CHAIN_ID = map[string]string{
		"BASE":    "8453",
		"POLYGON": "137",
	}
)
