package metablockchain

import (
	"errors"
)

const APIClientHttpMode = "http"

type ApiClient struct {
	Client    *Client
	ExplorerApiClient *ExplorerApiClient
	APIChoose string
}

func NewApiClient(wm *WalletManager) error {
	api := ApiClient{}

	if len(wm.Config.APIChoose) == 0 {
		wm.Config.APIChoose = APIClientHttpMode //默认采用rpc连接
	}
	api.APIChoose = wm.Config.APIChoose
	if api.APIChoose == APIClientHttpMode {
		api.Client = NewClient(wm.Config.NodeAPI, false, wm.Symbol() )
		api.ExplorerApiClient = NewExplorerApiClient(wm.Config.ExplorerAPI, false, wm.Symbol())
	}

	wm.ApiClient = &api

	return nil
}

// 获取当前最高区块
func (c *ApiClient) getBlockHeight() (uint64, error) {
	var (
		currentHeight uint64
		err           error
	)
	if c.APIChoose == APIClientHttpMode {
		currentHeight, err = c.Client.getBlockHeight()
	}

	return currentHeight, err
}

// 获取地址余额
func (c *ApiClient) getBalance(did string) (*AddrBalance, error) {
	var (
		balance *AddrBalance
		err     error
	)

	if c.APIChoose == APIClientHttpMode {
		balance, err = c.Client.getBalance(did)
	}

	return balance, err
}

// 获取地址对应的did
func (c *ApiClient) getDidByAddress(address string) (string, error) {
	var (
		err     error
	)

	result := ""

	if c.APIChoose == APIClientHttpMode {
		result, err = c.Client.GetDidByAddress(address)
	}

	if result==""{
		return "", errors.New("found did faild, address=" + address)
	}

	return result, err
}

// 获取地址对应的did
func (c *ApiClient) getEra(toaddress string, amount uint64, memo string, blocknumber uint64) (string, error) {
	var (
		err     error
	)

	result := ""

	if c.APIChoose == APIClientHttpMode {
		result, err = c.Client.GetEra(toaddress, amount, memo, blocknumber)
	}

	if result==""{
		return "", err
	}

	return result, err
}

func (c *ApiClient) getBlockByHeight(height uint64) (*Block, error) {
	var (
		block *Block
		err   error
	)
	if c.APIChoose == APIClientHttpMode {
		block, err = c.Client.getBlockByHeight(height)
		if err!=nil {
			return nil, err
		}

		for _, transaction := range block.Transactions {
			err = c.ExplorerApiClient.setTransactionStatus(&transaction)
			if err!=nil {
				return nil, err
			}
		}
	}

	return block, err
}

func (c *ApiClient) sendTransaction(rawTx string) (string, error) {
	var (
		txid string
		err  error
	)
	if c.APIChoose == APIClientHttpMode {
		txid, err = c.Client.sendTransaction(rawTx)
	}

	return txid, err
}

func (c *ApiClient) getRuntimeVersion() (*RuntimeVersion, error){
	var (
		result    *RuntimeVersion
		err         error
	)
	//if c.APIChoose == APIClientHttpMode {
	//	result, err = c.RpcClient.GetRuntimeVersion()
	//}

	return result, err
}

//获取当前最新高度
func (c *ApiClient) getMostHeightBlock() (*Block, error) {
	var (
		mostHeightBlock *Block
		err             error
	)
	if c.APIChoose == APIClientHttpMode {
		mostHeightBlock, err = c.Client.getMostHeightBlock()
	}

	return mostHeightBlock, err
}

func (c *ApiClient) getTxMaterial() (*TxMaterial, error) {
	var (
		txMaterial *TxMaterial
		err         error
	)
	if c.APIChoose == APIClientHttpMode {
		txMaterial, err = c.Client.getTxMaterial()
	}

	return txMaterial, err
}
