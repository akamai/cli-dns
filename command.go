// Copyright 2018-2020. Akamai Technologies, Inc
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

// V1 Records
var (
	baseCmdFlags = []cli.Flag{}

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
			"order":       {true, "uint16"},
			"preference":  {true, "uint16"},
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
			"active":                 {false, "bool"},
			"algorithm":              {true, "int"},
			"flags":                  {true, "int"},
			"iterations":             {true, "int"},
			"name":                   {true, "string"},
			"next-hashed-owner-name": {true, "string"},
			"salt":                   {true, "string"},
			"ttl":                    {false, "int"},
			"type-bitmaps":           {true, "string"},
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
			"port":     {true, "uint16"},
			"priority": {true, "int"},
			"target":   {true, "string"},
			"ttl":      {false, "int"},
			"weight":   {true, "uint16"},
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
			Flags: append(baseCmdFlags, []cli.Flag{
				cli.BoolFlag{
					Name:  "non-interactive",
					Usage: "Run in non-interactive mode",
				},
				cli.BoolFlag{
					Name:  "force-multiple",
					Usage: "Force removal of multiple matched records",
				},
			}...),
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
				case "uint16":
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

	commands = append(commands, addRecordCmd)
	commands = append(commands, rmRecordCmd)

	// V2 Base flags

	var baseV2BaseFlags = []cli.Flag{
		cli.BoolFlag{
			Name:   "json",
			Usage:  "Output as JSON",
			EnvVar: "AKAMAI_CLI_DNS_" + "JSON",
		},
		cli.StringFlag{
			Name:  "output",
			Usage: "Output command results to `FILE`",
		},
	}

	baseV2CmdFlags := append(baseV2BaseFlags, []cli.Flag{
		cli.BoolFlag{
			Name:   "suppress",
			Usage:  "Suppress command result output. Overrides other output related flags",
			EnvVar: "AKAMAI_CLI_DNS_" + "SUPPRESS",
		},
	}...)

	baseSetCmdFlags := append(baseV2CmdFlags, []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Recordset `NAME`",
		},
		cli.StringFlag{
			Name:  "type",
			Usage: "Recordset `TYPE`",
		},
	}...)

	// V2 Recordsets
	commands = append(commands, cli.Command{
		Name:        "list-recordsets",
		Description: "Retreive list of zone Recordsets",
		ArgsUsage:   "<zonename>",
		Action:      cmdListRecordsets,
		Flags: append(baseV2BaseFlags, []cli.Flag{
			cli.StringSliceFlag{
				Name:  "type",
				Usage: "List recordset(s) matching `TYPE`. Multiple flags allowed",
			},
			cli.StringFlag{
				Name:  "sortby",
				Usage: "List returned recordsets sorted by `SORTBY`",
			},
			cli.StringFlag{
				Name:  "search",
				Usage: "Filter returned recordsets by `SEARCH` criteria",
			},
		}...),
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "create-recordsets",
		Description: "Create multiple zone Recordsets from `FILE`",
		ArgsUsage:   "<zonename>",
		Action:      cmdCreateRecordsets,
		Flags: append(baseV2CmdFlags, []cli.Flag{
			cli.StringFlag{
				Name:  "file",
				Usage: "`FILE` path to JSON formatted recordset content",
			},
		}...),
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "update-recordsets",
		Description: "Update multiple zone Recordsets from `FILE`",
		ArgsUsage:   "<zonename>",
		Action:      cmdUpdateRecordsets,
		Flags: append(baseV2CmdFlags, []cli.Flag{
			cli.BoolFlag{
				Name:  "overwrite",
				Usage: "Replace ALL Recordsets",
			},
			cli.StringFlag{
				Name:  "file",
				Usage: "`FILE` path to JSON formatted recordset content",
			},
		}...),
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "retrieve-recordset",
		Description: "Retrieve recordset",
		ArgsUsage:   "<zonename>",
		Action:      cmdRetrieveRecordset,
		Flags: append(baseV2BaseFlags, []cli.Flag{
			cli.StringFlag{
				Name:  "name",
				Usage: "Recordset `NAME`",
			},
			cli.StringFlag{
				Name:  "type",
				Usage: "Recordset `TYPE`",
			},
		}...),
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "create-recordset",
		Description: "Create a new recordset",
		ArgsUsage:   "<zonename>",
		Action:      cmdCreateRecordset,
		Flags: append(baseSetCmdFlags, []cli.Flag{
			cli.IntFlag{
				Name:  "ttl",
				Usage: "Recordset `TTL`",
			},
			cli.StringSliceFlag{
				Name:  "rdata",
				Usage: "Recordset `RDATA`. Multiple flags allowed.",
			},
			cli.StringFlag{
				Name:  "file",
				Usage: "`FILE` path to JSON formatted recordset content",
			},
		}...),
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "update-recordset",
		Description: "Update existing recordset",
		ArgsUsage:   "<zonename>",
		Action:      cmdUpdateRecordset,
		Flags: append(baseSetCmdFlags, []cli.Flag{
			cli.IntFlag{
				Name:  "ttl",
				Usage: "Recordset `TTL`",
			},
			cli.StringSliceFlag{
				Name:  "rdata",
				Usage: "Record `RDATA`. Multiple flags allowed.",
			},
			cli.StringFlag{
				Name:  "file",
				Usage: "`FILE` path to JSON formatted recordset content. Allows multiple recordsets.",
			},
		}...),
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "delete-recordset",
		Description: "Delete recordset",
		ArgsUsage:   "<zonename>",
		Action:      cmdDeleteRecordset,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name",
				Usage: "Recordset `NAME`",
			},
			cli.StringFlag{
				Name:  "type",
				Usage: "Recordset `TYPE`",
			},
		},
		BashComplete: akamai.DefaultAutoComplete,
	})

	// V1 Zone
	commands = append(commands, cli.Command{
		Name:        "retrieve-zone",
		Description: "Fetch and display a zone",
		ArgsUsage:   "<hostname>",
		Action:      cmdRetrieveZone,
		Flags: append(baseCmdFlags, []cli.Flag{
			cli.BoolFlag{
				Name:   "json",
				Usage:  "Output as JSON",
				EnvVar: "AKAMAI_CLI_DNS_" + "JSON",
			},
			cli.BoolFlag{
				Name:  "output",
				Usage: "Output to `FILE`",
			},
			cli.StringSliceFlag{
				Name:  "filter",
				Usage: "Only show record types matching `TYPE`",
			},
		}...),
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "update-zone",
		Description: "Update a zone",
		ArgsUsage:   "<hostname>",
		Action:      cmdUpdateZone,
		Flags: append(baseCmdFlags, []cli.Flag{
			cli.BoolTFlag{
				Name:   "json",
				Usage:  "Input is in JSON format",
				EnvVar: "AKAMAI_CLI_DNS_" + "JSON",
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
		}...),
		BashComplete: akamai.DefaultAutoComplete,
	})

	// V2 Zones
	baseZoneCmdFlags := append(baseV2CmdFlags, []cli.Flag{
		cli.StringFlag{
			Name:  "type",
			Usage: "Zone `TYPE`",
		},
		cli.StringSliceFlag{
			Name:  "master",
			Usage: "Secondary Zone `MASTER`. Multiple flags may be specified",
		},
		cli.StringFlag{
			Name:  "comment",
			Usage: "Zone `COMMENT`",
		},
		cli.BoolFlag{
			Name:  "signandserve",
			Usage: "Primary or Secondary Zone `SIGNANDSERVE` flag",
		},
		cli.StringFlag{
			Name:  "algorithm",
			Usage: "Zone signandserve `ALGORITHM`",
		},
		cli.StringFlag{
			Name:  "tsigname",
			Usage: "TSIG key `NAME`",
		},
		cli.StringFlag{
			Name:  "tsigalgorithm",
			Usage: "TSIG key `ALGORITHM`",
		},
		cli.StringFlag{
			Name:  "tsigsecret",
			Usage: "TSIG key `SECRET`",
		},
		cli.StringFlag{
			Name:  "target",
			Usage: "Alias Zone `TARGET`",
		},
		cli.StringFlag{
			Name:  "endcustomerid",
			Usage: "`ENDCUSTOMERID`",
		},
		cli.StringFlag{
			Name:  "file",
			Usage: "Read JSON formatted input from `FILE`",
		},
	}...)

	commands = append(commands, cli.Command{
		Name:        "list-zoneconfig",
		Description: "List zone configuration(s)",
		Action:      cmdListZoneconfig,
		Flags: append(baseV2BaseFlags, []cli.Flag{
			cli.StringSliceFlag{
				Name:  "contractid",
				Usage: "Contract `ID`. Multiple flags allowed",
			},
			cli.StringSliceFlag{
				Name:  "type",
				Usage: "Zone `TYPE`. Multiple flags allowed",
			},
			cli.StringFlag{
				Name:  "search",
				Usage: "Zone search `VALUE`",
			},
			cli.BoolFlag{
				Name:  "summary",
				Usage: "List zone names and type",
			},
		}...),
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "create-zoneconfig",
		ArgsUsage:   "<zonename>",
		Description: "Create zone from configuration",
		Action:      cmdCreateZoneconfig,
		Flags: append(baseZoneCmdFlags, []cli.Flag{
			cli.StringFlag{
				Name:  "contractid",
				Usage: "Contract `ID`",
			},
			cli.StringFlag{
				Name:  "groupid",
				Usage: "Group `ID`",
			},
			cli.BoolFlag{
				Name:  "initialize",
				Usage: "Generate default SOA and NS Records",
			},
		}...),
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "retrieve-zoneconfig",
		Description: "Fetch and display zone configuration",
		ArgsUsage:   "<zonename>",
		Action:      cmdRetrieveZoneconfig,
		Flags: append(baseV2BaseFlags, []cli.Flag{
			cli.BoolFlag{
				Name:  "dns",
				Usage: "Retrieve Zone Master File",
			},
		}...),

		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "update-zoneconfig",
		Description: "Update a zone",
		ArgsUsage:   "<zonename>",
		Action:      cmdUpdateZoneconfig,
		Flags: append(baseZoneCmdFlags, []cli.Flag{
			cli.BoolFlag{
				Name:  "dns",
				Usage: "Input is Zone Master File",
			},
		}...),
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "submit-bulkzones",
		Description: "Submit Bulk Zones request",
		Action:      cmdSubmitBulkZones,
		Flags: append(baseV2CmdFlags, []cli.Flag{
			cli.StringFlag{
				Name:  "contractid",
				Usage: "Contract `ID`. Required for create.",
			},
			cli.StringFlag{
				Name:  "groupid",
				Usage: "Group `ID`. Optional for create.",
			},
			cli.BoolFlag{
				Name:  "bypasszonesafety",
				Usage: "Bypass zone safety check. Optional for delete.",
			},
			cli.BoolFlag{
				Name:  "create",
				Usage: "Bulk zone create operation.",
			},
			cli.BoolFlag{
				Name:  "delete",
				Usage: "Bulk zone delete operation.",
			},
			cli.StringFlag{
				Name:  "file",
				Usage: "Read JSON formatted input from `FILE`",
			},
		}...),
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "status-bulkzones",
		Description: "Query Bulk Zones Request Status",
		Action:      cmdStatusBulkZones,
		Flags: append(baseV2BaseFlags, []cli.Flag{
			cli.StringFlag{
				Name:  "requestid",
				Usage: "Request Id",
			},
			cli.BoolFlag{
				Name:  "create",
				Usage: "Bulk zone create operation.",
			},
			cli.BoolFlag{
				Name:  "delete",
				Usage: "Bulk zone delete operation.",
			},
		}...),
		BashComplete: akamai.DefaultAutoComplete,
	})

	commands = append(commands, cli.Command{
		Name:        "result-bulkzones",
		Description: "Query Bulk Zones Result Summary",
		Action:      cmdResultBulkZones,
		Flags: append(baseV2BaseFlags, []cli.Flag{
			cli.StringFlag{
				Name:  "requestid",
				Usage: "Request Id",
			},
			cli.BoolFlag{
				Name:  "create",
				Usage: "Bulk zone create operation.",
			},
			cli.BoolFlag{
				Name:  "delete",
				Usage: "Bulk zone delete operation.",
			},
		}...),
		BashComplete: akamai.DefaultAutoComplete,
	})

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
