# Output
aepctl provides several output formats:

1. __Table__ view with fixed columns. Usually a summarized output of the
   response and in most cases the default output format.
2. __Wide__ is like Table but with more columns (or more rows).   
3. __NVP__ (Name/Value/Path) displays all values in three columns.
4. __PV__ (Path/Value) displays all values int two columns.
5. __JSON__ pretty prints the complete response in JSON. This format does not support paging.
6. __Raw__ prints the response without any formatting. This format does not support paging.

Select the desired output format with the `--output` flag or the short form
`-o`. Please use one of the following notations:

* `--output json`
* `--output=json`
* `-o json`
* `-ojson`

All have the same meaning: Print out in indented JSON.

## Table

Table is usually the default output format. It presents selected values in a fixed column format.

### Example

```terminal
aepctl get behaviors

TITLE              VERSION
Ad Hoc Schema      1.22.3
Time-series Schema 1.22.3
Record Schema      1.22.3
```

The long form with dedicated output flag is not necessary:

```terminal
aepctl get behaviors -o table
```

Some commands with hierarchical responses combine the output with a tree, but it is still a table.

```terminal
aepctl get data-type http://schema.org/GeoCoordinates

STRUCTURE                        TYPE   XDM TYPE
Geo Coordinates                  object object
├─> #/definitions/geocoordinates object object
└── definitions
    └── geocoordinates
        ├── _id                  string string
        └── _schema              object object
            ├── description      string string
            ├── elevation        number number
            ├── latitude         number number
            └── longitude        number number
```


## Wide

Is like table but with more columns. Some commands don't provide more data, thus, the output will not differ.

### Example

```terminal
aepctl get behaviors -o wide

ID                                        TITLE              VERSION
https://ns.adobe.com/xdm/data/adhoc       Ad Hoc Schema      1.22.3
https://ns.adobe.com/xdm/data/time-series Time-series Schema 1.22.3
https://ns.adobe.com/xdm/data/record      Record Schema      1.22.3
```

## NVP

Use this format to print all relevant values. Each row shows the name, value and
path of the derived JSON values. 

### Example

```terminal
aepctl get data-type http://schema.org/GeoCoordinates -o nvp

NAME                  VALUE                                            PATH
$id                   http://www.iptc.org/rating
meta:altId            _www.iptc.org.rating
meta:resourceType     datatypes
version               1.22.3
title                 Rating
type                  object
description           The rating of the show. Based on [www.iptc.or...
title                 Rating Source Link                               definitions.rating.properties._iptc4xmpExt.properties.RatingSourceLink
type                  string                                           definitions.rating.properties._iptc4xmpExt.properties.RatingSourceLink
format                uri                                              definitions.rating.properties._iptc4xmpExt.properties.RatingSourceLink
description           Link to the site and optionally the page of t... definitions.rating.properties._iptc4xmpExt.properties.RatingSourceLink
meta:xdmType          string                                           definitions.rating.properties._iptc4xmpExt.properties.RatingSourceLink
meta:xdmField         iptc4xmpExt:RatingSourceLink                     definitions.rating.properties._iptc4xmpExt.properties.RatingSourceLink
title                 Rating Value                                     definitions.rating.properties._iptc4xmpExt.properties.RatingValue
type                  string                                           definitions.rating.properties._iptc4xmpExt.properties.RatingValue
description           Rating value as issued by the rating source.     definitions.rating.properties._iptc4xmpExt.properties.RatingValue
meta:xdmType          string                                           definitions.rating.properties._iptc4xmpExt.properties.RatingValue
meta:xdmField         iptc4xmpExt:RatingValue                          definitions.rating.properties._iptc4xmpExt.properties.RatingValue
[0]                   RatingSourceLink                                 definitions.rating.properties._iptc4xmpExt.required
[1]                   RatingValue                                      definitions.rating.properties._iptc4xmpExt.required
type                  object                                           definitions.rating.properties._iptc4xmpExt
meta:xdmType          object                                           definitions.rating.properties._iptc4xmpExt
meta:xedConverted     true                                             definitions.rating.properties._iptc4xmpExt
[0]                   _iptc4xmpExt                                     definitions.rating.required
$ref                  #/definitions/rating                             allOf[0]
type                  object                                           allOf[0]
meta:xdmType          object                                              
meta:status           stable
$schema               http://json-schema.org/draft-06/schema#
repo:createdDate      1621354385078                                    meta:registryMetadata
repo:lastModifiedDate 1621354385078                                    meta:registryMetadata
eTag                  b8162bf4abb18bc4516cc3221dc5c12d583401895cfd3... meta:registryMetadata
meta:globalLibVersion 1.22.3                                           meta:registryMetadata
meta:createdDate      2020-08-10
```

