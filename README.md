<p align="center"><img src="docs/img/squirrel.png" width="360"></p>
<p align="center">
  <a href="https://circleci.com/gh/HeySquirrel/tribe"><img src="https://circleci.com/gh/HeySquirrel/tribe.svg?style=svg" alt="Build Status"></img></a>
  <a href="https://codeclimate.com/github/HeySquirrel/tribe"><img src="https://codeclimate.com/github/HeySquirrel/tribe/badges/gpa.svg" alt="Code Climate"></img></a>
  <a href="https://codeclimate.com/github/HeySquirrel/tribe/coverage"><img src="https://codeclimate.com/github/HeySquirrel/tribe/badges/coverage.svg" /></a>
</p>

# tribe
Quickly answer the question “Why the @*$% does this code exist?”

## Features
- Integration with Jira or CA Agile Central to quickly access historical work items or issues
- Frequent contributors
- Commits across the last year of the file

## Installation

### Homebrew
```sh
$ brew tap HeySquirrel/tribe
$ brew install tribe
```

### Developer

```sh
$ go get -u github.com/HeySquirrel/tribe
```


## Usage

```sh
$ # Why do these lines of code exist?
$ tribe blame -L100,105 model/user.rb

$ # See basic information about your work items or issues
$ tribe show HIL-78

$ How risky is it to make a change to this file?
$ tribe risk app/models/user.rb

Risk for 'app/models/user.rb'

         1 month ago - Last commit
                 491 - Commit count
                   7 - Commits last six months
                  20 - Work items
                  37 - Contributors
                0.99 - Risk
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

## Like the idea?

:star: this repo to show support. Let us know you liked it on [Twitter](https://twitter.com/heysquirrelco).

## License
[MIT](https://github.com/HeySquirrel/tribe/blob/master/LICENSE)
