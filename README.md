# Akamai CLI for Fast DNS

An [Akamai CLI](https://developer.akamai.com/cli) package for managing DNS Zones using Fast DNS.

## Getting Started

### Installing

To install this package, use Akamai CLI:

```sh
$ akamai install dns
```

## Usage

```
$  akamai dns [--edgerc] [--section] <command> [sub-command]`
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
  "token": "9218376e14c2797e0d06e8d2f918d45f",
  "zone": {
    "name": "example.com",
    "a": [
      {
        "name": "www",
        "ttl": 3600,
        "active": true,
        "target": "192.0.2.1"
      }
    ],
    "aaaa": [
      {
        "name": "www",
        "ttl": 3600,
        "active": true,
        "target": "2001:db8:0:0:0:0:0:1"
      }
    ]
  }
}
```

### Update a Zone

Update a zone using `akamai dns update-zone`. This command allows you to input either
a [Fast DNS JSON payload](https://developer.akamai.com/api/luna/config-dns/resources.html#addormodifyazone), or a standard DNS zone file.

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
$ akamai dns add-record CNAME example.org --name www --target example.org --ttl 3600
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