# release-automaton

```
automaton n.

a machine that performs a function according to a predetermined set of coded instructions,
especially one capable of a range of programmed responses to different circumstances.
```

## Test Release Repo

- https://github.com/appscodelabs/release-automaton-demo

## Release Scenarios

### Alpha/Beta Release

- Don't update latest tag on firebase.json file in website repos.

### Patch Releases

- Make sure cherry picks are done ahead of time.

## Version Bump Script

Increment minor version (and reset patch to 0) for semver values in a product release file:

- `./hack/scripts/bump-release-minor.sh kubedb`
- `./hack/scripts/bump-release-minor.sh voyager`

Supported products:

- `voyager`
- `kubevault`
- `kubedb`
- `ace`
- `stash`
- `kubestash`

Alternative via Makefile:

- `make bump-release-minor PRODUCT=kubedb`
