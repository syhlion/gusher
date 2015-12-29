package main

import (
	"os"

	"github.com/codegangsta/cli"
)

const (
	APP_VER = "0.8.3"
)

//進入點
func main() {

	gopusher := cli.NewApp()
	gopusher.Name = "gusher"
	gopusher.Version = APP_VER
	gopusher.Commands = []cli.Command{
		CmdStart,
		InitStart,
	}

	gopusher.Run(os.Args)

}
