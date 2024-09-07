package okto

const (
	NATIVE_TOKEN_ADDR = "0x0000000000000000000000000000000000000000"
)

var (
	// TODO: use okto /aupported_tokens and /supported_networks api's
	// do not hardcode networks and tokens
	// for now, all networks returnd by /supported_networks do not work. ex: solana, osmosis
	// an array. do not handle each network separately. do not use enum to treat as first class attributes. okotron is network agnostic
	SUPPORTED_TOKENS = []string{
		"ETH",
		"MATIC",
		// "USDC",
		// "USDT",
		"WMATIC",
		"DAI",
		"WETH",
	}

	SUPPORTED_NETWORKS = map[string][]string{
		"ETH":    {"BASE"},
		"MATIC":  {"POLYGON"},
		// "USDC":   {"BASE", "POLYGON"},
		// "USDT":   {"POLYGON"},
		"WMATIC": {"POLYGON"},
		"DAI":  {"POLYGON"},
		"WETH": {"POLYGON"},
	}

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
		// "USDC":  6,
		//"USDT":  6,
		"WMATIC": 18,
		"DAI": 18,
		"WETH": 18,
	}

	// see the okto supported tokens and supported networks
	// these values are from okto supported tokens api
	TOKEN_TO_NETWORK_TO_ADDRESS = map[string]map[string]string{
		"ETH": {
			"BASE": NATIVE_TOKEN_ADDR,
		},
		"MATIC": {
			"POLYGON": NATIVE_TOKEN_ADDR,
		},
		// "USDC": {
		// 	"BASE":    "0xd9aaec86b65d86f6a7b5b1b0c42ffa531710b6ca",
		// 	"POLYGON": "0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
		// },
		// "USDT": {
		// 	"POLYGON": "0xc2132d05d31c914a87c6611c10748aeb04b58e8f",
		// },
		"WMATIC": {
			"POLYGON": "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
		},
		"DAI": {
			"POLYGON": "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063",
		},
		"WETH": {
			"POLYGON": "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619",
		},
	}
)
