# Tasks

The project uses [mise](https://mise.jdx.dev) for installing dev tools, declaring environment variables, and
defining canonical tasks.

### Key files

* `mise.toml` pins the dev tools used; further, it declares the environment variables required for the application and
  infrastructure.
* `mise.local.toml` should be used to provide the actual environment values.
* `.mise/tasks` contains all the tasks used in the project, including, testing, building, and infrastructure
  provisioning. Each task is simply a bash script with mise annotations. Tasks can be executed using `mise run`.
