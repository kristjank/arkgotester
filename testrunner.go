package main

import (
	"fmt"
	"time"

	"github.com/kristjank/ark-go/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func runTests() {
	testRecord := createTestRecord()
	testRecord.Save()

	testRecord.TestStarted = time.Now()
	for xx := 0; xx < viper.GetInt("env.txIterations"); xx++ {
		ArkAPIClient = ArkAPIClient.SetActiveConfiguration(core.DEVNET)
		if viper.GetBool("env.singlePeerTest") {
			log.Info("Single peer mode test active. Peer: ", viper.GetString("env.singlePeerIp"))
			fmt.Println("Single peer mode test active. Peer: ", viper.GetInt("env.singlePeerPort"))

			peer := core.Peer{}
			peer.IP = viper.GetString("env.singlePeerIp")
			peer.Port = viper.GetInt("env.singlePeerPort")
			ArkAPIClient = core.NewArkClientFromPeer(peer)
			//ArkAPIClient = core.NewArkClientFromIP(viper.GetString("env.singlePeerIp"))
		}

		payload := core.TransactionPayload{}

		for i := 0; i < viper.GetInt("env.txPerPayload"); i++ {
			tx := core.CreateTransaction(viper.GetString("account.recepient"),
				int64(i+1),
				viper.GetString("env.txDescription"),
				viper.GetString("account.passphrase"), viper.GetString("account.secondPassphrase"))
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

func runTestsParallel() {
	testRecord := createTestRecord()
	testRecord.Save()

	testRecord.TestStarted = time.Now()
	for xx := 0; xx < viper.GetInt("env.txIterations"); xx++ {
		go func() {
			ArkAPIClient = ArkAPIClient.SetActiveConfiguration(core.DEVNET)
			if viper.GetBool("env.singlePeerTest") {
				log.Info("Single peer mode test active. Peer: ", viper.GetString("env.singlePeerIp"))
				fmt.Println("Single peer mode test active. Peer: ", viper.GetString("env.singlePeerIp"))
				ArkAPIClient = core.NewArkClientFromIP(viper.GetString("env.singlePeerIp"))
			}

			payload := core.TransactionPayload{}

			for i := 0; i < viper.GetInt("env.txPerPayload"); i++ {
				tx := core.CreateTransaction(viper.GetString("account.recepient"),
					int64(i+1),
					viper.GetString("env.txDescription"),
					viper.GetString("account.passphrase"), viper.GetString("account.secondPassphrase"))
				payload.Transactions = append(payload.Transactions, tx)
			}

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
		}()
	}
	testRecord.TestStopped = time.Now()
	testRecord.Update()
	log.Info("The call took %v to run.\n", testRecord.TestStopped.Sub(testRecord.TestStarted))

}
