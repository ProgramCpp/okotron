package okto

import (
	"math/big"
	// "strings"
	"time"

	// "github.com/ethereum/go-ethereum/accounts/abi"
	// "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// const RPC_URL = "https://polygon-rpc.com"

// ERC20ApproveABI is the ERC20 function signature for "approve(address,uint256)"
const ERC20ApproveABI = `[{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"type":"function"}]`

// Function to approve ERC20 token transfer using eth_sendTransaction
func ApproveTokenTransfer(
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
	// parsedABI, err := abi.JSON(strings.NewReader(ERC20ApproveABI))
	// if err != nil {
	// 	return errors.Wrap(err, "Failed to parse the ERC20 ABI")
	// }

	// // Encode the approve function call
	// data, err := parsedABI.Pack("approve", common.HexToAddress(spenderAddress), amountToApprove)
	// if err != nil {
	// 	return errors.Wrap(err, "Failed to encode the function call")
	// }

	// hardcoding data for now. this supports base and polygon networks
	// approval for 100 * 10^18 for spender address 0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE
	// upto 100 of any erc20 e.g. 100 USDC, 100 WETH
	resData, err := RawTxn(authToken, RawTxPayload{
		NetworkName: networkName,
		Transaction: Transaction{
			From:  fromAddress,
			To:    contractAddress,
			Data:  "0x095ea7b30000000000000000000000001231deb6f5749ef6ce6943a275a1d3e7486f4eae0000000000000000000000000000000000000000000000056bc75e2d63100000",
			Value: "0x0",
		},
	})
	if err != nil {
		return errors.Wrap(err, "error calling okto raw txn")
	}

	// poll for success
	// maximum 10 calls with backoff. maximum wait of 110s(arithmatic progression of multiples of 2)
	for i := 1 ; i <= 10; i++{
		time.Sleep(time.Duration(i * 2) * time.Second)
		err = RawTxnStatus(authToken, resData.JobId)
		if err != nil && errors.Is(err, TXN_IN_PROGRESS){
			continue
		} else if err != nil {
			return errors.Wrap(err, "txn approval was unsuccessful")
		} else {
			return nil // success
		}
	}

	return errors.New("approval txn failed after max wait duration")
}
