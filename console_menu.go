package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/kristjank/ark-go/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//////////////////////////////////////////////////////////////////////////////
//GUI RELATED STUFF
func pause() {
	color.Set(color.FgHiYellow)
	fmt.Println("")
	fmt.Print("Press 'ENTER' key to continue... ")
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	ConsoleReader.ReadString('\n')
}

func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()

}

func printNetworkInfo() {
	color.Set(color.FgHiCyan)

	fmt.Println("Connected on ", core.EnvironmentParams.Network.Token, " peer:", core.BaseURL, "| ARKGoTester version", ArkGoTesterVersion)
	log.Info("Connected on ", core.EnvironmentParams.Network.Token, " peer:", core.BaseURL, "| ARKGoTester version", ArkGoTesterVersion)

}

func printBanner() {
	color.Set(color.FgHiGreen)
	dat, _ := ioutil.ReadFile("cfg/banner.txt")
	fmt.Print(string(dat))
}

func printMenu() {
	log.Info("--------- MAIN MENU ----------------")
	clearScreen()
	printBanner()
	printNetworkInfo()
	color.Set(color.FgHiYellow)
	fmt.Println("")
	fmt.Println("\t1-Deliver payload [", viper.GetInt("env.txPerPayload"), "/", viper.GetInt("env.txIterations"), "]")
	fmt.Println("\t2-Deliver payload random fees [", viper.GetInt("env.dynamicFeeMin"), "--", viper.GetInt("env.dynamicFeeMax"), "]")
	fmt.Println("\t8-Check delivery confirmations (latest run)")
	fmt.Println("\t9-List DB tests")
	fmt.Println("\t0-Exit")
	fmt.Println("")
	fmt.Print("\tSelect option [0-9]:")
	color.Unset()
}
