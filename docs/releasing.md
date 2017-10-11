# Releasing

This document contains details about the process for building release binaries for example-api.

## Versioning

We use [semantic versioning](http://semver.org/) which specifies version up to 3 levels.

* `MAJOR` version when you make incompatible API changes,
* `MINOR` version when you add functionality in a backwards-compatible manner, and
* `PATCH` version when you make backwards-compatible bug fixes.

## Process

First, you will need to cut a new version. There is a script for that which tags a version.


```sh
./cut-release.sh v0.0.0 --no-dry-run
```

Once you have cut a new release, you should ssh to the private registry server and build a release there. Let's say the private IP address of docker registry server is `ssh -F ssh.config 10.0.1.173`, then you would do:

```sh
ssh -F ssh.config 10.0.1.x
cd example-api
git pull --rebase && ./build-release.sh v0.0.0 --no-dry-run -y
```

After completing, the build script will have created a docker container and uploaded the image to private docker registry.
