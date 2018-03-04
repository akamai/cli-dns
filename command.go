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
	"strings"

	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var (
	recordOptions = map[string]map[string]struct {
		required bool
		flagType string
	}{
		"A": {
			"active": {false, "bool"},
			"name":   {true, "string"},
			"target": {true, "string"},
			"ttl":    {false, "int"},
		},
		"AAAA": {
			"active": {false, "bool"},
			"name":   {true, "string"},
			"target": {true, "string"},
			"ttl":    {false, "int"},
		},
		"AFSDB": {
			"active":  {false, "bool"},
			"name":    {true, "string"},
			"subtype": {true, "int"},
			"target":  {true, "string"},
			"ttl":     {false, "int"},
		},
		"CNAME": {
			"active": {false, "bool"},
			"name":   {true, "string"},
			"target": {true, "string"},
			"ttl":    {false, "int"},
		},
		"DNSKEY": {
			"active":    {false, "bool"},
			"algorithm": {true, "int"},
			"flags":     {true, "int"},
			"key":       {true, "string"},
			"name":      {true, "string"},
			"protocol":  {true, "int"},
			"ttl":       {false, "int"},
		},
		"DS": {
			"active":      {false, "bool"},
			"algorithm":   {true, "int"},
			"digest":      {true, "string"},
			"digest-type": {true, "int"},
			"keytag":      {true, "int"},
			"name":        {true, "string"},
			"ttl":         {false, "int"},
		},
		"HINFO": {
			"active":   {false, "bool"},
			"hardware": {true, "string"},
			"name":     {true, "string"},
			"software": {true, "string"},
			"ttl":      {false, "int"},
		},
		"LOC": {
			"active": {false, "bool"},
			"name":   {true, "string"},
			"target": {true, "string"},
			"ttl":    {false, "int"},
		},
		"MX": {
			"active":   {false, "bool"},
			"name":     {true, "string"},
			"priority": {true, "int"},
			"target":   {true, "string"},
			"ttl":      {false, "int"},
		},
		"NAPTR": {
			"active":      {false, "bool"},
			"flags":       {true, "string"},
			"name":        {true, "string"},
			"order":       {true, "uint"},
			"preference":  {true, "uint"},
			"regexp":      {true, "string"},
			"replacement": {true, "string"},
			"service":     {true, "string"},
			"ttl":         {false, "int"},
		},
		"NS": {
			"active": {false, "bool"},
			"name":   {true, "string"},
			"target": {true, "string"},
			"ttl":    {false, "int"},
		},
		"NSEC3": {
			"active":     {false, "bool"},
			"algorithm":  {true, "int"},
			"flags":      {true, "int"},
			"iterations": {true, "int"},
			"name":       {true, "string"},
			"next-hashed-owner-name": {true, "string"},
			"salt":         {true, "string"},
			"ttl":          {false, "int"},
			"type-bitmaps": {true, "string"},
		},
		"NSEC3PARAM": {
			"active":     {false, "bool"},
			"algorithm":  {true, "int"},
			"flags":      {true, "int"},
			"iterations": {true, "int"},
			"name":       {true, "string"},
			"salt":       {true, "string"},
			"ttl":        {false, "int"},
		},
		"PTR": {
			"active": {false, "bool"},
			"name":   {true, "string"},
			"target": {true, "string"},
			"ttl":    {false, "int"},
		},
		"RP": {
			"active":  {false, "bool"},
			"mailbox": {true, "string"},
			"name":    {true, "string"},
			"ttl":     {false, "int"},
			"txt":     {true, "string"},
		},
		"RRSIG": {
			"active":       {false, "bool"},
			"algorithm":    {true, "int"},
			"expiration":   {true, "string"},
			"inception":    {true, "string"},
			"keytag":       {true, "int"},
			"labels":       {true, "int"},
			"name":         {true, "string"},
			"original-ttl": {true, "int"},
			"signature":    {true, "string"},
			"signer":       {true, "string"},
			"ttl":          {false, "int"},
			"type-covered": {true, "string"},
		},
		"SOA": {
			"contact":      {true, "string"},
			"expire":       {true, "int"},
			"minimum":      {true, "uint"},
			"originserver": {true, "string"},
			"refresh":      {true, "int"},
			"retry":        {true, "int"},
			"serial":       {false, "uint"},
			"ttl":          {false, "int"},
		},
		"SPF": {
			"active": {false, "bool"},
			"name":   {true, "string"},
			"target": {true, "string"},
			"ttl":    {false, "int"},
		},
		"SRV": {
			"active":   {false, "bool"},
			"name":     {true, "string"},
			"port":     {true, "uint"},
			"priority": {true, "int"},
			"target":   {true, "string"},
			"ttl":      {false, "int"},
			"weight":   {true, "uint"},
		},
		"SSHFP": {
			"active":           {false, "bool"},
			"algorithm":        {true, "int"},
			"fingerprint":      {true, "string"},
			"fingerprint-type": {true, "int"},
			"name":             {true, "string"},
			"ttl":              {false, "int"},
		},
		"TXT": {
			"active": {false, "bool"},
			"name":   {true, "string"},
			"target": {true, "string"},
			"ttl":    {false, "int"},
		},
	}
)

