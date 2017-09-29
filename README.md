<p align="center">
  <a href="https://circleci.com/gh/HeySquirrel/tribe"><img src="https://circleci.com/gh/HeySquirrel/tribe.svg?style=svg" alt="Build Status"></img></a>
</p>

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

You can have as many `workitemservers` as you want. Tribe will search through all the servers defined in this section looking for matched work items in your commits.

### Rally Details
* As the code is currently implemented, you will need to obtain an API Key from Rally. You can access your API Key at - https://rally1.rallydev.com/login/ on the API KEYS tab.
* The `matcher` for your Rally subscription depends on how your workspaces are setup in Rally. The starting letters of your Artifacts can be changed by your workspace administrator. See above for an example matcher.
* Work items may not display correctly if the work item has been deleted or in a closed project or you don't have permissions to read that work item.

### JIRA Details
* If the JIRA server is public, you can leave off the username/password from it's configuration.
* Work items may not display correctly if the work item has been deleted or you don't have permissions to read that work item.


## Contribution

1. Fork this repo
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new pull request

## Installation

### Developer

```sh
$ go get -u github.com/HeySquirrel/tribe
```

## License
[MIT](https://github.com/HeySquirrel/tribe/blob/master/LICENSE)
