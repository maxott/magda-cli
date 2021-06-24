# Command-line tool for Magda

[Magda](http://magda.io) is a federated data catalog system developed by [Data61's Engineering & Design Group](https://data61.csiro.au/en/Our-Research/Programs-and-Facilities/Engineering-and-design).

Magda has an extensive [REST API](https://demo.dev.magda.io/api/v0/apidocs/index.html) and a rich web-based user interface. However, to support simple data operation from the command line (and to learn golang), I wrote this simple command-line tool. It only covers the subset of the Magda API I have been using. However, I'm mor than happy to receive pull requests to extend it's functionality or fix bugs.

## Install

Right now, we don't provide pre-compiled executable, but if you have [go installed](https://golang.org/doc/install), you can easily build & install it with:

```
go install github.com/maxott/magda-cli@latest
```

## Usage

Simply type `magda-cli -h` for a listing of the supported commands, or `magda-cli --help-man | nroff -man | more` for a more [detailed description](./man.md).

Some of the common information, like `--host` will be read from environment variables when they are set. Check for square bracketed names in the help text, like `MAGDA_HOST` for the `--host` flag:

```
% magda-cli -h
usage: magda [<flags>] <command> [<args> ...]

Managing records & schemas in Magda.

Flags:
  -H, --host=HOST    DNS name/IP of Magda host [MAGDA_HOST]
  ...
```

The most commonly set environment variables are:

```
MAGDA_HOST=magda.example.com
MAGDA_AUTH_ID=6a7...
MAGDA_AUTH_KEY=GL...=
MAGDA_TENANT_ID=ab...
```

Like many cli tools we are following the `command sub-command` pattern. For instance, the following will show (read) the content of a record:

```
magda-cli record read -i recordID
```

The sub commands are primarily following the CRUD verbs with some additional functionality if the API supports it.

This tool tries to stick as much as possible to the Magda API and often simply prints what is being returned by that API.

### Shell Completion

I'm using [Kingpin](https://github.com/alecthomas/kingpin), so you can setup shell completion with:

```
# If you're using Bash
eval "$(magda-cli --completion-script-bash)"
# If you're using Zsh
eval "$(magda-cli --completion-script-zsh)"
```

Don't forget to check out the completion for flags by pressing `<TAB>` after typing `--` as in:

```
% magda-cli --<TAB>
--authID --authKey --help --host ...
```

## Installation

This is my first golang program so I'm still trying to work out on how to best distribute it. At this point, you'll just need to follow the usual pattern:

```
go mod tidy
go build .
go install
```

## Building distributions

We are using [GoReleaser](https://goreleaser.com). 

To quote from the [manual](https://goreleaser.com/quick-start/): GoReleaser will use the latest Git tag of your repository. Create a tag and push it to GitHub:

```
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
```

To deploy releases, we also need a `GITHUB_TOKEN` environment variable, which should contain a valid GitHub token with the `repo` scope. You can create the token [here](https://github.com/settings/tokens/new).

```
export GITHUB_TOKEN="YOUR_GH_TOKEN"
```

Now you can run GoReleaser at the root of your repository:

```
goreleaser release
```

If you don't want to create a tag yet, you can also run GoReleaser without publishing based on the latest commit by using the `--snapshot` flag:

```
goreleaser --snapshot --skip-publish --rm-dist
```

