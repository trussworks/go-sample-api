# Application Packages (pkg)

`pkg` contains everything necessary
to serve the `bork` application.
Its entry-point is defined in [server](./server).

This doc defines its various subdirectories.
Note that directories prefixed with `app`
don't have any special meaning
other than providing uniqueness from standard library
or popular package naming,
and thus easier imports from the caller.

## [APIs](./apis)

`apis` hold API definitions for the application.
This layer sets up and fulfills the contract with calling clients.

The responsibility of this layer is to:

1. Accept a request and marshal it into Go domain model equivalents.
2. Call a service.
3. Write a response.

Note that APIs can be defined in different architectures/protocols
(REST, grpc, etc.),
but should be composable without affecting the business logic of the application.
The idea behind this is that APIs can evolve
and update to modern standards
without requires heavy changes to the rest of the code.
Current definitions include:

* [HTTP Handler REST](./apis/v1/http/handlers)

## [Configuration](./appconfig)

`appconfig` is mostly a set of constants
for referencing variables in a configuration implementation,
and providing consistent access.
It also contains some functional help
for working with environments.

Note that it doesn't reference non-standard libraries
and is agnostic to configuration implementation, like viper.

## [Context](./appcontext)

`appcontext` holds application specific
getters and setters for Go's [context](https://golang.org/pkg/context/).
Context is one of the few layers shareable across layers of the codebase.
It's most commonly used for injecting request/response information
from the calling client.
A common example is a `User` model for accessing role/authorization information.

As a general rule, `context.Context` should be passed down through the
layers of your program, as this is the conventional Go way to address
"cross cutting concerns", e.g. cancellation, logging, distributed
tracing, or other types of instrumentation (which in other languages
might be addressed via Thread Locals or similar constructs).

However, for any given function, concrete 1st order
dependencies for the operation of that function should be called out
as an explicit parameter, rather than buried in a call to the
not-type-safe `context.Value`.

## [Errors](./apperrors)

`apperrors` is a set of errors that adhere to Go's built-in `error` type.
These represent a set of known errors that can occur in the application
and are often one of two return values,
coinciding with domain [models](#domain-modelsmodels).

`apperrors` provide a way for calling functions
to act on an error return.
Using basic errors
(ex. `errors.New()`)
either forces the caller to fail
(ex. throw a 500 status code)
or attempt to recover based on type-unsafe string comparisons.
Instead,
with typed errors,
the caller can assert against the type
and make decisions.
An example in this code base
is determining HTTP status codes
based on service layer errors
in the [handler base](./apis/v1/http/handlers/handler_base.go).

## [Integration Tests](./integration)

Integration tests provide a way
to test that the rest of the layers of code
are set up properly.
Since application layers are composable
and independently testable,
we need a way to test that the current application configuration runs.

Integration tests should mock out as little as possible,
and run against a fully equipped server.
In this codebase,
[the integration test setup](./integration/integration_test.go)
shows an example running a test `http.Handler` based server.

Integration tests are slow and complex,
thus should be minimal--
attempting to only test main code routes
and integration availability.
Wide coverage tests,
should be left to layer based unit tests.

## [Domain Models](./models)

`models` contain types
for passing data across multiple layers of the codebase.
They are one of two returns values
(coinciding with [errors](#errorsapperrors)).
They are **used** in different ways
(ex. marshaling to JSON, retrieving from databases,
performing business functions against),
but do not provide function **within** themselves.
The simple provide structured data
for manipulation by other layers.

The goal of this layer,
is to restrict functional requirements
being shared across multiple layers of code.
This forces simple contracts across layers,
and in turn decouples functional implementation
from data.
[APIs](#apisapis) and [data source](#data-sourcesources)
become much easier to maintain and upgrade
independently of each other.

## [Server](./server)

`server` is the entry-point for serving the application.
Its focus is to pull in configurations,
set implementations from [APIs](#apisapis),
[services](#servicesservices), and [data sources](#data-sourcesources)
and execute their setup functions
(ex. serve a route).
`server` is the only place where configurations
and implementations are directly referenced.

`server` also contains some other application orchestration,
such as adding [auth middleware](./server/httpserver/auth.go)
and a logger.

Note that `server` should attempt to "fail fast"
when it doesn't have the requirements to run.
For example,
if the application can't set up a database with configurations,
[check them](./server/httpserver/config.go)
and exit.
Try not to let the lower levels of the codebase
fail on knowable startup issues.

## [Services](./services)

`services` contains the business logic of the codebase.
The goal is that this layer
can be run independent of APIs and data sources.
It's purpose is
to perform an action for a user.
This is where user value is derived in the server application.
This is opposed to other layers of the codebase,
which tend to perform or adhere to technical constraints.

For example, API code may answer "how do I adhere to the dog REST interface?",
a data source may say "how do I retrieve these dogs from postgres?",
but the service answers the question of
"how do I get the dogs that user 'x' owns?".

`services` defines a purely Go typed contract for parameters and returns,
leaving data retrieval and marshaling to other layers.
This allows it to be independently assessed,
without knowledge outside of Go's compile time context.

## [Data Sources](./sources)

`sources` provides packages
to retrieve/set data
to/from domain models.
The most common example,
is a adapter for a relational database,
such as Postgres.

The goal of this layer,
is to separate source implementation,
such as client API code or query language
from the rest of the codebase.
The server package will then set up a source,
based on application configuration.
Within the context of a single request,
the service layer will call functions on the source,
but is agnostic to its implementation.
