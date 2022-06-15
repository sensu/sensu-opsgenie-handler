# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased
## Changed
- Update to use sensu-plugin-sdk 0.16

## [0.9.0] - 2021-01-27

### Changed
- Changed tags from hard-coded list to templates settable via argument
- Moved test files into tests directory to clean up top level directory
- Updated all modules with 'go get -u' and 'go mod tidy'
- README updates and cleanup

### Added
- Lint GitHub Action

## [0.8.0] - 2020-12-01

### Changed
- Minor README fixes
- Updated SDK to 0.11.0
- Updated module dependencies

### Updates brought in from [upstream](https://github.com/betorvs/sensu-opsgenie-handler)
- Removed flag `OPSGENIE_APIURL` now we use constants from opsgenie sdk-v2.
- Removed `opsgenie_priority` annotation. Should use: `"sensu.io/plugins/sensu-opsgenie-handler/config/priority": "P3"`.
- Added flag `--region` to choose opsgenie region. Can be configured using environment variable too `OPSGENIE_REGION`. This feature replaces old `OPSGENIE_APIURL`.
- Added flag `--priority` to change Opsgenie default priority. String field. Expected: "P1", "P2", "P3", "P4" and "P5".
- Changed opsgenie sdk to [v2](https://github.com/opsgenie/opsgenie-go-sdk-v2).
- Changed withAnnotations to parse all annotations, and exclude if it contains `sensu.io/plugins/sensu-opsgenie-handler/config`, to send to Opsgenie.
- Added more tests
- Added `--allowLabels` to parse all Labels and send to Opsgenie.
- Added `--fullDetails` to add all kind of details in Opsgenie (changed from the upstream based on my opinions).

## [0.7.0] - 2020-08-14

### Changed
- Updated SDK to 0.8.0
- Set secret bool to true for authkey

## [0.6.2] - 2020-07-02

### Changed
- Fixed go.sum for goreleaser to work

## [0.6.1] - 2020-07-02

### Changed
- Increased the default description limit to 15000, matching OpsGenie published limit
- Increased the default message limit to 130, matching OpsGenie published limit
- Clened up a duplicate import of github.com/opsgenie/opsgenie-go-sdk/alertsv

## [0.6.0] - 2020-05-15

### Changed
- Updated README to reference secrets, add environment variables, and other cleanup
- Fixed Priority annotations to properly reflect entity annotations precedence over check annotations
- Removed last vestiges of withAnnotations

## [0.5.0] - 2020-05-14

### Added
- Details field for parity with Sensu Enterprise handler
- Expansion of embedded \n newlines in description
- More test coverage

### Changed
- Minor README fixes
- Changed source from "sensuGo" to "Sensu Go"
- Fixed bug where evenPriority annotations weren't being checked correctly

## [0.4.1] - 2020-05-11

### Changed
- Minor README fix

## [0.4.0] - 2020-05-11

### Removed
- withAnnotations flags and its code to create description field

### Added
- descriptionTemplate and descriptionLimit flags to allow customizing the description field

## [0.3.2] - 2020-05-01

### Added
- includeEventInNote bool flag to include JSON of event in Note, false by default

## [0.3.1] - 2020-05-01

### Change
- Formatting cleanup and goreportcard badge
- Use Sensu SDK templating engine

## [0.3.0] - 2020-05-01

### Changed
- Fixed alias to be set string entity name/check name, as opposed to matching message title

## [0.2.1] - 2020-04-30

### Changed
- Updated .bonsai.yml

## [0.2.0] - 2020-04-30

### Changed
- Move to Go Modules
- Use new Sensu SDK
- Move from Travis to GitHub Actions
- Reorganize README
- Added timestamps to all the test events

### Added
- Made message a configurable template, including a length limit

## [0.1.0] - 2020-03-31

### Changed
- Changed behaviour from opsgenie-handler to not add new alerts as a note and send them as alert.

### Removed
- Removed addNote func
- Removed eventKey func
- Removed eventTags func
- Removed goreleaser goos freeebsd and arch arm64

## [0.0.10] - 2020-02-23

### Added
- Parse Entity.Name, Check.Name, Entity.Namespace, Entity.Class as tags to Opsgenie
- Add OPSGENIE_SENSU_DASHBOARD variable to add new field in Description with "source: Sensu Dashboard URL"


## [0.0.9] - 2020-02-08

### Added
- Added the Sensu event json dump to the OpsGenie `Note` field.
- Added more tests

## [0.0.8] - 2020-01-20

### Changed
- change from dep to go mod
- gometalinter to golangci-lint
- correct goreleaser

## [0.0.7] - 2019-12-09

### Added
- Correct issue [#6](https://github.com/betorvs/sensu-opsgenie-handler/issues/6): `trim additional ending slash in --url argument`
- add script test.all.events.sh

## [0.0.6] - 2019-11-24

### Added
- Add `OPSGENIE_ANNOTATIONS` to parse annotations to include that information inside the alert.
- Update README.

## [0.0.5] - 2019-10-15

### Added
- Add `OPSGENIE_APIURL` to change OpsGenie API URL
- Updated Gopkg.lock file.
- Changed travis go version.

## [0.0.4] - 2019-08-26

### Added
- Add bonsai configuration

## [0.0.3] - 2019-08-02

### Added
- Add OpsGenie Priority as annotations inside check annotations.
- Add Get, Close and Add Note functions to manage alerts already open. 

## [0.0.2] - 2019-07-10

### Added
- Add OpsGenie Priority as annotations inside sensu-agent to override default Alert Event Priority in OpsGenie.

## [0.0.1] - 2019-07-10

### Added
- Initial release
