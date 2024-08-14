package generator

import (
	"crypto/ecdsa"
	"log"
	"net/http"
	"prometheus/pkg/utils"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/sha3"
)

type Body struct {
	Prefix  string `json:"prefix"`
	Suffix  string `json:"suffix"`
	Threads int    `json:"threads"`
}

func Generate(c *gin.Context) {
	var b Body
	c.BindJSON(&b)

	log.Println("Prefix:", b.Prefix)
	log.Println("Suffix:", b.Suffix)

	for {
		privKey, err := crypto.GenerateKey()
		utils.Check(err)
		privKeyBytes := crypto.FromECDSA(privKey)

		pubKey := privKey.Public()
		pubKeyECDSA, _ := pubKey.(*ecdsa.PublicKey)
		pubKeyBytes := crypto.FromECDSAPub(pubKeyECDSA)

		// addr := crypto.PubkeyToAddress(*pubKeyECDSA)

		hash := sha3.NewLegacyKeccak256()
		hash.Write(pubKeyBytes[1:])

		genAddr := hexutil.Encode(hash.Sum(nil)[12:])
		if strings.HasSuffix(genAddr, b.Suffix) && strings.HasPrefix(genAddr, "0x"+b.Prefix) {
			// if strings.HasSuffix(genAddr, s) {
			log.Printf("Private key: %s\n", hexutil.Encode(privKeyBytes)[2:])
			log.Printf("Public key: %s\n", hexutil.Encode(pubKeyBytes)[4:])
			// log.Printf("Public address (from ECDSA): %s\n", addr)
			log.Printf("Public address (hash of public key):\t%s\n\n", genAddr)
			c.JSON(http.StatusOK, gin.H{
				"message": b,
				"status":  "OK",
			})
		}
	}

}
