package cli

import (
	"flag"
	"fmt"
	"github.com/behance/go-logging/log"
	"io"
)

type ZkSetAcl struct {
	AclList
}

func (zk *ZkSetAcl) Parse(args []string) (exec CommandExec, err error) {
	log.Debugf(" %+v", args)
	flags := flag.NewFlagSet("zk acls", flag.ExitOnError)
	zk.FlagSet(flags)
	if err = flags.Parse(args); err != nil {
		return nil, err
	}
	return zk, nil
}

func (zk *ZkSetAcl) Usage(writer io.Writer) {
	flags := flag.NewFlagSet("zk getAcls", flag.ExitOnError)
	zk.FlagSet(flags)
	fmt.Fprintln(writer, "setAcl ")
	flags.SetOutput(writer)
	flags.PrintDefaults()
	fmt.Fprintln(writer, `
	Example - add acls to wide open
	./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw@172.20.0.2:2181/foo7 --debug setAcl -acl 'perms=all scheme=digest id=user:pw'
	`)
}

func (zk *ZkSetAcl) Execute(runtime *Runtime) (interface{}, error) {
	log.Debugf("zkGetAcl.Execute %+v", runtime)
	_, stat, err := runtime.Client.GetACL(runtime.ZkPath)
	if err != nil {
		return nil, err
	}
	stat, err = runtime.Client.SetACL(runtime.ZkPath, zk.AclList, stat.Version)
	if err != nil {
		return nil, err
	}
	log.Debugf("create %s -> %#v", runtime.ZkPath, stat)
	return stat, nil
}
