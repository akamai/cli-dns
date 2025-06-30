package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func renderRecordsetTable(zone string, record *dns.GetRecordResponse) string {
	return fmt.Sprintf(`
    Zone: %s
    Name: %s
    Type: %s
    TTL: %d
    Rdata:
      %s
      `,
		zone, record.Name, record.RecordType, record.TTL, strings.Join(record.Target, "\n "))
}

func renderRecordsetListTable(zone string, recordsets []dns.RecordSet) string {
	var out strings.Builder
	out.WriteString("\nZone Recordsets\n\n")
	table := tablewriter.NewWriter(&out)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_LEFT})
	table.SetHeader([]string{"NAME", "TYPE", "TTL", "RDATA"})
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetBorder(false)
	table.SetCaption(true, fmt.Sprintf("Zone: %s", zone))

	if len(recordsets) == 0 {
		rowData := []string{"No recordsets found", " ", " "}
		table.Append(rowData)
	} else {
		for _, set := range recordsets {
			name := set.Name
			typeVal := set.Type
			ttl := strconv.Itoa(set.TTL)
			//rdata := strings.Join(set.Rdata, ", ")
			for i, rdata := range set.Rdata {
				if i == 0 {
					table.Append([]string{name, typeVal, ttl, rdata})
				} else {
					table.Append([]string{" ", " ", " ", rdata})
				}
			}
		}
	}
	table.Render()
	return out.String()
}

func renderZoneconfigTable(zone *dns.GetZoneResponse, c *cli.Context) string {

	//bold := color.New(color.FgWhite, color.Bold)
	outString := ""
	outString += fmt.Sprintln(" ")
	outString += fmt.Sprintln("Zone Configuration")
	outString += fmt.Sprintln(" ")

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	table.SetHeader([]string{"ZONE", "ATTRIBUTE", "VALUE"})
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetBorder(false)

	if zone == nil {
		rowData := []string{"No zone info to display", " ", " "}
		table.Append(rowData)
	} else {
		zname := zone.Zone
		ztype := zone.Type
		table.Append([]string{zname, "Type", ztype})
		if len(zone.Comment) > 0 {
			table.Append([]string{" ", "Comment", zone.Comment})
		}
		if len(zone.ContractID) > 0 {
			table.Append([]string{" ", "ContractId", zone.ContractID})
		}
		if strings.ToUpper(ztype) == "SECONDARY" {
			if len(zone.Masters) > 0 {
				masters := strings.Join(zone.Masters, " ,")
				table.Append([]string{" ", "Masters", masters})
			}
			if zone.TSIGKey != nil {
				if len(zone.TSIGKey.Name) > 0 {
					table.Append([]string{" ", "TsigKey:Name", zone.TSIGKey.Name})
				}
				if len(zone.TSIGKey.Algorithm) > 0 {
					table.Append([]string{" ", "TsigKey:Algorithm", zone.TSIGKey.Algorithm})
				}
				if len(zone.TSIGKey.Secret) > 0 {
					table.Append([]string{" ", "TsigKey:Secret", zone.TSIGKey.Secret})
				}
			}
		}
		if strings.ToUpper(ztype) == "PRIMARY" || strings.ToUpper(ztype) == "SECONDARY" {
			table.Append([]string{" ", "SignAndServe", fmt.Sprintf("%t", zone.SignAndServe)})
			if len(zone.SignAndServeAlgorithm) > 0 {
				table.Append([]string{" ", "SignAndServeAlgorithm", fmt.Sprintf("%s", zone.SignAndServeAlgorithm)})
			}
		}
		if strings.ToUpper(ztype) == "ALIAS" {
			table.Append([]string{" ", "Target", zone.Target})
			table.Append([]string{" ", "AliasCount", strconv.FormatInt(zone.AliasCount, 10)})
		}
		table.Append([]string{" ", "ActivationState", zone.ActivationState})
		if len(zone.LastActivationDate) > 0 {
			table.Append([]string{" ", "LastActivationDate", zone.LastActivationDate})
		}
		if len(zone.LastModifiedDate) > 0 {
			table.Append([]string{" ", "LastModifiedDate", zone.LastModifiedDate})
		}
		table.Append([]string{" ", "VersionId", zone.VersionID})
	}
	table.Render()
	outString += fmt.Sprintln(tableString.String())

	return outString
}

