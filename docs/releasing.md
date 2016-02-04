# Releasing

First, cut a release using `cut-release.sh` script. If you make a mistake, you can re-run the script as it will force push tags.

```
./cut-release.sh v0.0.0 --no-dry-run
```

Second, build a release using `build-release.sh` script. It will build a Docker image, tag it and push it to S3 bucket bucket. It might take a while depending on your connection speed as gzipped image is about 270MB.

```
./build-release.sh v0.0.0 --no-dry-run
```
