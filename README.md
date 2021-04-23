# Go Sample API Application

The bork application
is a sample application
to demonstrate Go API patterns
and provide a lightweight template
for a new application.

Rather than attempting to create
a monolithic CRUD/MVP framework based app,
this application centers a set of "services"
which contain business logic for the application.
API and data layers are then orchestrated
to serve and retrieve data defined by the services.

**Note**, this application is mostly a guide
for developing applications with a particular set of patterns.
Adopting it requires some orchestrating
based on your business context
(ex. plugging in an auth solution).
It attempts to illustrate concepts
that have been repeated in multiple Truss environments,
while implying space for domain specific changes.
Some basic uses of this repo are:

1. Forking some boilerplate code commonly written at application bootstrap.
2. Looking for Go [API patterns](#patterns) to employ.
3. Adding sample API or data sources.
   Want to show what a GraphQL server looks like
   as opposed to standard library handlers?
   [APIs](./pkg/apis) would be a great spot.
   Maybe there's a new sql adapter--
   take a look at [Data Sources](./pkg/sources).
4. Testing out a new pattern on a minimal dependency application.
   While a production application may have compiled layers or history
   and third party dependencies,
   this API mostly requires Go and Postgres
   and has minimal "business" constraints.

## Patterns

The main patterns of the API are:

* **Service package for business logic.**
  Business logic is scoped to a single package
  using single functions per instruction of work.
  Since application development is
  led by client user value,
  the goal is to focus on adding value
  to this layer of the code.
  Feature oriented work
  will tend to be mirrored in this layer.
* **Composable API and data source layers.**
  The API and data source layers
  are treated as tools to serve the logic of the application.
  Rather than overloading those layers with logic
  (such as monolithic handlers),
  they serve small, scoped sets of responsibilities.
  Separating these packages
  allows the tooling of the application to progress
  without affecting the business logic (ex. migrating from REST to GraphQL).
  A goal  is to spend less time in these layers
  and automate out of as much of the boilerplate
  (ex. generating code)
  to allow for focus on the service layers.
* **Cross package domain models.**
  Few resources are shared across layers of the application
  and these resources are almost entirely scoped
  to method-less data structures
  and Go basic types or standard library definitions.
  These include:
  [models](./pkg/models),
  [errors](./pkg/apperrors),
  and [context](./pkg/appcontext).
  Reducing the function of these types
  heavily reduces coupling of the application's components.
* **Internally defined, minimal, dependency/parameter interfaces.**
  When defining dependencies for inter-layer communication,
  interfaces are defined internally to the **calling package**.
  These interfaces are scoped to fewest methods necessary
  and often defined as typed functions.
  These implementations of interfaces
  are set up in a main [server package](./pkg/server/httpserver)
  and explicitly defined.
  Elsewhere in the application,
  functions are almost exclusively referred to
  as interfaces or typed functions.
  Paired with minimal domain models,
  this reduces layer coupling
  and allows implementation packages
  to implicitly satisfy caller requirements.

Note that many of these patterns pull from
[SOLID](https://en.wikipedia.org/wiki/SOLID) principles
and attempt to implement them in line with Go's language design.

## Libraries

The application also consists of a few dependable libraries
used in Truss codebases,
such as:

* [Testify](https://github.com/stretchr/testify)
  for testing suites and sub-test structuring.
* [sqlx](https://github.com/jmoiron/sqlx)
  for marshaling sql into Go structs.
* [mux](https://github.com/gorilla/mux)
  mostly used for simple sub-routing.
* [viper](https://github.com/spf13/viper)
  for app configuration.
* [clock](https://github.com/facebookgo/clock/) It's a clock interface and mock clock.
  Mostly used for `time.Now()` mocking in tests.
* [zap](https://github.com/uber-go/zap)
  A structured logger with nice syntax for adding fields.

For more information on each package,
see [codebase layout](#codebase-layout).

## Codebase Layout

### [Command Line Utilities](./cmd)

`cmd` is the entry-point
for **running** the application from the **command line**.
It's lightweight and holds very few responsibilities.
Namely:

1. Read the command line arguments.
2. Initialize environment variables.
3. Execute the application.

In a production application,
the may also hold some helper scripts
for running maintenance or testing tasks
on the application.

### [Application Packages](./pkg)

`pkg` holds the necessary code
for **serving** the application.
This is distinct from `cmd`
in that is should be executable from multiple environments
including, `cmd`, integration testing or other Go code.

In relation to `cmd`, it takes `viper.Viper` as it's only argument
(in the server package),
and has no awareness of the command line or executing environment.

## Application Setup

### Setup: Developer Setup

There are a number of things you'll need at a minimum
to be able to check out,
develop,
and run this project.

* Install [Homebrew](https://brew.sh)
  * Use the following command
    `/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"`
* Install Go with Homebrew:
    `brew install go`
      * **Note**:
        If you have previously modified your PATH
        to point to a specific version of go,
        make sure to remove that.
        This would be either in your `.bash_profile` or `.bashrc`,
        and might look something like
        `PATH=$PATH:/usr/local/opt/go@1.12/bin`.
* Ensure you are using the latest version of bash for this project:
  * Install it with Homebrew:
    `brew install bash`
  * Update list of shells that users can choose from:

    ```bash
        [[ $(cat /etc/shells | grep /usr/local/bin/bash) ]] \
        || echo "/usr/local/bin/bash" | sudo tee -a /etc/shells
    ```

  * If you are using bash as your shell
    (and not zsh, fish, etc)
    and want to use the latest shell as well,
    then change it (optional): `chsh -s /usr/local/bin/bash`
  * Ensure that `/usr/local/bin` comes before `/bin`
    on your `$PATH` by running `echo $PATH`.
    Modify your path by editing `~/.bashrc` or `~/.bash_profile`
    and changing the `PATH`.
    Then source your profile with `source ~/.bashrc` or `~/.bash_profile`
    to ensure that your terminal has it.
* **Note**:
    If you have previously used Golang, please make sure none
    of them are pinned to an old version by running `brew list --pinned`.
    If they are pinned, please run `brew unpin <formula>`.
    You can upgrade these formulas instead of installing by running
    `brew upgrade <formula`.

### Setup: Git

Use your work email when making commits to our repositories.
The simplest path to correctness is setting global config:

  ```bash
  git config --global user.email "trussel@truss.works"
  git config --global user.name "Trusty Trussel"
  ```

If you drop the `--global` flag,
these settings will only apply to the current repo.
If you ever re-clone that repo or clone another repo,
you will need to remember to set the local config again.
You won't.
Use the global config. :-)

For web-based Git operations,
GitHub will use your primary email unless you choose
"Keep my email address private".
If you don't want to set your work address as primary,
please [turn on the privacy setting](https://github.com/settings/emails).

Note that with 2-factor-authentication enabled,
in order to push local code to GitHub through HTTPS,
you need to [create a personal access token](https://gist.github.com/ateucher/4634038875263d10fb4817e5ad3d332f)
and use that as your password.

### Setup: Project Checkout

You can checkout this repository by running
`git clone git@github.com:trussworks/go-sample.git`.
Please check out the code in a directory like
`~/Projects/go-sample` and NOT in your `$GOPATH`. As an example:

  ```bash
  mkdir -p ~/Projects
  git clone git@github.com:trussworks/go-sample.git
  cd go-sample
  ```

You will then find the code at `~/Projects/go-sample`.
You can check the code out anywhere EXCEPT inside your `$GOPATH`.
So this is customization that is up to you.

### Setup: direnv

* Install direnv:
    `brew install direnv`
* Set environment with:
    `direnv allow`

### Setup: Pre-Commit

* Install pre-commit:
    `brew install pre-commit`
* Run `pre-commit install`
    to install a pre-commit hook
    into `./git/hooks/pre-commit`.
* Next install the pre-commit hook libraries
    with `pre-commit install-hooks`.

### Golang cli app

To build the cli application in your local filesystem:

```sh
go build -a -o bin/bork ./cmd/bork
```

You can then access the tool with the `bork` command.

#### air

Alternatively to manually building the app,
you can run the server and reload files with air:

```bash
# outside of a Go project,
# to install globally
go get -u github.com/cosmtrek/air
air
```

### Run the Database

To use the database, run:

```bash
    docker run --name bork-postgres \
        --publish $PGPORT:5432 \
        -e POSTGRES_DB=$PGDATABASE \
        -e POSTGRES_PASSWORD=$PGPASS \
        -d postgres
```

Install flyway to run migrations:

```brew install flyway```

Then run them with:

```flway migrate```

To add a new migration, add a new file to the `migrations` directory
following the standard
`V_${last_migration_version + 1}_your_migration_name_here.sql`

## Testing

Run tests with:

```bash
    go test ./pkg/...
```

## Development and Debugging

### APIs

The APIs reside at `localhost:8080` when running.
To run a test request,
you can send a GET to the health check endpoint:
`curl localhost:8080/api/v1/healthcheck`

### Logging

Logging in this project is intentionally restricted to a single structured log
line per request. The log line will be either Info or Error Level and will have
a variety of key/value pairs associated with it. The primary reason for doing
things this way is to aid in debugging issues in production by bucketing
all info by request. This allows you to use line-oriented tools to quickly
determine if there are commonalities between errors.

The line logging is performed by the middleware in
`pkg/server/httpserver/logger.go`, in code you interact with it with these
methods from `appcontext`.

```go
// LogRequestField adds a zap.Field to the line logged at the end of the request
LogRequestField(ctx context.Context, field zap.Field)
```

You can call this anytime you have a piece of information that you might like
associated with this request. e.g.:

* `user_id`
* `auth_type`
* `number_of_dogs`

```go
// LogRequestError adds a message to the request log line and also sets it to
// log at the Error level
LogRequestError(ctx context.Context, err error)
```

You should call this at most once in a request if you want the request to log
at the Error level. The error will be logged with the zap standard `error` key.
In general, these should generally correspond with requests that have `5XX`
status codes. Error log levels may trigger notifications and require
investigation.
