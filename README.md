# aepctl

aepctl is a command line tool for the [Adobe Experience
Platform](https://experienceleague.adobe.com/docs/experience-platform/landing/home.html)
implementing a part of the [REST
API](https://www.adobe.io/apis/experienceplatform/home/api-reference.html).

This is the initial state of this project and no release is available.

# Overview

aepctl is a complement to the existing web interface and has been developed for
advanced users as well as developers. In combination with activated syntax
completion, aepctl accelerates the execution of repeating tasks, prototyping and
learning of the APIs.

# Status of Implementation

At the moment the following APIs are implemented:

* [Access Control API](https://www.adobe.io/apis/experienceplatform/home/api-reference.html#!acpdr/swagger-specs/access-control.yaml)
* [Offer Decisioning](https://experienceleague.adobe.com/docs/offer-decisioning/using/api-reference/getting-started.html?lang=en#api-reference)

Planed APIs:

* Catalog Service API

# Installation
At the moment you have to build the binary on your own:

    make build

# License
aepctl is released under the Apache 2.0 license.