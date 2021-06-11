# Schema Registry

The Schema Registry manages the schema library of the AEP. The command
uses the following pattern:

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

The following verbs are supported by Schema Registry:

* `audit` ([Audit](#Audit))
* `delete` ([Delete](#Delete))
* `export` ([Export](#Export))
* `get` ([Get Resource](#Get-Resource) and [Get Stats](#Get-Stats))
* `import` ([Import](#Import))
* `ls` or `list`([List](#List))

# Audit
The `audit`command returns the audit log for the passed ID.

```terminal
aepctl get resource_id
```

## Views

The default view shows a table with the columns TIME, USER, VERSION and ACTION:

```terminal
aepctl audit https://ns.adobe.com/experienceplatform/schemas/f1fr70ac3d2b8f3fac09e542cf5f444fa8c1ea247d44897b

TIME                USER                             VERSION ACTION
05-24-2021 02:00:31 0D5858275FF35C910A495FC9@AdobeID 1.5     add
                                                             remove
05-24-2021 01:59:13 0D5858275FF35C910A495FC9@AdobeID 1.4     add
05-24-2021 01:59:12 0D5858275FF35C910A495FC9@AdobeID 1.3     add
                                                             add
05-24-2021 01:43:43 0D5858275FF35C910A495FC9@AdobeID 1.2     remove
05-24-2021 01:42:18 0D5858275FF35C910A495FC9@AdobeID 1.1     update
                                                             add
```

The flag `-o wide` adds the column VALUE, showing the JSON representation of the change:

```terminal
aepctl audit https://ns.adobe.com/experienceplatform/schemas/f1fr70ac3d2b8f3fac09e542cf5f444fa8c1ea247d44897b -o wide

TIME                USER                             VERSION ACTION PATH                  VALUE
05-24-2021 02:00:31 0D5858275FF35C910A495FC9@AdobeID 1.5     add    /allOf/3              {
                                                                                            "meta:xdmType": "object",
                                                                                            "type": "object",
                                                                                            "＄ref": "https://ns.adobe.com/xdm/context/profile-person-details"
                                                                                          }
                                                             remove /meta:usageCount      0
05-24-2021 01:59:13 0D5858275FF35C910A495FC9@AdobeID 1.4     add    /meta:immutableTags/0 "union"
05-24-2021 01:59:12 0D5858275FF35C910A495FC9@AdobeID 1.3     add    /meta:immutableTags   []
                                                             add    /meta:usageCount      0
05-24-2021 01:43:43 0D5858275FF35C910A495FC9@AdobeID 1.2     remove /meta:usageCount      0
05-24-2021 01:42:18 0D5858275FF35C910A495FC9@AdobeID 1.1     update /title                "My Profiles"
                                                             add    /meta:usageCount      0
```

See [Output](output.md) for other output formats.

# Delete

The `delete` command deletes the resources with the passed IDs.
It supports the following resources:

* Class
```terminal
aepctl delete class class_id
```
* Data type
```terminal
aepctl delete datatype datatype_id
```
* Descriptor
```terminal
aepctl delete descriptor descriptor_id
```
* Field group
```terminal
aepctl delete fieldgroup fieldgroup_id
```
* Schemas
```terminal
aepctl delete schema schema_id
```

The following example deletes the schema with the ID `https://ns.adobe.com/experienceplatform/schemas/349c2e69db34a49fcbde727bc249540ea8141b12cdfa8d0` :

```terminal
aepctl delete schema https://ns.adobe.com/experienceplatform/schemas/349c2e69db34a49fcbde727bc249540ea8141b12cdfa8d0
```

This command supports multiple arguments:
```terminal
aepctl delete class class_id1 class_id2 class_idn
```
Alternative:
```terminal
aepctl delete classes class_id1 class_id2 class_idn
```

If the command has been executed successfully then it just returns with no
output. Otherwise, it will return an error message.

# Export

The `export` command exports the resource with the passed ID including
dependencies. In combination with `import` it is possible to copy resources from
one account or sandbox to another. Predefined resources by Adobe cannot be
exported.

A simple export to standard out looks like this:
```terminal
aepctl export https://ns.adobe.com/experienceplatform/schemas/b1032dfc8f8598c6c59884028da0cbde649d9c4dd00e7e88
```

Redirect the output to a file:
```terminal
aepctl export https://ns.adobe.com/experienceplatform/schemas/b1032dfc8f8598c6c59884028da0cbde649d9c4dd00e7e88 > file.json
```

Export to different account:

```terminal
aepctl export https://ns.adobe.com/experienceplatform/schemas/b1032dfc8f8598c6c59884028da0cbde649d9c4dd00e7e88 | aepctl import --config name_of_configuration
```

The flag `--config` provides the name of a second configuration. Use the command
`aepctl configure --config name_of_configuration` to create and manage
additional configurations.

# Get Resource

The command `get RESOURCE` returns a resource by ID and uses the following pattern:

```terminal
aepctl get RESOURCE RESOURCE_ID
```

It supports the following resources:

* Behavior
```terminal
 aepctl get behavior https://ns.adobe.com/xdm/data/adhoc
 ```
* Class
```terminal
aepctl get class https://ns.adobe.com/xdm/context/profile
```
* Field group
```terminal
aepctl get fieldgroup https://ns.adobe.com/xdm/context/identitymap
```
* Data types
```terminal
aepctl get datatype https://ns.adobe.com/xdm/context/product
```
* Descriptor
```terminal
aepctl get descriptor 397dc205d9064f3fd9b06ee53355f76b238f95ae4fbf8752
```
* Union
```terminal
aepctl get union https://ns.adobe.com/xdm/context/profile__union
```

## Views

The default view shows a table with the columns STRUCTURE, TYPE and XDM Types.
Following tables explain the symbols and names used in the structure tree.

|Tree Item | Meaning |
|-----|--------------|
|`├─>`| References a definition. IDs starting with `#` are defined in the same JSON object |
|`├──` or `└──`| Links a child node. |
|`*`| Indicates an array. |
|`+`| Indicates additional properties. |

Some information is stored in dedicated meta-nodes. Please visit the provided
links taken from [Understanding JSON
Schema](http://json-schema.org/understanding-json-schema/index.html) for more
information :

|Tree Item | Meaning | Link |
|-----|--------------|------|
|`├─ definitions`| Contains the local definitions of embedded resources. | [Schema Compositions](http://json-schema.org/understanding-json-schema/reference/combining.html?highlight=definitions)
|`├─ extends`| Lists all extended resources. | [Extensions](https://github.com/adobe/xdm/blob/master/docs/extensions.md) |
|`├─ required`| Shows all required properties. | [Required Properties](https://json-schema.org/understanding-json-schema/reference/object.html?highlight=required#required-properties) |

This following example gets a data type with the ID
https://ns.adobe.com/xdm/common/geo. The STRUCTURE column shows a tree of the
definition:

```terminal
aepctl get datatype https://ns.adobe.com/xdm/common/geo

STRUCTURE                                TYPE    XDM TYPE
Geo                                      object  object
├─> http://schema.org/GeoCoordinates     object  object
├─> #/definitions/geo                    object  object
├── definitions
│   └── geo
│       ├── city                         string  string
│       ├── countryCode                  string  string
│       ├── dmaID                        integer int
│       ├── msaID                        integer int
│       ├── postalCode                   string  string
│       └── stateProvince                string  string
└── extends
    └── http://schema.org/GeoCoordinates
```

The flag `--full` resolves all references to a complete definition. The
same example with that flag looks like this:

```terminal
bin/aepctl get datatype https://ns.adobe.com/xdm/common/geo --full
STRUCTURE           TYPE    XDM TYPE
Geo                 object  object
├── _id             string  string
├── _schema         object  object
│   ├── description string  string
│   ├── elevation   number  number
│   ├── latitude    number  number
│   └── longitude   number  number
├── city            string  string
├── countryCode     string  string
├── dmaID           integer int
├── msaID           integer int
├── postalCode      string  string
└── stateProvince   string  string
```

The flag `-o wide` adds more information to each node:

```terminal
aepctl get datatype https://ns.adobe.com/xdm/common/geo -o wide

STRUCTURE                                TYPE          XDM TYPE
Geo                                      object        object
├─> http://schema.org/GeoCoordinates     object        object
├─> #/definitions/geo                    object        object
├── definitions
│   └── geo
│       ├── city
│       │                                description   The name of the city.
│       │                                meta:xdmField xdm:city
│       │                                meta:xdmType  string
│       │                                title         City
│       │                                type          string
│       ├── countryCode
│       │                                description   The two-character [ISO 3166-1 alpha-2](https://datahub.io/core/country-list) code for the country.
│       │                                meta:xdmField xdm:countryCode
│       │                                meta:xdmType  string
│       │                                pattern       ^[A-Z]{2}$
│       │                                title         Country code
│       │                                type          string
│       ├── dmaID
│       │                                description   The Nielsen media research designated market area.
│       │                                meta:xdmField xdm:dmaID
│       │                                meta:xdmType  int
│       │                                title         Designated market area
│       │                                type          integer
│       ├── msaID
│       │                                description   The metropolitan statistical area in the United States where the observation occurred.
│       │                                meta:xdmField xdm:msaID
│       │                                meta:xdmType  int
│       │                                title         Metropolitan statistical area
│       │                                type          integer
│       ├── postalCode
│       │                                description   The postal code of the location. Postal codes are not available for all countries. In some countries, this will only contain part of the postal code.
│       │                                meta:xdmField xdm:postalCode
│       │                                meta:xdmType  string
│       │                                title         Postal code
│       │                                type          string
│       └── stateProvince
│                                        description   The state, or province portion of the observation. The format follows the [ISO 3166-2 (country and subdivision)][http://www.unece.org/cefact/locode/subdivisions.html] standard.
│                                        examples[0]   US-CA
│                                        examples[1]   DE-BB
│                                        examples[2]   JP-13
│                                        meta:xdmField xdm:stateProvince
│                                        meta:xdmType  string
│                                        pattern       ([A-Z]{2}-[A-Z0-9]{1,3}|)
│                                        title         State or province
│                                        type          string
└── extends
    └── http://schema.org/GeoCoordinates
```

See [Output](output.md) for other output formats.

# Get Sample

The command `get sample RESOURCE_ID` generates some example data in JSON format for the passed `RESOURCE_ID`:

```terminal
aepctl get sample https://ns.adobe.com/xdm/common/geo
{
  "_id": "/uri-reference",
  "_schema": {
    "description": "string",
    "elevation": 16641.31,
    "latitude": -10.74,
    "longitude": -119.82
  },
  "city": "string",
  "countryCode": "US",
  "dmaID": 952,
  "msaID": 7700,
  "postalCode": "string",
  "stateProvince": "US-CA"
}
```

# Get Stats

The command `get stats` retrieves the tenant ID with some additional information:

```terminal
aepctl get stats

ORG                               TENANT             # SCHEMAS # MIXINS # DATATYPES # CLASSES # UNIONS
B06A75B93BF479EC1A495A73@AdobeOrg experienceplatform 46        21       25          3         6
```

See [Output](output.md) for other output formats.

Implements [GET /stats](https://www.adobe.io/apis/experienceplatform/home/api-reference.html#/Stats/ims_org_stats)

# Import

The `import` command imports the passed resources, usually the results of the [`export`](#Export) command.

Import a file:
```terminal
aepctl import file.json
```

If the command has been executed successfully then it just returns without
output. Otherwise, it will return an error message.

Import multiple files:
```terminal
aepctl import file1.json file2.json filen.json
```

Import from `export`:
```terminal
aepctl export https://ns.adobe.com/experienceplatform/schemas/b1032dfc8f8598c6c59884028da0cbde649d9c4dd00e7e88 | aepctl import --config name_of_configuration
```

The flag `--config` provides the name of a second configuration. Use the command
`aepctl configure --config name_of_configuration` to create and manage
additional configurations.


# List

The `ls` command returns a list of resources, either predefined, custom or both.
It supports the following resources:

* Behaviors
```terminal
 aepctl ls behaviors
 ```
* Classes
```terminal
aepctl ls classes
```
* Data types
```terminal
aepctl ls datatypes
```
* Descriptors
```terminal
aepctl ls descriptors
```
* Field groups
```terminal
aepctl ls fieldgroups
```
* Schemas
```terminal
aepctl ls schemas
```

* Unions
```terminal
aepctl ls unions
```

For the resources class, field group, data type, descriptor and schema `aepctl`
returns all custom definitions by default. In order to get the predefined
resources use the flag `--predefined` and to show both with the flag `--all`.

List all resources, predefined and custom, for example:

```terminal
 aepctl ls classes --all
 ```

List predefined resources, for example:

```terminal
 aepctl ls schemas --predefined
 ```

Behaviors are always predefined and unions can only be custom. Thus, both commands don't support
either `--predefined` nor `--all`.

## Views

The default view shows a table with the columns ID, TITLE and VERSION:

```terminal
aepctl ls classes --predefined -o wide
ID                                                                               TITLE                                            VERSION
https://ns.adobe.com/experience/journeyOrchestration/stepEvents/journey          Journey Orchestration Class                      1.22.3
https://ns.adobe.com/experience/journeyOrchestration/stepEvents/journeyStepEvent Journey Step Event                               1.22.3
https://ns.adobe.com/experience/decisioning/option                               Decision Option                                  1.22.3
…
```

The flag `--show` will add more information about the referenced resources,
definitions and extensions. Please take a look at [Get Resource](#Get-Resource)
for more information about the tree structure.

```terminal
aepctl ls classes --predefined --show
STRUCTURE                                                                     TYPE    XDM TYPE
Journey Orchestration Class                                                   object  object
├─> https://ns.adobe.com/xdm/data/record                                      object  object
├─> #/definitions/journeyClass                                                object  object
├── definitions
│   └── journeyClass
└── extends
    └── https://ns.adobe.com/xdm/data/record
Journey Step Event                                                            object  object
├─> https://ns.adobe.com/xdm/data/time-series                                 object  object
├─> #/definitions/journeyStepEventClass                                       object  object
├── definitions
│   └── journeyStepEventClass
└── extends
    └── https://ns.adobe.com/xdm/data/time-series
…
```

See [Output](output.md) for other output formats.

## Flags
|Flag | Type | Default | Example | Usage |
|-----|------|---------|---------|-------|
| --limit | uint | | 2 | Limits the number of returned results per request. Should be used in combination with `-o json` or `-o raw`. Otherwise, it just leads to more requests. |
| --orderby | string |  | title | Sorts the response by specified fields (separated by \",\"). |
| --start | string |  | 1607575965330 | Offests the start of returned results and is used for paging. |
