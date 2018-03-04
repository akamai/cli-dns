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
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v1"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/dshafik/gozone"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdUpdateZone(c *cli.Context) error {
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
	akamai.StopSpinnerOk()

	var reader io.Reader
	if c.IsSet("file") {
		reader, err = os.Open(c.String("file"))
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Unable to open file (%s)", c.String("file")), 1)
		}
	} else {
		reader = bufio.NewReader(os.Stdin)
	}

	if c.IsSet("overwrite") {
		newZone := dns.NewZone(hostname)
		newZone.Token = zone.Token
		zone = newZone
	}

	if !c.IsSet("dns") {
		contents, err := ioutil.ReadAll(reader)
		if err != nil {
			return cli.NewExitError("An error occurred reading input", 1)
		}

		if c.IsSet("overwrite") {
			err = json.Unmarshal(contents, zone)
			if err != nil {
				return cli.NewExitError(color.RedString("An error occurred parsing input:\n"+err.Error()), 1)
			}
		} else {
			newZone := dns.NewZone(hostname)
			err = json.Unmarshal(contents, newZone)
			if len(newZone.Zone.A) > 0 {
				for _, record := range newZone.Zone.A {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Aaaa) > 0 {
				for _, record := range newZone.Zone.Aaaa {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Afsdb) > 0 {
				for _, record := range newZone.Zone.Afsdb {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Cname) > 0 {
				for _, record := range newZone.Zone.Cname {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Dnskey) > 0 {
				for _, record := range newZone.Zone.Dnskey {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Ds) > 0 {
				for _, record := range newZone.Zone.Ds {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Hinfo) > 0 {
				for _, record := range newZone.Zone.Hinfo {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Loc) > 0 {
				for _, record := range newZone.Zone.Loc {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Mx) > 0 {
				for _, record := range newZone.Zone.Mx {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Naptr) > 0 {
				for _, record := range newZone.Zone.Naptr {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Ns) > 0 {
				for _, record := range newZone.Zone.Ns {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Nsec3) > 0 {
				for _, record := range newZone.Zone.Nsec3 {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Nsec3param) > 0 {
				for _, record := range newZone.Zone.Nsec3param {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Ptr) > 0 {
				for _, record := range newZone.Zone.Ptr {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Rp) > 0 {
				for _, record := range newZone.Zone.Rp {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Rrsig) > 0 {
				for _, record := range newZone.Zone.Rrsig {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if newZone.Zone.Soa != nil && newZone.Zone.Soa.Serial > 0 {
				zone.Zone.Soa = newZone.Zone.Soa
			}
			if len(newZone.Zone.Spf) > 0 {
				for _, record := range newZone.Zone.Spf {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Srv) > 0 {
				for _, record := range newZone.Zone.Srv {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Sshfp) > 0 {
				for _, record := range newZone.Zone.Sshfp {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
			if len(newZone.Zone.Txt) > 0 {
				for _, record := range newZone.Zone.Txt {
					zone.RemoveRecord(record)
					zone.AddRecord(record)
				}
			}
		}
	}

	if c.IsSet("dns") {
		var record gozone.Record
		scanner := gozone.NewScanner(reader)

		i := 0
		for {
			i++
			err := scanner.Next(&record)
			if err != nil {
				if err.Error() != "EOF" {
					return cli.NewExitError("An error occurred parsing input:\n"+err.Error(), 1)
				}
				break
			}

			if strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname) == record.DomainName {
				return cli.NewExitError(color.RedString("%s Record on line %d does not match hostname \"%s\"", record.Type.String(), i, hostname), 1)
			}

			switch record.Type {
			case gozone.RecordType_A:
				newRecord := dns.NewARecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Target = record.Data[0]
				zone.RemoveRecord(newRecord)
				zone.AddRecord(newRecord)
			case gozone.RecordType_AAAA:
				newRecord := dns.NewAaaaRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Target = record.Data[0]
				zone.AddRecord(newRecord)
			case gozone.RecordType_AFSDB:
				newRecord := dns.NewAfsdbRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Subtype, _ = strconv.Atoi(record.Data[0])
				newRecord.Target = record.Data[1]
				zone.AddRecord(newRecord)
			case gozone.RecordType_CNAME:
				newRecord := dns.NewCnameRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Target = record.Data[0]
				if !strings.HasSuffix(newRecord.Target, ".") {
					newRecord.Target += "."
				}
				zone.AddRecord(newRecord)
			case gozone.RecordType_DNSKEY:
				newRecord := dns.NewDnskeyRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Flags, _ = strconv.Atoi(record.Data[0])
				newRecord.Protocol, _ = strconv.Atoi(record.Data[1])
				newRecord.Algorithm, _ = strconv.Atoi(record.Data[2])
				newRecord.Key = strings.Trim(record.Data[3], `"`)
				zone.AddRecord(newRecord)
			case gozone.RecordType_DS:
				newRecord := dns.NewDsRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Keytag, _ = strconv.Atoi(record.Data[0])
				newRecord.Algorithm, _ = strconv.Atoi(record.Data[1])
				newRecord.DigestType, _ = strconv.Atoi(record.Data[2])
				newRecord.Digest = record.Data[3]
				zone.AddRecord(newRecord)
			case gozone.RecordType_HINFO:
				newRecord := dns.NewHinfoRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Hardware = strings.Trim(record.Data[0], `"`)
				newRecord.Software = strings.Trim(record.Data[1], `"`)
				zone.AddRecord(newRecord)
			case gozone.RecordType_LOC:
				newRecord := dns.NewLocRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Target = record.Data[0]
				zone.AddRecord(newRecord)
			case gozone.RecordType_MX:
				newRecord := dns.NewMxRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Priority, _ = strconv.Atoi(record.Data[0])
				newRecord.Target = record.Data[1]
				if !strings.HasSuffix(newRecord.Target, ".") {
					newRecord.Target += "."
				}
				zone.AddRecord(newRecord)
			case gozone.RecordType_NAPTR:
				newRecord := dns.NewNaptrRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				order, _ := strconv.Atoi(record.Data[0])
				newRecord.Order = uint16(order)
				preference, _ := strconv.Atoi(record.Data[1])
				newRecord.Preference = uint16(preference)
				newRecord.Flags = strings.Trim(record.Data[2], `"`)
				newRecord.Service = strings.Trim(record.Data[3], `"`)
				newRecord.Regexp = strings.Trim(record.Data[4], `"`)
				newRecord.Replacement = strings.Trim(record.Data[5], `"`)
				zone.AddRecord(newRecord)
			case gozone.RecordType_NS:
				newRecord := dns.NewNsRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Target = record.Data[0]
				if !strings.HasSuffix(newRecord.Target, ".") {
					newRecord.Target += "."
				}
				zone.AddRecord(newRecord)
			case gozone.RecordType_NSEC3:
				newRecord := dns.NewNsec3Record()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Algorithm, _ = strconv.Atoi(record.Data[0])
				newRecord.Flags, _ = strconv.Atoi(record.Data[1])
				newRecord.Iterations, _ = strconv.Atoi(record.Data[2])
				newRecord.Salt = record.Data[3]
				newRecord.NextHashedOwnerName = record.Data[4]
				newRecord.TypeBitmaps = strings.Join(record.Data[5:], " ")
				zone.AddRecord(newRecord)
			case gozone.RecordType_NSEC3PARAM:
				newRecord := dns.NewNsec3paramRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Algorithm, _ = strconv.Atoi(record.Data[0])
				newRecord.Flags, _ = strconv.Atoi(record.Data[1])
				newRecord.Iterations, _ = strconv.Atoi(record.Data[2])
				newRecord.Salt = record.Data[3]
				zone.AddRecord(newRecord)
			case gozone.RecordType_PTR:
				newRecord := dns.NewPtrRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Target = record.Data[0]
				if !strings.HasSuffix(newRecord.Target, ".") {
					newRecord.Target += "."
				}
				zone.AddRecord(newRecord)
			case gozone.RecordType_RP:
				newRecord := dns.NewRpRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Mailbox = record.Data[0]
				if len(record.Data) > 0 {
					newRecord.Txt = record.Data[1]
				}
				zone.AddRecord(newRecord)
			case gozone.RecordType_RRSIG:
				newRecord := dns.NewRrsigRecord()
				newRecord.TTL = int(record.TimeToLive)
				newRecord.TypeCovered = record.Data[0]
				newRecord.Algorithm, _ = strconv.Atoi(record.Data[1])
				newRecord.Labels, _ = strconv.Atoi(record.Data[2])
				newRecord.OriginalTTL, _ = strconv.Atoi(record.Data[3])
				newRecord.Expiration = record.Data[4]
				newRecord.Inception = record.Data[5]
				newRecord.Keytag, _ = strconv.Atoi(record.Data[6])
				newRecord.Signer = record.Data[7]
				newRecord.Signature = record.Data[8]
				zone.AddRecord(newRecord)
			case gozone.RecordType_SOA:
				newRecord := dns.NewSoaRecord()
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Originserver = record.Data[0]
				newRecord.Contact = record.Data[1]
				serial, _ := strconv.Atoi(record.Data[2])
				newRecord.Serial = uint(serial)
				newRecord.Refresh, _ = strconv.Atoi(record.Data[3])
				newRecord.Retry, _ = strconv.Atoi(record.Data[4])
				newRecord.Expire, _ = strconv.Atoi(record.Data[5])
				minimum, _ := strconv.Atoi(record.Data[6])
				newRecord.Minimum = uint(minimum)
				zone.AddRecord(newRecord)
			case gozone.RecordType_SPF:
				newRecord := dns.NewSpfRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Target = strings.Trim(record.Data[0], `"`)
				zone.AddRecord(newRecord)
			case gozone.RecordType_SRV:
				newRecord := dns.NewSrvRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Priority, _ = strconv.Atoi(record.Data[0])
				weight, _ := strconv.Atoi(record.Data[1])
				newRecord.Weight = uint16(weight)
				port, _ := strconv.Atoi(record.Data[2])
				newRecord.Port = uint16(port)
				newRecord.Target = record.Data[3]
				if !strings.HasSuffix(newRecord.Target, ".") {
					newRecord.Target += "."
				}
				zone.AddRecord(newRecord)
			case gozone.RecordType_SSHFP:
				newRecord := dns.NewSshfpRecord()
				newRecord.Algorithm, _ = strconv.Atoi(record.Data[0])
				newRecord.FingerprintType, _ = strconv.Atoi(record.Data[1])
				newRecord.Fingerprint = record.Data[2]
				zone.AddRecord(newRecord)
			case gozone.RecordType_TXT:
				newRecord := dns.NewTxtRecord()
				newRecord.Name = strings.TrimSuffix(strings.TrimSuffix(record.DomainName, hostname+"."), hostname)
				newRecord.Active = true
				newRecord.TTL = int(record.TimeToLive)
				newRecord.Target = strings.Trim(record.Data[0], `"`)
				zone.AddRecord(newRecord)
			}
		}
	}

	akamai.StartSpinner("Updating zone...", fmt.Sprintf("Updating zone...... [%s]", color.GreenString("OK")))
	err = zone.Save()
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError("Error saving zone: "+err.Error(), 1)
	}
	akamai.StopSpinnerOk()

	return nil
}
