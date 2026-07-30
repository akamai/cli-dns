package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/urfave/cli"
)

func renderRecordsetTable(zone string, record *dns.GetRecordResponse) string {
	var sb strings.Builder
	table := tablewriter.NewTable(&sb)
	table.Header([]string{"ZONE", "NAME", "TYPE", "TTL", "RDATA"})
	rdata := strings.Join(record.Target, ", ")
	row := []string{
		zone,
		record.Name,
		record.RecordType,
		strconv.Itoa(record.TTL),
		rdata,
	}
	table.Bulk([][]string{row})
	table.Render()
	return sb.String()
}

// Recordsets list table format
func renderRecordsetListTable(recordsets []dns.RecordSet) string {
	var sb strings.Builder
	sb.WriteString("Zone Recordsets:\n")
	table := tablewriter.NewTable(&sb)
	table.Header([]string{"NAME", "TYPE", "TTL", "RDATA"})

	var rows [][]string
	if len(recordsets) == 0 {
		rows = append(rows, []string{"No recordsets found", " ", " ", " "})
	} else {
		for _, set := range recordsets {
			name := set.Name
			typeVal := set.Type
			ttl := strconv.Itoa(set.TTL)
			for i, rdata := range set.Rdata {
				if i == 0 {
					rows = append(rows, []string{name, typeVal, ttl, rdata})
				} else {
					rows = append(rows, []string{" ", " ", " ", rdata})
				}
			}
		}
	}
	table.Bulk(rows)
	table.Render()
	return sb.String()
}

// Zone table format
func renderZoneconfigTable(zone *dns.GetZoneResponse) string {
	var sb strings.Builder
	sb.WriteString("Zone Configuration:\n")

	table := tablewriter.NewTable(&sb, tablewriter.WithConfig(tablewriter.Config{
		Row: tw.CellConfig{
			Alignment: tw.CellAlignment{Global: tw.AlignLeft},
		},
		Header: tw.CellConfig{
			Alignment: tw.CellAlignment{Global: tw.AlignLeft},
		},
	}))

	table.Header([]string{"ZONE", "ATTRIBUTE", "VALUE"})

	var rows [][]string
	if zone == nil {
		rows = append(rows, []string{"No zone info to display", " ", " "})
	} else {
		zname := zone.Zone
		ztype := zone.Type
		rows = append(rows, []string{zname, "Type", ztype})
		if len(zone.Comment) > 0 {
			rows = append(rows, []string{" ", "Comment", zone.Comment})
		}
		if len(zone.ContractID) > 0 {
			rows = append(rows, []string{" ", "ContractId", zone.ContractID})
		}
		if strings.ToUpper(ztype) == "SECONDARY" {
			if len(zone.Masters) > 0 {
				masters := strings.Join(zone.Masters, " ,")
				rows = append(rows, []string{" ", "Masters", masters})
			}
			if zone.TSIGKey != nil {
				if len(zone.TSIGKey.Name) > 0 {
					rows = append(rows, []string{" ", "TsigKey:Name", zone.TSIGKey.Name})
				}
				if len(zone.TSIGKey.Algorithm) > 0 {
					rows = append(rows, []string{" ", "TsigKey:Algorithm", zone.TSIGKey.Algorithm})
				}
				if len(zone.TSIGKey.Secret) > 0 {
					rows = append(rows, []string{" ", "TsigKey:Secret", zone.TSIGKey.Secret})
				}
			}
		}
		if strings.ToUpper(ztype) == "PRIMARY" || strings.ToUpper(ztype) == "SECONDARY" {
			rows = append(rows, []string{" ", "SignAndServe", fmt.Sprintf("%t", zone.SignAndServe)})
			if len(zone.SignAndServeAlgorithm) > 0 {
				rows = append(rows, []string{" ", "SignAndServeAlgorithm", zone.SignAndServeAlgorithm})
			}
		}
		if strings.ToUpper(ztype) == "ALIAS" {
			rows = append(rows, []string{" ", "Target", zone.Target})
			rows = append(rows, []string{" ", "AliasCount", strconv.FormatInt(zone.AliasCount, 10)})
		}
		rows = append(rows, []string{" ", "ActivationState", zone.ActivationState})
		if len(zone.LastActivationDate) > 0 {
			rows = append(rows, []string{" ", "LastActivationDate", zone.LastActivationDate})
		}
		if len(zone.LastModifiedDate) > 0 {
			rows = append(rows, []string{" ", "LastModifiedDate", zone.LastModifiedDate})
		}
		if len(zone.LastModifiedBy) > 0 {
			rows = append(rows, []string{" ", "LastModifiedBy", zone.LastModifiedBy})
		}
		rows = append(rows, []string{" ", "VersionId", zone.VersionID})
	}

	table.Bulk(rows)
	table.Render()
	return sb.String()
}

