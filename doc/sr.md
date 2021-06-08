# Schema Registry

The Schema Registry manages the schema library of the AEP. The command
invocation uses the following pattern:

```terminal
aepctl (verb) (noun)
```

E.g., if you want to list the custom schemas 

```terminal
aepctl list schemas
```

Some commands support a preferred short notation, e.g. `ls` for `list`:

```terminal
aepctl ls schemas
```

Flags change and control the execution, e.g. use `--predefined` to show all
resources provided by Adobe:

```terminal
aepctl ls schemas --predefined
```

## List Resources

The Schema Registry `ls` command supports the following resources:

* Behaviors
```terminal
 aepctl ls behaviors
 ```
* Classes
```terminal
aepctl ls classes
```
* Field Groups
```terminal
aepctl ls fieldgroups
```
* Data types
```terminal
aepctl ls data-types
```
* Descriptors
```terminal
aepctl ls descriptors
```
* Unions
```terminal
aepctl ls unions
```

The resources classes, field groups, 

## Stats

Get the tenant ID with some additional information:

```terminal
aepctl get stats

ORG                               TENANT             # SCHEMAS # MIXINS # DATATYPES # CLASSES # UNIONS
B06A75B93BF479EC1A495A73@AdobeOrg experienceplatform 46        21       25          3         6
```

See [Output](output.md) for other output formats.

Implements [GET /stats](https://www.adobe.io/apis/experienceplatform/home/api-reference.html#/Stats/ims_org_stats)

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

