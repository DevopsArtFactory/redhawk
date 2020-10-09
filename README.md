# redhawk
redhawk is an open source resource audit tool. When you want to get all list of infrastructure resources in cloud provider,
then you could use redhawk to easily get list of all resources. For better security, it is important to check infrastructures.

## Important Notice
- Currently redhawk only supports AWS resources.

## Install
* macOS user
```bash
```

* Linux user
```bash
```

* Windows user
  - file: 
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
  scan-resources Scan infrastructure resources in AWS

Other Commands:
  completion     Output shell completion for the given shell (bash or zsh)
  version        Print the version information

Usage:
  redhawk [flags] [options]

Use "redhawk <command> --help" for more information about a given command.
```

## Contribution Guide
- Check [CONTRIBUTING.md](CONTRIBUTING.md)
