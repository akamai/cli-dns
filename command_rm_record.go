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
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v1"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli"
)

func cmdRmRecord(c *cli.Context) error {
	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dns.Config = config

	hostname := c.Args().First()
	if hostname == "" {
		color.Red("hostname is required")
		cli.ShowAppHelpAndExit(c, 1)
		return cli.NewExitError("hostname is required", 1)
	}

	if !c.IsSet("name") {
		color.Red("name is required")
		cli.ShowAppHelpAndExit(c, 1)
		return cli.NewExitError("name is required", 1)
	}

	akamai.StartSpinner(
		fmt.Sprintf("Removing %s record...", c.Command.Name),
		fmt.Sprintf("Removing %s record...... [%s]", c.Command.Name, color.GreenString("OK")),
	)

	zone, err := dns.GetZone(hostname)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError("Unable to fetch zone", 1)
	}

	recordType := c.Command.Name
	options := make(map[string]interface{})

	for option, settings := range recordOptions[recordType] {
		if !c.IsSet(option) {
			continue
		}

		switch option {
		case "active":
			active := true
			inactive := c.Bool("inactive")

			if inactive {
				active = false
			}
			options["active"] = active
		default:
			switch settings.flagType {
			case "string":
				options[option] = c.String(option)
			case "int":
				options[option] = c.Int(option)
			case "uint":
				options[option] = c.Uint(option)
			case "uint16":
				options[option] = uint16(c.Uint(option))
			case "bool":
				options[option] = c.Bool(option)
			}
		}
	}

	if c.IsSet("target") && (recordType == "CNAME" || recordType == "MX" || recordType == "NS" || recordType == "PTR" || recordType == "SRV") {
		if !strings.HasSuffix(options["target"].(string), ".") {
			target := options["target"].(string)
			target += "."
			options["target"] = target
		}
	}

	records := zone.FindRecords(recordType, options)
	if len(records) == 0 {
		akamai.StopSpinnerWarnOk()
		fmt.Fprintln(c.App.Writer, color.CyanString("No records found"))
		return nil
	}

	if len(records) == 1 {
		err := zone.RemoveRecord(records[0])
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError("Unable to remove record", 1)
		}

		err = zone.Save()
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError("Unable to remove record", 1)
		}

		akamai.StopSpinnerOk()
		return nil
	}

	if len(records) > 1 {
		if !akamai.IsInteractive(c) && !c.Bool("force-multiple") {
			akamai.StopSpinnerFail()
			return cli.NewExitError("Multiple records found, no records removed", 1)
		}

		if c.Bool("force-multiple") {
			for _, record := range records {
				err = zone.RemoveRecord(record)
				if err != nil {
					akamai.StopSpinnerFail()
					return cli.NewExitError("Unable to remove records", 1)
				}
			}

			err = zone.Save()
			if err != nil {
				akamai.StopSpinnerFail()
				return cli.NewExitError("Unable to remove records", 1)
			}

			akamai.StopSpinnerOk()
			fmt.Fprintf(akamai.App.ErrWriter, "%d records successfully removed\n", len(records))
			return nil
		}

		akamai.StopSpinnerWarn()
		color.Cyan("Multiple records found:")
		newZone := dns.NewZone(hostname)
		newZone.Token = zone.Token
		for _, record := range records {
			newZone.AddRecord(record)
		}
		renderZoneTable(newZone, c)
		if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
			return cli.NewExitError(color.RedString("Use --force-multiple to remove multiple records"), 1)
		}

		fmt.Fprint(c.App.Writer, "\nWhich records would you like to remove? [comma-separated, leave empty to quit]: ")
		answer := ""
		fmt.Scanln(&answer)
		if strings.TrimSpace(answer) == "" {
			fmt.Fprintln(c.App.Writer, color.CyanString("No records selected"))
			return nil
		}

		akamai.StartSpinner(
			fmt.Sprintf("Removing %s records...", c.Command.Name),
			fmt.Sprintf("Removing %s records...... [%s]", c.Command.Name, color.GreenString("OK")),
		)
		recordKeys := strings.Split(answer, ",")
		for key := range recordKeys {
			recordKeys[key] = strings.TrimSpace(recordKeys[key])
			var i int
			if i, err = strconv.Atoi(recordKeys[key]); err != nil || i > len(records) {
				akamai.StopSpinnerFail()
				return cli.NewExitError(color.RedString("Invalid record ID: %s", recordKeys[key]), 1)
			}
			err = zone.RemoveRecord(records[i-1])
			if err != nil {
				akamai.StopSpinnerFail()
				return cli.NewExitError(color.RedString("Unable to remove record %s: %s", recordKeys[key], err.Error()), 1)
			}
		}

		err = zone.Save()
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError("Unable to remove records", 1)
		}

		akamai.StopSpinnerOk()
		return nil
	}

	return nil
}
