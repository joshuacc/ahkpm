# ahkpm

The [AutoHotKey][ahk] package manager - a work in progress

## Commands

- [ ] `ahkpm init`: Interactively create an `ahkpm.json` file in the current directory
- [ ] `ahkpm install github.com/user/repo`: Installs a package (other hosts to come later)
- [ ] `ahkpm update <name>`: Updates installed packages to the newest version
- [ ] `ahkpm list`: List local packages and their installed versions
- [ ] `ahkpm help`: Provides documentation
- [ ] `ahkpm selfupdate`: Upgrade to the latest version of ahkpm

Install an AutoHotKey script (or scripts) from GitHub: `ahkpm github.com/user/repo`

## ahkpm.json

```json
{
    "name": "my-project",
    "git": "github.com/user/my-project",
    "website": "https://my-project-website.com",
    "dependencies": {
        "github.com/user/repo": "tag or commit hash"
    }
}
```

[ahk]:https://www.autohotkey.com/