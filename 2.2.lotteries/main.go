package main

import (
	"fmt"
	"encoding/hex"
	"encoding/binary"
	"bytes"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"math"
)

func main() {
	//float64からintに変換
	limit := int(math.Pow(2, 8))
	targetHash := "db81b4d58595fbbbb592d3661a34cdca14d7ab379441400cbfa1b78bc447c365"

	for i := 0; i <limit ; i++ {

		buf := new(bytes.Buffer)
		num := uint8(i) //intからuint8に変換
		//uint8から[]byteに変換
		err := binary.Write(buf, binary.LittleEndian, num)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}

		hash := sha3.NewKeccak256()

		var hashByte []byte
		//hash.Write([]byte)
		hash.Write(buf.Bytes())
		hashByte = hash.Sum(hashByte)

		if hex.EncodeToString(hashByte) == targetHash {
			fmt.Println(i)
			break
		}
	}
}