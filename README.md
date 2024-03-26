aws-describer
=============

Tool that displays information on various AWS resources in a tabular format combined with information on other related resources.
Covers all valid regions by default.

Usage
-----

```text
NAME:
   aws-describer - AWS resources describer CLI

USAGE:
   aws-describer [global options] command [command options] [arguments...]

VERSION:
   0.0.0

DESCRIPTION:
   A cli application to join and list AWS resources with various other resources

COMMANDS:
   completion  Generate completion scripts: bash|zsh|pwsh
   ec2         Invoke EC2 API and list resources
   iam         Invoke IAM API and list resources
   s3          Invoke S3 API and list resources

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

EC2 commands

```text
NAME:
   aws-describer ec2 - Invoke EC2 API and list resources

USAGE:
   aws-describer ec2 command

DESCRIPTION:
   Invoke EC2 API and list resources in various output formats

COMMANDS:
   get-instances        List EC2 instance info
   get-images           List EC2 image info
   get-security-groups  List EC2 security group info
   get-vpcs             List EC2 VPC info
   get-subnets          List EC2 subnet info
   get-route-tables     List EC2 route table info

OPTIONS:
   --help, -h  show help
```

IAM commands

```text
NAME:
   aws-describer iam - Invoke IAM API and list resources

USAGE:
   aws-describer iam command

DESCRIPTION:
   Invoke IAM API and list resources in various output formats

COMMANDS:
   get-users     List IAM user info
   get-groups    List IAM group info
   get-roles     List IAM role info
   get-policies  List IAM policy info

OPTIONS:
   --help, -h  show help
```

S3 commands

```text
NAME:
   aws-describer s3 - Invoke S3 API and list resources

USAGE:
   aws-describer s3 command

DESCRIPTION:
   Invoke S3 API and list resources in various output formats

COMMANDS:
   get-buckets  List S3 bucket info

OPTIONS:
   --help, -h  show help
