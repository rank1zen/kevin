# Tasks

The project uses ([`mise`](https://mise.jdx.dev)) for installing dev tools, declaring environment variables, and
defining canonical tasks.

### Key files

* `mise.toml` pins the dev tools used, as well as the environment variables required; all environment variables are
  required.
* `mise.local.toml` should be used to provide the actual environment values. `mise` will error if a value is missing.
* `.mise/tasks` contains all the tasks used in the project, including, testing, building, and infrastructure
  provisioning.
