# Identity Service

The Identity Service handles the multiple identities for a user profile. Data
from different channels comes with its own identities, e.g. data from CRM
contains hashed email addresses (first party PII) while data from web analytics
comes with cookie IDs (device ID). Each channel uses its specific namespace for
IDs in order to address the kind and source. Data from Adobe solutions usually
uses the namespace code `ECID` (first party cookie). The mentioned CRM data
would use the namespace `Email_LC_SHA256` (SHA256 with lower-cased email
addresses). Please visit [Identity Service
Overview](https://experienceleague.adobe.com/docs/experience-platform/identity/home.html?lang=en)
for more information.

The command uses the following pattern:

```terminal
aepctl (verb) (noun)
```

E.g., if you want to list the available namespaces

```terminal
aepctl list namespaces
```

Some commands support a preferred short notation, e.g. `ls` for `list`:

```terminal
aepctl ls namespaces
```

The following verbs are supported by Identity Service:

* `create` ([Create Namespaces](#Create-Namespaces))
* `get` ([Get Namespace](#Get-Namespace), [Get XID](#Get-XID) and [Get Cluster of IDs](#Get-Cluster-of-IDs))
* `import` ([Import](#Import))
* `ls` or `list`([List Namespaces](#List-Namespaces))
* `update` ([Update Namespaces](#Update-Namespaces))

# Create Namespaces

The `create` command creates a new namespace and requires a payload in JSON
format with the details. 

`name`, `code` and `idType` are mandatory, `description` is optional. For more
information see the official
[documentation](https://www.adobe.io/apis/experienceplatform/home/api-reference.html#/Identity_Namespace/CreateNameSpace).

```json
{
  "name": "My Loyalty Member",
  "code": "MyLoyalty",
  "description": "My Loyalty Program Member ID",
  "idType": "Cross_device"
}
```

This payload can be provided in a file, e.g. `namespace.json` in the folder
`examples/create`:

```terminal
aepctl create namespace examples/create/namespace.json
```

If no file name is provided then `aepctl` reads from the standard input stdin,
hence heredoc is supported, too:

```terminal
aepctl create namespace << EOF
{
  "name": "My Loyalty Member",
  "code": "MyLoyalty",
  "description": "My Loyalty Program Member ID",
  "idType": "Cross_device"
}
EOF
```

In case of error a message will be printed out, otherwise no output will be
produced. Use the flag `--respponse` to print the returned information.

It is also possible to create multiple namespaces with one call:

```terminal
aepctl create namespaces namespace_loyalty.json namespace_custom.json namespace_hashed.json
```

Use the flag `--ignore` to skip over errors during the execution. Otherwise,
`aepctl` would stop the execution with the first error.

# Get Namespace

The `get` command returns information for the passed namespace ID:

```terminal
aepctl get namespace 4
```

will return the following for the ECID namespace with the ID 4 :

```terminal
ID CODE NAME ID TYPE DESCRIPTION
4  ECID ECID COOKIE  Adobe Experience Cloud ID
```

The flag `-o wide` adds the columns STATUS, TYPE, UPDATED and DESCRIPTION:

```terminal
aepctl get namespace 4 -o wide
ID CODE NAME ID TYPE STATUS TYPE     UPDATED             DESCRIPTION
4  ECID ECID COOKIE  ACTIVE Standard 04 Mar 19 09:33 CET Adobe Experience Cloud ID
```

See [Output](output.md) for other output formats. 

## Shared Namespaces

Provide an IMS Organization ID with the flag `--ims-org ID` to get the shared namespace with a different organization:

```terminal
aepctl get namespace 4 --ims-org 792D3C635C5CDF980A395CB2@AdobeOrg
```

Or

```terminal
aepctl get namespace --ims-org 792D3C635C5CDF980A395CB2@AdobeOrg 4
```

# Get XID

The AEP supports composite IDs, e.g. the namespace ECID in combination with the
ID 67504705834917073766149225290256685656, which can be represented by a
single-value ID with the name XID, e.g. A28eOco1-QqGQERvuJjKVoEd.

It is possible to calculate the single-value XID, either from the namespace code
or namespace ID in combination with the related ID. Use the following command to
retrieve the XID, the default namespace code is ECID:

```terminal
aepctl get xid 67504705834917073766149225290256685656
```

With the following output. See [Output](output.md) for other output formats. 

```terminal
XID
Bm_wQb_rOn2uXVgamU3Sb1ZV
```

Provide either a namespace code, e.g. `Email`, or the namespace ID, which is 6 for `Email`. Use `aepctl ls namespaces` to list all namespaces. The following examples return the same results:

```terminal
aepctl get xid --namespace Email max@mustermann.de
```

With the namespace code:

```terminal
aepctl get xid --ns-id 6 max@mustermann.de
```

See [Output](output.md) for other output formats. 

# Get Cluster of IDs

The identity service builds clusters of related identities. Use either the
single-value XID or the combination of namespace and ID.

The following command returns the related XIDs:

```terminal
aepctl get ids A28eOco1-QqGQERvuJjKVoEe
```

The output could look like this:

```terminal
SRC                      XID
A28eOco1-QqGQERvuJjKVoEe A28eOco1-QqGQERvuJjKVoEe
                         A2_9t-I0hvDrL4hk_uuc6Rde
                         A2_wQb_rOn2uXVkamU3Sb1ZYV
                         BUF8bWJ1bmdlbnMrgkUW9c79bzf6Bp2ht616y35dyiX
                         BkFuK4TuGkp9ciuCRRb1LTgy
                         CkF9r3spqD0SqaqydEBYbaKRveKvkuOOOzjjytRMyS0ri4wK4B
                         A2_MthJFcLjbZK1mvuK3oLVe
```

The following command returns the related namespace code + ID combinations
(using the alias cluster for ids):

```terminal
aepctl get cluster --namespace ECID 67504705834917073766149225290256685657
```

The output could look like this:

```terminal
NAMESPACE       ID
ECID            46494157307721956142054018333203263300
Email_LC_SHA256 7b29a83d12a9aab27440586da291bde292e38e3b38e3cad44cc92d2b8b8c0ae1
Phone           +49110112
Email           max@mustermann.de
ECID            67504705834917073766149225290256685657
```

The following command returns the related namespace ID + ID combinations:

```terminal
aepctl get ids --ns-id 4 67504705834917073766149225290256685657
```

The output could look like this:

```terminal
NAMESPACE ID
4         46494157307721956142054018333203263300
11        7b29a83d12a9aab27440586da291bde292e38e3b38e3cad44cc92d2b8b8c0ae1
7         +49110112
6         max@mustermann.de
4         67504705834917073766149225290256685657
```

See [Output](output.md) for other output formats. 

## Get multiple Clusters of IDs

The following command returns multiple clusters. It uses the plural `clusters`
and executes one HTTP POST command:

```terminal
aepctl get clusters A28eOco1-QqGQERvuJjKVoEc A2_MthJFcLjbZK1mvuK3oLVd
```

Returns:

```terminal
SRC                      XID
A28eOco1-QqGQERvuJjKVoEc A28eOco1-QqGQERvuJjKVoEc
                         BUF8bWJ1bmdlbnMrgkUW9c79bzf6Bp2ht616y35dyiY
                         BkFuK4TuGkp9ciuCRRb1LTgz
                         CkF9r3spqD0SqaqydEBYbaKRveKvkuOOOzjjytRMyS0ri4wK4A
                         A2_MthJFcLjbZK1mvuK3oLVd
A2_MthJFcLjbZK1mvuK3oLVd A28eOco1-QqGQERvuJjKVoEc
                         BUF8bWJ1bmdlbnMrgkUW9c79bzf6Bp2ht616y35dyiY
                         BkFuK4TuGkp9ciuCRRb1LTgz
                         CkF9r3spqD0SqaqydEBYbaKRveKvkuOOOzjjytRMyS0ri4wK4A
                         A2_MthJFcLjbZK1mvuK3oLVd
```

The flag `--namespace` is not supported, please use the namespace ID:

```terminal
aepctl get clusters --ns-id 4  67504705834917073766149225290256685657 --ns-id 6 max@mustermann.de
```

```terminal
SRC                                    NAMESPACE ID
67504705834917073766149225290256685657 4         46494157307721956142054018333203263300
                                       11        7b29a83d12a9aab27440586da291bde292e38e3b38e3cad44cc92d2b8b8c0ae0
                                       6         max@mustermann.de
                                       4         67504705834917073766149225290256685657
                                       7         +49110112
max@mustermann.de                      4         46494157307721956142054018333203263300
                                       11        7b29a83d12a9aab27440586da291bde292e38e3b38e3cad44cc92d2b8b8c0ae0
                                       6         max@mustermann.de
                                       4         67504705834917073766149225290256685657
                                       7         +49110112
```

The last `--ns-id` sets the namespace for all following IDs. The example will return 2 clusters:

```terminal
aepctl get clusters --ns-id 4  67504705834917073766149225290256685657 46494157307721956142054018333203263300
```

# List Namespaces

The `ls` command returns a list of all declared namespaces:

```terminal
aepctl ls namespaces
```

The default view shows a table with the columns ID, CODE, NAME and ID TYPE

```terminal
ID       CODE              NAME                         ID TYPE
0        CORE              CORE                         COOKIE
4        ECID              ECID                         COOKIE
411      AdCloud           AdCloud                      COOKIE
11       Email_LC_SHA256   Emails (SHA256, lowercased)  Email
17       Phone_E.164       Phone (E.164)                Phone
…
```

The flag `-o wide` adds the columns STATUS, TYPE, UPDATED and DESCRIPTION:

```terminal
aepctl ls namespaces -o wide
ID       CODE     NAME      ID TYPE    STATUS TYPE        UPDATED              DESCRIPTION
0        CORE     CORE      COOKIE     ACTIVE Standard    04 Mar 19 09:33 CET  Adobe Audience Manger UUID
4        ECID     ECID      COOKIE     ACTIVE Standard    04 Mar 19 09:33 CET  Adobe Experience Cloud ID
411      AdCloud  AdCloud   COOKIE     ACTIVE Standard    04 Mar 19 09:33 CET  Adobe AdCloud - ID Syncing Partner
…
```

See [Output](output.md) for other output formats. 

## Shared Namespaces

Provide an IMS Organization ID with the flag `--ims-org ID` to show the shared namespaces with a different organization.

```terminal
aepctl ls namespaces --ims-org 792D3C635C5CDF980A395CB2@AdobeOrg
```

# Update Namespaces

The `update` command replaces an existing namespace with the passed payload,
hence all fields have to be provided. See [Create
Namespaces](#create-namespaces) for the exact JSON payload.


This payload can be provided in a file, e.g. `namespace.json` in the folder
`examples/update`:


```terminal
aepctl update namespace --id 10885436 examples/update/namespace.json
```

In heredoc:

```terminal
aepctl update namespace --id 10885436 << EOF
{
  "name": "My Loyalty Member",
  "code": "MyLoyalty",
  "description": "My Loyalty Program Member ID",
  "idType": "Cross_device"
}
EOF
```

In case of error a message will be printed out, otherwise no output will be
produced. Use the flag `--respponse` to print the returned information.