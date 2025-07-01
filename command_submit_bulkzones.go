// Copyright 2020. Akamai Technologies, Inc
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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/akamai/cli-dns/edgegrid"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

const (
	MaxUint     = ^uint(0)
	MaxInt      = int(MaxUint >> 1)
	httpMaxBody = MaxInt
)

type BulkZonesResponse struct {
	RequestID      string
	ExpirationDate string
}

func cmdSubmitBulkZones(c *cli.Context) error {

	// Initialize context and Edgegrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return fmt.Errorf("session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	var (
		outputPath     string
		contractid     string
		groupid        string
		inputPath      string
		bulkDeleteList *dns.ZoneNameListResponse
		newBulkZones   *dns.BulkZonesCreate
		op             string = "create"
		bypass         bool
		maxNumZones    int = 1000
	)

	fmt.Println("Preparing bulk zones submit request ", "")
	queryArgs := dns.ZoneQueryString{}

	if c.IsSet("contractid") {
		contractid = c.String("contractid")
		queryArgs.Contract = contractid
	} else if c.IsSet("create") {
		return cli.NewExitError(color.RedString("contractid is required"), 1)
	}
	if c.IsSet("groupid") {
		groupid = c.String("groupid")
		if groupid != "" {
			queryArgs.Group = groupid
			fmt.Println("Using groupid:", groupid)
		} else {
			fmt.Println("groupid flag set but empty; ignoring")
		}
	} else {
		fmt.Println("groupid flag not set; proceeding without groupid")
	}
	if c.IsSet("bypassZoneSafety") && c.Bool("bypassZoneSafety") {
		bypass = true
	}

	// Validate that only one operation is selected (create or delete)
	if (c.IsSet("create") && c.IsSet("delete")) || (!c.IsSet("create") && !c.IsSet("delete")) {
		return cli.NewExitError(color.RedString("Either create or delete arg is required. "), 1)
	}

	// Creating object based on operation type
	if c.IsSet("delete") {
		op = "delete"
		bulkDeleteList = &dns.ZoneNameListResponse{}
	} else {
		newBulkZones = &dns.BulkZonesCreate{}
	}
	if op == "create" {
		if bypass {
			fmt.Println("Warning: bypassZoneSafety arg ignored")
		}
	} else {
		if c.IsSet("contractid") {
			fmt.Println("Warning: contractid arg ignored")
		}
		if c.IsSet("groupid") {
			fmt.Println("Warning: groupid arg ignored")
		}
	}
	if c.IsSet("file") {
		inputPath = c.String("file")
		inputPath = filepath.FromSlash(inputPath)
	} else {
		return cli.NewExitError(color.RedString(" Bulk create JSON source file must be specified"), 1)
	}
	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}

	val, ok := os.LookupEnv("AKAMAI_ZONES_BATCH_SIZE")
	if ok {
		batchsize, err := strconv.Atoi(val)
		if err != nil {
			return cli.NewExitError(color.RedString(" Environ variable AKAMAI_ZONEBATCH has invalid value"), 1)
		}
		maxNumZones = batchsize
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		return cli.NewExitError(color.RedString("Failed to read input file"), 1)
	}

	if op == "create" {
		err = json.Unmarshal(data, newBulkZones)
	} else {
		err = json.Unmarshal(data, bulkDeleteList)
	}
	if err != nil {
		return cli.NewExitError(color.RedString("Failed to parse json file content into bulk zones object"), 1)
	}

	/*var (
		submitStatusList []*dns.CreateBulkZonesResponse
		deleteResp       *dns.DeleteBulkZonesResponse
	)*/

	submitStatusList := make([]*dns.BulkZonesResponse, 0)

	fmt.Println("Submitting Bulk Zones request  ", "")

	// Handling bulk create in batches
	if op == "create" {
		ZonesMax := len(newBulkZones.Zones)
		numZones := ZonesMax
		if ZonesMax > maxNumZones {
			numZones = maxNumZones
		}
		bulkZonesList := make([]*dns.BulkZonesCreate, 0)
		bulkZones := &dns.BulkZonesCreate{Zones: make([]dns.ZoneCreate, 0)}

		for _, zone := range newBulkZones.Zones {
			bulkZones.Zones = append(bulkZones.Zones, zone)
			numZones--
			ZonesMax--
			if numZones == 0 {
				bulkZonesList = append(bulkZonesList, bulkZones)
				bulkZones = &dns.BulkZonesCreate{Zones: make([]dns.ZoneCreate, 0)}
				numZones = ZonesMax
				if ZonesMax > maxNumZones {
					numZones = maxNumZones
				}
			}

		}
		if len(bulkZones.Zones) > 0 {
			bulkZonesList = append(bulkZonesList, bulkZones)
		}

		// Submit each batch to API
		for _, zonesRequest := range bulkZonesList {
			req := dns.CreateBulkZonesRequest{
				BulkZones:       zonesRequest,
				ZoneQueryString: queryArgs,
			}
			resp, err := dnsClient.CreateBulkZones(ctx, req)
			if err != nil {
				return cli.NewExitError(color.RedString("bulk zone submit request failed: %s", err), 1)
			}
			submitStatusList = append(submitStatusList, &dns.BulkZonesResponse{
				RequestID:      resp.RequestID,
				ExpirationDate: resp.ExpirationDate,
			})
		}
	} else {
		// Handle bulk zone DELETE request
		resp, err := dnsClient.DeleteBulkZones(ctx, dns.DeleteBulkZonesRequest{
			ZonesList:          bulkDeleteList,
			BypassSafetyChecks: &bypass,
		})
		if err != nil {
			return cli.NewExitError(color.RedString(fmt.Sprintf("Bulk Zone Request submit failed. Error: %s", err.Error())), 1)
		}
		submitStatusList = append(submitStatusList, &dns.BulkZonesResponse{
			RequestID:      resp.RequestID,
			ExpirationDate: resp.ExpirationDate,
		})
	}

	// Format results in JSON or table format
	results := ""
	if c.IsSet("json") && c.Bool("json") {
		jsonBytes, err := json.MarshalIndent(submitStatusList, "", " ")
		if err != nil {
			return cli.NewExitError(color.RedString("unable to marshal"), 1)
		}
		results = string(jsonBytes)
	} else {
		results = renderBulkZonesRequestStatusTable(submitStatusList, c)
	}

	// Write output either to output file or console
	if len(outputPath) < 1 {
		outputPath = fmt.Sprintf("Bulk_Submit_Request_Status_%d.json", time.Now().Unix())
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return cli.NewExitError(color.RedString("failed to create output file: %s", err), 1)
	}
	defer file.Close()

	if _, err := file.WriteString(results); err != nil {
		return cli.NewExitError(color.RedString("failed to write output"), 1)
	}

	if c.IsSet("suppress") && c.Bool("suppress") {
		return nil
	}

	fmt.Fprintln(c.App.Writer, "")
	fmt.Fprintln(c.App.Writer, results)
	return nil
}
