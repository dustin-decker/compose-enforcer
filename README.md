# ComposeEnforcer

Compose enforcer validates that a provided docker-compose file meets restrictions set.

It can be used via CLI from the compiled binary, or you can use the public interfaces
provided. A common usecase is within a continuous deployment pipeline.

Currently it is developed against [docker compose v3.2 spec](https://docs.docker.com/compose/compose-file/) and supports:

  - volumes
  - networks
  - secrets
  - resources

Under the hood it uses exported docker compose interfaces so it should *hopefully*
be easy to maintain with time.
