# Release Guide
If you want to release a new version of `redhawk`, you have to follow this guide

## Administrator
- If you want to release, then you have to become an administrator for redhawk.
- List of administrators are specified in the `./hack/release/check_permission.go`

## How to release
- You have to unset all AWS related environment variable
  - AWS_ACCESS_KEY_ID
  - AWS_SECRET_ACCESS_KEY
  - AWS_SESSION_TOKEN
- Release files will be created in `out` directory and uploaded to s3://benx-devops-files
- You have to check `version` and set the new version to `VERSION` environment variable

```bash
$ unset AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN
$ redhawk version
v0.0.1
$ export VERSION=v0.0.2 
$ make release
```

