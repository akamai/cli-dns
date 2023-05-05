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
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v1"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdCreateRecord(c *cli.Context) error {
	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dns.Config = config

	hostname := c.Args().First()
	if hostname == "" {
		return cli.NewExitError(color.RedString("hostname is required"), 1)
	}

	err = validateFields(c.Command.Name, c)
	if err != nil {
		return err
	}

	akamai.StartSpinner(
		fmt.Sprintf("Adding new %s record...", c.Command.Name),
		fmt.Sprintf("Adding new %s record...... [%s]", c.Command.Name, color.GreenString("OK")),
	)
	zone, err := dns.GetZone(hostname)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError("Unable to fetch zone", 1)
	}

	var active bool
	if !c.IsSet("inactive") {
		active = true
	}

	var record dns.DNSRecord

	switch c.Command.Name {
	case "A":
		record = dns.NewARecord()
		record.(*dns.ARecord).Active = active
		record.(*dns.ARecord).Name = c.String("name")
		record.(*dns.ARecord).Target = c.String("target")
		record.(*dns.ARecord).TTL = c.Int("ttl")
	case "AAAA":
		record = dns.NewAaaaRecord()
		record.(*dns.AaaaRecord).Active = active
		record.(*dns.AaaaRecord).Name = c.String("name")
		record.(*dns.AaaaRecord).Target = c.String("target")
		record.(*dns.AaaaRecord).TTL = c.Int("ttl")
	case "AFSDB":
		record = dns.NewAfsdbRecord()
		record.(*dns.AfsdbRecord).Active = active
		record.(*dns.AfsdbRecord).Name = c.String("name")
		record.(*dns.AfsdbRecord).Subtype = c.Int("subtype")
		record.(*dns.AfsdbRecord).Target = c.String("target")
		record.(*dns.AfsdbRecord).TTL = c.Int("ttl")
	case "CNAME":
		record = dns.NewCnameRecord()
		record.(*dns.CnameRecord).Active = active
		record.(*dns.CnameRecord).Name = c.String("name")
		record.(*dns.CnameRecord).Target = c.String("target")
		record.(*dns.CnameRecord).TTL = c.Int("ttl")

		if !strings.HasSuffix(c.String("target"), ".") {
			record.(*dns.CnameRecord).Target += "."
		}
	case "DNSKEY":
		record = dns.NewDnskeyRecord()
		record.(*dns.DnskeyRecord).Active = active
		record.(*dns.DnskeyRecord).Algorithm = c.Int("algorithm")
		record.(*dns.DnskeyRecord).Flags = c.Int("flags")
		record.(*dns.DnskeyRecord).Key = c.String("key")
		record.(*dns.DnskeyRecord).Name = c.String("name")
		record.(*dns.DnskeyRecord).Protocol = c.Int("protocol")
		record.(*dns.DnskeyRecord).TTL = c.Int("ttl")
	case "DS":
		record = dns.NewDsRecord()
		record.(*dns.DsRecord).Active = active
		record.(*dns.DsRecord).Algorithm = c.Int("algorithm")
		record.(*dns.DsRecord).Digest = c.String("digest")
		record.(*dns.DsRecord).DigestType = c.Int("digest-type")
		record.(*dns.DsRecord).Keytag = c.Int("keytag")
		record.(*dns.DsRecord).Name = c.String("name")
		record.(*dns.DsRecord).TTL = c.Int("ttl")
	case "HINFO":
		record = dns.NewHinfoRecord()
		record.(*dns.HinfoRecord).Active = active
		record.(*dns.HinfoRecord).Hardware = c.String("hardware")
		record.(*dns.HinfoRecord).Name = c.String("name")
		record.(*dns.HinfoRecord).Software = c.String("software")
		record.(*dns.HinfoRecord).TTL = c.Int("ttl")
	case "LOC":
		record = dns.NewLocRecord()
		record.(*dns.LocRecord).Active = active
		record.(*dns.LocRecord).Name = c.String("name")
		record.(*dns.LocRecord).Target = c.String("target")
		record.(*dns.LocRecord).TTL = c.Int("ttl")
	case "MX":
		record = dns.NewMxRecord()
		record.(*dns.MxRecord).Active = active
		record.(*dns.MxRecord).Name = c.String("name")
		record.(*dns.MxRecord).Priority = c.Int("priority")
		record.(*dns.MxRecord).Target = c.String("target")
		record.(*dns.MxRecord).TTL = c.Int("ttl")

		if !strings.HasSuffix(c.String("target"), ".") {
			record.(*dns.MxRecord).Target += "."
		}
	case "NAPTR":
		record = dns.NewNaptrRecord()
		record.(*dns.NaptrRecord).Active = active
		record.(*dns.NaptrRecord).Flags = c.String("flags")
		record.(*dns.NaptrRecord).Name = c.String("name")
		record.(*dns.NaptrRecord).Order = uint16(c.Uint("order"))
		record.(*dns.NaptrRecord).Preference = uint16(c.Uint("preference"))
		record.(*dns.NaptrRecord).Regexp = c.String("regexp")
		record.(*dns.NaptrRecord).Replacement = c.String("replacement")
		record.(*dns.NaptrRecord).Service = c.String("service")
		record.(*dns.NaptrRecord).TTL = c.Int("ttl")
	case "NS":
		record = dns.NewNsRecord()
		record.(*dns.NsRecord).Active = active
		record.(*dns.NsRecord).Name = c.String("name")
		record.(*dns.NsRecord).Target = c.String("target")
		record.(*dns.NsRecord).TTL = c.Int("ttl")

		if !strings.HasSuffix(c.String("target"), ".") {
			record.(*dns.NsRecord).Target += "."
		}
	case "NSEC3":
		record = dns.NewNsec3Record()
		record.(*dns.Nsec3Record).Active = active
		record.(*dns.Nsec3Record).Algorithm = c.Int("algorithm")
		record.(*dns.Nsec3Record).Flags = c.Int("flags")
		record.(*dns.Nsec3Record).Iterations = c.Int("iterations")
		record.(*dns.Nsec3Record).Name = c.String("name")
		record.(*dns.Nsec3Record).NextHashedOwnerName = c.String("next-hashed-owner-name")
		record.(*dns.Nsec3Record).Salt = c.String("salt")
		record.(*dns.Nsec3Record).TTL = c.Int("ttl")
		record.(*dns.Nsec3Record).TypeBitmaps = c.String("type-bitmaps")
	case "NSEC3PARAM":
		record = dns.NewNsec3paramRecord()
		record.(*dns.Nsec3paramRecord).Active = active
		record.(*dns.Nsec3paramRecord).Algorithm = c.Int("algorithm")
		record.(*dns.Nsec3paramRecord).Flags = c.Int("flags")
		record.(*dns.Nsec3paramRecord).Iterations = c.Int("iterations")
		record.(*dns.Nsec3paramRecord).Name = c.String("name")
		record.(*dns.Nsec3paramRecord).Salt = c.String("salt")
		record.(*dns.Nsec3paramRecord).TTL = c.Int("ttl")
	case "PTR":
		record = dns.NewPtrRecord()
		record.(*dns.PtrRecord).Active = active
		record.(*dns.PtrRecord).Name = c.String("name")
		record.(*dns.PtrRecord).Target = c.String("target")
		record.(*dns.PtrRecord).TTL = c.Int("ttl")

		if !strings.HasSuffix(c.String("target"), ".") {
			record.(*dns.PtrRecord).Target += "."
		}
	case "RP":
		record = dns.NewRpRecord()
		record.(*dns.RpRecord).Active = active
		record.(*dns.RpRecord).Mailbox = c.String("mailbot")
		record.(*dns.RpRecord).Name = c.String("name")
		record.(*dns.RpRecord).TTL = c.Int("ttl")
		record.(*dns.RpRecord).Txt = c.String("txt")
	case "RRSIG":
		record = dns.NewRrsigRecord()
		record.(*dns.RrsigRecord).Active = active
		record.(*dns.RrsigRecord).Algorithm = c.Int("algorithm")
		record.(*dns.RrsigRecord).Expiration = c.String("expiration")
		record.(*dns.RrsigRecord).Inception = c.String("inception")
		record.(*dns.RrsigRecord).Keytag = c.Int("keytag")
		record.(*dns.RrsigRecord).Labels = c.Int("labels")
		record.(*dns.RrsigRecord).Name = c.String("name")
		record.(*dns.RrsigRecord).OriginalTTL = c.Int("original-ttl")
		record.(*dns.RrsigRecord).Signature = c.String("signature")
		record.(*dns.RrsigRecord).Signer = c.String("signer")
		record.(*dns.RrsigRecord).TTL = c.Int("ttl")
		record.(*dns.RrsigRecord).TypeCovered = c.String("type-covered")
	case "SOA":
		record = dns.NewSoaRecord()
		record.(*dns.SoaRecord).Contact = c.String("contact")
		record.(*dns.SoaRecord).Expire = c.Int("expire")
		record.(*dns.SoaRecord).Minimum = c.Uint("minimum")
		record.(*dns.SoaRecord).Originserver = c.String("originserver")
		record.(*dns.SoaRecord).Refresh = c.Int("refresh")
		record.(*dns.SoaRecord).Retry = c.Int("retry")
		record.(*dns.SoaRecord).Serial = c.Uint("serial")
		record.(*dns.SoaRecord).TTL = c.Int("ttl")
	case "SPF":
		record = dns.NewSpfRecord()
		record.(*dns.SpfRecord).Active = active
		record.(*dns.SpfRecord).Name = c.String("name")
		record.(*dns.SpfRecord).Target = c.String("target")
		record.(*dns.SpfRecord).TTL = c.Int("ttl")
	case "SRV":
		record = dns.NewSrvRecord()
		record.(*dns.SrvRecord).Active = active
		record.(*dns.SrvRecord).Name = c.String("name")
		record.(*dns.SrvRecord).Port = uint16(c.Uint("port"))
		record.(*dns.SrvRecord).Priority = c.Int("priority")
		record.(*dns.SrvRecord).Target = c.String("target")
		record.(*dns.SrvRecord).TTL = c.Int("ttl")
		record.(*dns.SrvRecord).Weight = uint16(c.Uint("weight"))

		if !strings.HasSuffix(c.String("target"), ".") {
			record.(*dns.SrvRecord).Target += "."
		}
	case "SSHFP":
		record = dns.NewSshfpRecord()
		record.(*dns.SshfpRecord).Active = active
		record.(*dns.SshfpRecord).Algorithm = c.Int("algorithm")
		record.(*dns.SshfpRecord).Fingerprint = c.String("fingerprint")
		record.(*dns.SshfpRecord).FingerprintType = c.Int("fingerprint-type")
		record.(*dns.SshfpRecord).Name = c.String("name")
		record.(*dns.SshfpRecord).TTL = c.Int("ttl")
	case "TXT":
		record = dns.NewTxtRecord()
		record.(*dns.TxtRecord).Active = active
		record.(*dns.TxtRecord).Name = c.String("name")
		record.(*dns.TxtRecord).Target = c.String("target")
		record.(*dns.TxtRecord).TTL = c.Int("ttl")
	}

	zone.AddRecord(record)
	err = zone.Save()
	if err != nil {
		akamai.StopSpinnerFail()
		fmt.Printf("%#v\n", err)
		return cli.NewExitError(err.Error(), 1)
	}

	akamai.StopSpinnerOk()
	return nil
}

func validateFields(recordType string, c *cli.Context) error {
	var missing []string = make([]string, 0)
	for option, settings := range recordOptions[recordType] {
		if settings.required {
			if !c.IsSet(option) {
				missing = append(missing, option)
			}
		}
	}

	if len(missing) != 0 {
		error := "Missing required options: \n"
		for _, option := range missing {
			error += "  " + fmt.Sprintf("--%s", option) + "\n"
		}

		return cli.NewExitError(error, 1)
	}

	return nil
}
