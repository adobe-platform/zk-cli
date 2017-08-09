package cli

import (
	"flag"
	"io"
	"fmt"
	"strings"
	"errors"
	"github.com/samuel/go-zookeeper/zk"
	"time"
	"github.com/behance/go-logging/log"
)

type Runtime struct {
	zkHostRaw string
	Debug     bool
	ZkConnect string
	ZkPath string
	ZkServers []string
	Client    *zk.Conn
}

const RPC_RETRIES = 5
const RPC_TIMEOUT = time.Second * 5

func auth(raw string) (servers []string, acl []zk.ACL, err error) {
	acl = make([]zk.ACL,0)
	up := strings.Split(raw, "@")
	comp := raw
	if len(up) > 1 {

		comp = up[1]
		credList := strings.Split(up[0], ",")
		for _, cred := range credList {
			log.Debugf("unsplit cred %s", cred)
			// user:password or digest:user:password
			credPieces := strings.Split(cred, ":")
			log.Debugf("cred %s pieces %v",cred,credPieces)
			if len(credPieces) < 2 {
				return servers, acl, fmt.Errorf("Can detect creds in >%s< outer >%s<", cred, up[0])
			} else if len(credPieces) == 2 {
				acl = append(acl, zk.DigestACL(zk.PermAll, credPieces[0], credPieces[1])...)
			}  else if len(credPieces) == 3 {
				acl = append(acl, []zk.ACL{ zk.ACL{Scheme:credPieces[0], ID: strings.Join(credPieces[1:], ":")} }...)
			}  else {
				return servers, acl, fmt.Errorf("Unknown credentials >%s<", cred)
			}
		}
		servers = strings.Split(up[1], ",")
		log.Debugf("servers: %s",servers)
		return servers,acl,nil
	}

	servers = strings.Split(comp, ",")
	return servers, acl, nil
}
func ParseZKURI(zkURI string) (servers []string, chroot string, acl []zk.ACL, err error) {
	servers = []string{}

	// this is to use the canonical zk://host1:ip,host2/zkChroot format
	strippedZKConnect := strings.TrimPrefix(zkURI, "zk://")
	parts := strings.Split(strippedZKConnect, "/")
	log.Debugf("parts: %v",parts)
	if len(parts) == 2 {
		if parts[1] == "" {
			return nil, "", nil,errors.New("ZK chroot must not be the root path \"/\"!")
		}
		chroot = "/" + parts[1]
		servers, acl, err = auth(parts[0])

	} else if len(parts) == 1 {
		servers,acl,err = auth(parts[0])
	} else {
		return nil, "", nil,errors.New("ZK URI must be of the form " +
		"zk://$host1:$port1,$host2:$port2/path/to/zk/chroot")
	}
	for _, zk := range servers {
		if len(strings.Split(zk, ":")) != 2 {
			return nil, "", nil,errors.New("ZK URI must be of the form " +
			"zk://$host1:$port1,$host2:$port2/path/to/zk/chroot")
		}
	}
	return servers, chroot,acl, err
}

// FlagSet - Set up the flags
func (runtime *Runtime) FlagSet(name  string) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ExitOnError)

	flags.StringVar(&runtime.zkHostRaw, "zk-hosts", "", "Zookeeper URI of the form zk://user:passord@host1:port1,host2:port2/chroot/path  REQUIRED")
	flags.BoolVar(&runtime.Debug, "debug", false, "Turn on debug")
	return flags
}
// Usage - emit the usage
func (runtime *Runtime) Usage(writer io.Writer) {
	flags := runtime.FlagSet("<global options help>")
	flags.SetOutput(writer)
	flags.PrintDefaults()
}
// Parse - Process command line arguments
func (runtime *Runtime) Parse(args []string) (CommandExec, error) {
	flags := runtime.FlagSet("<global options> ")
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	if runtime.Debug {
		log.SetLevel(log.DebugLevel)
	}

	if runtime.zkHostRaw == "" {
		return nil, fmt.Errorf("Missing zk-hosts")
	}
	zkServers, zkChroot, acl, err := ParseZKURI(runtime.zkHostRaw)
	if err != nil {
		return nil, err
	}

	runtime.ZkServers = zkServers
	runtime.ZkPath = zkChroot
	c, _, err := zk.Connect(zkServers, RPC_TIMEOUT)
	if err != nil {
		return nil,err
	}
	for _,auth := range acl{
		c.AddAuth(auth.Scheme,[]byte(auth.ID))
	}
	runtime.Client = c

	//log.Debugf("Runtime <global flags> ok")
	// No exec returned
	return nil, nil
}

