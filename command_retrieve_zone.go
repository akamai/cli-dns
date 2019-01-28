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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v1"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func cmdRetrieveZone(c *cli.Context) error {
	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dns.Config = config

	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("hostname is required"), 1)
	}

	hostname := c.Args().First()
	akamai.StartSpinner(
		"Fetching zone...",
		fmt.Sprintf("Fetching zone...... [%s]", color.GreenString("OK")),
	)
	zone, err := dns.GetZone(hostname)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString("Zone not found "), 1)
	}

	outputZone := dns.NewZone(hostname)
	outputZone.Token = zone.Token
	if c.IsSet("json") && c.Bool("json") && !c.IsSet("filter") {
		json, err := json.MarshalIndent(zone, "", "  ")
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to display zone"), 1)
		}
		akamai.StopSpinnerOk()

		fmt.Fprintln(c.App.Writer, string(json))
		return nil
	}

	filter := make(map[string]bool)
	if c.IsSet("filter") {
		for _, recordType := range c.StringSlice("filter") {
			filter[strings.ToUpper(recordType)] = true
		}
	}

	if _, ok := filter["A"]; (!c.IsSet("filter") || ok) && len(zone.Zone.A) > 0 {
		outputZone.Zone.A = zone.Zone.A
	}
	if _, ok := filter["AAAA"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Aaaa) > 0 {
		outputZone.Zone.Aaaa = zone.Zone.Aaaa
	}
	if _, ok := filter["AFSDB"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Afsdb) > 0 {
		outputZone.Zone.Afsdb = zone.Zone.Afsdb
	}
	if _, ok := filter["CNAME"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Cname) > 0 {
		outputZone.Zone.Cname = zone.Zone.Cname
	}
	if _, ok := filter["DNSKEY"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Dnskey) > 0 {
		outputZone.Zone.Dnskey = zone.Zone.Dnskey
	}
	if _, ok := filter["DS"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Ds) > 0 {
		outputZone.Zone.Ds = zone.Zone.Ds
	}
	if _, ok := filter["HINFO"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Hinfo) > 0 {
		outputZone.Zone.Hinfo = zone.Zone.Hinfo
	}
	if _, ok := filter["LOC"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Loc) > 0 {
		outputZone.Zone.Loc = zone.Zone.Loc
	}
	if _, ok := filter["MX"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Mx) > 0 {
		outputZone.Zone.Mx = zone.Zone.Mx
	}
	if _, ok := filter["NAPTR"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Naptr) > 0 {
		outputZone.Zone.Naptr = zone.Zone.Naptr
	}
	if _, ok := filter["NS"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Ns) > 0 {
		outputZone.Zone.Ns = zone.Zone.Ns
	}
	if _, ok := filter["NSEC3"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Nsec3) > 0 {
		outputZone.Zone.Nsec3 = zone.Zone.Nsec3
	}
	if _, ok := filter["NSEC3PARAM"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Nsec3param) > 0 {
		outputZone.Zone.Nsec3param = zone.Zone.Nsec3param
	}
	if _, ok := filter["PTR"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Ptr) > 0 {
		outputZone.Zone.Ptr = zone.Zone.Ptr
	}
	if _, ok := filter["RP"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Rp) > 0 {
		outputZone.Zone.Rp = zone.Zone.Rp
	}
	if _, ok := filter["RRSIG"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Rrsig) > 0 {
		outputZone.Zone.Rrsig = zone.Zone.Rrsig
	}
	if _, ok := filter["SOA"]; (!c.IsSet("filter") || ok) && zone.Zone.Soa != nil {
		outputZone.Zone.Soa = zone.Zone.Soa
	}

	if _, ok := filter["SPF"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Spf) > 0 {
		outputZone.Zone.Spf = zone.Zone.Spf
	}
	if _, ok := filter["SRV"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Srv) > 0 {
		outputZone.Zone.Srv = zone.Zone.Srv
	}
	if _, ok := filter["SSHFP"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Sshfp) > 0 {
		outputZone.Zone.Sshfp = zone.Zone.Sshfp
	}
	if _, ok := filter["TXT"]; (!c.IsSet("filter") || ok) && len(zone.Zone.Txt) > 0 {
		outputZone.Zone.Txt = zone.Zone.Txt
	}

	if c.IsSet("json") && c.Bool("json") {
		json, err := json.MarshalIndent(outputZone, "", "  ")
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to display zone"), 1)
		}
		akamai.StopSpinnerOk()

		fmt.Fprintln(c.App.Writer, string(json))
		return nil
	} else {
		akamai.StopSpinnerOk()

		fmt.Fprintln(c.App.Writer, "")
		renderZoneTable(outputZone, c)
	}

	return nil
}

