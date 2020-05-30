# Bork Go Application

The bork application
is a sample application
to demonstrate Go API patterns.

## Codebase Layout

### [cmd](./cmd)

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

### Migrating the Database

To add a new migration, add a new file to the `migrations` directory
following the standard
`V_${last_migration_version + 1}_your_migration_name_here.sql`

## Development and Debugging

### APIs

The APIs reside at `localhost:8080` when running.
To run a test request,
you can send a GET to the health check endpoint:
`curl localhost:8080/api/v1/healthcheck`
