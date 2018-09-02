package main

import (
	"os"
	"github.com/urfave/cli"
	"fmt"
	"log"
	qcli "github.com/qnib/go-rfxbridge/cli"
	"github.com/qnib/go-rfxbridge/http"

)


func Run(ctx *cli.Context) {
	fmt.Printf("%s Args:%s", ctx.Command.Name, ctx.Args())
}



func main() {
	app := cli.NewApp()
	app.Name = "Deamon to expose RFXCOM devices to others"
	app.Usage = "go-rfxbridge [options]"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen-addr",
			Value: "0.0.0.0:8081",
			Usage: "IP:PORT to bind http server",
			EnvVar: "RFXBRIDGE_HTTP_ADDR",
		},
		cli.StringFlag{
			Name:  "dev-map",
			Usage: "Comma separated mapping of DevName:RFX_ID",
			EnvVar: "RFXBRIDGE_DEV_MAP",
		},
		cli.StringFlag{
			Name:  "usb",
			Value: "/dev/ttyUSB0",
			Usage: "USB device of RFSCOM",
			EnvVar: "RFXBRIDGE_USB",
		},
		cli.BoolFlag{
			Name: "debug",
			Usage: "Be more verbose..",
			EnvVar: "RFXBRIDGE_DEBUG",

		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Serve HTTP endpoint.",
			Action:  http.RunServer,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:   "debug",
					Usage:  "Be more verbose..",
					EnvVar: "RFXBRIDGE_RFX_DEBUG",
				},
			},
		},{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "Query device status",
			Action:  qcli.Request,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
