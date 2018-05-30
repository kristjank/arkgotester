package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/kristjank/ark-go/arkcoin"
	"github.com/kristjank/ark-go/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

func getRandomSender() (string, string) {
	return viper.GetString("account.passphrase1"), viper.GetString("account.secondPassphrase1")
}

func fillTransactions(randomFees bool) {
	f, _ := os.Create("passphrases.log")
	defer f.Close()

	if viper.GetBool("env.singlePeerTest") {
		log.Info("Single peer mode test active. Peer: ", viper.GetString("env.singlePeerIp"))
		fmt.Println("Single peer mode test active. Peer: ", viper.GetInt("env.singlePeerPort"))

		peer := core.Peer{}
		peer.IP = viper.GetString("env.singlePeerIp")
		peer.Port = viper.GetInt("env.singlePeerPort")
		ArkAPIClient = core.NewArkClientFromPeer(peer)
	}

	testRecord := createTestRecord()
	testRecord.Save()

	testRecord.TestStarted = time.Now()
	for xx := 0; xx < viper.GetInt("env.txIterations"); xx++ {
		payload := core.TransactionPayload{}
		senderP1, senderP2 := getRandomSender()
		for i := 0; i < viper.GetInt("env.txPerPayload"); i++ {
			fee := 0
			if randomFees {
				fee = random(viper.GetInt("env.dynamicFeeMin"), viper.GetInt("env.dynamicFeeMax"))
			} else {
				fee = 10000000
			}

			recepientAddress, recepientPassword := getWallet(getRandomPassword())
			f.WriteString(arkcoin.NewPrivateKeyFromPassword(senderP1, arkcoin.ActiveCoinConfig).PublicKey.Address() + "-" + recepientAddress + "-" + recepientPassword + "\n")
			log.Info("Creating random recepient ", recepientAddress, recepientPassword)
			tx := core.CreateTransaction(recepientAddress,
				int64(i+1),
				viper.GetString("env.txDescription"),
				senderP1,
				senderP2,
				int64(fee))
			//f.WriteString(recepientPassword + "\n")
			payload.Transactions = append(payload.Transactions, tx)
		}

		log.Info("Sending transactions, nr of tx: ", len(payload.Transactions))

		testIterRecord := createTestIterationRecord(testRecord.ID)
		testIterRecord.Save()
		res, httpresponse, err := ArkAPIClient.PostTransaction(payload)
		testIterRecord.IterationStopped = time.Now()

		if res.Success {
			log.Info("Success,", httpresponse.Status, xx)
			testIterRecord.TestStatus = "SUCCESS"
			testIterRecord.TxIDs = res.TransactionIDs
		} else {
			testIterRecord.TestStatus = "FAILED"
			if httpresponse != nil {
				log.Error(res.Message, res.Error, xx)
			}
			log.Error(err.Error(), res.Error)
		}
		testIterRecord.Update()
	}
	testRecord.TestStopped = time.Now()
	testRecord.Update()
	log.Info("The call took %v to run.\n", testRecord.TestStopped.Sub(testRecord.TestStarted))
}

func createDelegates() {
}
