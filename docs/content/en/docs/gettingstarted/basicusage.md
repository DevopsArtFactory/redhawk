---
title: "Basic Usage"
description: "test post"
date: 2020-01-28T00:34:51+09:00
draft: false
weight: -3
---

## Basic Usage

{{< notice warning "Possible Provider" >}}
Currently, we are supporting AWS provider. So, you need AWS credentials set up in the local environment.
{{< /notice >}}


## Commands
Audit:
* [redhawk list](#redhawk-list) - to gather data of infrastructure resources.

### Redhawk List 
- In order to find resources, you need to specify resources with `--resources`.
- By default, `redhawk` will show you the result on Stdout. 
- You can simply change output format with `--output, -o` .
```
Usage:
  redhawk list [flags] [options]

Example:
  # Find resources in detailed version
  - redhawk list --resources=ec2 --detail

  # Find ec2 resource only and region set to Seoul
  - redhawk list --resources=ec2 --region=ap-northeast-2

  # Find ec2,iam resources and region set to all(all region)
  - redhawk list --resources=ec2,iam --all

  # Find ec2,iam resources and region set to us-west-2 and output set to csv
  - redhawk list --resources=ec2,iam --region=us-west-2 -o csv

```

### Result Sample
```
$ redhawk list --resources=ec2,iam
INFO[0000] start scanning resources
PROVIDER: aws
==============================================
SERVICE   NAME           ID                    STATUS    TYPE        AZ                LAUNCHED
EC2       ec2            i-0e2d26c29388283     stopped   t3a.small   ap-northeast-2a   2020-09-15 13:49:41 &#43;0000 UTC
==============================================
SERVICE     NAME               USER_COUNT   USERS                        GROUP_POLICIES
IAM_GROUP   art_devops_black   3            asbubam|jupiter.song|gslee   SelfManageMFA|RotateKeys|assume-art-prod-admin-policy|ForceMFA
IAM_GROUP   devops             2            juyoung.song|gildong.hong    AmazonEC2FullAccess
==============================================
SERVICE    NAME                                       TRUST_ENTITIES                       ROLE_LAST_ACTIVITY
IAM_ROLE   AWSServiceRoleForAWSCloud9                 cloud9.amazonaws.com
IAM_ROLE   AWSServiceRoleForSupport                   support.amazonaws.com
IAM_ROLE   AWSServiceRoleForAmazonElasticFileSystem   elasticfilesystem.amazonaws.com
IAM_ROLE   jenkins                                    ec2.amazonaws.com
IAM_ROLE   admin                                      ec2.amazonaws.com
IAM_ROLE   AWSServiceRoleForEC2Spot                   spot.amazonaws.com
IAM_ROLE   codebuild-deployment                       codebuild.amazonaws.com
IAM_ROLE   AWSServiceRoleForElasticLoadBalancing      elasticloadbalancing.amazonaws.com
IAM_ROLE   AWSServiceRoleForOrganizations             organizations.amazonaws.com
IAM_ROLE   AWSServiceRoleForGlobalAccelerator         globalaccelerator.amazonaws.com
IAM_ROLE   app-hello                                  ec2.amazonaws.com
IAM_ROLE   app-hello2                                 ec2.amazonaws.com
IAM_ROLE   AWSServiceRoleForAutoScaling               autoscaling.amazonaws.com
IAM_ROLE   AWSServiceRoleForTrustedAdvisor            trustedadvisor.amazonaws.com
IAM_ROLE   hello-iam-role                             ec2.amazonaws.com
==============================================
SERVICE    NAME                MFA                                   GROUP_COUNT   ACCESS_KEY_LAST_USED                CREATED
IAM_USER   readonly@art.com                                          1                                                 2020-06-12 10:12:05 &#43;0000 UTC
IAM_USER   gildong.hong                                              1                                                 2020-08-30 18:57:43 &#43;0000 UTC
IAM_USER   gslee               arn:aws:iam::816736805842:mfa/gslee   1             2020-10-11 14:06:00 &#43;0000 UTC   2020-06-12 10:12:05 &#43;0000 UTC
```