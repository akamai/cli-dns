# Release Notes

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







