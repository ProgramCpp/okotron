package telegram

// primary commands
const (
	CMD_LOGIN       = "/login"
	CMD_PORTFOLIO   = "/portfolio"
	CMD_TRANSFER    = "/transfer"
	CMD_SWAP        = "/swap"
	CMD_LIMIT_ORDER = "/limitorder"
	CMD_COPY_TRADE  = "/copytrade"
)

// sub commands that can be executed after the primary commands
const (
	// TODO: fix the weird naming convention to pass context from commands and sub commands
	CMD_LOGIN_CMD_SETUP_PROFILE = "/login/setup-profile"

	CMD_TRANSFER_CMD_FROM_TOKEN   = "/transfer/from-token"
	CMD_TRANSFER_CMD_FROM_NETWORK = "/transfer/from-network"
	CMD_TRANSFER_CMD_QUANTITY     = "/transfer/quantity"
	CMD_TRANSFER_CMD_ADDRESS      = "/transfer/address"

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

	CMD_COPY_TRADE_CMD_ADDRESS = "/copy-trade/address"
)
