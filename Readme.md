# Thomas' Automation Tools

## Installation

```shell
go install github.com/trichner/tb@latest
```

## TL;DR

```bash
printf "a,b,c\nhello,2,3" | tb csv2json | jq .
```

```bash
echo '{"a":1, "b":true}' | tb json2sheet
```

```bash
tb sheet2json --spreadsheet-url=<sheetUrl>
```

## Bash 'command not found'

```shell
# binaries can be found here:
ls "$(go env GOPATH)/bin/"

# add them to your path if not already done
export PATH=$PATH:$(go env GOPATH)/bin
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
