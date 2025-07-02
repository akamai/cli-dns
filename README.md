<!--esi 
<esi:assign name="dac_stylesheets" value="['/stylesheets/screen.css']" />
<esi:assign name="dac_footer_js" value="['/javascripts/app/toc.js']" />
<esi:eval src="/internal/header.html?http" dca="esi" />
-->
# Akamai CLI for Edge DNS

[![Go Report Card](https://goreportcard.com/badge/github.com/akamai/cli-dns)](https://goreportcard.com/report/github.com/akamai/cli-dns) [![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fakamai%2Fcli-dns.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fakamai%2Fcli-dns?ref=badge_shield)

An [Akamai CLI](https://developer.akamai.com/cli) package for managing DNS Zones using Edge DNS (formerly known as Fast DNS).

## Getting Started

### Installing

To install this package, use Akamai CLI:

```sh
$ akamai install dns
```

You may also use this as a stand-alone command by downloading the
[latest release binary](https://github.com/akamai/cli-dns/releases)
for your system, or by cloning this repository and compiling it yourself.

### Compiling from Source

If you want to compile the package from source, you will need Go 1.23 or later installed:

1. Create a clone of the target repository:
   `git clone https://github.com/akamai/cli-dns.git`
2. Change to the package directory and compile the binary:
  - Linux/macOS/*nix: `go build -o akamai-dns`
  - Windows: `go build -o akamai-dns.exe`

## Command Summary

### Usage

```
$  akamai dns [--edgerc] [--section] [--accountkey] <command> [sub-command]
```

or 

```
$  akamai-dns [--edgerc] [--section] [--accountkey] <command> [sub-command]
```

### Description

   Manage DNS Zones with Edge DNS



### Global Flags

```
   --edgerc value      Location of the credentials file (default: "/home/elynes/.edgerc") [$AKAMAI_EDGERC]
   --section value     Section of the credentials file (default: "dns") [$AKAMAI_EDGERC_SECTION]
   --accountkey value  Account switch key [$AKAMAI_EDGERC_ACCOUNT_KEY]
```

## Built-In Commands

```
  add-record [Deprecated]
  rm-record [Deprecated]
  list-recordsets
  create-recordsets
  update-recordsets
  retrieve-recordset
  create-recordset
  update-recordset
  delete-recordset
  retrieve-zone [Deprecated]
  update-zone [Deprecated]
  list-zoneconfig
  create-zoneconfig
  retrieve-zoneconfig
  update-zoneconfig
  submit-bulkzones
  status-bulkzones
  result-bulkzones
  list
  help
```

Commands are grouped into several categories.

*-zoneconfig, *-recordsets and *-recordset commands provide the ability to manage zone configurations directly, as well as manage recordsets individually or in groupings. These commands should be preferred.

*-bulkzones commands provide the ability to submit and monitor bulk zone operations.

*-zone and *-record commands provide a more constrained scope of control, treating the zone and records as a single entity. 

### Listing Zone Configurations

A list of existing zone configurations can be retrieved by using the `akamai dns list-zoneconfig` command.

The complete command line is:

```
   akamai dns list-zoneconfig  [--json] [--output] [--contractid] [--type] [--search] [--summary] 

Flags: 
   --json           Output as JSON [$AKAMAI_CLI_DNS_JSON]
   --output FILE    Output command results to FILE
   --contractid ID  Contract ID. Multiple flags allowed
   --type TYPE      Zone TYPE. Multiple flags allowed
   --search VALUE   Zone search VALUE
   --summary        List zone names and type
```

To list Primary zones and generate results in json format, the `--type` and `--json` flags would be used. For example:

```
$akamai dns list-zoneconfig --type primary --json
```

would result in the following output:

```
{
  "Zones": [
    {
      "zone": "example.com",
      "type": "PRIMARY",
      "signAndServe": false,
      "contractId": "1-ABC123",
      "activationState": "NEW",
      "lastModifiedBy": "jsmith",
      "lastModifiedDate": "2020-06-05T21:05:04.298125Z",
      "versionId": "60a1f29b-85e8-44e0-a921-2bcae8728f75"
    }
  ]
}
```

To generate a summary list of zones containing specific text, the `--search` and `--summary` flags would be used. For example:

```
$ akamai dns list-zoneconfig --search example --summary
```

would result in the following output:

```
Zone List Summary
 
                   ZONE                     TYPE     ACTIVATION STATE   CONTRACT ID  
                                                                                     
  example.com                              PRIMARY         NEW           1-3CV382    
``` 

### Retrieving a Zone Configuration

An existing zone configuration can be retrieved by using the `akamai dns retrieve-zoneconfig` command.

The complete command line is:

```
   akamai dns retrieve-zoneconfig <zonename> [--json] [--output] [--dns]

Flags: 
   --json         Output as JSON [$AKAMAI_CLI_DNS_JSON]
   --output FILE  Output command results to FILE
   --dns          Retrieve Zone Master File
```

To retrieve a zoneconfig and output result in json format, the `--json` flag would be used. For example:

```
$ akamai dns retrieve-zoneconfig example.com --json
```

would result in the following output:

```
{
  "zone": "example.com",
  "type": "PRIMARY",
  "comment": "Primary zone",
  "signAndServe": false,
  "contractId": "1-ABC123",
  "activationState": "PENDING",
  "lastModifiedBy": "jsmith",
  "lastModifiedDate": "2020-06-09T18:06:01.266155Z",
  "versionId": "a0b4730e-fbbe-40ad-96b3-ac6a4cbadb1e"
} 
```

To retrieve the master zone file and output to the console, an example would be:

```
$ akamai dns retrieve-zoneconfig xxx_primary_test.com --dns
Retrieving Zone
Retrieving Zone  ... [OK]

Assembling Zone Content
Assembling Zone Content ... [OK]


;; File Generated at 2020-10-21T21:22:58.672933Z
;; Last Modified at 2020-10-19T14:54:56.407659Z[UTC]
;; Version Identifier 166415f9-0744-4acc-9ca1-84173df6f69d
xxx_primary_test.com.        86400   IN      NS      a4-65.akam.net.
xxx_primary_test.com.        86400   IN      NS      a13-65.akam.net.
xxx_primary_test.com.        86400   IN      NS      a3-64.akam.net.
xxx_primary_test.com.        86400   IN      NS      a2-64.akam.net.
xxx_primary_test.com.        86400   IN      NS      a12-65.akam.net.
xxx_primary_test.com.        86400   IN      NS      a1-98.akam.net.
xxx_primary_test.com.        86400   IN      SOA     a1-98.akam.net. hostmaster.xxx_primary_test.com. 2020052603 3600 600 604800 300
arec.xxx_primary_test.com.   1800    IN      A       1.2.4.8
```

### Creating a Zone Configuration

A zone configuration can be created by using the `akamai dns create-zoneconfig` command. The configuration can be provided as command line values or json file.

The complete command line is:

```
akamai dns create-zoneconfig <zonename> [--json] [--suppress] [--output] [--type] [--master] [--comment] [--signandserve] [--algorithm] [--tsigname] [--tsigalgorithm] [--tsigsecret] [--target] [--endcustomerid] [--file] [--contractid] [--groupid] [--initialize] 

Flags: 
   --json                         Output as JSON [$AKAMAI_CLI_DNS_JSON]
   --suppress                     Suppress command result output. Overrides other output related flags [$AKAMAI_CLI_DNS_SUPPRESS]
   --output FILE                  Output command results to FILE
   --type TYPE                    Zone TYPE
   --master MASTER                Secondary Zone MASTER. Multiple flags may be specified
   --comment COMMENT              Zone COMMENT
   --signandserve SIGNANDSERVE    Primary or Secondary Zone SIGNANDSERVE flag
   --algorithm ALGORITHM          Zone signandserve ALGORITHM
   --tsigname NAME                TSIG key NAME
   --tsigalgorithm ALGORITHM      TSIG key ALGORITHM
   --tsigsecret SECRET            TSIG key SECRET
   --target TARGET                Alias Zone TARGET
   --endcustomerid ENDCUSTOMERID  ENDCUSTOMERID
   --file FILE                    Read JSON formatted input from FILE
   --contractid ID                Contract ID
   --groupid ID                   Group ID
   --initialize                   Generate default SOA and NS Records
```

To create a zone, the desired fields and values would be provided. For example, to create a simple primary zone with default SOA and NS records via command line:

```
$ akamai dns create-zoneconfig example_primary.com --type primary --contractid 1-ABC123 --initialize
```

would create the zone with the following output:

```
Zone Configuration
 
           ZONE               ATTRIBUTE                      VALUE                  
                                                                                    
  example_primary.com      Type               PRIMARY                               
                           
                           ContractId         1-ABC123
                                                                                    
                           SignAndServe       false                                 
                                                                                    
                           ActivationState    PENDING                               
                                                                                    
                           LastModifiedDate   2020-06-09T18:06:01.266155Z           
                                                                                    
                           VersionId          a0b4730e-fbbe-40ad-96b3-ac6a4cbadb1e  
```

To create a secondary zone with a Tsig Key and comment and but suppress output via file input, the following command would be specified:

```
$ akamai dns create-zoneconfig example_secondary.com --file zone_create.json --suppress
```

where zone_create.json contains:

```
{
  "zone": "example_secondary.com",
  "type": "SECONDARY",
  "comment": "secondary zone",
  "masters": [
    "10.0.1.1"
  ],
  "signAndServe": false,
  "tsigKey": {
    "name": "testtsig",
    "algorithm": "hmac-md5",
    "secret": "p/jzrJpXOLf4mPUtx/z+Sw=="
  },
  "contractId": "1-ABC123" 
}
```

returns no output

### Updating a Zone Configuration

A zone configuration can be updated by using the `akamai dns update-zoneconfig` command. The updated configuration can be provided as command line values or json file. Note: Only updated fields and values need to be specified if using command line input.

The complete command line is:

```
akamai dns update-zoneconfig <zonename> [--json] [--suppress] [--output] [--type] [--master] [--comment] [--signandserve] [--algorithm] [--tsigname] [--tsigalgorithm] [--tsigsecret] [--target] [--endcustomerid] [--file] [--contractid] [--dns]

Flags:
   --json                         Output as JSON [$AKAMAI_CLI_DNS_JSON]
   --suppress                     Suppress command result output. Overrides other output related flags [$AKAMAI_CLI_DNS_SUPPRESS]
   --output FILE                  Output command results to FILE
   --type TYPE                    Zone TYPE
   --master MASTER                Secondary Zone MASTER. Multiple flags may be specified
   --comment COMMENT              Zone COMMENT
   --signandserve SIGNANDSERVE    Primary or Secondary Zone SIGNANDSERVE flag
   --algorithm ALGORITHM          Zone signandserve ALGORITHM
   --tsigname NAME                TSIG key NAME
   --tsigalgorithm ALGORITHM      TSIG key ALGORITHM
   --tsigsecret SECRET            TSIG key SECRET
   --target TARGET                Alias Zone TARGET
   --endcustomerid ENDCUSTOMERID  ENDCUSTOMERID
   --contractid ID                Contract ID
   --file FILE                    Read JSON formatted input from FILE
   --dns                          Input is Zone Master File
```

For example, to  update the previously created primary zone and add a comment via the command line:

```
$ akamai dns update-zoneconfig example_primary.com --type primary --comment "This is a comment"
```

would update the zone and provide the following output:

```
Zone Configuration

           ZONE               ATTRIBUTE                      VALUE

  example_primary.com      Type               PRIMARY

                           Comment            This is a comment

                           ContractId         1-ABC123

                           SignAndServe       false

                           ActivationState    PENDING

                           LastModifiedDate   2020-06-09T18:06:01.266155Z

                           VersionId          a0b4730e-fbbe-40ad-96b3-ac6a4cbadb1e
```

To update the master zone file (previously retrieved), an example would be:

```
$ akamai dns update-zoneconfig xxx_primary_test.com --dns --file ./master_file
Preparing zone for update
Updating Master Zone File
Updating Master Zone File ... [OK]
```

where ./master_file contains the following:

```
;; File Generated at 2020-10-21T21:22:58.672933Z
;; Last Modified at 2020-10-19T14:54:56.407659Z[UTC]
;; Version Identifier 166415f9-0744-4acc-9ca1-84173df6f69d
xxx_primary_test.com.        86400   IN      NS      a4-65.akam.net.
xxx_primary_test.com.        86400   IN      NS      a13-65.akam.net.
xxx_primary_test.com.        86400   IN      NS      a3-64.akam.net.
xxx_primary_test.com.        86400   IN      NS      a2-64.akam.net.
xxx_primary_test.com.        86400   IN      NS      a12-65.akam.net.
xxx_primary_test.com.        86400   IN      NS      a1-98.akam.net.
xxx_primary_test.com.        86400   IN      SOA     a1-98.akam.net. hostmaster.xxx_primary_test.com. 2020052604 3600 600 604800 300
arec.xxx_primary_test.com.   1800    IN      A       1.2.3.4
```

### Listing a Zone's Recordsets

The recordsets in a zone can be [selectively] listed by using the `akamai dns list-recordsets` command. 

The complete command line is:

```
$ akamai dns list-recordsets <zonename> [--json] [--output] [--type] [--sortby] [--search] 

Flags: 
   --json           Output as JSON [$AKAMAI_CLI_DNS_JSON]
   --output FILE    Output command results to FILE
   --type TYPE      List recordset(s) matching TYPE. Multiple flags allowed
   --sortby SORTBY  List returned recordsets sorted by SORTBY
   --search SEARCH  Filter returned recordsets by SEARCH criteria
```

To list a zone's recordsets selectively by type, the following command line would be used:

```
$ akamai dns list-recordsets example.com --type soa --type ns
```

With the following output:

```
Zone Recordsets
 
              NAME                TYPE    TTL                                             RDATA                                           
                                                                                                                                           
  example.com                     NS     86400   a1-98.akam.net.                                                                           
                                                                                                                                           
                                                 a12-65.akam.net.                                                                          
                                                                                                                                           
                                                 a13-65.akam.net.                                                                          
                                                                                                                                           
                                                 a2-64.akam.net.                                                                           
                                                                                                                                           
                                                 a3-64.akam.net.                                                                           
                                                                                                                                           
                                                 a4-65.akam.net.                                                                           
                                                                                                                                           
  example.com                     SOA    86400   a1-98.akam.net. hostmaster.example.com. 2020060510 3600 600 604800 300  
                                                                                                                                           
Zone: example.com
```

A similar example, sorting by recordset type and outputting in json format, would be the following:

```
$ akamai dns list-recordsets example.com --sortby type --json
```

and result on the following result: 

```
{
  "Recordsets": [
    {
      "name": "a_example.com",
      "type": "A",
      "ttl": 900,
      "rdata": [
        "10.0.0.10"
      ]
    },
    {
      "name": "aaaa_example.com",
      "type": "AAAA",
      "ttl": 900,
      "rdata": [
        "8001:ab8:85b3:0:0:8a1e:370:7225"
      ]
    },
    {
      "name": "example.com",
      "type": "NS",
      "ttl": 86400,
      "rdata": [
        "a1-98.akam.net.",
        "a12-65.akam.net.",
        "a4-65.akam.net."
      ]
    },
    {
      "name": "example.com",
      "type": "SOA",
      "ttl": 86400,
      "rdata": [
        "a1-98.akam.net. hostmaster.egl_clidns_primary_test_1.com. 2020060510 3600 600 604800 300"
      ]
    }
  ]
}
```

### Creating Multiple Zone Recordsets

The command `akamai dns create-recordsets` is used to create multiple recordsets in one command invocation.

The complete command line is:

```
$ akamai dns create-recordsets <zonename> [--json] [--suppress] [--output] [--file] 

Flags: 
   --json         Output as JSON [$AKAMAI_CLI_DNS_JSON]
   --suppress     Suppress command result output. Overrides other output related flags [$AKAMAI_CLI_DNS_SUPPRESS]
   --output FILE  Output command results to FILE
   --file FILE    FILE path to JSON formatted recordset content 
```

For example, the following command:

```
$ akamai dns create-recordsets rs_example.com --file new_recordsets.json --json
```

where the file file `new_recordsets.json` contains:

```
{
    "recordsets": [
        {
            "name": "example.rs_example.com",
            "type": "CNAME",
            "ttl": 1200,
            "rdata": [
                "www.example.com"
            ]
        },
       {
            "name": "a_rs_example.com",
            "type": "A",
            "ttl": 900,
            "rdata": ["10.0.0.20"]
       }
   ]
}
```

would result in the following output:

```
{
  "Recordsets": [
    {
      "name": "example.rs_example.com",
      "type": "CNAME",
      "ttl": 1200,
      "rdata": [
        "www.example.com."
      ]
    },
    {
      "name": "rs_example.com",
      "type": "SOA",
      "ttl": 86400,
      "rdata": [
        "a1-98.akam.net. hostmaster.rs_example.com. 2020060513 3600 600 604800 300"
      ]
    },
    {
      "name": "rs_example.com",
      "type": "NS",
      "ttl": 86400,
      "rdata": [
        "a1-98.akam.net.",
        "a12-65.akam.net.",
        "a3-64.akam.net.",
        "a4-65.akam.net."
      ]
    },
    {
      "name": "a_rs_example.com",
      "type": "A",
      "ttl": 900,
      "rdata": [
        "10.0.0.20"
      ]
    }
  ]
}
```

### Updating multiple zone Recordsets 

The command `akamai dns update-recordsets` is used to update multiple recordsets in one command invocation. Note: The default operation of the update-recordsets command is to specifically replace the recordsets in the provided file if they exist. The `overwrite' flag will REPLACE ALL existing recordsets in the zone.

The complete command line is:

```
$ akamai dns update-recordsets <zonename> [--json] [--suppress] [--output] [--file] 

Flags: 
   --json         Output as JSON [$AKAMAI_CLI_DNS_JSON]
   --suppress     Suppress command result output. Overrides other output related flags [$AKAMAI_CLI_DNS_SUPPRESS]
   --overwrite    Replace ALL Recordsets
   --output FILE  Output command results to FILE
   --file FILE    FILE path to JSON formatted recordset content
```

The following incorrect example updates the recordsets with the same file used in the previous create-recordsets example WITH the `--overwrite` flag. 

```
$ akamai dns update-recordsets --file new_recordsets.json --overwrite
```

resulting in the following error output. [The expected result being to remove the SOA record!] :

```
Updating Recordsets ... [FAIL]

Recordset update failed. Error: Zone "example.com" validation failed: [SOA record set is required for zone example.com]
```

### Retrieving a Recordset

The command `akamai dns retrieve-recordset` is used to retrieve a single recordset.

The complete command line is:

```
$ akamai dns retrieve-recordset <zonename> [--json] [--output] [--name] [--type] 

Flags: 
   --json         Output as JSON [$AKAMAI_CLI_DNS_JSON]
   --output FILE  Output command results to FILE
   --name NAME    Recordset NAME
   --type TYPE    Recordset TYPE
```

A recordset retrieval example is the following:

```
$ akamai dns retrieve-recordset egl_clidns_primary_test_1.com --name a_rs_example.com --type A
```

would result in the following output:

```
                     NAME                       TYPE   TTL     RDATA    
                                                                        
  a_rs_example.com                               A     900   10.0.0.20  
                                                                        
Zone: a_rs_example.com
```

The following example would direct output to a file:


```
$ akamai dns retrieve-recordset egl_clidns_primary_test_1.com --name a_rs_example.com --type A --output ./recordset_a.json
```

resulting in ./recordset_a.json would containing:

```
{
        "name": "a_rs_example.com",
        "type": "A",
        "ttl": 900,
        "rdata": ["10.0.0.20"]
}
```

### Creating a Recordset

The command `akamai dns create-recordset` is used to create a single recordset. Fields and values can be provided on the command line or input file.

The complete command line is:

```
$ akamai dns create-recordset <zonename> [--json] [--suppress] [--output] [--name] [--type] [--ttl] [--rdata] [--file] 

Flags: 
   --json         Output as JSON [$AKAMAI_CLI_DNS_JSON]
   --suppress     Suppress command result output. Overrides other output related flags [$AKAMAI_CLI_DNS_SUPPRESS]
   --output FILE  Output command results to FILE
   --name NAME    Recordset NAME
   --type TYPE    Recordset TYPE
   --ttl TTL      Recordset TTL (default: 0)
   --rdata RDATA  Recordset RDATA. Multiple flags allowed.
   --file FILE    FILE path to JSON formatted recordset content
```

To create a recordset via command line, the following would be invoked:

```
$ akamai dns create-recordset example.com --name a_record_example.com --type A --ttl 600 --rdata 10.1.1.1 --json
```

and result in the following output:

```
{
  "name": "a_record_example.com",
  "type": "A",
  "ttl": 600,
  "rdata": [
    "10.1.1.1"
  ]
}
``` 

### Updating a Recordset

The command `akamai dns update-recordset` is used to update a single existing recordset. Fields and values can be provided on the command line or input file.

The complete command line is:

```
$ akamai dns update-recordset <zonename> [--json] [--suppress] [--output] [--name] [--type] [--ttl] [--rdata] [--file] 

Flags: 
   --json         Output as JSON [$AKAMAI_CLI_DNS_JSON]
   --suppress     Suppress command result output. Overrides other output related flags [$AKAMAI_CLI_DNS_SUPPRESS]
   --output FILE  Output command results to FILE
   --name NAME    Recordset NAME
   --type TYPE    Recordset TYPE
   --ttl TTL      Recordset TTL (default: 0)
   --rdata RDATA  Record RDATA. Multiple flags allowed.
   --file FILE    FILE path to JSON formatted recordset content. Allows multiple recordsets.
```

An example recordset update using an input file would be as follows:

```
$ akamai dns update-recordset example.com --file ./a_recordset.json --suppress
```

where the file `./a_recordset.json` contains:

```
{
  "name": "a_record_example.com",
  "type": "A",
  "ttl": 900,
  "rdata": [
    "10.1.1.1"
  ]
}
```

 and return no output.

### Deleting a Recordset

The command to delete a single recordset is `akamai dns delete-recordset`.

The complete command line is:

```
$ akamai dns delete-recordset <zonename> [--name] [--type] 

Flags: 
   --name NAME    Recordset NAME
   --type TYPE    Recordset TYPE
```

### Submit Bulk Zone Request

To submit a create or delete bulk zone operation request, use the following command:

The complete command line is:

```
$ akamai dns submit-bulkzones  [--json] [--output] [--suppress] [--contractid] [--groupid] [--bypasszonesafety] [--create] [--delete] [--file]

Flags:
   --json              Output as JSON [$AKAMAI_CLI_DNS_JSON]
   --output FILE       Output command results to FILE
   --suppress          Suppress command result output. Overrides other output related flags [$AKAMAI_CLI_DNS_SUPPRESS]
   --contractid ID     Contract ID. Required for create.
   --groupid ID        Group ID. Required for create.
   --bypasszonesafety  Bypass zone safety check. Optional for delete.
   --create            Bulk zone create operation.
   --delete            Bulk zone delete operation.
   --file FILE         Read JSON formatted input from FILE
```

NOTE: The CLI currently limits the number of zones in a submit request to 1000. If an invocation presents more than 1000 zones, the zones will be submitted in batches of 1000 and multiple Request Ids will be returned. The batch size can be changed by setting the environment variable AKAMAI_ZONES_BATCH_SIZE.

An example create submit request  would be as follows:

```
$ akamai dns submit-bulkzones --create --contractid 1-3CV382 --groupid 18432 -file bulkcreate.json
Preparing bulk zones submit request
Submitting Bulk Zones request
Submitting Bulk Zones request  ... [OK]

Assembling Bulk Zone Response Content
Assembling Bulk Zone Response Content ... [OK]

Writing Request Status to bulkSubmitRequest.f3fcbf11-1b03-420e-9e2b-88cd0096fa62
Writing Request Status to bulkSubmitRequest.f3fcbf11-1b03-420e-9e2b-88cd0096fa62 ... [OK]



Bulk Zones Request Submission Status

  Request Id        f3fcbf11-1b03-420e-9e2b-88cd0096fa62

  Expiration Date   2020-10-27T13:23:20.042484Z
```

where the file `./bulkcreate.json` contains:

```
{
    "zones": [
        {
            "zone": "one.xxx_testbulk.net",
            "type": "secondary",
            "comment": "testing bulk operations",
            "masters": [
                "1.2.3.4",
                "1.2.3.10"
            ]
        },
        {
            "zone": "two.xxx_testbulk.net",
            "type": "secondary",
            "comment": "testing bulk operations",
            "masters": [
                "1.2.3.40",
                "1.2.3.100"
            ]
        }
    ]
}
```

An example delete submit request would be as follows:

```
$ akamai dns submit-bulkzones --delete --file bulkdelete.json
Preparing bulk zones submit request
Submitting Bulk Zones request
Submitting Bulk Zones request  ... [OK]

Assembling Bulk Zone Response Content
Assembling Bulk Zone Response Content ... [OK]

Writing Request Status to bulkSubmitRequest.f3fcbf11-1b03-420e-9e2b-88cd0096fa62
Writing Request Status to bulkSubmitRequest.f3fcbf11-1b03-420e-9e2b-88cd0096fa62 ... [OK]



Bulk Zones Request Submission Status

  Request Id        f3fcbf11-1b03-420e-9e2b-88cd0096fa62

  Expiration Date   2020-10-27T13:23:20.042484Z
```

where the file `./bulkdelete.json` contains:

```
{
    "zones": [
            "one.xxx_testbulk.net",
            "two.xxx_testbulk.net"
    ]
}
```

### Get Bulk Zone Request Status 

To retrieve the current status of a create or delete bulk zone operation request, use the following command:

The complete command line is:

```
$ akamai dns status-bulkzones  [--json] [--output] [--create] [--delete] [--requestid]

Flags:
   --json         Output as JSON [$AKAMAI_CLI_DNS_JSON]
   --output FILE      Output command results to FILE
   --requestid value  Request Id. Multiple args allowed.
   --create           Bulk zone create operation.
   --delete           Bulk zone delete operation.
```

An example status check for a request would be as follows:

```
$ akamai dns status-bulkzones -create -requestid 309679b5-1ab1-4837-9666-0019d1be891e -json
Preparing bulk zones status request(s)
Submitting Bulk Zones request 
Submitting Bulk Zones request  ... [OK]

Assembling Bulk Zone Response Content
Assembling Bulk Zone Response Content ... [OK]


[
    {
        "requestId": "309679b5-1ab1-4837-9666-0019d1be891e",
        "zonesSubmitted": 2,
        "successCount": 2,
        "failureCount": 0,
        "isComplete": true,
        "expirationDate": "2020-10-26T18:03:21.072004Z"
    }
]
```

An example status check for multiple request ids would be as follows:

```
$ akamai dns status-bulkzones -create -requestid 15bc138f-8d82-451b-80b7-a56b88ffc474 -requestid 0c22641b-7a30-44be-8fdd-092bf875f3bc
Preparing bulk zones status request(s)
Submitting Bulk Zones request
Submitting Bulk Zones request  ... [OK]

Assembling Bulk Zone Response Content
Assembling Bulk Zone Response Content ... [OK]



Bulk Zones Request Status

  Request Id   15bc138f-8d82-451b-80b7-a56b88ffc474

               Zones Submitted                        2

               Success Count                          0

               Failure Count                          2

               Complete                               true

               Expiration Date                        2020-10-28T17:10:04.515792Z

 Request Id   0c22641b-7a30-44be-8fdd-092bf875f3bc

               Zones Submitted                        2

               Success Count                          0

               Failure Count                          2

               Complete                               true

               Expiration Date                        2020-10-28T17:10:05.613474Z
```

### Get Bulk Zone Request Result

To retrieve the result of a create or delete bulk zone operation request, use the following command:

The complete command line is:

```
$ akamai dns status-bulkzones  [--json] [--output] [--create] [--delete] [--requestid]

Flags:
   --json             Output as JSON [$AKAMAI_CLI_DNS_JSON]
   --output FILE      Output command results to FILE
   --requestid value  Request Id. Multiple args allowed.
   --create           Bulk zone create operation.
   --delete           Bulk zone delete operation.

```

An example result retrieval for a request would be as follows:

```
$ akamai dns result-bulkzones -delete -requestid f3fcbf11-1b03-420e-9e2b-88cd0096fa62 -json
Preparing bulk zones result request(s)
Submitting Bulk Zones request
Submitting Bulk Zones request  ... [OK]

Assembling Bulk Zone Response Content
Assembling Bulk Zone Response Content ... [OK]


[
    {
        "requestId": "f3fcbf11-1b03-420e-9e2b-88cd0096fa62",
        "successfullyDeletedZones": [
            "one.xxx_testbulk.net",
            "two.xxx_testbulk.net"
        ],
        "failedZones": []
    }
]
```

An example result retrieval for multiple request ids would be as follows:

```
$ akamai dns result-bulkzones -create -requestid 15bc138f-8d82-451b-80b7-a56b88ffc474 -requestid 0c22641b-7a30-44be-8fdd-092bf875f3bc -json
Preparing bulk zones result request(s)
Submitting Bulk Zones request
Submitting Bulk Zones request  ... [OK]

Assembling Bulk Zone Response Content
Assembling Bulk Zone Response Content ... [OK]


[
  {
    "requestId": "15bc138f-8d82-451b-80b7-a56b88ffc474",
    "successfullyCreatedZones": [],
    "failedZones": [
      {
        "zone": "one.xxx_testbulk.net",
        "failureReason": "ZONE_ALREADY_EXISTS"
      }
    ]
  },
  {
    "requestId": "0c22641b-7a30-44be-8fdd-092bf875f3bc",
    "successfullyCreatedZones": [],
    "failedZones": [
      {
        "zone": "two.xxx_testbulk.net",
        "failureReason": "ZONE_ALREADY_EXISTS"
      }
    ]
  }
]
```

### Retrieving a Zone

To retrieve a Zone use the `retrieve-zone` command:

```
$ akamai dns retrieve-zone example.org
```

To filter to specific record types use one or more `--filter <TYPE>` flags. For example,
to show just `A` and `AAAA` records:

```
$ akamai dns retrieve-zone example.org --filter A --filter AAAA
```

You can also output the result as JSON, by adding the `--json` flag:

```sh
$ akamai dns retrieve-zone example.org --filter A --filter AAAA --json
```
```json
{
 "records": [
  {
   "name": "a.example.org",
   "type": "A",
   "ttl": 900,
   "rdata": [
    "10.0.0.10",
    "10.0.0.20"
   ]
  },
  {
   "name": "test.example.org",
   "type": "AAAA",
   "ttl": 600,
   "rdata": [
    "2001:db8:0:0:0:0:0:1"
   ]
  }
 ],
 "zone": {
  "zone": "example.org",
  "type": "PRIMARY",
  "comment": "primary zone config",
  "signAndServe": true,
  "signAndServeAlgorithm": "RSA_SHA256",
  "contractId": "1-1AB123",
  "activationState": "PENDING",
  "lastActivationDate": "2025-07-02T07:57:34.146414Z",
  "lastModifiedBy": "xxxxx",
  "lastModifiedDate": "2025-07-02T09:28:38.141705Z",
  "versionId": "a0b4730e-fbbe-40ad-96b3-ac6a4cbadb1e"
 }
}
```

### Update a Zone

Update a zone using `akamai dns update-zone`. This command allows you to input either
a [Edge DNS JSON payload](https://developer.akamai.com/api/luna/config-dns/resources.html#addormodifyazone), or a standard DNS zone file.

By default, this will **append** the records to the zone.

You can either specify a file or redirect content via STDIN:

```sh
$ akamai dns update-zone example.org -f new-records.zone.json
```

is identical to:

```sh
$ cat new-records.zone.json | akamai dns update-zone example.org
```

To use DNS zone format, specify the `--dns` flag:

```sh
$ akamai dns update-zone --dns -f new-records.zone
```

#### Overwriting a Zone

If you want to overwrite the existing zone entirely, add the `--overwrite` flag:

```sh
$ akamai dns update-zone example.org --overwrite -f example.org.zone.json
```

### Add a New Record

To add a new DNS record use `akamai dns add-record <record type>`. Each setting for the record is a flag, for example to add a `CNAME` record:

```
$ akamai dns add-record CNAME example.org --name www --rdata example.org --ttl 3600
```

### Remove a Record

Use `akamai dns rm-record <record type>` to remove one or more records matching the given flags.

```
$ akamai dns rm-record A example.org --name www
```


By default the command will ask you to verify which records to remove if more than one matches the given criteria. You can force it to remove all matching records by passing in the `--force-multiple` flag.

If the command is run in a non-interactive terminal, **or** the `--non-interactive` flag is passed in, without the `--force-multiple` flag the command will remove records if only one match is found, otherwise it will exit with status code `1`.


## License

This package is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for details.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fakamai%2Fcli-dns.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fakamai%2Fcli-dns?ref=badge_large)

<!--esi 
<esi:eval src="/internal/footer.html?http" dca="esi" />
-->
