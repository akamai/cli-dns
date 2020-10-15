// Copyright 2018. Akamai Technologies, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"os"

	akamai "github.com/akamai/cli-common-golang"
)

var (
	VERSION = "0.4.0"
)

func main() {
	akamai.CreateApp(
		"dns",
		"A CLI for Edge DNS",
		"Manage DNS Zones with Edge DNS",
		VERSION,
		"dns",
		commandLocator,
	)

	setHelpTemplates()
	akamai.App.Run(os.Args)
}
