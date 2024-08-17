package telegram

// primary commands
const (
	CMD_LOGIN       = "/login"
	CMD_PORTFOLIO   = "/portfolio"
	CMD_SWAP        = "/swap"
	CMD_LIMIT_ORDER = "/limitorder"
)

// sub commands that can be executed after the primary commands
const (
	// TODO: fix the weird naming convention to pass context from commands and sub commands
	CMD_LOGIN_CMD_SETUP_PROFILE = "/login/setup-profile"

	CMD_SWAP_CMD_FROM_TOKEN   = "/swap/from-token"
	CMD_SWAP_CMD_FROM_NETWORK = "/swap/from-network"
	CMD_SWAP_CMD_TO_TOKEN     = "/swap/to-token"
	CMD_SWAP_CMD_TO_NETWORK   = "/swap/to-network"
	CMD_SWAP_CMD_QUANTITY     = "/swap/quantity"

	CMD_LIMIT_ORDER_CMD_BUY_OR_SELL  = "/limit-order/buy-or-sell"
	CMD_LIMIT_ORDER_CMD_FROM_TOKEN   = "/limit-order/from-token"
	CMD_LIMIT_ORDER_CMD_FROM_NETWORK = "/limit-order/from-network"
	CMD_LIMIT_ORDER_CMD_TO_TOKEN     = "/limit-order/to-token"
	CMD_LIMIT_ORDER_CMD_TO_NETWORK   = "/limit-order/to-network"
	CMD_LIMIT_ORDER_CMD_QUANTITY     = "/limit-order/quantity"
	CMD_LIMIT_ORDER_CMD_PRICE        = "/limit-order/price"
)

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
