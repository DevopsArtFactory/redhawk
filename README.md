# redhawk
redhawk is an open source resource audit tool. When you want to get all list of infrastructure resources in cloud provider,
then you could use redhawk to easily get list of all resources. For better security, it is important to check infrastructures.

## Important Notice
- Currently redhawk only supports AWS resources.

## Install
* macOS user
```bash
# Using cURL
curl -Lo redhawk https://devopsartfactory.s3.ap-northeast-2.amazonaws.com/redhawk/releasees/latest/redhawk-linux-amd64
sudo install redhawk /usr/local/bin/
redhawk version

# Using brew
brew tap devopsartfactory/devopsart
brew install redhawk
redhawk version
```

* Linux user
```bash
curl -Lo redhawk https://devopsartfactory.s3.ap-northeast-2.amazonaws.com/redhawk/releasees/latest/redhawk-linux-amd64
sudo install redhawk /usr/local/bin/
redhawk version
```

* Windows user
  - file: https://devopsartfactory.s3.ap-northeast-2.amazonaws.com/redhawk/releasees/latest/redhawk-windows-amd64.exe
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

## Commands 
```bash
Opensource cloud resources audit and management tool

checking all resources in cloud provider
  list        list infrastructure resources in AWS

Other Commands:
  completion  Output shell completion for the given shell (bash or zsh)
  version     Print the version information

Usage:
  redhawk [flags] [options]

Use "redhawk <command> --help" for more information about a given command.
```

## Contribution Guide
- Check [CONTRIBUTING.md](CONTRIBUTING.md)
