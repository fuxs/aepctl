# Schema Registry

The Schema Regsitry manages the schema library of the AEP. The general aeptctl
command invocation uses the following pattern:

```terminal
aepctl (verb) sr (noun)
```

E.g., if you want to list the custom schems 

```terminal
aepctl get sr schemas
```

## Stats

The default table view returns the tenant ID with some additional information.

Command
```terminal
aepctl get sr stats
```

Output
```terminal
ORG                               TENANT             # SCHEMAS # MIXINS # DATATYPES # CLASSES # UNIONS
B06A75B93BF479EC1A495A73@AdobeOrg experienceplatform 46        21       25          3         6
```

### Show All Data

Use the `--output=nvp` flag for the complete response in a generice Name/Value/Path view 

```terminal
aepctl get sr stats --output=nvp
```
or `--output=json` for pretty printed JSON.

```terminal
aepctl get sr stats --output=json
```
### REST API

Implements the [GET /stats](https://www.adobe.io/apis/experienceplatform/home/api-reference.html#/Stats/ims_org_stats) command.

## Schemas

The default table view returns a list of custom schema titles and versions.

Command
```terminal
aepctl get sr schemas
```

Output
```terminal
TITLE                             VERSION
My Destinations                   1.20.4
My Destinations Segment Mapping   1.20.4
My Destinations Namespace Mapping 1.20.4
```

The core schemas provided by Adobe can be retrieved with the argument `global`

```terminal
aepctl get sr schemas global
```

By default this command returns a summary of each resourc. For full resouce
descriptions use the flag `--full`. This will always return a Name/Value/Path
visualization due to the different structures of the returned respnses.

```terminal
# custom schemas
aepctl get sr schemas --full
# predefined schemas
aepctl get sr schemas global --full
```

### Flags

|Flag | Type | Default | Example | Usage |
|-----|------|---------|---------|-------|
| --limit | uint | | 2 | Limit the number of results to be displayed |
| --properties | string |  | meta:abstract==true | Filter results on any top-level attribute

