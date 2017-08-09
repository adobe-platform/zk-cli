# 
zkcli is a pure golang zookeeper cli capable of executing commands with requiring auth in a very tiny footprint

Currently, this cli can create, rm, get, setAcl, getAcl, ls

> ls & rm support recursive

## Features
- Pure go build.  You don't need to install go
- Only 16M docker image!
- Can do ACLs and auth.
## Build

```
$  `git clone https://github.com/adobe-platform/zk-cli
$  cd zk-cli
$  make upload-container DOCKER_CONFIG=`cd ~/.docker-hub-ethos/;pwd` PUBLISH_TAG=adobeplatform
```
## Usage
```
docker run -it adobeplatform/zk-cli:0.2

USAGE

         ./zk-cli-linux-amd64 <global-options>  {create|rm|get|getAcl|setAcl|help} [<action options>|help]

COMMANDS:

create
  -acl value
          perms=(all|(r|w|c|d|a)" scheme=(world|digest|auth|ip)  id=(anyone|id).
                        if scheme is digest and id has ':' then id value is assumed to be cleartext user:password which digested with sha1
  -input-data string
        data for create node operation
  -input-file string
        data read from file for create node operation

                Example: with password

                ./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw@172.20.0.2:2181/foo6 --debug create -acl 'perms=all scheme=digest id=user:pw2' --input-data "something else"

                Wide open:
                ./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw@172.20.0.2:2181/foo7 --debug create -acl 'perms=all scheme=world id=anyone' --input-data "wide open"


rm
  -recursive
        convert the node value to a string

        Example:  ./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw2@172.20.0.2:2181/foo6 -debug rm  --recursive

get
  -as-hex
        convert the node value to a hex
  -as-string
        convert the node value to a string

        Example:  ./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw2@172.20.0.2:2181/foo6 -debug get --as-string
        INFO[0000] result "something else"

        Example - no auth
        ./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw2@172.20.0.2:2181/foo7 -debug get --as-string
        INFO[0000] result "wide open"

getAcl  

        Example: ./zk-cli-linux-amd64 --zk-hosts zk://172.20.0.2:2181/foo7 -debug getAcl
         {
                "Perms": 31,
                "Scheme": "world",
                "ID": "anyone"
         }

setAcl
  -acl value
        perms=(all|(r|w|c|d|a)" scheme=(world|digest|auth|ip)  id=(anyone|id).
	                        if scheme is digest and id has ':' then id value is assumed to be cleartext user:password which digested with sha1

        Example - add acls to wide open
	        ./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw@172.20.0.2:2181/foo7 --debug setAcl -acl 'perms=all scheme=digest id=user:pw'


GLOBAL OPTIONS:
 The zk target path is taken --zk-hosts for most commands

  -debug
          Turn on debug
	    -zk-hosts string
	            Zookeeper URI of the form zk://user:passord@host1:port1,host2:port2/chroot/path  REQUIRED
```
