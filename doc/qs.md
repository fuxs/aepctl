# Query Service

The Query Service offers an SQL interface to the AEP.

The command uses the following pattern:

```terminal
aepctl (verb) (noun)
```

E.g., if you want to list all queries

```terminal
aepctl list queries
```

Some commands support a preferred short notation, e.g. `ls` for `list`:

```terminal
aepctl ls namespaces
```

The following verbs are supported by Query Service:

* `get` ([Get Connection](#Get-Connection))
* `ls` or `list` ([List Queries](#List-Queries))
* `psql` ([PSQL](#PSQL))

# Get Connection

The `get connection` command returns all parameters for a connection with the
interactive PosgreSQL terminal `psql`:

```terminal
aepctl get connection
```

The default view shows a table with the columns PATH and Value

```terminal
PATH      VALUE
username  127075E95BF479EC0A495C73@AdobeOrg
dbName    mysandbox:all
host      experienceplatform.platform-query.adobe.io
version   1
port      80
token     eyJhbGciOiJSUzI1NiIsIng1dSI6Imltc19uYTEta2V5LTEuY2VyIn0…
```

See [Output](output.md) for other output formats.

The command `aepctl psql` uses this information to start the `psql` client for
you, see section [PSQL](#PSQL) for more information.

# List Queries

The `ls queries` command returns a list of all queries:

```terminal
aepctl ls queries
```

The default view shows a table with the columns ID, NAME and LAST MODIFIED

```terminal
ID                                   NAME                 LAST MODIFIED
58bcd85a-f712-4896-bf6d-50d0d6f9a761 SELECT sum("crmid")  01 Sep 21 11:06 CEST
ddad3647-3d7c-4d95-84bb-d2d938c53623 SELECT "date_key",   01 Sep 21 11:04 CEST
1815c06a-3b8d-4fd2-9d85-bbf8a4fe8be8 SELECT "adwh_dim_nam 01 Sep 21 11:03 CEST
…
```

The flag `-o wide` adds the column SQL:

```terminal
ID                                   NAME          LAST MODIFIED       SQL
87907703-b1fe-4e41-be6e-8f017bdd6eeb First Call    11 Dec 20 18:30 CET show tables
78764d94-c458-438a-8326-21f80d43f248 Get Events    11 Dec 20 18:31 CET select * from demo_system_event_dataset_for_website_global_v1_1
0a105db0-283b-4dc9-80ad-6b3936134ee1 -             11 Dec 20 18:31 CET select date_format( timestamp , 'yyyy-MM-dd') AS Day,
                                                                     count(*) AS productViews
                                                              from   website_events
                                                              where  _experienceplatform.demoEnvironment.category IN ('Wurstwaren', 'Eisenwaren')
                                                              and    _eventType = 'productView'
                                                              group by Day
                                                              limit 10
```

The number of results can be very high. The command uses paging, by default each
call (page) returns 100 entries and the next call will be executed
automatically. Use the flag `--limit` to specify the number of results per call,
e.g. `--limit=10` will return 10 entries. This would cause 10 times more calls
to the AEP compared to the default of 100. Usually this flag is used in
combination with `--paging=false` which disables paging and leads to only one
call.

`aepctl ls queries` prints the result of each call immediately, leading to
different column widths for each page. Use the flag `--flush=false` for equal
column width.

In order to sort the queries by time use one of the following flags:

* `--order=+created` ascending creation date
* `--order=-created` descending creation date
* `--order=+updated` ascending modification date
* `--order=-updated` descending modification date

Specify the start date for the order flag with `--start`, e.g.
`=2021-01-01T00:00:00Z` to get all queries from the beginning of 2021.

The number of results can be reduced by the flag `--filter` with comma separated
definitions:

Supported properties:

* created
* updated
* state
* id
* referenced_datasets
* userId
* sql
* templateId
* templateName
* client
* scheduleId
* scheduleRunId

Supported operators:

* `>` greater than
* `<` less than
* `>=` greater than or equal to
* `<=` less than or equal to
* `==` equal to
* `!=` not equal to
* `~` contains

Examples:

* `id==78764d94-c458-438a-8326-21f80d43f249`
* `client=API,created>=2021-08-01T18:30:00Z` (comma means AND)
* `sql~SELECT c1` (SQL text must contain no comma)

Other flags to control the number of results:
* `--exclude-deleted=false` to show soft deleted queries (`true` is default and
  doesn't have to be set).
* `--exclude-hidden=false` to show system generated queries (`true` is default
  and doesn't have to be set).

See [Output](output.md) for other output formats.

# PSQL

The command `aepctl psql` runs an PostgreSQL interactive terminal connected to
the AEP instance. It requires an installed `psql` command (OS X: run `brew
install postgres` for installation). If the command is not on the path or if a
specific version is required then use the flag `--command=/usr/local/bin/psql`

The flag `--print` just prints the command with all parameters. Use the command
`aepctl get connection` for a structured output of all parameters.