func renderZoneListTable(zones []dns.ZoneResponse) string {
	outString := ""
	outString += fmt.Sprintln(" ")
	outString += fmt.Sprintln("Zone List")
	outString += fmt.Sprintln(" ")
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	table.SetHeader([]string{"ZONE", "ATTRIBUTE", "VALUE"})
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetBorder(false)

	if len(zones) == 0 {
		rowData := []string{"No zones found", " ", " "}
		table.Append(rowData)
	} else {
		for _, zone := range zones {
			zname := zone.Zone
			ztype := zone.Type
			table.Append([]string{zname, "Type", ztype})
			if len(zone.Comment) > 0 {
				table.Append([]string{" ", "Comment", zone.Comment})
			}
			if strings.ToUpper(ztype) == "SECONDARY" {
				if len(zone.Masters) > 0 {
					masters := strings.Join(zone.Masters, " ,")
					table.Append([]string{" ", "Masters", masters})
				}
				if zone.TSIGKey != nil {
					table.Append([]string{" ", "TsigKey:Name", zone.TSIGKey.Name})
					table.Append([]string{" ", "TsigKey:Algorithm", zone.TSIGKey.Algorithm})
					table.Append([]string{" ", "TsigKey:Secret", zone.TSIGKey.Secret})
				}
			}
			if strings.ToUpper(ztype) == "PRIMARY" || strings.ToUpper(ztype) == "SECONDARY" {
				table.Append([]string{" ", "SignAndServe", fmt.Sprintf("%t", zone.SignAndServe)})
				if len(zone.SignAndServeAlgorithm) > 0 {
					table.Append([]string{" ", "SignAndServeAlgorithm", fmt.Sprintf("%s", zone.SignAndServeAlgorithm)})
				}
			}
			if strings.ToUpper(ztype) == "ALIAS" {
				table.Append([]string{" ", "Target", zone.Target})
				table.Append([]string{" ", "AliasCount", strconv.FormatInt(zone.AliasCount, 10)})
			}
			table.Append([]string{" ", "ActivationState", zone.ActivationState})
			table.Append([]string{" ", "LastActivationDate", zone.LastActivationDate})
			table.Append([]string{" ", "LastModifiedDate", zone.LastModifiedDate})
			table.Append([]string{" ", "VersionId", zone.VersionID})
			table.Append([]string{" ", " ", " "})
		}
	}
	table.Render()
	outString += fmt.Sprintln(tableString.String())

	return outString
}

/*var b strings.Builder
b.WriteString("\nZone List\n\n")

t := tablewriter.NewWriter(&b)
t.SetHeader([]string{"ZONE", "ATTRIBUTE", "VALUE"})
t.SetAutoWrapText(false)
t.SetRowLine(true)
t.SetBorder(false)
t.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})

if len(zones) == 0 {
	t.Append([]string{"No zones found", " ", " "})
} else {
	for _, z := range zones {
		t.Append([]string{z.Zone, "Type", z.Type})
		if z.Comment != "" {
			t.Append([]string{" ", "Comment", z.Comment})
		}
		if strings.ToUpper(z.Type) == "SECONDARY" {
			if len(z.Masters) > 0 {
				t.Append([]string{" ", "Masters", strings.Join(z.Masters, ",")})
			}
			if z.TSIGKey != nil {
				t.Append([]string{" ", "TsigKey:Name", z.TSIGKey.Name})
				t.Append([]string{" ", "TsigKey:Algorithm", z.TSIGKey.Algorithm})
				t.Append([]string{" ", "TsigKey:Secret", z.TSIGKey.Secret})
			}
		}
		if strings.ToUpper(z.Type) == "PRIMARY" || strings.ToUpper(z.Type) == "SECONDARY" {
			t.Append([]string{" ", "SignAndServe", fmt.Sprintf("%t", z.SignAndServe)})
			if z.SignAndServeAlgorithm != "" {
				t.Append([]string{" ", "SignAndServeAlgorithm", z.SignAndServeAlgorithm})
			}
		}
		if strings.ToUpper(z.Type) == "ALIAS" {
			t.Append([]string{" ", "Target", z.Target})
			t.Append([]string{" ", "AliasCount", strconv.FormatInt(z.AliasCount, 10)})
		}

		t.Append([]string{" ", "ActivationState", z.ActivationState})
		t.Append([]string{" ", "LastActivationDate", z.LastActivationDate})
		t.Append([]string{" ", "LastModifiedDate", z.LastModifiedDate})
		t.Append([]string{" ", "VersionId", z.VersionID})
		t.Append([]string{" ", " ", " "})
	}
}

t.Render()
return b.String()*/

func renderZoneSummaryListTable(zones []dns.ZoneResponse) string {
	var b strings.Builder
	b.WriteString("\n Zone List Summary\n\n")

	t := tablewriter.NewWriter(&b)
	t.SetHeader([]string{"ZONE", "TYPE", "ACTIVATION STATE", "CONTRACT ID"})
	t.SetAutoWrapText(false)
	t.SetRowLine(true)
	t.SetBorder(false)
	t.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})

	if len(zones) == 0 {
		t.Append([]string{"No zones found", " ", " ", " "})
	} else {
		for _, z := range zones {
			t.Append([]string{z.Zone, z.Type, z.ActivationState, z.ContractID})
		}
	}
	t.Render()
	return b.String()

}

