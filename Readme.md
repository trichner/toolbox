# Thomas' Automation Tools

## Quickstart

```
#--- initial setup
brew install go

#setup the 'GOPATH'
export GOPATH=$HOME/workspaces/go

# fetch all dependencies
go get

# install the binaries
go install ./...

# binaries can be found here:
ls $GOPATH/bin/

# add them to your path
export PATH=$PATH:$GOPATH/bin

#--- update the binaries whenever code changes

go install ./...
```

## TL;DR

```bash
tb sheet2json --spreadsheet-url=<sheetUrl>
```

```bash
tb csv2json < some.csv | jq .
```

```bash
echo '{"a":1, "b":true}' | tb json2sheet
```

## Authentication

### Google APIs

1. create an OAuth consent screen as
   documented [here](https://support.google.com/cloud/answer/6158849?hl=en)
2. create client credentials for a 'Desktop App' for said consent screen
3. store the client credentials in a `client_secret.json` in the working directory

### GitHub

*WARN: this is not particularly safe!*

1. create a Personal Access Token (PAT)
2. put it into `~/.config/github/token.txt` or `export GITHUB_TOKEN=<your PAT>`

**TIP:** you can also load it from the file
via `export GITHUB_TOKEN=$(cat ~/.config/github/token.txt | tr -d '\n')`

### Jira

*WARN: this is not particularly safe!*

1. provision an API token
2. create a `credentials.json` file, looking like
    ```
    {
    "username": "octo@example.com",
    "token": "s0meT0kn"
    }
    ```
3. put the file into `~/.config/jira/credentials.json`

### Freshdesk

1. create an API key
2. set the `FRESHDESK_TOKEN` environment variable, e.g. `export FRESHDESK_TOKEN=<your token>`