// Zone list table format
func renderZoneListTable(zones []dns.ZoneResponse) string {
	var sb strings.Builder
	sb.WriteString("Zone List:\n")

	table := tablewriter.NewTable(&sb)
	table.Header([]string{"ZONE", "ATTRIBUTE", "VALUE"})

	var rows [][]string
	if len(zones) == 0 {
		rows = append(rows, []string{"No zones found", " ", " "})
	} else {
		for _, zone := range zones {
			zname := zone.Zone
			ztype := zone.Type
			rows = append(rows, []string{zname, "Type", ztype})
			if len(zone.Comment) > 0 {
				rows = append(rows, []string{" ", "Comment", zone.Comment})
			}
			if strings.ToUpper(ztype) == "SECONDARY" {
				if len(zone.Masters) > 0 {
					masters := strings.Join(zone.Masters, " ,")
					rows = append(rows, []string{" ", "Masters", masters})
				}
				if zone.TSIGKey != nil {
					rows = append(rows, []string{" ", "TsigKey:Name", zone.TSIGKey.Name})
					rows = append(rows, []string{" ", "TsigKey:Algorithm", zone.TSIGKey.Algorithm})
					rows = append(rows, []string{" ", "TsigKey:Secret", zone.TSIGKey.Secret})
				}
			}
			if strings.ToUpper(ztype) == "PRIMARY" || strings.ToUpper(ztype) == "SECONDARY" {
				rows = append(rows, []string{" ", "SignAndServe", fmt.Sprintf("%t", zone.SignAndServe)})
				if len(zone.SignAndServeAlgorithm) > 0 {
					rows = append(rows, []string{" ", "SignAndServeAlgorithm", zone.SignAndServeAlgorithm})
				}
			}
			if strings.ToUpper(ztype) == "ALIAS" {
				rows = append(rows, []string{" ", "Target", zone.Target})
				rows = append(rows, []string{" ", "AliasCount", strconv.FormatInt(zone.AliasCount, 10)})
			}
			rows = append(rows, []string{" ", "ActivationState", zone.ActivationState})
			rows = append(rows, []string{" ", "LastActivationDate", zone.LastActivationDate})
			rows = append(rows, []string{" ", "LastModifiedDate", zone.LastModifiedDate})
			rows = append(rows, []string{" ", "LastModifiedBy", zone.LastModifiedBy})
			rows = append(rows, []string{" ", "VersionId", zone.VersionID})
			rows = append(rows, []string{" ", " ", " "})
		}
	}
	table.Bulk(rows)
	table.Render()
	return sb.String()
}

// Zone list summary format
func renderZoneSummaryListTable(zones []dns.ZoneResponse) string {
	var sb strings.Builder
	sb.WriteString("Zone List Summary:\n")

	table := tablewriter.NewTable(&sb)
	table.Header([]string{"ZONE", "TYPE", "ACTIVATION STATE", "CONTRACT ID"})

	var rows [][]string
	if len(zones) == 0 {
		rows = append(rows, []string{"No zones found", " ", " ", " "})
	} else {
		for _, z := range zones {
			rows = append(rows, []string{z.Zone, z.Type, z.ActivationState, z.ContractID})
		}
	}
	table.Bulk(rows)
	table.Render()
	return sb.String()
}

// Zone table format
func renderZoneTable(zone *dns.GetZoneResponse, records []dns.RecordSet, c *cli.Context) {
	var sb strings.Builder
	table := tablewriter.NewTable(&sb)
	table.Header([]string{"Field", "value"})

	var rows [][]string
	rows = append(rows, []string{"Zone", zone.Zone})
	rows = append(rows, []string{"Type", zone.Type})
	rows = append(rows, []string{"Masters", strings.Join(zone.Masters, ", ")})
	rows = append(rows, []string{"Comment", zone.Comment})
	rows = append(rows, []string{"Contract ID", zone.ContractID})
	rows = append(rows, []string{"SignAndServe", fmt.Sprintf("%v", zone.SignAndServe)})
	rows = append(rows, []string{"Target", zone.Target})
	rows = append(rows, []string{"EndCustomerID", zone.EndCustomerID})
	rows = append(rows, []string{"Activation State", zone.ActivationState})
	rows = append(rows, []string{"Last Modified By", zone.LastModifiedBy})
	rows = append(rows, []string{"Last Modified Date", zone.LastModifiedDate})
	rows = append(rows, []string{"Version ID", zone.VersionID})

	if zone.TSIGKey != nil {
		rows = append(rows, []string{"TSIG Name", zone.TSIGKey.Name})
		rows = append(rows, []string{"TSIG Algorithm", zone.TSIGKey.Algorithm})
		rows = append(rows, []string{"TSIG Secret", zone.TSIGKey.Secret})
	}
	table.Bulk(rows)
	table.Render()
	_, _ = fmt.Fprintln(c.App.Writer, sb.String())

	if len(records) > 0 {
		_, _ = fmt.Fprintln(c.App.Writer, "DNS Records: ")

		var recSb strings.Builder
		recordsTable := tablewriter.NewTable(&recSb)
		recordsTable.Header([]string{"Name", "Type", "TTL", "Data"})

		var recRows [][]string
		for _, rec := range records {
			for _, data := range rec.Rdata {
				recRows = append(recRows, []string{
					rec.Name, rec.Type, fmt.Sprintf("%d", rec.TTL), data,
				})
			}
		}
		recordsTable.Bulk(recRows)
		recordsTable.Render()
		_, _ = fmt.Fprintln(c.App.Writer, recSb.String())
	}
}

