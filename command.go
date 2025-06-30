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
	"github.com/urfave/cli"
)

// V1 Records
/*var (
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
)*/

func GetCommands() []cli.Command {
	var commands []cli.Command

	/*recordMap := []string{
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
	}*/

	// V11 Recordsets
	baseV11BaseFlags := []cli.Flag{
		cli.BoolFlag{
			Name:   "json",
			Usage:  "Output as JSON",
			EnvVar: "AKAMAI_CLI_DNS_" + "JSON",
		},
		cli.StringFlag{
			Name:  "output",
			Usage: "Output command results to FILE",
		},
	}

	baseV11CmdFlags := append(baseV11BaseFlags,
		cli.BoolFlag{
			Name:   "suppress",
			Usage:  "Suppress command result output",
			EnvVar: "AKAMAI_CLI_DNS_" + "SUPPRESS",
		})

	baseSetCmdFlags := append(baseV11CmdFlags,
		cli.StringFlag{
			Name:  "name",
			Usage: "Recordset NAME",
		},
		cli.StringFlag{
			Name:  "type",
			Usage: "Recordset TYPE",
		},
	)

	commands = append(commands, cli.Command{
		Name:        "list-recordsets",
		Description: "Retreive list of zone Recordsets",
		ArgsUsage:   "<zonename>",
		Action:      cmdListRecordsets,
		Flags: append(baseV11BaseFlags,
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
		),
	})

	commands = append(commands, cli.Command{
		Name:        "create-recordsets",
		Description: "Create multiple zone Recordsets from `FILE`",
		ArgsUsage:   "<zonename>",
		Action:      cmdCreateRecordsets,
		Flags: append(baseV11CmdFlags,
			cli.StringFlag{
				Name:  "file",
				Usage: "`FILE` path to JSON formatted recordset content",
			},
		),
	})

	commands = append(commands, cli.Command{
		Name:        "update-recordsets",
		Description: "Update multiple zone Recordsets from `FILE`",
		ArgsUsage:   "<zonename>",
		Action:      cmdUpdateRecordsets,
		Flags: append(baseV11CmdFlags,
			cli.BoolFlag{
				Name:  "overwrite",
				Usage: "Replace ALL Recordsets",
			},
			cli.StringFlag{
				Name:  "file",
				Usage: "`FILE` path to JSON formatted recordset content",
			},
		),
	})

	commands = append(commands, cli.Command{
		Name:        "retrieve-recordset",
		Description: "Retrieve recordset",
		ArgsUsage:   "<zonename>",
		Action:      cmdRetrieveRecordset,
		Flags: append(baseV11BaseFlags,
			cli.StringFlag{
				Name:  "name",
				Usage: "Recordset `NAME`",
			},
			cli.StringFlag{
				Name:  "type",
				Usage: "Recordset `TYPE`",
			},
		),
	})

	commands = append(commands, cli.Command{
		Name:        "create-recordset",
		Description: "Create a new recordset",
		ArgsUsage:   "<zonename>",
		Action:      cmdCreateRecordset,
		Flags: append(baseSetCmdFlags,
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
		),
	})

	commands = append(commands, cli.Command{
		Name:        "update-recordset",
		Description: "Update existing recordset",
		ArgsUsage:   "<zonename>",
		Action:      cmdUpdateRecordset,
		Flags: append(baseSetCmdFlags,
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
		),
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
	})

	// V11 Zones
	//Zone level flags
	baseZoneCmdFlags := append(baseV11CmdFlags,
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
	)

	commands = append(commands, cli.Command{
		Name:        "list-zoneconfig",
		Description: "List zone configuration(s)",
		Action:      cmdListZoneconfig,
		Flags: append(baseV11BaseFlags,
			cli.StringFlag{
				Name:  "contractid",
				Usage: "Contract `ID`",
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
		),
	})

	commands = append(commands, cli.Command{
		Name:        "retrieve-zoneconfig",
		Description: "Fetch and display zone configuration",
		ArgsUsage:   "<zonename>",
		Action:      cmdRetrieveZoneconfig,
		Flags: append(baseV11BaseFlags,
			cli.BoolFlag{
				Name:  "dns",
				Usage: "Retrieve Zone Master File",
			},
		),
	})

	commands = append(commands, cli.Command{
		Name:        "create-zoneconfig",
		ArgsUsage:   "<zonename>",
		Description: "Create zone from configuration",
		Action:      cmdCreateZoneconfig,
		Flags: append(baseZoneCmdFlags,
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
		),
	})

	commands = append(commands, cli.Command{
		Name:        "update-zoneconfig",
		Description: "Update a zone",
		ArgsUsage:   "<zonename>",
		Action:      cmdUpdateZoneconfig,
		Flags: append(baseZoneCmdFlags,
			cli.BoolFlag{
				Name:  "dns",
				Usage: "Input is Zone Master File",
			},
		),
	})

	commands = append(commands, cli.Command{
		Name:        "submit-bulkzones",
		Description: "Submit Bulk Zones request",
		Action:      cmdSubmitBulkZones,
		Flags: append(baseV11CmdFlags,
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
		),
	})

	commands = append(commands, cli.Command{
		Name:        "status-bulkzones",
		Description: "Query Bulk Zones Request Status",
		Action:      cmdStatusBulkZones,
		Flags: append(baseV11BaseFlags,
			cli.StringSliceFlag{
				Name:  "requestid",
				Usage: "Request Id. Multiple args allowed.",
			},
			cli.BoolFlag{
				Name:  "create",
				Usage: "Bulk zone create operation.",
			},
			cli.BoolFlag{
				Name:  "delete",
				Usage: "Bulk zone delete operation.",
			},
		),
	})

	commands = append(commands, cli.Command{
		Name:        "result-bulkzones",
		Description: "Query Bulk Zones Result Summary",
		Action:      cmdResultBulkZones,
		Flags: append(baseV11BaseFlags,
			cli.StringSliceFlag{
				Name:  "requestid",
				Usage: "Request Id. Multiple args allowed.",
			},
			cli.BoolFlag{
				Name:  "create",
				Usage: "Bulk zone create operation.",
			},
			cli.BoolFlag{
				Name:  "delete",
				Usage: "Bulk zone delete operation.",
			},
		),
	})

	return commands

}
