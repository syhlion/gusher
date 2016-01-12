package main

import (
	"os"

	"github.com/codegangsta/cli"
)

const (
	APP_VER = "2.0.0"
)

//進入點
func main() {
	gopusher := cli.NewApp()
	gopusher.Name = "gusher"
	gopusher.Version = APP_VER
	gopusher.Commands = []cli.Command{
		CmdStart,
		CmdInitConfig,
	}
	gopusher.Run(os.Args)

}
