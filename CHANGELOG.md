# Changelog

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