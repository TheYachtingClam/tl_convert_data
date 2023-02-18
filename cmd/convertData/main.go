//
//  TERALYTIC CONFIDENTIAL
//  _________________
//   2021- 2022 TERALYTIC
//   All Rights Reserved.
//
//   NOTICE:  All information contained herein is, and remains
//   the property of TERALYTIC and its suppliers,
//   if any.  The intellectual and technical concepts contained
//   herein are proprietary to TERALYTIC
//   and its suppliers and may be covered by U.S. and Foreign Patents,
//   patents in process, and are protected by trade secret or copyright law.
//   Dissemination of this information or reproduction of this material
//   is strictly forbidden unless prior written permission is obtained
//   from TERALYTIC.
//

package main

import (
	"os"

	_ "net/http/pprof"

	"github.com/apex/log"
	clih "github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/discard"

	jsonh "github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/text"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Name = "convertData"
	app.Usage = "TERALYTIC convert data"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "set the logging level",
			Value:   "info",
			EnvVars: []string{"LOG_LEVEL"},
		},
		&cli.StringFlag{
			Name:    "log-format",
			Usage:   "set the logging format",
			Value:   "text",
			EnvVars: []string{"LOG_FORMAT"},
		},
	}

	app.Commands = []*cli.Command{
		convertCmd,
	}
	app.Before = func(c *cli.Context) error {
		if logLevel := c.String("log-level"); logLevel != "" {
			if level, err := log.ParseLevel(logLevel); err == nil {
				log.SetLevel(level)
			}
		}

		if c.Bool("silent") {
			log.SetHandler(discard.Default)
		} else {
			switch c.String("log-format") {
			case "json":
				log.SetHandler(jsonh.Default)

			case "logfmt":
				log.SetHandler(logfmt.Default)

			case "text":
				log.SetHandler(text.Default)

			default:
				log.SetHandler(clih.Default)
			}
		}

		return nil
	}
	

	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
}
