// Copyright 2020. Akamai Technologies, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

var (
	VERSION = "0.6.0"
)

func main() {
	app := cli.NewApp()
	app.Name = "akamai-dns"
	app.Usage = "CLI DNS"
	app.Version = VERSION

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "edgerc",
			Usage:  "Path to the .edgerc file",
			EnvVar: "AKAMAI_EDGERC",
		},
		cli.StringFlag{
			Name:   "section",
			Usage:  "Section in the .edgerc file",
			EnvVar: "AKAMAI_EDGERC_SECTION",
		},
		cli.StringFlag{
			Name:   "accountkey, account-key",
			Usage:  "Account switch key",
			EnvVar: "AKAMAI_EDGERC_ACCOUNT_KEY",
		},
	}

	app.Commands = GetCommands()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
