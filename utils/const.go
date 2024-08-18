package utils

var (
	// TODO: use okto /aupported_tokens and /supported_networks api's
	// do not hardcode networks and tokens
	// for now, all networks returnd by /supported_networks do not work. ex: solana, osmosis
	// an array. do not handle each network separately. do not use enum to treat as first class attributes. okotron is network agnostic
	SUPPORTED_TOKENS   = []string{"ETH", "MATIC", "USDC", "USDT"}
	SUPPORTED_NETWORKS = map[string][]string{
		"ETH":   {"BASE"},
		"MATIC": {"POLYGON"},
		"USDC":  {"BASE", "POLYGON"},
		"USDT":  {"POLYGON"},
	}
)