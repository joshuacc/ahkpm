Installs specified package(s). If none, reinstalls all packages in ahkpm.json.
	
For example, `ahkpm install github.com/user/repo@1.0.0` will download version
1.0.0 of the package into the `ahkpm-modules` folder as well as save the package
name and version to ahkpm.json for future use.

You may also use package name shorthands, such as `gh:user/repo`.

Running `ahkpm install` without specifying a package name will download all
dependencies specified in ahkpm.json into the `ahkpm-modules` folder.