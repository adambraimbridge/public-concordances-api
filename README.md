# Public API for Concordances (public-concordances-api)
__Provides a public API for Concordances stored in a Neo4J graph database__

## Build & deployment etc:
_NB You will need to tag a commit in order to build, since the UI asks for a tag to build / deploy_
* [Jenkins view](http://ftjen10085-lvpr-uk-p:8181/view/JOBS-public-concordances-api/)
* [Build and publish to forge](http://ftjen10085-lvpr-uk-p:8181/job/public-concordances-api-build)
* [Deploy to Test](http://ftjen10085-lvpr-uk-p:8181/view/JOBS-public-concordances-api/job/public-concordances-api-deploy-test/)
* [Deploy to Prod](http://ftjen10085-lvpr-uk-p:8181/view/JOBS-public-concordances-api/job/public-concordances-api-deploy-to-prod/)

## Installation & running locally
* `go get -u github.com/Financial-Times/public-concordances-api`
* `cd $GOPATH/src/github.com/Financial-Times/public-concordances-api`
* `go test ./...`
* `go install`
* `$GOPATH/bin/public-concordances-api --neo-url={neo4jUrl} --port={port} --log-level={DEBUG|INFO|WARN|ERROR}`
_Both arguments are optional.
--neo-url defaults to http://localhost:7474/db/data, which is the out of box url for a local neo4j instance.
--port defaults to 8080._
* `curl http://localhost:8080/concordances/143ba45c-2fb3-35bc-b227-a6ed80b5c517 | json_pp`
Or using [httpie](https://github.com/jkbrzt/httpie)
* `http GET http://localhost:8080/concordances/143ba45c-2fb3-35bc-b227-a6ed80b5c517`

## API definition
Based on the following [google doc](https://docs.google.com/a/ft.com/document/d/1onyyb-XoByB00RQNZvjNoL_IsO_eHKe-vOpUuAVHyJE)

## Healthchecks
Healthchecks: [http://localhost:8080/__health](http://localhost:8080/__health)
