# Contributing to ahkpm

## People stuff

Every human being is to be treated with dignity and respect. Act accordingly.

## Workflow stuff

If you're making large changes to the codebase or documentation,
open an issue to discuss it first to ensure that it is likely to be accepted.
It is very disheartening as writing a bunch of code only to find that a different approach is needed.
We don't want to waste anyone's time or effort.

## Technical stuff

To contribute to ahkpm's codebase, you will need the following:

- A computer running Microsoft Windows
- [AutoHotkey](https://www.autohotkey.com/) installed (Optional, but recommended)
- [Go 1.19 or later](https://go.dev/) installed.
- [Mage](https://magefile.org/) installed
- [golangci-lint](https://golangci-lint.run/) installed (Optional, but recommended for in-editor feedback)

After your system meets the requirements above:

- Fork the ahkpm repo
- Clone your fork
- Within the ahkpm directory, run `go get ./...` to install all dependencies
- Make your changes
- Use the following mage commands as needed
  - `mage build`: Compiles source code and outputs `ahkpm.exe` into the `bin` directory
  - `mage lint`: Runs lint checks against the source code
  - `mage test`: Runs all unit tests
  - `mage verify`: Runs all of the above
- Push your changes to your fork
- Open up a pull request with your changes

## License stuff

Any contributions made to ahkpm will be under the MIT license which covers the project as a whole.