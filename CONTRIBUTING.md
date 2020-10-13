# Contributing Guide
If you want to contribute to `redhawk`, then please follow this guide

## Process
1. Set up local development environment
2. Change something...
3. Run test
4. Create PR

### 1. Set up local development environment
- You have to install...
    - golang ( >= 1.14 )
- Set GOPATH 

```bash
$ cd $GOPATH/src
$ mkdir -p github.com/DevopsArtFactory
$ cd github.com/DevopsArtFactory
$ git clone https://github.com/DevopsArtFactory/redhawk.git
$ cd redhawk
```

### 2. Change something
- Change codes
- If you create new function, then please **make unit test**.
- Please run `make format` in order to do formatting

### 3. Run test
- `make linters`: This will check the rules for clean code.
- `make test`: Run unit test

### 4. Create PR
- Thank you so much for your Pull Request!!
