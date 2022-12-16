# Security

This manual provides some hints how to improve the security of the deployment.


## 1. Use Containers

The container image of the IAT endpoint is provided on [DockerHub (external URL)](https://hub.docker.com/r/jonasprimbs/oidc-e2ea-server).

Follow the instructions there to deploy the container or use the Docker Compose file [here](./docker-compose.yaml) to run a predefined composition with [Keycloak (external URL)](https://www.keycloak.org/) as OpenID Provider.