// Bulk zone request status format
func renderBulkZonesRequestStatusTable(submitStatusList []*dns.BulkZonesResponse) string {
	var sb strings.Builder
	sb.WriteString("Bulk Zones Request Submission Status:\n")

	table := tablewriter.NewTable(&sb)
	var rows [][]string

	for i, submitStatus := range submitStatusList {
		rows = append(rows, []string{"Request Id", submitStatus.RequestID})
		rows = append(rows, []string{"Expiration Date", submitStatus.ExpirationDate})
		if i == len(submitStatusList)-1 {
			rows = append(rows, []string{"", ""})
		}
	}
	table.Bulk(rows)
	table.Render()
	return sb.String()
}

// Bulk zone status format
func renderBulkZonesStatusTable(submitStatusList []*dns.BulkStatusResponse) string {
	var sb strings.Builder
	sb.WriteString("Bulk Zones Request Status:\n")

	table := tablewriter.NewTable(&sb)
	var rows [][]string

	for _, submitStatus := range submitStatusList {
		rows = append(rows, []string{"Request Id", submitStatus.RequestID})
		rows = append(rows, []string{"Zones Submitted", strconv.Itoa(submitStatus.ZonesSubmitted)})
		rows = append(rows, []string{"Success Count", strconv.Itoa(submitStatus.SuccessCount)})
		rows = append(rows, []string{"Failure Count", strconv.Itoa(submitStatus.FailureCount)})
		rows = append(rows, []string{"Complete", fmt.Sprintf("%t", submitStatus.IsComplete)})
		rows = append(rows, []string{"Expiration Date", submitStatus.ExpirationDate})
	}
	table.Bulk(rows)
	table.Render()
	return sb.String()
}

// Bulk zone result format
func renderBulkZonesResultTable(resultRespList interface{}) string {
	var sb strings.Builder
	var requestid string
	var succzones []string
	var failzones []dns.BulkFailedZone
	var rows [][]string

	if resultList, ok := resultRespList.([]*dns.GetBulkZoneCreateResultResponse); ok {
		op := "Created"
		sb.WriteString(fmt.Sprintf("Bulk Zones %s Request Results\n", op))
		table := tablewriter.NewTable(&sb)

		for _, crreq := range resultList {
			requestid = crreq.RequestID
			succzones = crreq.SuccessfullyCreatedZones
			failzones = crreq.FailedZones
			rows = append(rows, []string{"Request Id", requestid})
			if len(succzones) == 0 {
				rows = append(rows, []string{fmt.Sprintf("Successfully %s Zones", op), "None"})
			} else {
				for i, zn := range succzones {
					if i == 0 {
						rows = append(rows, []string{fmt.Sprintf("Successfully %s Zones", op), zn})
					} else {
						rows = append(rows, []string{"", zn})
					}
				}
			}
			if len(failzones) == 0 {
				rows = append(rows, []string{fmt.Sprintf("Failed %s Zones", op), "None"})
			} else {
				for i, fzn := range failzones {
					if i == 0 {
						rows = append(rows, []string{fmt.Sprintf("Failed %s Zones", op), fzn.Zone + ": " + fzn.FailureReason})
					} else {
						rows = append(rows, []string{"", fzn.Zone + ": " + fzn.FailureReason})
					}
				}
			}
		}
		table.Bulk(rows)
		table.Render()
		return sb.String()
	}

	resultList, ok := resultRespList.([]*dns.GetBulkZoneDeleteResultResponse)
	if !ok {
		return "Unable to create result table"
	}

	op := "Deleted"
	sb.WriteString(fmt.Sprintf("Bulk Zones %s Request Results\n", op))
	table := tablewriter.NewTable(&sb)

	for _, delreq := range resultList {
		requestid = delreq.RequestID
		succzones = delreq.SuccessfullyDeletedZones
		failzones = delreq.FailedZones

		rows = append(rows, []string{"Request Id", requestid})
		if len(succzones) == 0 {
			rows = append(rows, []string{fmt.Sprintf("Successfully %s Zones", op), "None"})
		} else {
			for i, zn := range succzones {
				if i == 0 {
					rows = append(rows, []string{fmt.Sprintf("Successfully %s Zones", op), zn})
				} else {
					rows = append(rows, []string{"", zn})
				}
			}
		}
		if len(failzones) == 0 {
			rows = append(rows, []string{fmt.Sprintf("Failed %s Zones", op), "None"})
		} else {
			for i, fzn := range failzones {
				if i == 0 {
					rows = append(rows, []string{fmt.Sprintf("Failed %s Zones", op), fzn.Zone + ": " + fzn.FailureReason})
				} else {
					rows = append(rows, []string{"", fzn.Zone + ": " + fzn.FailureReason})
				}
			}
		}
	}

	table.Bulk(rows)
	table.Render()
	return sb.String()
}
