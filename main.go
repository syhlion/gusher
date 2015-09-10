package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/syhlion/gusher/cmd"
)

const (
	APP_VER = "0.6.0"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

//進入點
func main() {

	gopusher := cli.NewApp()
	gopusher.Name = "gusher"
	gopusher.Version = APP_VER
	gopusher.Commands = []cli.Command{
		cmd.CmdStart,
		cmd.InitStart,
	}

	gopusher.Run(os.Args)

}
