package main

import (
	"log"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"fmt"
	"context"
	"github.com/ethereum/go-ethereum/common"
)

var apiKey string = ""

func main(){
	number := big.NewInt(4688441)
	//Ropstenネットワークに接続
	client := ConnectClient()
	blockInfo, err := client.BlockByNumber(context.Background(), number)
	if err != nil{
		log.Fatal("can not get block info:", err)
	}

	previousBlockHash := blockInfo.Header().ParentHash
	timestamp := blockInfo.Header().Time
	//32byteになるようにpaddingをする
	paddedTimestamp := common.LeftPadBytes(timestamp.Bytes(), 32)

	hash := sha3.NewKeccak256()
	var buf []byte
	//hash.Write([]byte)
	var data []byte
	data = append(data, previousBlockHash[:]...)
	data = append(data, paddedTimestamp...)

	hash.Write(data)
	buf = hash.Sum(buf)

	//keccak256の最後の1byte分がuint8に当たる
	answer := buf[len(buf)-1]

	fmt.Println(answer)
}

func ConnectClient() *ethclient.Client {
	client, err := ethclient.Dial("https://ropsten.infura.io/" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	return client
}