# Releasing versions

The repo has 2 packages :
* github.com/outscale/goutils/sdk
* github.com/outscale/goutils/k8s

Each has its own versions.

To create new version of github.com/outscale/goutils/sdk:

```shell
TAG=sdk/vX.Y.Z git tag $TAG -m "🔖 $TAG" && git push origin $TAG
```

To create new version of github.com/outscale/goutils/k8s:

```shell
TAG=k8s/vX.Y.Z git tag $TAG -m "🔖 $TAG" && git push origin $TAG
```

No releases are created, Go only needs tags.