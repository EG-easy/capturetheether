package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/ethclient"
)

var apiKey string = ""
var privateKey string = ""
var contractAddress string = ""

func main() {
	sendExternalRawTransaction("your nick name")
}

func sendExternalRawTransaction(nickName string) (transaction string) {
	//Ropstenネットワークに接続
	client, err := ethclient.Dial("https://ropsten.infura.io/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}

	//PrivateKeyを読み込む
	privateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	//　PrivateKeyからPublickeyを生成
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}

	//PublicKeyから、送金主アドレスを生成
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//Ropstenネットワークから、Nonce情報を読み取る
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	//トークン送金Transactionをテストネット送るためのgasLimit、
	value := big.NewInt(0) //（オプション）後で使用する関数NewTransactionの引数で必要になるため設定。Transactionと同時に送るETHの量を設定できます。
	gasLimit := uint64(2000000)

	//ロプステンネットワークから、現在のgasPriceを取得。トランザクションがマイニングされずに放置されることを防ぐ。
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	//トークンコントラクトアドレスを指定
	tokenAddress := common.HexToAddress(contractAddress)
	//ERC20のどの関数を使用するか指定。https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendtransaction
	transferFnSignature := []byte("setNickname(bytes32)")
	//hash化し、先頭から4バイトまで取得。これで使用する関数を指定したことになる。
	hash := sha3.NewKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	NickName := []byte(nickName)

	//0埋め
	paddedNickName := common.RightPadBytes(NickName, 32)

	//トランザクションで送るデータを作成
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedNickName...)

	/***** Preparing signed transaction *****/
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	//サインしたトランザクションをRopstenNetworkに送る。
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Signed tx sent: %s", signedTx.Hash().Hex())

	return signedTx.Hash().Hex()
}