`[0]` and `[1]` in the NAME column stand for the index in an array.

## PV

Use this format to print all relevant values. Each row shows the path and the value of the derived JSON values.

## Example
```terminal
aepctl get data-type http://schema.org/GeoCoordinates -o pv

PATH                                                                               VALUE
$id                                                                                http://schema.org/GeoCoordinates
meta:altId                                                                         _schema.org.GeoCoordinates
meta:resourceType                                                                  datatypes
version                                                                            1.22.3
title                                                                              Geo Coordinates
type                                                                               object
description                                                                        The geographic coordinates of a place. Based on [schema...
definitions.geocoordinates.properties._id.title                                    Coordinates ID
definitions.geocoordinates.properties._id.type                                     string
definitions.geocoordinates.properties._id.format                                   uri-reference
definitions.geocoordinates.properties._id.description                              The unique identifier of the coordinates.
definitions.geocoordinates.properties._id.meta:xdmType                             string
definitions.geocoordinates.properties._id.meta:xdmField                            @id
definitions.geocoordinates.properties._schema.properties.description.title         Description
definitions.geocoordinates.properties._schema.properties.description.type          string
definitions.geocoordinates.properties._schema.properties.description.description   A description of what the coordinates identify.
definitions.geocoordinates.properties._schema.properties.description.meta:xdmType  string
definitions.geocoordinates.properties._schema.properties.description.meta:xdmField schema:description
definitions.geocoordinates.properties._schema.properties.elevation.title           Elevation
definitions.geocoordinates.properties._schema.properties.elevation.type            number
definitions.geocoordinates.properties._schema.properties.elevation.description     The specific elevation of the defined coordinate. The v...
definitions.geocoordinates.properties._schema.properties.elevation.meta:xdmType    number
definitions.geocoordinates.properties._schema.properties.elevation.meta:xdmField   schema:elevation
definitions.geocoordinates.properties._schema.properties.latitude.title            Latitude
definitions.geocoordinates.properties._schema.properties.latitude.type             number
definitions.geocoordinates.properties._schema.properties.latitude.minimum          -90
definitions.geocoordinates.properties._schema.properties.latitude.maximum          90
definitions.geocoordinates.properties._schema.properties.latitude.description      The signed vertical coordinate of a geographic point.
definitions.geocoordinates.properties._schema.properties.latitude.meta:xdmType     number
definitions.geocoordinates.properties._schema.properties.latitude.meta:xdmField    schema:latitude
definitions.geocoordinates.properties._schema.properties.longitude.title           Longitude
definitions.geocoordinates.properties._schema.properties.longitude.type            number
definitions.geocoordinates.properties._schema.properties.longitude.minimum         -180
definitions.geocoordinates.properties._schema.properties.longitude.maximum         180
definitions.geocoordinates.properties._schema.properties.longitude.description     The signed horizontal coordinate of a geographic point.
definitions.geocoordinates.properties._schema.properties.longitude.meta:xdmType    number
definitions.geocoordinates.properties._schema.properties.longitude.meta:xdmField   schema:longitude
definitions.geocoordinates.properties._schema.type                                 object
definitions.geocoordinates.properties._schema.meta:xdmType                         object
definitions.geocoordinates.properties._schema.meta:xedConverted                    true
allOf[0].$ref                                                                      #/definitions/geocoordinates
allOf[0].type                                                                      object
allOf[0].meta:xdmType                                                              object
meta:extensible                                                                    true
meta:xdmType                                                                       object
meta:status                                                                        stable
$schema                                                                            http://json-schema.org/draft-06/schema#
meta:registryMetadata.repo:createdDate                                             1621354385104
meta:registryMetadata.repo:lastModifiedDate                                        1621354385104
meta:registryMetadata.eTag                                                         8c81a430f048e7efbde01d4f20be9342fd8baae7d1a649244caa82d7b9515727
meta:registryMetadata.meta:globalLibVersion                                        1.22.3
meta:registryMetadata.meta:usageCount                                              3
meta:createdDate                                                                   2020-08-10
```

## JSON

Prints out the raw JSON response with indention. This format doesn't support
paging.

### Example

```terminal
 aepctl get sandboxes types -o json

{
  "sandboxTypes": [
    "development",
    "production"
  ]
}
```

## Raw

Prints out the response without any formatting. This format does not support
paging.

### Example

```terminal
 aepctl get sandboxes types -o json

{"sandboxTypes":["development","production"]}
```
