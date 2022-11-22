Installs any packages you specify at the command line

Running `ahkpm install` without specifying a package name will download all
dependencies specified in ahkpm.json into the `ahkpm-modules` folder.

Packages may be specified as either `<packageName>@<version>` or as just
`<packageName>`.

For example, `ahkpm install github.com/user/repo@1.0.0` will download version
1.0.0 of the package into the `ahkpm-modules` folder as well as save the package
name and version to `ahkpm.json` for future use.

You may also use package name shorthands, such as `gh:user/repo`.

For versions you may specify a range such as `1.x.x` or `1.2.x`.

If you do not specify a version, ahkpm will attempt to find the latest valid
semantic version. If no valid semantic version of the package is available,
it will fall back to `branch:main`. If there is no `main` branch, it will
fall back to `branch:master`. There are no further fallbacks.
