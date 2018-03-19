Public API for Concordances (public-concordances-api)
===

__Provides a public API for Concordances stored in a Neo4J graph database__

This Go app provides REST endpoint for the retrieval of Concordance data given one or more parameters. This app
directly depends upon a connection to a neo4j database.

## Installation & running locally
* `go get -u github.com/Financial-Times/public-concordances-api`
* `cd $GOPATH/src/github.com/Financial-Times/public-concordances-api`
* `go test ./...`
* `go install`
* `$GOPATH/bin/public-concordances-api --neo-url={neo4jUrl} --port={port}`

_Both arguments are optional.
--neo-url defaults to http://localhost:7474/db/data, which is the out of box url for a local neo4j instance.
--port defaults to 8080._
 
Alternatively, use socks proxy to forward to tunnel url via ssh and hit a cluster's neo4j directly.

* `ssh -L 1234:localhost:8080 core@{cluster}.ft.com` Maps your local port 1234 to the cluster's port 8080(vulcan)
* `$GOPATH/bin/public-concordances-api --neo-url=http://localhost:{your port}/__neo4j-{red/blue}/db/data --port={port}` 

## API Endpoints
Based on the following [google doc](https://docs.google.com/a/ft.com/document/d/1onyyb-XoByB00RQNZvjNoL_IsO_eHKe-vOpUuAVHyJE)

    - `GET /concordances?conceptId={thingUri}` - Returns a list of all identifiers for given concept
    - `GET /concordances?conceptId={thingUri}&conceptId={thingUri}...` - Returns a list of all identifiers for each concept provided   
    - `GET /concordances?authority={identifierUri}&identifierValue{identifierValue}` - Returns the apiUrl that matches the corresponding identifier 
    - `GET /concordances?authority={identifierUri}&idenifierValue={identifierValue}&idenifierValue={identifierValue}...` -Returns a list of all apiUrl's for the corresponding identifiers

## Admin endpoints (CoCo) 

  - `https://{host}/__public-concordances-api/__health`
  - `https://{host}/__public-concordances-api/__build-info`
  - `https://{host}/__public-concordances-api/__gtg`

Note: All API endpoints in CoCo require Authentication.
See [service panic guide] (https://sites.google.com/a/ft.com/universal-publishing/ops-guides/concordances-read) on how to access cluster credentials.  


## Error handling
The service expects at least 1 conceptId or (authority + identifierValue pair) parameter and will respond with an Error HTTP status code if these are not provided.
The service will respond with Error HTTP codes if both a conceptId is presented with an authority parameter or if an identifierValue is presented without the authority parameter.
The service will never respond with Error HTTP status codes if none of the conceptId's or identifierValues are present in concordance,
instead it will return an empty array of Concepts or Identifiers.

Sourcecode is [here] (https://github.com/Financial-Times/public-concordances-api)
Run Book is [here] (https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/public-concordances-api)

