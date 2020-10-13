---
title: "Configuration"
description: "test post"
date: 2020-01-28T00:34:56+09:00
draft: false
weight: -2
---

## Configurations
- You can create configuration file for redhawk.
- In order to apply the configuration in the custom config file, then please use `--config=<path/to/config.yaml>`
```
# config.yaml
#Basic Configurations
provider: aws
accounts:
  # Account : name of account alias
  # Role Arn : Assume role that is used for original account to assume to access the cross account
  - name: prerpod
    role_arn: arn:aws:iam::11111...

# Regions to scan
regions:
  - ap-northeast-2
 # - us-east-1
 # - us-east-2
 # - us-west-1
 # - us-west-2

resources:
#  - name: ec2
#  - name: security_group
#  - name: route53
#  - name: rds
#  - name: s3
  - name: iam
```


## Configuration Detail