```

Example
-------

Combine and output security group information for EC2 instances.

```text
$ aws-describer ec2 get-instances --join sg --output compressed
+---------------------+------------------+-----------------------+----------------+----------------------+-------------------+---------------+------------+----------+--------+---------------+-----------------------------------+------------------+
| InstanceId          | InstanceName     | VpcId                 | VpcName        | SecurityGroupId      | SecurityGroupName | FlowDirection | IpProtocol | FromPort | ToPort | AddressType   | CidrBlock                         | AvailabilityZone |
+---------------------+------------------+-----------------------+----------------+----------------------+-------------------+---------------+------------+----------+--------+---------------+-----------------------------------+------------------+
| i-xxxxxxxxxxxxxxxxx | test-instance-01 | vpc-11111111111111111 | vpc-01         | sg-xxxxxxxxxxxxxxxxx | sg-01             | Ingress       | icmp       |       -1 |     -1 | Ipv4          | aaa.bbb.ccc.ddd/32                | ap-northeast-1a  |
|                     |                  |                       |                |                      |                   | Ingress       | icmp       |       -1 |     -1 | SecurityGroup | sg-bbbbbbbbbbbbbbbbb/000000000000 | ap-northeast-1a  |
|                     |                  |                       |                |                      |                   | Ingress       | tcp        |        0 |  65535 | SecurityGroup | sg-bbbbbbbbbbbbbbbbb/000000000000 | ap-northeast-1a  |
|                     |                  |                       |                |                      |                   | Ingress       | tcp        |       22 |     22 | Ipv4          | aaa.bbb.ccc.ddd/32                | ap-northeast-1a  |
|                     |                  |                       |                |                      |                   | Ingress       | tcp        |       80 |     80 | Ipv4          | aaa.bbb.ccc.ddd/32                | ap-northeast-1a  |
|                     |                  |                       |                |                      |                   | Ingress       | tcp        |       80 |     80 | SecurityGroup | sg-ccccccccccccccccc/000000000000 | ap-northeast-1a  |
|                     |                  |                       |                |                      |                   | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0                         | ap-northeast-1a  |
+---------------------+------------------+-----------------------+----------------+----------------------+-------------------+---------------+------------+----------+--------+---------------+-----------------------------------+------------------+
| i-yyyyyyyyyyyyyyyyy | test-instance-02 | vpc-22222222222222222 | vpc-02         | sg-yyyyyyyyyyyyyyyyy | sg-02             | Ingress       | tcp        |       22 |     22 | Ipv4          | aaa.bbb.ccc.ddd/32                | ap-northeast-1a  |
|                     |                  |                       |                |                      |                   | Ingress       | tcp        |       80 |     80 | Ipv4          | aaa.bbb.ccc.ddd/32                | ap-northeast-1a  |
|                     |                  |                       |                |                      |                   | Ingress       | tcp        |      443 |    443 | Ipv4          | aaa.bbb.ccc.ddd/32                | ap-northeast-1a  |
|                     |                  |                       |                |                      |                   | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0                         | ap-northeast-1a  |
+---------------------+------------------+-----------------------+----------------+----------------------+-------------------+---------------+------------+----------+--------+---------------+-----------------------------------+------------------+
| i-zzzzzzzzzzzzzzzzz | test-instance-03 | vpc-33333333333333333 | vpc-03         | sg-zzzzzzzzzzzzzzzzz | sg-03             | Ingress       | tcp        |       22 |     22 | Ipv4          | aaa.bbb.ccc.ddd/32                | us-east-2a       |
|                     |                  |                       |                |                      |                   | Ingress       | tcp        |       80 |     80 | Ipv4          | aaa.bbb.ccc.ddd/32                | us-east-2a       |
|                     |                  |                       |                |                      |                   | Ingress       | tcp        |      443 |    443 | Ipv4          | aaa.bbb.ccc.ddd/32                | us-east-2a       |
|                     |                  |                       |                |                      |                   | Egress        |         -1 |        0 |      0 | Ipv4          | 0.0.0.0/0                         | us-east-2a       |
|                     |                  |                       |                | sg-aaaaaaaaaaaaaaaaa | sg-04             | Egress        | tcp        |     3306 |   3306 | SecurityGroup | sg-ddddddddddddddddd/000000000000 | us-east-2a       |
+---------------------+------------------+-----------------------+----------------+----------------------+-------------------+---------------+------------+----------+--------+---------------+-----------------------------------+------------------+
```

Complete list of IAM user policies.

```text
$ aws-describer iam get-users --join assoc --document
+----------+----------------------------------------+------------+----------------------+-------------------------------------------------------------------------+
| UserName | AttachedBy                             | PolicyType | PolicyName           | PolicyDocument                                                          |
+----------+----------------------------------------+------------+----------------------+-------------------------------------------------------------------------+
| user1    | arn:aws:iam::000000000000:user/user1   | Attached   | AddministratorAccess | SKIPPED                                                                 |
+----------+----------------------------------------+------------+----------------------+-------------------------------------------------------------------------+
| user2    | arn:aws:iam::000000000000:user/user2   | -          | -                    | -                                                                       |
+          +----------------------------------------+------------+----------------------+-------------------------------------------------------------------------+
|          | arn:aws:iam::000000000000:group/group1 | Attached   | ReadOnlyAccess       | SKIPPED                                                                 |
+          +                                        +------------+----------------------+-------------------------------------------------------------------------+
|          |                                        | Inline     | inline-policy-01     | {                                                                       |
|          |                                        |            |                      |   "Version": "2012-10-17",                                              |
|          |                                        |            |                      |   "Statement": [                                                        |
|          |                                        |            |                      |     {                                                                   |
|          |                                        |            |                      |       "Effect": "Allow",                                                |
|          |                                        |            |                      |       "Action": [                                                       |
|          |                                        |            |                      |         "s3:*",                                                         |
|          |                                        |            |                      |       ],                                                                |
|          |                                        |            |                      |       "Resource": "*"                                                   |
|          |                                        |            |                      |     }                                                                   |
|          |                                        |            |                      |   ]                                                                     |
|          |                                        |            |                      | }                                                                       |
+----------+----------------------------------------+------------+----------------------+-------------------------------------------------------------------------+
| user3    | arn:aws:iam::000000000000:user/user3   | Attached   | attached-policy-01   | {                                                                       |
|          |                                        |            |                      |   "Version": "2012-10-17",                                              |
|          |                                        |            |                      |   "Statement": [                                                        |
|          |                                        |            |                      |     {                                                                   |
|          |                                        |            |                      |       "Effect": "Allow",                                                |
|          |                                        |            |                      |       "Action": [                                                       |
|          |                                        |            |                      |         "firehose:*"                                                    |
|          |                                        |            |                      |       ],                                                                |
|          |                                        |            |                      |       "Resource": [                                                     |
|          |                                        |            |                      |         "arn:aws:firehose:ap-northeast-1:000000000000:deliverystream/*" |
|          |                                        |            |                      |       ]                                                                 |
|          |                                        |            |                      |     }                                                                   |
|          |                                        |            |                      |   ]                                                                     |
|          |                                        |            |                      | }                                                                       |
+----------+----------------------------------------+------------+----------------------+-------------------------------------------------------------------------+
```

Todo
----

- Write test
- Write document

Author
------

[nekrassov01](https://github.com/nekrassov01)

License
-------

[MIT](https://github.com/nekrassov01/aws-describer/blob/main/LICENSE)
