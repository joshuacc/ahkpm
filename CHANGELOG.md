# Changelog

## 0.7.0

- Fixed a bug where `ahkpm init` did not allow the value `"UNLICENSED"`. Added clarifying text, and changed the value to `"NO LICENSE"` to distinguish it from the license named "The Unlicense."
- `ahkpm init` now defaults the `version` field to `1.0.0` to avoid the unintuitive behavior of `0.x.x` versions
- `ahkpm include` now correctly calculates relative script paths when the `--file` flag points to a file in a different folder
- Made `ahkpm include` output compatible with AutoHotkey 2
- Added support for ahkpm scripts to be stored in the `scripts` object in `ahkpm.json`
- The AutoHotkey version is no longer reported by `--version`, but is available with `--ahk-version`.

## 0.6.0

- Added `ahkpm search` command to find packages
- Added `--defaults` flag for `ahkpm init` to allow bypassing prompts
- Added `ahkpm u` as an alias for `ahkpm update`
- Added `ahkpm i` as an alias for `ahkpm install`
- Added `ahkpm include` command to automatically generate `#Include` directive
- `ahkpm install` now supports omitting the version from packages
- `ahkpm update` now supports the `--all` flag

## 0.5.0

- ahkpm now supports version ranges such as `1.x.x`.
- `ahkpm install` now makes the smallest possible change to the dependency tree
- Added a new `ahkpm version` command to bump package version
- `ahkpm install` now supports specifying multiple dependencies
- Added support for `gh:` shorthand for GitHub dependencies in `ahkpm install` and `ahkpm update`
- Added `ahkpm list` to display table of top level dependencies
- The command to get the version of ahkpm has moved. Now use `ahkpm --version` instead of `ahkpm version`

## 0.4.0

- Added `ahkpm update` to update package(s) to latest allowed version
- Fixed a bug where the local package cache was taking precedence over the remote
- `ahkpm version` now also attempts to display the version of AutoHotkey installed
- Set up a new documentation site at [ahkpm.dev][https://ahkpm.dev].

## 0.3.0

- Resolve transitive dependencies
- Add `ahkpm.lock` to prevent unexpected code changes in dependencies
- Validate dependency arguments passed to `ahkpm install`

## 0.2.0

- Remove some unnecessary binary bloat
- Add `ahkpm version` command
- Make ahkpm available with a Windows installer
- Fix GitHub url validation
- Add friendly error messages if package or version can't be found

## 0.1.0

- Initial version