func renderBulkZonesRequestStatusTable(submitStatusList []*dns.BulkZonesResponse, c *cli.Context) string {

	//bold := color.New(color.FgWhite, color.Bold)
	outString := ""
	outString += fmt.Sprintln(" ")
	outString += fmt.Sprintln("Bulk Zones Request Submission Status")
	outString += fmt.Sprintln(" ")
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetBorder(false)

	for i, submitStatus := range submitStatusList {
		table.Append([]string{"Request Id", submitStatus.RequestID})
		table.Append([]string{"Expiration Date", submitStatus.ExpirationDate})
		if i == len(submitStatusList)-1 {
			table.Append([]string{"", ""})
		}
	}
	table.Render()
	outString += fmt.Sprintln(tableString.String())

	return outString
}

func renderBulkZonesStatusTable(submitStatusList []*dns.BulkStatusResponse, c *cli.Context) string {

	//bold := color.New(color.FgWhite, color.Bold)
	outString := ""
	outString += fmt.Sprintln(" ")
	outString += fmt.Sprintln("Bulk Zones Request Status")
	outString += fmt.Sprintln(" ")
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetBorder(false)

	for _, submitStatus := range submitStatusList {
		table.Append([]string{"Request Id", submitStatus.RequestID, ""})
		table.Append([]string{"", "Zones Submitted", strconv.Itoa(submitStatus.ZonesSubmitted)})
		table.Append([]string{"", "Success Count", strconv.Itoa(submitStatus.SuccessCount)})
		table.Append([]string{"", "Failure Count", strconv.Itoa(submitStatus.FailureCount)})
		table.Append([]string{"", "Complete", fmt.Sprintf("%t", submitStatus.IsComplete)})
		table.Append([]string{"", "Expiration Date", submitStatus.ExpirationDate})
	}
	table.Render()
	outString += fmt.Sprintln(tableString.String())

	return outString
}

func renderBulkZonesResultTable(resultRespList interface{}, c *cli.Context) string {

	//bold := color.New(color.FgWhite, color.Bold)
	var requestid string
	var succzones []string
	var failzones []dns.BulkFailedZone
	op := "Created"
	tableHeader := "Bulk Zones %s Request Results"

	outString := ""
	outString += fmt.Sprintln(" ")
	outString += fmt.Sprintln(fmt.Sprintf(tableHeader, op))
	outString += fmt.Sprintln(" ")
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetBorder(false)

	if resultList, ok := resultRespList.([]*dns.GetBulkZoneCreateResultResponse); ok {
		for _, crreq := range resultList {
			requestid = crreq.RequestID
			succzones = crreq.SuccessfullyCreatedZones
			failzones = crreq.FailedZones
			table.Append([]string{"Request Id", requestid, "", ""})
			table.Append([]string{"", fmt.Sprintf("Successfully %s Zones", op), "", ""})
			if len(succzones) == 0 {
				table.Append([]string{"", "", "None", ""})
			} else {
				for _, zn := range succzones {
					table.Append([]string{"", "", zn, ""})
				}
			}
			table.Append([]string{"", fmt.Sprintf("Failed %s Zones", op), "", ""})
			if len(failzones) == 0 {
				table.Append([]string{"", "", "None", ""})
			} else {
				for _, fzn := range failzones {
					table.Append([]string{"", "", fzn.Zone, fzn.FailureReason})
				}
			}
		}
		table.Render()
		outString += fmt.Sprintln(tableString.String())

		return outString
	}
	resultList, ok := resultRespList.([]*dns.GetBulkZoneDeleteResultResponse)
	if !ok {
		return "Unable to create result table"
	}
	for _, delreq := range resultList {
		requestid = delreq.RequestID
		succzones = delreq.SuccessfullyDeletedZones
		failzones = delreq.FailedZones
		op = "Deleted"
		table.Append([]string{"Request Id", requestid, "", ""})
		table.Append([]string{fmt.Sprintf("", "Successfully %s Zones", op), "", ""})
		if len(succzones) == 0 {
			table.Append([]string{"", "", "None", ""})
		} else {
			for _, zn := range succzones {
				table.Append([]string{"", "", zn, ""})
			}
		}
		table.Append([]string{fmt.Sprintf("", "Failed %s Zones", op), "", ""})
		if len(succzones) == 0 {
			table.Append([]string{"", "", "None", ""})
		} else {
			for _, fzn := range failzones {
				table.Append([]string{"", "", fzn.Zone, fzn.FailureReason})
			}
		}
	}

	table.Render()
	outString += fmt.Sprintln(tableString.String())

	return outString
}
