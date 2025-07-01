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

func GetCommands() []cli.Command {
	var commands []cli.Command

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
		Name:        "retrieve-zone",
		Description: "Retrieve a zone's configuration and records",
		ArgsUsage:   "<zonename>",
		Action:      cmdRetrieveZone,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "json",
				Usage: "Output zone in JSON format",
			},
			cli.StringSliceFlag{
				Name:  "filter",
				Usage: "Filter by record type",
			},
		},
	})

	commands = append(commands, cli.Command{
		Name:        "update-zone",
		Description: "Update a zone using either a recordsets JSON file or a DNS master zone file",
		ArgsUsage:   "<zonename>",
		Action:      cmdUpdateZone,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file, f",
				Usage: "Path to input file (JSON for recordsets or DNS master file)",
			},
			cli.BoolFlag{
				Name:  "dns",
				Usage: "Use this flag if input file is a DNS master zone file",
			},
			cli.StringFlag{
				Name:  "output, o",
				Usage: "Optional output path for updated zone details",
			},
			cli.BoolFlag{
				Name:  "overwrite",
				Usage: "Overwrite all recordsets instead of merging with existing",
			},
			cli.BoolFlag{
				Name:  "json",
				Usage: "Output zone response in JSON format",
			},
			cli.BoolFlag{
				Name:  "suppress",
				Usage: "Suppress output to console",
			},
		},
	})

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

	commands = append(commands, cli.Command{
		Name:        "add-record",
		Description: "Create or update a DNS recordset in a zone",
		ArgsUsage:   "<type> <zonename>",
		Action:      cmdAddRecord,
		Flags: append(baseSetCmdFlags,
			cli.StringSliceFlag{
				Name:  "target",
				Usage: "Record target (RDATA), multiple flags allowed",
			},
			cli.IntFlag{
				Name:  "ttl",
				Usage: "Recordset TTL in seconds",
			},
		),
	})

	commands = append(commands, cli.Command{
		Name:        "rm-record",
		Description: "Remove a DNS recordset from a zone",
		ArgsUsage:   "<record type> <zonename>",
		Action:      cmdRmRecord,
		Flags: append(baseV11CmdFlags,
			cli.StringFlag{
				Name:  "name",
				Usage: "Record name to delete (eg: --name www)",
			},
			cli.BoolFlag{
				Name:  "force-multiple",
				Usage: "Force delete all matching records without confirmation",
			},
			cli.BoolFlag{
				Name:  "non-interactive",
				Usage: "Run in non-interactive mode (e.g. CI). Fails if multiple matches and not forced.",
			},
		),
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
