# Tags


| Master Branch |
| ------------- |
| [![Build Status](https://travis-ci.com/arundo/platform-deleteme.svg?token=aUDFrwWnpp2hjLD7LeQC&branch=master)](https://travis-ci.com/arundo/platform-deleteme) |
| [![Coverage Status](https://coveralls.io/repos/github/arundo/platform-deleteme/badge.svg?branch=master&t=i0c3z7)](https://coveralls.io/github/arundo/platform-deleteme?branch=master) |

**[Table of Contents]**
- [Synopsis](#synopsis)
- [Motivation](#motivation)
- [Code Example](#code-example)
- [Endpoints](#endpoints)
- [Getting Started](#getting-started)
  - [Configuration and/or Environment Variables](#configuration-andor-environment-variables)
  - [Prerequisites](#prerequisites)
  - [Installing](#installing)
- [Running the tests](#running-the-tests)
  - [Break down into unit tests](#break-down-into-unit-tests)
  - [Break down into end to end tests](#break-down-into-end-to-end-tests)
- [Deployment](#deployment)
- [Built With](#built-with)
- [License](#license)
- [Links](#links)

This service provides CRUD endpoints for tag meta data. This is a standard CRUD service for tag meta data.

## Synopsis

Users of the Arundo platform are streaming their delete me into storage and need additional meta data for a tag. This service provides CRUD endpoints for users that need tag meta data for their tenant in Arundo's tag system.

## Motivation

It is important to provide users access additional data associated with their delete me. This service provides CRUD endpoints for tag meta data.

## Code Example

Below is a code snippet in GoLang that shows how to use the service.

```
package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func main() {

	url := "http://develop.arundo.com/v0/deletemes"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Basic <token>")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}
```

Response
```
{
    "data": [
        {
            "guid": "32646335-3561-6616-6263-643761316465",
            "description": "Random Tag",
            "valueType": "float",
            "tagId": "a_tag",
            "owner": "user@arundo.com",
            "deviceId": "EdgeDeviceA",
            "displayName": "Tag Display Name",
            "lastActivity": null,
            "insertedAt": "2018-12-18T19:53:34.447678Z",
            "updatedAt": "2018-12-18T19:53:34.447678Z",
            "lastActiveValue": ""
        },
        {
            "guid": "38623132-6233-6231-6534-336431646461",
            "description": "A Random Tag",
            "valueType": "float",
            "tagId": "a-tag",
            "owner": "user@arundo.com",
            "deviceId": "EdgeDeviceA",
            "displayName": "TagDisplayName",
            "lastActivity": "2018-12-17T19:40:17.510912Z",
            "insertedAt": "2018-12-17T19:32:17.510912Z",
            "updatedAt": "2018-12-17T19:32:17.510912Z",
            "lastActiveValue": "2"
        }
    ]
}
```

## Endpoints

+ **[<code>POST</code> /v1/deletemes](https://github.com/arundo/platform-deleteme/blob/master/main.go#L83)**
+ **[<code>GET</code> /v1/deletemes](https://github.com/arundo/platform-deleteme/blob/master/main.go#L84)**
+ **[<code>GET</code> /v1/deletemes/:id](https://github.com/arundo/platform-deleteme/blob/master/main.go#L85)**
+ **[<code>PUT</code> /v1/deletemes/:id](https://github.com/arundo/platform-deleteme/blob/master/main.go#L86)**
+ **[<code>DELETE</code> /v1/deletemes](https://github.com/arundo/platform-deleteme/blob/master/main.go#L87)**
+ **[<code>DELETE</code> /v1/deletemes/:id](https://github.com/arundo/platform-deleteme/blob/master/main.go#L88)**


+ **[<code>POST</code> /v1/ideletemes/:company](https://github.com/arundo/platform-deleteme/blob/master/main.go#L72)**
+ **[<code>GET</code> /v1/ideletemes/:company](https://github.com/arundo/platform-deleteme/blob/master/main.go#L73)**
+ **[<code>PUT</code> /v1/ideletemes/:company/:id](https://github.com/arundo/platform-deleteme/blob/master/main.go#L74)**

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project in the cloud.

### Configuration and/or Environment Variables

The location of the configuration and/or environment variables for to get this project up and running.

[env vars in the deployment.yaml](https://github.com/arundo/platform-deleteme/blob/master/.chart/templates/deployment.yaml#L37)

### Prerequisites

You will need a computer with GoLang 1.11 or compatible installed.

### Installing

A step by step series of steps that tell you have to get a development env running.

Clone the repo
```
git clone https://github.com/arundo/platform-deleteme.git
```

Create a `.env` file with the following environment variables

```
AUTH0_DOMAIN=arundo-develop.auth0.com
AUTH0_AUDIENCE=https://develop.arundo.com
GIN_MODE=debug
GORM_LOG_MODE=true
POSTGRES_SSL_MODE=require
POSTGRES_USERNAME=
POSTGRES_PASSWORD=
POSTGRES_HOSTNAME=
POSTGRES_DATABASE_NAME=
POSTGRES_PORT=5432
PORT=5000
QUERY_LIMIT=1000
TAGS_PER_PAGE=100
TAG_MONITOR_API_KEY=
```

Execute the following command, and depending on your environment variables you may or may not see any messages in the terminal:
`go run main.go`

## Running the tests

WIP
Explain how to run the automated tests for this system

### Break down into unit tests

Explain what these tests test and why

```
Give an example
```

### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

## Deployment

This repo takes advantage of the centralized [Travis CICD repo](https://github.com/arundo/cicd-travis). CICD creates both incubator and stable helm charts.

[incubator](https://github.com/arundo/charts/tree/master/incubator/fabric/platform-deleteme) helm chart
[stable](https://github.com/arundo/charts/tree/master/stable/fabric/platform-deleteme) helm chart

#### Configmap

User needs access to the kubernetes cluster the service will be running on. Create the following yaml locally, name it `config.yaml`, with the details filled out.

```
# Please edit the object below. Lines beginning with a '#' will be ignored,
# and an empty file will abort the edit. If an error occurs while saving this file will be
# reopened with the relevant failures.
#
apiVersion: v1
data:
  auth0_audience: https://develop.arundo.com
  auth0_domain: arundo-develop.auth0.com
  enable_ip_whitelist: "true"
  gin_mode: debug
  gorm_log_mode: "true"
  ip_whitelist: 127.0.0.1,127.0.0.2
  port: "5000"
  postgres_database_name: fabric_service_deletemes
  postgres_hostname:
  postgres_password:
  postgres_port: "5432"
  postgres_ssl_mode: require
  postgres_username: 
  query_limit: "100000"
  tag_monitor_api_key:
  deletemes_per_page: "100"
kind: ConfigMap
metadata:
  name: deletemes
  namespace: fabric
```

Execute the following command:
`kubectl apply -f config.yaml`

Response should be:
```
configmap "deletemes" created
```

This service now has the configmaps it needs to be installed on the kubernetes cluster.

#### Helm Install

Once the configmap and secret has been setup the next step is to install the helm chart for this service. To do this make sure the helm repo list contains the correct chart repo, then execute the following command:
```
helm upgrade --install --namespace fabric platform-deleteme arundo/platform-deleteme
```

Standard helm args apply. This service has now been installed on the kubernetes cluster.

## Built With

* [GoLang](https://golang.org/doc/) - GoLang
* [Travis-CI](https://travis-ci.com/) - Build service

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Links

* [Arundo-Docs](https://arundo-docs.gitbook.io/data-management/data/streaming-v2)
* [changelog](https://gist.github.com/alexbednarczyk/2519c770d0ff3e9d0d71bb84ce1ef056#file-changelog-md)
