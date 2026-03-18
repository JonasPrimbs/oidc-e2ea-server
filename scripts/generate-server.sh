#!/bin/bash

docker run --rm -v ../:/local/ openapitools/openapi-generator-cli:v7.20.0 generate -i /local/docs/swagger.yaml -g go-server -o /local/src
