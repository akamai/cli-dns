# Release Notes

## Version 0.8.0 

### Features/Enhancements

* Upgrade to Edgegrid v13.2.0
* Update documentation and --help for recordset commands to enforce mutually exclusive options (filename or recordname+recordtype)
* Add LastModifiedBy field to Zoneconfig table output
* Standardize command status/failure messaging ([OK/FAIL]) and table rendering format
* Update project license file
* Migrate to go 1.26.4

## Version 0.7.0 (January 2026)

### Features/Enhancements

* Add support for arm64 architecture
* Upgrade to Edgegrid v12.3.0 
* Migrate to go 1.25.0

## Version 0.6.0 (July 4, 2025)

### Features/Enhancements

* Upgrade to Edgegrid v11.0.0 
* Session based authentication
* Migrate to go 1.23
* The add-record, retrieve-zone, rm-record and update-zone command now use the Edge DNS API v2

* add-record command
    - The command checks if the DNS record already exists before creating it.
    - If the record does not exist, it is created.
    - If the record exists with the same type:
        - It is updated only if the TTL or RDATA values have changed.
        - If there are no changes, no update is performed and a message is shown.
    - uses the --rdata flag instead of --target flag


* retrieve-zone command
    - Fetches detailed information about the specified DNS zone.
    - If the zone is of type ALIAS, it displays only the zone details.
    - Otherwise, fetches all DNS recordsets associated with the zone.
    - Supports filtering the recordsets by DNS record type via a --filter flag.
    - Supports output in either human-readable table format or JSON format

* rm-record command
    - Enables deletion of DNS records from a zone using record type and name.

* update-zone command
    - Updates DNS recordsets for the specified zone using input JSON or master zone file.
    - Validates zone; disallows ALIAS zones for recordset updates.
    - Accepts input from file (--file) or STDIN.
    - Supports master zone file upload with --dns flag.
    - Can overwrite all existing recordsets (--overwrite) or merge changes.

## Version 0.5.0 (May 10, 2023)

### Features/Enhancements

* Add M1 support
* Migrate to go 1.18







