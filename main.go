package main

import (
	"github.com/codegangsta/cli"
	"github.com/syhlion/gopusher/cmd"
	"os"
	"runtime"
)

const (
	APP_VER = "0.4.2"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

//進入點
func main() {

	gopusher := cli.NewApp()
	gopusher.Name = "gopusher"
	gopusher.Version = APP_VER
	gopusher.Commands = []cli.Command{
		cmd.CmdStart,
	}

	gopusher.Run(os.Args)

}
