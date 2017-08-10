package cli

import (
	"flag"
	"fmt"
	"github.com/behance/go-logging/log"
	"io"
	"io/ioutil"
)

//CreateZkNode -- create a node with data
type CreateZkNode struct {
	AclList
	File string
	Data string
	byt  []byte
}

// Usage - create usage
func (zk *CreateZkNode) Usage(writer io.Writer) {
	flags := flag.NewFlagSet("zk acls", flag.ExitOnError)
	zk.FlagSet(flags)
	fmt.Fprintln(writer, "create")
	flags.SetOutput(writer)
	flags.PrintDefaults()
	fmt.Fprintln(writer, `
		Example: with password

		./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw@172.20.0.2:2181/foo6 --debug create -acl 'perms=all scheme=digest id=user:pw2' --input-data "something else"

		Wide open:
		./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw@172.20.0.2:2181/foo7 --debug create -acl 'perms=all scheme=world id=anyone' --input-data "wide open"

		`)
}
func (zk *CreateZkNode) FlagSet(flags *flag.FlagSet) *flag.FlagSet {
	zk.AclList.FlagSet(flags)
	flags.StringVar(&zk.File, "input-file", "", "data read from file for create node operation")
	flags.StringVar(&zk.Data, "input-data", "", "data for create node operation")
	return flags
}

// Parse - create flag parse
func (zk *CreateZkNode) Parse(args []string) (exec CommandExec, err error) {
	log.Debugf(" %+v", args)
	flags := flag.NewFlagSet("zk acls", flag.ExitOnError)
	zk.FlagSet(flags)
	if err = flags.Parse(args); err != nil {
		return nil, err
	}
	if zk.File != "" {
		bb, err := ioutil.ReadFile(zk.File)
		if err != nil {
			return nil, err
		}
		zk.byt = bb
	} else if zk.Data != "" {
		zk.byt = []byte(zk.Data)
	}
	return zk, nil

}

// Execute -- create the node
func (zk *CreateZkNode) Execute(runtime *Runtime) (interface{}, error) {
	return runtime.Client.Create(runtime.ZkPath, zk.byt, 0, zk.AclList)
}
