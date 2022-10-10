# ahkpm - The AutoHotKey Package Manager.

[AutoHotKey][ahk] is a powerful tool for scripting automations on Windows,
but managing dependencies for your scripts is painful.
`ahkpm` intends to bring modern package management to AutoHotkey,
making it easier than ever to automate away the drudgery.

## Commands

- [ ] `ahkpm init`: Interactively create an `ahkpm.json` file in the current directory
- [ ] `ahkpm install github.com/user/repo`: Installs a package (other hosts to come later)
- [ ] `ahkpm update <name>`: Updates installed packages to the newest version
- [ ] `ahkpm list`: List local packages and their installed versions
- [x] `ahkpm help`: Provides documentation
- [ ] `ahkpm selfupdate`: Upgrade to the latest version of ahkpm

Install an AutoHotKey script (or scripts) from GitHub: `ahkpm github.com/user/repo`

## ahkpm.json

```json
{
  "name": "my-project",
  "version": "0.0.1",
  "description": "A brief description",
  "repository": "github.com/user/my-project",
  "website": "example.com",
  "license": "MIT",
  "issueTracker": "github.com/user/my-project/issues",
  "author": {
    "name": "joshuacc",
    "email": "",
    "website": "joshuaclanton.dev"
  },
  "dependencies": {
    "github.com/user/repo": "branch, tag or commit hash"
  }
}
```

[ahk]:https://www.autohotkey.com/