var commandLocator akamai.CommandLocator = func() ([]cli.Command, error) {
	var commands []cli.Command

	recordMap := []string{
		"A",
		"AAAA",
		"AFSDB",
		"CNAME",
		"DNSKEY",
		"DS",
		"HINFO",
		"LOC",
		"MX",
		"NAPTR",
		"NS",
		"NSEC3",
		"NSEC3PARAM",
		"PTR",
		"RP",
		"RRSIG",
		"SOA",
		"SPF",
		"SRV",
		"SSHFP",
		"TXT",
	}

	addRecordCmd := cli.Command{
		Name:        "add-record",
		ArgsUsage:   "<record type> <hostname>",
		Description: "Add a new record to the zone",
		HideHelp:    true,
		Action: func(c *cli.Context) error {
			if recordType := c.Args().First(); recordType == "" {
				cli.ShowAppHelp(c)
				return cli.NewExitError(color.RedString("You must specify a record type"), 1)
			}
			os.Args[2] = strings.ToUpper(c.Args().First())
			return akamai.App.Run(os.Args)
		},
	}
	rmRecordCmd := cli.Command{
		Name:        "rm-record",
		ArgsUsage:   "<record type> <hostname>",
		Description: "Remove a record from the zone",
		HideHelp:    true,
		Action: func(c *cli.Context) error {
			if recordType := c.Args().First(); recordType == "" {
				cli.ShowAppHelp(c)
				return cli.NewExitError(color.RedString("You must specify a record type"), 1)
			}
			os.Args[2] = strings.ToUpper(c.Args().First())
			return akamai.App.Run(os.Args)
		},
	}

	for _, recordType := range recordMap {
		addCmd := cli.Command{
			Name:         recordType,
			ArgsUsage:    "<hostname>",
			Description:  fmt.Sprintf("Add a new %s record to the zone", recordType),
			Action:       cmdCreateRecord,
			HideHelp:     true,
			BashComplete: akamai.DefaultAutoComplete,
		}

		rmCmd := cli.Command{
			Name:         recordType,
			ArgsUsage:    "<hostname>",
			Description:  fmt.Sprintf("Remove a %s record from the zone", recordType),
			Action:       cmdRmRecord,
			HideHelp:     true,
			BashComplete: akamai.DefaultAutoComplete,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "non-interactive",
					Usage: "Run in non-interactive mode",
				},
				cli.BoolFlag{
					Name:  "force-multiple",
					Usage: "Force removal of multiple matched records",
				},
			},
		}

		for option, settings := range recordOptions[recordType] {
			var flag cli.Flag
			switch option {
			case "ttl":
				flag = cli.IntFlag{
					Name:  option,
					Value: 7200,
				}
			case "active":
				flag = cli.BoolTFlag{
					Name: "active",
				}
				addCmd.Flags = append(addCmd.Flags, flag)
				rmCmd.Flags = append(rmCmd.Flags, flag)

				flag = cli.BoolFlag{
					Name: "inactive",
				}
			case "target":
				flag = cli.StringFlag{
					Name:  option,
					Usage: "Target `ADDRESS`",
					//EnvVar: "AKAMAI_DNS_" + strings.ToUpper(option),
				}
			default:
				switch settings.flagType {
				case "string":
					flag = cli.StringFlag{
						Name: option,
						//EnvVar: "AKAMAI_DNS_" + strings.ToUpper(option),
					}
				case "int":
					flag = cli.IntFlag{
						Name: option,
						//EnvVar: "AKAMAI_DNS_" + strings.ToUpper(option),
					}
				case "uint":
					flag = cli.UintFlag{
						Name: option,
						//EnvVar: "AKAMAI_DNS_" + strings.ToUpper(option),
					}
				case "bool":
					flag = cli.BoolFlag{
						Name: option,
						//EnvVar: "AKAMAI_DNS_" + strings.ToUpper(option),
					}
				}
			}

			addCmd.Flags = append(addCmd.Flags, flag)
			rmCmd.Flags = append(rmCmd.Flags, flag)
		}

		addRecordCmd.Subcommands = append(addRecordCmd.Subcommands, addCmd)
		rmRecordCmd.Subcommands = append(rmRecordCmd.Subcommands, rmCmd)
	}

	commands = append(commands, cli.Command{
		Name:        "retrieve-zone",
		Description: "Fetch and display a zone",
		ArgsUsage:   "<hostname>",
		Action:      cmdRetrieveZone,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "json",
				Usage: "Output as JSON",
			},
			cli.BoolFlag{
				Name:  "output",
				Usage: "Output to `FILE`",
			},
			cli.StringSliceFlag{
				Name:  "filter",
				Usage: "Only show record types matching `TYPE`",
			},
		},
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "update-zone",
		Description: "Update a zone",
		ArgsUsage:   "<hostname>",
		Action:      cmdUpdateZone,
		Flags: []cli.Flag{
			cli.BoolTFlag{
				Name:  "json",
				Usage: "Input is in JSON format",
			},
			cli.BoolFlag{
				Name:  "dns",
				Usage: "Input is in DNS Zone format",
			},
			cli.BoolFlag{
				Name:  "overwrite",
				Usage: "Overwrite all existing records (default: false)",
			},
			cli.StringFlag{
				Name:  "file",
				Usage: "Read input from `FILE`",
			},
		},
		BashComplete: akamai.DefaultAutoComplete,
	})
	commands = append(commands, addRecordCmd)
	commands = append(commands, rmRecordCmd)
	commands = append(commands,
		cli.Command{
			Name:        "list",
			Description: "List commands",
			Action:      akamai.CmdList,
		},
		cli.Command{
			Name:         "help",
			Description:  "Displays help information",
			ArgsUsage:    "[command] [sub-command]",
			Action:       akamai.CmdHelp,
			BashComplete: akamai.DefaultAutoComplete,
		},
	)

	return commands, nil
}
