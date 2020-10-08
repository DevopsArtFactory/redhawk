# redhawk
redhawk is a command line tool for beNX developer. This command line interface will help you do work more easily and fast.
You can assume to the different AWS account or generate rds authentication token with one command.

## Important Notice
- For all commands, redhawk will check whether the access key is expired or not (180 days)
- If your access key is older than `180 days`, run `redhawk renew-credential`.
  - It will change your `aws_access_key_id` and `aws_secret_access_key` in $HOME/.aws/credentials.

## Install
* macOS user
```bash
$ brew tap DevopsArtFactory/benx
$ brew install redhawk
$ redhawk version -v info
```

* Linux user
```bash
$ curl -Lo redhawk https://benx-devops-files.s3.ap-northeast-2.amazonaws.com/redhawk/releases/latest/redhawk-linux-amd64
$ sudo install redhawk /usr/local/bin/
$ redhawk version -v info
```

* Windows user
  - file : https://benx-devops-files.s3.ap-northeast-2.amazonaws.com/redhawk/releases/latest/redhawk-windows-amd64.exe
  - Simply download it and place it in your PATH as redhawk.exe.
  
  
 ## Auto completion
- zsh 
  - This is recommended.
 ```bash
$ echo "source <(redhawk completion zsh)" >> ~/.zshrc
$ source  ~/.zshrc
```

- bash 
 ```bash
$ echo "source <(redhawk completion bash)" >> ~/.bash_rc or ~/.bash_profile
$ source  ~/.bashrc
```

## Setting Configuration
- Configuration file should be in `$HOME/.aws/config.yaml`.
- You can easily create your configuration file with `redhawk init`
- `redhawk init` will create a configuration for default profile.
```bash
$ redhawk init
? Your base account:  gslee@bighitcorp.com
- profile: default
  name: gslee@bighitcorp.com
  assume_roles:
    preprod: ""
    prod: ""
  databases:
    preprod:
    - ""
    prod:
    - ""

? Are you sure to generate configuration file?  yes
New configuration file is successfully generated in $HOME//Users/gslee/.aws/config.yaml
```

## Alias for assume role
- You can set alias with alias list.
- **You cannot use `-` prefix for alias because golang will detect it as flag.**
```bash
$ vim ~/.aws/config.yaml
- profile: default
  name: gslee@bighitcorp.com
  alias:
    d: preprod
    p: prod
  assume_roles:
    preprod: arn:aws:iam::413573177390:role/userassume-benx-preprod-admin
    prod: arn:aws:iam::056370501645:role/userassume-benx-prod-admin

# This is equal to `redhawk setup preprod`
$ redhawk setup d
Assume Credentials copied to clipboard, please paste it.
```

## RDS IAM Authentication
- You can get RDS auth token in order to log in to the database.
- If you follow these steps, then you will get your token in the clipboard

```bash
$ redhawk get rds-token
[preprod-aurora prod]
? Choose the environment:   [Use arrows to move, type to filter]
> preprod-aurora
  prod

# If you select environment
[preprod-aurora prod]
? Choose the environment:  preprod-aurora
? Choose an instance: wet  [Use arrows to move, type to filter]
> wemember-aurora-wetapne2.cluster-capmencz1lhe.ap-northeast-2.rds.amazonaws.com
  weply-aurora-wetapne2.cluster-capmencz1lhe.ap-northeast-2.rds.amazonaws.com
  weverse-aurora-wetapne2.cluster-capmencz1lhe.ap-northeast-2.rds.amazonaws.com

# If you choose instance
$ redhawk get rds-token
[preprod-aurora prod]
? Choose the environment:  preprod-aurora
? Choose an instance: wemember-aurora-wetapne2.cluster-capmencz1lhe.ap-northeast-2.rds.amazonaws.com
Assume Role MFA token code: 712352
INFO[0084] Token is copied to clipboard.
```


## Commands 
```bash
$ redhawk help
Private command line tool for beNX

managing configuration of redhawk
  init             initialize redhawk command line tool

commands related to aws IAM credentials
  renew-credential recreates aws credential of profile

commands for controlling assume role
  setup            create assume credentials for multi-account
  who              check the account information of current shell

Other Commands:
  assume           do work about assume role
  completion       Output shell completion for the given shell (bash or zsh)
  get              Get token or information with redhawk
  version          Print the version information

Usage:
  redhawk [flags] [options]

Use "redhawk <command> --help" for more information about a given command.
```

## Contribution Guide
- Check [CONTRIBUTING.md](CONTRIBUTING.md)

## Release Guide
- Check [RELEASE.md](RELEASE.md)