func renderZoneTable(zone *dns.Zone, c *cli.Context) {
	bold := color.New(color.FgWhite, color.Bold)

	table := tablewriter.NewWriter(c.App.Writer)
	table.SetHeader([]string{"INDEX", "TYPE", "NAME", "TTL", "ACTIVE", "OPTIONS"})
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)

	i := 0
	if len(zone.Zone.A) > 0 {
		for _, record := range zone.Zone.A {
			i++
			options := fmt.Sprintf("%s: %s", bold.Sprintf("target"), record.Target)

			values := []string{strconv.Itoa(i), "A", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Aaaa) > 0 {
		for _, record := range zone.Zone.Aaaa {
			i++
			options := fmt.Sprintf("%s: %s", bold.Sprintf("target"), record.Target)

			values := []string{strconv.Itoa(i), "AAAA", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Afsdb) > 0 {
		for _, record := range zone.Zone.Afsdb {
			i++
			options := fmt.Sprintf("%s: %s\n", bold.Sprintf("target"), record.Target)
			options += fmt.Sprintf("%s: %d", bold.Sprintf("subtype"), record.Subtype)

			values := []string{strconv.Itoa(i), "AFSDB", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Cname) > 0 {
		for _, record := range zone.Zone.Cname {
			i++
			options := fmt.Sprintf("%s: %s", bold.Sprintf("target"), record.Target)

			values := []string{strconv.Itoa(i), "CNAME", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Dnskey) > 0 {
		for _, record := range zone.Zone.Dnskey {
			i++
			options := fmt.Sprintf("%s: %d\n", bold.Sprintf("flags"), record.Flags)
			options += fmt.Sprintf("%s: %d\n", bold.Sprintf("protocol"), record.Protocol)
			options += fmt.Sprintf("%s: %d\n", bold.Sprintf("algorithm"), record.Algorithm)
			options += fmt.Sprintf("%s: %s", bold.Sprintf("key"), record.Key)

			values := []string{strconv.Itoa(i), "DNSKEY", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Ds) > 0 {
		for _, record := range zone.Zone.Ds {
			i++
			options := fmt.Sprintf("%s: %d\n", bold.Sprintf("keytag"), record.Keytag)
			options += fmt.Sprintf("%s: %d\n", bold.Sprintf("algorithm"), record.Algorithm)
			options += fmt.Sprintf("%s: %d\n", bold.Sprintf("digest-type"), record.DigestType)
			options += fmt.Sprintf("%s: %s", bold.Sprintf("digest"), record.Digest)

			values := []string{strconv.Itoa(i), "DS", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Hinfo) > 0 {
		for _, record := range zone.Zone.Hinfo {
			i++
			options := fmt.Sprintf("%s: %s\n", bold.Sprintf("hardware"), record.Hardware)
			options += fmt.Sprintf("%s: %s", bold.Sprintf("software"), record.Software)

			values := []string{strconv.Itoa(i), "HINFO", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Loc) > 0 {
		for _, record := range zone.Zone.Loc {
			i++
			options := fmt.Sprintf("%s: %s", bold.Sprintf("target"), record.Target)

			values := []string{strconv.Itoa(i), "LOC", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Mx) > 0 {
		for _, record := range zone.Zone.Mx {
			i++
			options := fmt.Sprintf("%s: %s\n", bold.Sprintf("target"), record.Target)
			options += fmt.Sprintf("%s: %d", bold.Sprintf("priority"), record.Priority)

			values := []string{strconv.Itoa(i), "MX", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Naptr) > 0 {
		for _, record := range zone.Zone.Naptr {
			i++
			options := fmt.Sprintf("%s: %d\n", bold.Sprintf("order"), record.Order)
			options += fmt.Sprintf("%s: %d\n", bold.Sprintf("preference"), record.Preference)
			options += fmt.Sprintf("%s: %s\n", bold.Sprintf("flags"), record.Flags)
			options += fmt.Sprintf("%s: %s\n", bold.Sprintf("service"), record.Service)
			options += fmt.Sprintf("%s: %s\n", bold.Sprintf("regexp"), record.Regexp)
			options += fmt.Sprintf("%s: %s", bold.Sprintf("replacement"), record.Replacement)

			values := []string{strconv.Itoa(i), "NAPTR", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Ns) > 0 {
		for _, record := range zone.Zone.Ns {
			i++
			options := fmt.Sprintf("%s: %s", bold.Sprintf("target"), record.Target)

			values := []string{strconv.Itoa(i), "NS", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Nsec3) > 0 {
		for _, record := range zone.Zone.Nsec3 {
			i++
			options := fmt.Sprintf("%s: %d\n", bold.Sprintf("algorithm"), record.Algorithm)
			options += fmt.Sprintf("%s: %d\n", bold.Sprintf("flags"), record.Flags)
			options += fmt.Sprintf("%s: %d\n", bold.Sprintf("iterations"), record.Iterations)
			options += fmt.Sprintf("%s: %s\n", bold.Sprintf("salt"), record.Salt)
			options += fmt.Sprintf("%s: %s\n", bold.Sprintf("next-hashed-owner-name"), record.NextHashedOwnerName)
			options += fmt.Sprintf("%s: %s", bold.Sprintf("type-bitmaps"), record.TypeBitmaps)

			values := []string{strconv.Itoa(i), "NSEC3", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Nsec3param) > 0 {
		for _, record := range zone.Zone.Nsec3param {
			i++
			options := fmt.Sprintf("%s: %d\n", bold.Sprintf("algorithm"), record.Algorithm)
			options += fmt.Sprintf("%s: %d\n", bold.Sprintf("flags"), record.Flags)
			options += fmt.Sprintf("%s: %d\n", bold.Sprintf("iterations"), record.Iterations)
			options += fmt.Sprintf("%s: %s", bold.Sprintf("salt"), record.Salt)

			values := []string{strconv.Itoa(i), "NSEC3PARAM", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Ptr) > 0 {
		for _, record := range zone.Zone.Ptr {
			i++
			options := fmt.Sprintf("%s: %s", bold.Sprintf("target"), record.Target)

			values := []string{strconv.Itoa(i), "PTR", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Rp) > 0 {
		for _, record := range zone.Zone.Rp {
			i++
			options := fmt.Sprintf("%s: %s\n", bold.Sprintf("mailbox"), record.Mailbox)
			options += fmt.Sprintf("%s: %s", bold.Sprintf("txt"), record.Txt)

			values := []string{strconv.Itoa(i), "RP", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Rrsig) > 0 {
		for _, record := range zone.Zone.Rrsig {
			i++
			options := fmt.Sprintf("%s: %s\n", bold.Sprintf("type-covered"), record.TypeCovered)
			options += fmt.Sprintf("%s: %d\n", bold.Sprintf("algorithm"), record.Algorithm)
			options += fmt.Sprintf("%s: %d\n", bold.Sprintf("original-TTL"), record.OriginalTTL)
			options += fmt.Sprintf("%s: %s\n", bold.Sprintf("expiration"), record.Expiration)
			options += fmt.Sprintf("%s: %s\n", bold.Sprintf("inception"), record.Inception)
			options += fmt.Sprintf("%s: %d\n", bold.Sprintf("keytag"), record.Keytag)
			options += fmt.Sprintf("%s: %s\n", bold.Sprintf("signer"), record.Signer)
			options += fmt.Sprintf("%s: %s\n", bold.Sprintf("signature"), record.Signature)
			options += fmt.Sprintf("%s: %d", bold.Sprintf("labels"), record.Labels)

			values := []string{strconv.Itoa(i), "RRSIG", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if zone.Zone.Soa != nil && zone.Zone.Soa.Serial > 0 {
		i++
		options := fmt.Sprintf("%s: %s\n", bold.Sprintf("originserver"), zone.Zone.Soa.Originserver)
		options += fmt.Sprintf("%s: %s\n", bold.Sprintf("contact"), zone.Zone.Soa.Contact)
		options += fmt.Sprintf("%s: %d\n", bold.Sprintf("serial"), zone.Zone.Soa.Serial)
		options += fmt.Sprintf("%s: %d\n", bold.Sprintf("refresh"), zone.Zone.Soa.Refresh)
		options += fmt.Sprintf("%s: %d\n", bold.Sprintf("retry"), zone.Zone.Soa.Retry)
		options += fmt.Sprintf("%s: %d\n", bold.Sprintf("expire"), zone.Zone.Soa.Expire)
		options += fmt.Sprintf("%s: %d", bold.Sprintf("minimum"), zone.Zone.Soa.Minimum)

		values := []string{strconv.Itoa(i), "SOA", "", strconv.Itoa(zone.Zone.Soa.TTL), fmt.Sprintf("%t", true), options}
		table.Append(values)
	}

	if len(zone.Zone.Spf) > 0 {
		for _, record := range zone.Zone.Spf {
			i++
			options := fmt.Sprintf("%s: %s\n", bold.Sprintf("target"), record.Target)

			values := []string{strconv.Itoa(i), "SPF", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Srv) > 0 {
		for _, record := range zone.Zone.Srv {
			i++
			options := fmt.Sprintf("%s: %s\n", bold.Sprintf("target"), record.Target)
			options += fmt.Sprintf("%s: %d\n", bold.Sprintf("priority"), record.Priority)
			options += fmt.Sprintf("%s: %d\n", bold.Sprintf("weight"), record.Weight)
			options += fmt.Sprintf("%s: %d", bold.Sprintf("port"), record.Port)

			values := []string{strconv.Itoa(i), "SRV", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Sshfp) > 0 {
		for _, record := range zone.Zone.Sshfp {
			i++
			options := fmt.Sprintf("%s: %d\n", bold.Sprintf("algorithm"), record.Algorithm)
			options += fmt.Sprintf("%s: %d\n", bold.Sprintf("fingerprint-type"), record.FingerprintType)
			options += fmt.Sprintf("%s: %s", bold.Sprintf("fingerprint"), record.Fingerprint)

			values := []string{strconv.Itoa(i), "SSHFP", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}
	if len(zone.Zone.Txt) > 0 {
		for _, record := range zone.Zone.Txt {
			i++
			options := fmt.Sprintf("%s: %s", bold.Sprintf("target"), record.Target)

			values := []string{strconv.Itoa(i), "TXT", record.Name, strconv.Itoa(record.TTL), fmt.Sprintf("%t", record.Active), options}
			table.Append(values)
		}
	}

	if table.NumLines() == 0 {
		fmt.Fprintln(akamai.App.Writer, color.CyanString("No records found"))
		return
	}

	table.SetCaption(true, fmt.Sprintf("Zone: %s; Token: %s", zone.Zone.Name, zone.Token))
	table.Render()
}
