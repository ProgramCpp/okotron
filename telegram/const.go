package telegram

// primary commands
const (
	CMD_LOGIN       = "/login"
	CMD_PORTFOLIO   = "/portfolio"
	CMD_SWAP        = "/swap"
	CMD_LIMIT_ORDER = "/limit-order"
)

// sub commands that can be executed after the primary commands
const (
	// TODO: fix the weird naming convention to pass context from commands and sub commands
	CMD_LOGIN_CMD_SETUP_PROFILE = "/login/setup-profile"
	CMD_ANY_CMD_FROM_TOKEN      = "/any/from-token"
	CMD_ANY_CMD_FROM_NETWORK    = "/any/from-network"
	CMD_ANY_CMD_TO_TOKEN        = "/any/to-token"
	CMD_ANY_CMD_TO_NETWORK      = "/any/to-network"
	CMD_ANY_CMD_QUANTITY        = "/any/quantity"
)

var (
	// TODO: use okto /aupported_tokens and /supported_networks api's
	// do not hardcode networks and tokens
	// for now, all networks returnd by /supported_networks do not work. ex: solana, osmosis
	// an array. do not handle each network separately. do not use enum to treat as first class attributes. okotron is network agnostic
	SUPPORTED_TOKENS   = []string{"APT", "ETH", "MATIC", "USDC", "USDT"}
	SUPPORTED_NETWORKS = map[string][]string{
		"APT":   {"APTOS"},
		"ETH":   {"BASE"},
		"MATIC": {"POLYGON"},
		"USDC":  {"BASE", "POLYGON"},
		"USDT":  {"POLYGON", "APTOS"},
	}
)