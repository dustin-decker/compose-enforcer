# ComposeEnforcer
### Compose enforcer validates that a provided docker-compose file meets restrictions set.

It can be used via CLI from the compiled binary, or you can use the public interfaces
provided. A common usecase is within a continuous deployment pipeline.

I am making no guarantees of API stability at this time.

Currently it is developed against the file spec v3.2 and supports:
  - volumes
  - networks
  - secrets
  - resources

Under the hood it uses docker compose interfaces so it should *hopefully*
be easy to maintain with time.
