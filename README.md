# go-pkg-lib

This repository contains packages that pantheon can use company wide in their go based services. Each package will have its own release through a bot called [release-please](https://github.com/googleapis/release-please) or check [this document](https://getpantheon.atlassian.net/wiki/spaces/VULCAN/pages/3082911827/Library+Package+Release+Process).

## Release process

### Commit message and automatic pull request

To follow a release process, we just need to take care of the commit message, rest will be taken care by the `release-please` workflow.

We need to follow [conventional commit message](https://www.conventionalcommits.org/en/v1.0.0/) however, with release-please it seems only [few are supported](https://www.conventionalcommits.org/en/v1.0.0/).

- `fix:`: which represents bug fixes, and correlates to a [SemVer](https://semver.org/) patch.
- `feat:` which represents a new feature, and correlates to a SemVer minor.
- `feat!:`, or `fix!:`, `refactor!:`, etc., which represent a breaking change (indicated by the !) and will result in a SemVer major.

We can only `Squash and Merge` to maintain the commit messages in the `main` branch. Once a PR is merged, the `release-please` workflow will create a new pull request automatically with a commit prefixed with `chore(main): release main`. 

The new PR will update `.release-please-menifest.json` file with the new tag it is going to release for the package and the package's CHANGELOG file.

Once merged, it will release a new tag for the package and package can be downloaded and used with following command:

```
go get github.com/pantheon-systems/go-pkg-lib/metric@v1.0.0
```

### Package release dependency chain

A good process for commit message would be to add a commit including changes on that package only. For example, if you have change in two packages such as `metrics` and `api` where `api` depends on the `metrics` package, first release the api package.

Make changes in the `metric` package, use `go workspace` to test the package change locally in the `api` package and once satisfied, raise the pull request for `metric` package following the proper commit message.

An example would be: 
```
git commit -m "feat(fesi-00): added support for labels with request counter."
```

Once the PR is merged, a new PR will be created to release a new version for `metric` package with the message in the change log. Once the new PR is merged, it will release a new version for the `metric` package.

The new version of `metric` package can be udpated in the `api` package's go.mod file and then it can be released in the same way.

If there are more packages that depends on the `metric` and requires updates simultaneously, we still need to include package specific files in a single commit, but the workflow will be able to update the automatic PR with all release at once.

In case of `api` package, you can change some other package as well and commit those changes into a separate commit message. Once merged, the automatic release PR created for `api` will be updated with this new package as well.

## Adding a new package

Added a new package will require changes in three files, `.release-please-menifest.json`, `release-please-config.json` and `.github/labeler.yml`. Having them in those files is enough to follow the same release pattern for the newly added package as well.
