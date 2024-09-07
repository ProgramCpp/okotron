package okto

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// const RPC_URL = "https://polygon-rpc.com"

// ERC20ApproveABI is the ERC20 function signature for "approve(address,uint256)"
const ERC20ApproveABI = `[{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"type":"function"}]`

// Function to approve ERC20 token transfer using eth_sendTransaction
func approveTokenTransfer(
	authToken string,
	networkName string,
	contractAddress string, // The ERC20 token contract address
	spenderAddress string, // The address that will be allowed to spend the tokens
	amountToApprove *big.Int, // The amount of tokens to approve
	fromAddress string, // The address sending the transaction (must be unlocked)
) error {

	// Connect to the Ethereum client
	// client, err := ethclient.Dial(RPC_URL)
	// if err != nil {
	// 	return errors.Wrap(err, "Failed to connect to the Ethereum client")
	// }

	// Load the ABI of the contract
	parsedABI, err := abi.JSON(strings.NewReader(ERC20ApproveABI))
	if err != nil {
		return errors.Wrap(err, "Failed to parse the ERC20 ABI")
	}

	// Encode the approve function call
	data, err := parsedABI.Pack("approve", common.HexToAddress(spenderAddress), amountToApprove)
	if err != nil {
		return errors.Wrap(err, "Failed to encode the function call")
	}

	resData, err := RawTxn(authToken, RawTxPayload{
		NetworkName: networkName,
		Transaction: Transaction{
			From:  fromAddress,
			To:    contractAddress,
			Data:  string(data),
			Value: "0x0",
		},
	})
	if err != nil {
		return errors.Wrap(err, "error calling okto raw txn")
	}

	// poll for success
	for {
		err = RawTxnStatus(authToken, resData.JobId)
		if err != nil {
			return errors.Wrap(err, "txn approval was unsuccessful")
		}
	}
}
