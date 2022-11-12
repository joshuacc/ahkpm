Updates the specified package(s), by checking `ahkpm.json`, determining the
latest version allowed by the version range listed there, and downloading it
to `ahkpm-modules`.

For example, if you have a dependency on `github.com/user/repo` with version
`branch:main`, running `ahkpm update github.com/user/repo` will update the
package to the latest commit on the main branch.

You may also use package name shorthands, such as `gh:user/repo`.