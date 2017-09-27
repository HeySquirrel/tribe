# tribe
Quickly answer the question “Why the @*$% does this code exist?”

## Features
- Integration with Jira or CA Agile Central to quickly access historical work items or issues
- Frequent contributors
- Commits across the last year of the file

## Usage

```sh
$ # Why do these lines of code exist
$ tribe blame -L100,105 model/user.rb

$ # See basic information about your work items or issues
$ tribe show HIL-78
```

## Configuration
The configuration for tribe is stored in $HOME/.tribe.json. Currently the only configuration is the work item servers you want tribe to understand. Below is an example format.

```json
{
  "workitemservers": {
    "rally1": {
      "type": "rally",
      "host": "https://<cool rally server>",
      "apikey": "<rally api key here>",
      "matcher": "(S|DE|F|s|de|f)[0-9][0-9]+"
    },
    "myjira": {
      "type": "jira",
      "host": "https://<cool jira server>",
      "username": "<jira username>",
      "password": "<jira password>",
      "matcher": "HIL-[0-9]+"
    }
  }
}
```

## Contribution

1. Fork this repo
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new Pull Request

## Installation

### Developer

```sh
$ go get -u github.com/heysquirrel/tribe
```
