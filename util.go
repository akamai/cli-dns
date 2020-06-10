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
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func setHelpTemplates() {
	cli.AppHelpTemplate =
		color.YellowString("Usage: \n") +
			`{{if or (or (eq .HelpName "akamai-dns add-record") (eq .HelpName "akamai dns add-record")) (or (eq .HelpName "akamai-dns rm-record") (eq .HelpName "akamai dns rm-record")) (or (eq .HelpName "akamai-dns retrieve-zone") (eq .HelpName "akamai dns retrieve-zone")) (or (eq .HelpName "akamai-dns update-zone") (eq .HelpName "akamai dns update-zone")) (or (eq .HelpName "akamai-dns list-recordsets") (eq .HelpName "akamai dns list-recordsets")) (or (eq .HelpName "akamai-dns retrieve-recordset") (eq .HelpName "akamai dns retrieve-recordset")) (or (eq .HelpName "akamai-dns create-recordset") (eq .HelpName "akamai dns create-recordset")) (or (eq .HelpName "akamai-dns update-recordset") (eq .HelpName "akamai dns update-recordset")) (or (eq .HelpName "akamai-dns delete-recordset") (eq .HelpName "akamai dns delete-recordset")) (or (eq .HelpName "akamai-dns list-zoneconfig") (eq .HelpName "akamai dns list-zoneconfig")) (or (eq .HelpName "akamai-dns create-zoneconfig") (eq .HelpName "akamai dns create-zoneconfig")) (or (eq .HelpName "akamai-dns retrieve-zoneconfig") (eq .HelpName "akamai dns retrieve-zoneconfig")) (or (eq .HelpName "akamai-dns update-zoneconfig") (eq .HelpName "akamai dns update-zoneconfig"))}}` +
			color.BlueString(`	{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .ArgsUsage}} {{.ArgsUsage}}{{end}}{{end}}`) +
			`{{else}}` +
			color.BlueString(`	{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}}{{range .VisibleFlags}} [--{{.Name}}]{{end}}{{end}}{{if .ArgsUsage}} {{.ArgsUsage}}{{end}}{{if .Commands}} <command> [sub-command]{{end}}{{end}}`) +
			`{{end}}` +

			"{{if .Description}}\n\n" +
			color.YellowString("Description:\n") +
			"   {{.Description}}" +
			"\n\n{{end}}" +

			"{{if .VisibleFlags}}" +
			color.YellowString("Global Flags:\n") +
			"{{range $index, $option := .VisibleFlags}}" +
			"{{if $index}}\n{{end}}" +
			"   {{$option}}" +
			"{{end}}" +
			"\n\n{{end}}" +

			"{{if .VisibleCommands}}" +
			`{{if or (or (eq .HelpName "akamai-dns add-record") (eq .HelpName "akamai dns add-record")) (or (eq .HelpName "akamai-dns rm-record") (eq .HelpName "akamai dns rm-record")) (or (eq .HelpName "akamai-dns retrieve-zone") (eq .HelpName "akamai dns retrieve-zone")) (or (eq .HelpName "akamai-dns update-zone") (eq .HelpName "akamai dns update-zone")) (or (eq .HelpName "akamai-dns list-recordsets") (eq .HelpName "akamai dns list-recordsets")) (or (eq .HelpName "akamai-dns retrieve-recordset") (eq .HelpName "akamai dns retrieve-recordset")) (or (eq .HelpName "akamai-dns create-recordset") (eq .HelpName "akamai dns create-recordset")) (or (eq .HelpName "akamai-dns update-recordset") (eq .HelpName "akamai dns update-recordset")) (or (eq .HelpName "akamai-dns delete-recordset") (eq .HelpName "akamai dns delete-recordset")) (or (eq .HelpName "akamai-dns list-zoneconfig") (eq .HelpName "akamai dns list-zoneconfig")) (or (eq .HelpName "akamai-dns create-zoneconfig") (eq .HelpName "akamai dns create-zoneconfig")) (or (eq .HelpName "akamai-dns retrieve-zoneconfig") (eq .HelpName "akamai dns retrieve-zoneconfig")) (or (eq .HelpName "akamai-dns update-zoneconfig") (eq .HelpName "akamai dns update-zoneconfig"))}}` +
			color.YellowString("Record Types:\n") +
			`{{else}}` +
			color.YellowString("Built-In Commands:\n") +
			`{{end}}` +
			"{{range .VisibleCategories}}" +
			"{{if .Name}}" +
			"\n{{.Name}}\n" +
			"{{end}}" +
			"{{range .VisibleCommands}}" +
			color.GreenString("  {{.Name}}") +
			"{{if .Aliases}} ({{ $length := len .Aliases }}{{if eq $length 1}}alias:{{else}}aliases:{{end}} " +
			"{{range $index, $alias := .Aliases}}" +
			"{{if $index}}, {{end}}" +
			color.GreenString("{{$alias}}") +
			"{{end}}" +
			"){{end}}\n" +
			"{{end}}" +
			"{{end}}" +
			"{{end}}\n" +

			"{{if .Copyright}}" +
			color.HiBlackString("{{.Copyright}}") +
			"{{end}}\n"

	cli.CommandHelpTemplate =
		color.YellowString("Name: \n") +
			"   {{.HelpName}}\n\n" +

			`{{if .Description}}` +
			color.YellowString("Description: \n") +
			"   {{.Description}}\n\n" +
			`{{end}}` +

			color.YellowString("Usage: \n") +
			color.BlueString("   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .ArgsUsage}}{{.ArgsUsage}}{{end}} {{if .VisibleFlags}}{{range .VisibleFlags}}[--{{.Name}}] {{end}}{{end}}{{end}}\n\n") +

			"{{if .Category}}" +
			color.YellowString("Type: \n") +
			"   {{.Category}}\n\n{{end}}" +

			"{{if .VisibleFlags}}" +
			color.YellowString("Flags: \n") +
			"{{range .VisibleFlags}}   {{.}}\n{{end}}\n{{end}}" +

			"{{if .Subcommands}}" +
			color.YellowString("Record Types: \n") +
			"{{range .Subcommands}}   {{.Name}}\n{{end}}{{end}}"

	cli.SubcommandHelpTemplate = cli.CommandHelpTemplate
}
