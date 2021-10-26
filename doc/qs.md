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

* `cancel` ([Cancel Query](#Cancel-Query) and [Cancel Scheduled Query
  Run](#Cancel-Scheduled-Query-Run))
* `create` ([Create Query](#Create-Query), [Create Query
  Template](#Create-Query-Template) and [Create Scheduled
  Query](#Create-Scheduled-Query))
* `delete` ([Delete Query](#Delete-Query), [Delete Query
  Template](#Delete-Query-Template) and [Delete Scheduled
  Query](#Delete-Scheduled-Query))
* `get` ([Get Connection](#Get-Connection), [Get Query](#Get-Query), [Get Query
  Template](#Get-Query-Template), [Get Scheduled Query](#Get-Scheduled-Query)
  and [Get Scheduled Query Run](#Get-Scheduled-Query-Run))
* `ls` or `list` ([List Queries](#List-Queries), [List Query
  Templates](#List-Query-Templates), [List Scheduled
  Queries](#List-Scheduled-Queries) and [List Scheduled Query
  Runs](#List-Scheduled-Query-Runs))
* `psql` ([PSQL](#PSQL))
* `trigger` ([Trigger Scheduled Query Run](#Trigger-Scheduled-Query-Run))
* `update` ([Update Query Template](#Update-Query-Template) and [Update
  Scheduled Query](#Update-Scheduled-Query))

# Cancel Query

The `cancel query` command supports the cancellation of one or more submitted
queries which haven't been executed.

```terminal
aepctl cancel query 17cc37a7-015f-4e5e-9b3b-787b9d90a050 12fec78c-eff3-4266-a106-e02dde48aa1c
```

In case of no errors the command returns without any output. Use the flag
`--response` to show the response from the server.

```terminal
aepctl cancel query 17cc37a7-015f-4e5e-9b3b-787b9d90a050 12fec78c-eff3-4266-a106-e02dde48aa1c --response
```
The command should return something like:
```terminal
{
  "message": "Query cancel request received.",
  "statusCode": 202
}
{
  "message": "Query cancel request received.",
  "statusCode": 202
}
```

If you pass multiple query IDs and an error occurs then the command will stop
the execution. Use the flag `--ignore` to execute the command for all IDs
ignoring any errors.

# Cancel Scheduled Query Run

The `cancel run` command cancels a specific run of a scheduled query. It
requires the ID of the scheduled query (e.g. list scheduled queries with `aepctl
ls schedules`) and the ID of the run (e.g. list runs with `aepctl ls runs
SCHEDULED_QUERY_ID`).

```terminal
aepctl cancel run 907075e95bf479ec0a495c73_68b9c64d-0dde-4db5-b9c6-4d0ddebdb5a7_my_scheduled_querywvo1ozznm5_bsngzg c2NoZWR1bGVkX18yMDIxLTEwLTI2VDA1OjMwOjAwKzAwOjAw
```

In case of no errors the command returns without any output. Use the flag
`--response` to show the response from the server.

# Create Query

The `create query` command creates a new query and requires a payload in JSON format.

`dbName` and `sql` are mandatory, `name` and `description` are optional. Please
see the [Query Service
API](https://www.adobe.io/experience-platform-apis/references/query-service/#operation/create_query)
documentation for a full list of JSON attributes.

```json
{
    "dbName": "myorg:all",
    "sql": "SELECT * FROM callcenter_interaction_analysis LIMIT 5;",
    "name": "Sample Query",
    "description": "A sample of a query."
}
```

This payload can be provided in a file, e.g. `query.json` in the folder `examples/create`:

```terminal
aepctl create query examples/create/query.json
```

If no file name is provided then aepctl reads from the standard input stdin,
hence heredoc is supported, too:

```terminal
aepctl create query << EOF
{
  "dbName": "myorg:all",
  "sql": "SELECT * FROM callcenter_interaction_analysis LIMIT 5;",
  "name": "Sample Query",
  "description": "A sample of a query."
}
EOF
```

In case of no errors the command returns without any output. Use the flag
`--response` to show the response from the server.

It is also possible to create multiple queries with one call:

```terminal
aepctl create query query1.json query2.json last_query.json
```

Use the flag `--ignore` to skip over errors during the execution. Otherwise,
aepctl would stop the execution in the case of the first error.

# Create Query Template

The `create template` command creates a new query-template and requires a
payload in JSON format.

```json
{
  "sql": "SELECT $key from $key1 where $key > $key2;",
  "queryParameters": {
    "key": "value",
    "key1": "value1",
    "key2": "value2"
  },
  "name": "Sample-Template"
}
```

This payload can be provided in a file, e.g. `query-template.json` in the folder `examples/create`:

```terminal
aepctl create template examples/create/query-template.json
```

If no file name is provided then aepctl reads from the standard input stdin,
hence heredoc is supported, too:

```terminal
aepctl create template << EOF
{
  "sql": "SELECT $key from $key1 where $key > $key2;",
  "queryParameters": {
    "key": "value",
    "key1": "value1",
    "key2": "value2"
  },
  "name": "Sample-Template"
}
EOF
```

In case of no errors the command returns without any output. Use the flag
`--response` to show the response from the server.

It is also possible to create multiple queries with one call:

```terminal
aepctl create template query-template1.json query-template2.json last_query-template.json
```

Use the flag `--ignore` to skip over errors during the execution. Otherwise,
aepctl would stop the execution in case of the first error.

# Create Scheduled Query

The `create schedule` command creates a new scheduled query and requires a
payload in JSON format.

```json
{
  "query": {
    "dbName": "myorg:all",
    "sql": "SELECT * FROM callcenter_interaction_analysis LIMIT 5;",
    "name": "My Scheduled Query",
    "description": "A sample scheduled query."
  },
  "schedule": {
      "schedule": "30 * * * *",
      "startDate": "2021-09-07T12:00:00Z"
  }
}
```

This payload can be provided in a file, e.g. `schedule.json` in the folder `examples/create`:

```terminal
aepctl create schedule examples/create/schedule.json
```

If no file name is provided then aepctl reads from the standard input stdin,
hence heredoc is supported, too:

```terminal
aepctl create schedule << EOF
{
  "query": {
    "dbName": "myorg:all",
    "sql": "SELECT * FROM callcenter_interaction_analysis LIMIT 5;",
    "name": "My Scheduled Query",
    "description": "A sample scheduled query."
  },
  "schedule": {
      "schedule": "30 * * * *",
      "startDate": "2021-09-07T12:00:00Z"
  }
}
EOF
```

In case of no errors the command returns without any output. Use the flag
`--response` to show the response from the server.

It is also possible to create multiple queries with one call:

```terminal
aepctl create query query1.json query2.json last_query.json
```

Use the flag `--ignore` to skip over errors during the execution. Otherwise,
aepctl would stop the execution in the case of the first error.

# Delete Query

The `delete query` command deletes one or more queries. This is a soft delete
and the `ls queries --exclude-deleted=false` can be used to show deleted
queries.

```terminal
aepctl delete query 17cc37a7-015f-4e5e-9b3b-787b9d90a050 12fec78c-eff3-4266-a106-e02dde48aa1c
```

In case of no errors the command returns without any output. Use the flag
`--response` to show the response from the server.

```terminal
aepctl delete query 17cc37a7-015f-4e5e-9b3b-787b9d90a050 12fec78c-eff3-4266-a106-e02dde48aa1c --response
```
The command should return something like:
```terminal
{
  "message": "Query soft delete successful.",
  "statusCode": 200
}
{
  "message": "Query soft delete successful.",
  "statusCode": 200
}
```
If you pass multiple query IDs and an error occurs then the command will stop
the execution. Use the flag `--ignore` to execute the command for all IDs
ignoring any errors.

# Delete Query Template

The `delete template` command deletes one or more query templates.

```terminal
aepctl delete template 169a0264-9946-4cf2-af40-231a55cc6d47 ab1f5dda-3e7c-4383-b670-593e6abe885c
```

In case of no errors the command returns without any output. Use the flag
`--response` to show the response from the server.

If you pass multiple query template IDs and an error occurs then the command will stop
the execution. Use the flag `--ignore` to execute the command for all IDs
ignoring any errors.

# Delete Scheduled Query

The `delete schedule` command deletes one or more scheduled queries.

```terminal
aepctl delete schedule 907075e95bf479ec0a495c73_68b9c64d-0dde-4db5-b9c6-4d0ddebdb5a7_my_scheduled_querywvo1ozznm5_bsngzg
```

In case of no errors the command returns without any output. Use the flag
`--response` to show the response from the server.

If you pass multiple schedules IDs and an error occurs then the command will
stop the execution. Use the flag `--ignore` to execute the command for all IDs
ignoring any errors.

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

# Get Query

The `get query` command returns the status of the query with the passed ID:

```terminal
aepctl get query 87907703-b1fe-4e41-be6e-8f017bdd6eeb
```

The default view shows a table with the columns ID, NAME, STATE and LAST MODIFIED:

```terminal
ID                                   NAME         STATE   LAST MODIFIED
87907703-b1fe-4e41-be6e-8f017bdd6eeb First Call   SUCCESS 30 Sep 21 06:16 CEST
```

The flag `-o wide` adds the column SQL:

```terminal
ID                                   NAME          LAST MODIFIED       SQL
87907703-b1fe-4e41-be6e-8f017bdd6eeb First Call    11 Dec 20 18:30 CET SHOW TABLES;
```

See [Output](output.md) for other output formats.

# Get Query Template

The `get template` command returns the query-template with the passed ID:

```terminal
aepctl get template ab1f5dda-3e7c-4383-b670-593e6abe885b
```

The default view shows a table with the columns ID, NAME, STATE and LAST MODIFIED:

```terminal
ID                                   NAME        STATE   LAST MODIFIED
a9941df4-be98-417a-ae9f-5ad8ad2d3c45 My Example  SUCCESS 30 Sep 21 06:16 CEST
```

The flag `-o wide` adds the column SQL:

```terminal
ID                                   NAME        STATE   LAST MODIFIED         SQL
a9941df4-be98-417a-ae9f-5ad8ad2d3c45 My Example  SUCCESS 30 Sep 21 06:16 CEST  SHOW TABLES;
```

See [Output](output.md) for other output formats.

# Get Scheduled Query

The `get schedule` command returns the status of the scheduled query with the passed ID:

```terminal
aepctl get schedule 907075e95bf479ec0a495c73_68b9c64d-0dde-4db5-b9c6-4d0ddebdb5a7_my_scheduled_querywvo1ozznm5_bsngzg
```

The default view shows a table with the columns ID, NAME, STATE and LAST MODIFIED:

```terminal
ID                                                                                                             NAME                            STATE   LAST MODIFIED
907075e95bf479ec0a495c73_68b9c64d-0dde-4db5-b9c6-4d0ddebdb5a7_my_scheduled_querywvo1ozznm5_bsngzg My Scheduled Query ENABLED 07 Sep 21 11:48 CEST
```

The flag `-o wide` adds the columns SCHEDULE, START DATE and SQL:

```terminal
907075e95bf479ec0a495c73_68b9c64d-0dde-4db5-b9c6-4d0ddebdb5a7_my_scheduled_querywvo1ozznm5_bsngzg My Scheduled Query ENABLED 07 Sep 21 11:48 CEST 30 * * * * 2021-09-08T12:00:00.000Z SELECT * FROM callcenter_interaction_analysis LIMIT 5;
```

See [Output](output.md) for other output formats.

# Get Scheduled Query Run

The `get rund` command returns the status of the scheduled query run with the
passed IDs for the scheduled query (first argument) and the run (second argument):

```terminal
aepctl get schedule 907075e95bf479ec0a495c73_68b9c64d-0dde-4db5-b9c6-4d0ddebdb5a7_my_scheduled_querywvo1ozznm5_bsngzg c2NoZWR1bGVkX18yMDIxLTA5LTEwVDExOjMwOjAwKzAwOjAw
```

The default view shows a table with the columns TASK ID, STATE, MESSAGE,
DURATION, START DATE, END DATE:

```terminal
TASK ID    STATE  MESSAGE           DURATION START DATE           END DATE
wvo1oZzNm5 FAILED Processing Failed 370      10 Sep 21 14:53 CEST 10 Sep 21 14:59 CEST
```

See [Output](output.md) for other output formats.

# List Queries

The `ls queries` command returns a list of all queries:

```terminal
aepctl ls queries
```

The default view shows a table with the columns ID, NAME and LAST MODIFIED:

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

# List Query Templates

The `ls templates` command returns a list of all query templates:

```terminal
aepctl ls templates
```

The default view shows a table with the columns ID, NAME and LAST MODIFIED:

```terminal
ID                                   NAME                                                          LAST MODIFIED
179a0264-9946-4bf2-af40-231a55cc6d46 Target Audience                                               05 Oct 21 19:02 CEST
ac1f5dda-3e7c-4283-b670-593e6abe885b High Intent Shoppers                                          20 Sep 21 22:40 CEST
97336c6d-d9e4-38d7-9452-ce479d633e12 Test Group                                                    06 Sep 21 14:01 CEST
…
```

The flag `-o wide` adds the column SQL:

```terminal
ID                                   NAME                                                          LAST MODIFIED    SQL
179a0264-9946-4bf2-af40-231a55cc6d46 Target Audience                                               05 Oct 21 19:02 CEST  SHOW TABLES;
ac1f5dda-3e7c-4283-b670-593e6abe885b High Intent Shoppers                                          20 Sep 21 22:40 CEST  SELECT * FROM hight_intent_shoppers WHERE value > 1000;
97336c6d-d9e4-38d7-9452-ce479d633e12 Test Group                                                    06 Sep 21 14:01 CEST  SELECT * FROM test_group;
…
```

The number of results can be reduced by the flag `--filter` with comma separated
definitions:

Supported properties:

* created (with operators `>` greater than and `<` less than)
* lastUpdatedBy (with operator `==`)
* name (with operator `==` equal to and `~` contains)
* userId (with operator `==`)

See [Output](output.md) for other output formats.

# List Scheduled Queries

The `ls schedules` command returns a list of all queries:

```terminal
aepctl ls schedules
```

The default view shows a table with the columns ID, NAME, STATE and LAST MODIFIED:

```terminal
ID                                                                                                             NAME                            STATE   LAST MODIFIED
907075e95bf479ec0a495c73_68b9c64d-0dde-4db5-b9c6-4d0ddebdb5a7_my_scheduled_querywvo1ozznm5_bsngzg My Scheduled Query ENABLED 07 Sep 21 11:48 CEST
907075e95bf479ec0a495c73_68b9c64d-0dde-4db5-b9c6-4d0ddebdb5a7_my_scheduled_queryigvl8ysbd3_tf6vgf My Scheduled Query ENABLED 07 Sep 21 13:21 CEST
…
```

The flag `-o wide` adds the columns SCHEDULE, START DATE and SQL:

```terminal
907075e95bf479ec0a495c73_68b9c64d-0dde-4db5-b9c6-4d0ddebdb5a7_my_scheduled_querywvo1ozznm5_bsngzg My Scheduled Query ENABLED 07 Sep 21 11:48 CEST 30 * * * * 2021-09-08T12:00:00.000Z SELECT * FROM callcenter_interaction_analysis LIMIT 5;
907075e95bf479ec0a495c73_68b9c64d-0dde-4db5-b9c6-4d0ddebdb5a7_my_scheduled_queryigvl8ysbd3_tf6vgf My Scheduled Query ENABLED 07 Sep 21 13:21 CEST 30 * * * * 2021-09-07T12:00:00.000Z SELECT * FROM store_interaction_analysis LIMIT 5;
```

The number of results can be reduced by the flag `--filter` with comma separated
definitions:

Supported properties:

* created (with operators `>` greater than and `<` less than)
* templateId (with operator `==`)
* userId (with operator `==`)

See [Output](output.md) for other output formats.

# List Scheduled Query Runs

The `ls runs` command returns a list of all runs for a scheduled query:

```terminal
aepctl ls runs 
```

The default view shows a table with the columns ID, STATE and CREATED:

```terminal
ID                                               STATE   CREATED
c2NoZWR1bGVkX18yMDIxLTA5LTA4VDExOjMwOjAwKzAwOjAw SUCCESS 08 Sep 21 13:30 CEST
c2NoZWR1bGVkX18yMDIxLTA5LTA4VDEyOjMwOjAwKzAwOjAw FAILED  08 Sep 21 14:30 CEST
c2NoZWR1bGVkX18yMDIxLTA5LTA4VDEzOjMwOjAwKzAwOjAw FAILED  08 Sep 21 15:30 CEST
…
```

The number of results can be reduced by the flag `--filter` with comma separated
definitions:

Supported properties:

* created (with operators `>` greater than and `<` less than)
* state (with operators `==` equal to and `!=` not equal to)
* externalTrigger (with operator `==`)

See [Output](output.md) for other output formats.

# PSQL

The command `aepctl psql` runs an PostgreSQL interactive terminal connected to
the AEP instance. It requires an installed `psql` command (OS X: run `brew
install postgres` for installation). If the command is not on the path or if a
specific version is required then use the flag `--command=/usr/local/bin/psql`

The flag `--print` just prints the command with all parameters. Use the command
`aepctl get connection` for a structured output of all parameters.

# Trigger Scheduled Query Run

The command `trigger run` triggers the immediate execution of one or multiple
scheduled queries.

```terminal
aepctl trigger run 907075e95bf479ec0a495c73_68b9c64d-0dde-4db5-b9c6-4d0ddebdb5a7_my_scheduled_querywvo1ozznm5_bsngzg
```

In case of no errors the command returns without any output. Use the flag
`--response` to show the response from the server.

If you pass multiple scheduled query IDs and an error occurs then the command
will stop the execution. Use the flag `--ignore` to execute the command for all
IDs ignoring any errors.

# Update Query Template

The `update` command replaces an existing query template with the passed
payload, hence all fields have to be provided. See [Create Query
Template](#Create-Query-Template) for the exact JSON payload.


This payload can be provided in a file, e.g. `query-template.json` in the folder
`examples/update`:


```terminal
aepctl update template --id 169a0264-9946-4cf2-af41-231a55cc6d46 examples/update/query-template.json
```

In heredoc:

```terminal
aepctl update template --id 169a0264-9946-4cf2-af41-231a55cc6d46 << EOF
{
  "sql": "SELECT $key from $key1 where $key > $key2;",
  "queryParameters": {
    "key": "value",
    "key1": "value1",
    "key2": "value2"
  },
  "name": "My updated query template"
}
EOF
```

In case of no errors the command returns without any output. Use the flag
`--response` to show the response from the server.

# Update Scheduled Query

The `update` command supports the replacement of selected values. The payload
body contains an array with JSON objects having the attributes `op` (mandatory
operator), `path` (optional path to attribute) and the new value `value`
(optional).

This payload can be provided in a file, e.g. `schedule.json` in the folder
`examples/update`:

```terminal
aepctl update schedule --id 907075e95bf479ec0a495c73_68b9c64d-0dde-4db5-b9c6-4d0ddebdb5a7_my_scheduled_querywvo1ozznm5_bsngzg examples/update/schedule.json
```

In heredoc:

```terminal
aepctl update schedule --id 907075e95bf479ec0a495c73_68b9c64d-0dde-4db5-b9c6-4d0ddebdb5a7_my_scheduled_querywvo1ozznm5_bsngzg << EOF
{
    "body": [
        {
            "op": "replace",
            "path": "/state",
            "value": "disable"
        }
    ]
}
EOF
```

In case of no errors the command returns without any output. Use the flag
`--response` to show the response from the server.