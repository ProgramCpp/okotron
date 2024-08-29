package okto

var (
	// TODO: get this mapping dynamically at app start and cache it
	// okto supports only a sub set of networks returned  by /supported-networks. ex, solana is not supported
	// these values are from okto /supported-networks call
	NETWORK_NAME_TO_CHAIN_ID = map[string]string{
		"BASE":    "8453",
		"POLYGON": "137",
	}

	// these values are from lifi supported tokens api
	TOKEN_TO_DECIMALS = map[string]int{
		"ETH":   18,
		"MATIC": 18,
		"USDC":  6,
		"USDT":  6,
	}

	// see the okto supported tokens and supported networks
	// these values are from okto supported tokens api
	TOKEN_TO_NETWORK_TO_ADDRESS = map[string]map[string]string{
		"ETH": {
			"BASE": " ", // space input for okto
		},
		"MATIC": {
			"POLYGON": " ",
		},
		"USDC": {
			"BASE":    "0xd9aaec86b65d86f6a7b5b1b0c42ffa531710b6ca",
			"POLYGON": "0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
		},
		"USDT": {
			"POLYGON": "0xc2132d05d31c914a87c6611c10748aeb04b58e8f",
		},
	}
)
