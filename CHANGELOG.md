# Changelog

## v0.2.1 (30/07/2021)

---

## v0.2.2 (16/12/2021)
## Changes

## Features

- https://github.com/vdjagilev/nmap-formatter/pull/58: #34: Remove blank vertical spaces in Markdown format

## Other

- https://github.com/vdjagilev/nmap-formatter/pull/57: #55: Remove href links from markdown template
- https://github.com/vdjagilev/nmap-formatter/pull/56: Refactor code with gocritic linter
- https://github.com/vdjagilev/nmap-formatter/pull/53: other: minor refactoring of show version functionality
- https://github.com/vdjagilev/nmap-formatter/pull/51: Fix "go install" url
- https://github.com/vdjagilev/nmap-formatter/pull/50: Update dependencies \& go.mod file
- https://github.com/vdjagilev/nmap-formatter/pull/49: #48: Update golang to 1.17

---

## Version v0.2.0 (30/07/2021)
## Changes

### Bugs

- #33: OSMatch can have multiple values

### Features

- #37: Add version to help output & version command

### Other

- #38: Use release drafter github action to generate new release content
- #41: Add lint
- #44: Add release autocreation
---

## Version v0.1.0 (19/07/2021)
Little improvements & some bugfixes https://github.com/vdjagilev/nmap-formatter/pull/36

# Bugfixes

* https://github.com/vdjagilev/nmap-formatter/issues/5: Fix CSV formatter skipping down hosts

# Features

* https://github.com/vdjagilev/nmap-formatter/issues/7: Rename down-hosts output option
* https://github.com/vdjagilev/nmap-formatter/issues/12: Increase unit test coverage
* #11: Implement json-pretty print option
* #8: Add new skip output options

# Other

* https://github.com/vdjagilev/nmap-formatter/issues/2: Move types formatter package
* https://github.com/vdjagilev/nmap-formatter/issues/10: Remove unused variable from `root.go` validate function
* https://github.com/vdjagilev/nmap-formatter/issues/6: Add build status badge
* https://github.com/vdjagilev/nmap-formatter/issues/3: Add table of contents to the README
* #28: Cleanup
* https://github.com/vdjagilev/nmap-formatter/issues/9: Add matrix strategy for ci workflow
* #20: Add codeclimate badge
* #4: Add nmap-formatter examples with jq tool



---

## Version v0.0.2 (29/06/2021)
# Changelog

* Minor changes in documentation
* Added assets
---

## Alpha release (28/06/2021)

---

## Initial release (28/06/2021)
Basic functionality:

* Convert:
  * Markdown
  * JSON
  * CSV
  * HTML
* Output to file
* Show hosts that are down
