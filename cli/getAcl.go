package cli

import (
	"flag"
	"fmt"
	"github.com/behance/go-logging/log"
	"io"
)

type ZkGetAcl struct {
}

func (zk *ZkGetAcl) Parse(args []string) (exec CommandExec, err error) {
	log.Debugf(" %+v", args)
	/*
		flags := flag.NewFlagSet("zk acls", flag.ExitOnError)
		zk.FlagSet(flags)
		if err = flags.Parse(args); err != nil {
			return nil,err
		}
	*/
	return zk, nil
}

func (zk *ZkGetAcl) Usage(writer io.Writer) {
	flags := flag.NewFlagSet("zk getAcls", flag.ExitOnError)
	// zk.FlagSet(flags)
	fmt.Fprintln(writer, "getAcl ")
	flags.SetOutput(writer)
	flags.PrintDefaults()
	fmt.Fprintln(writer, `
	Example: ./zk-cli-linux-amd64 --zk-hosts zk://172.20.0.2:2181/foo7 -debug getAcl
	 {
    		"Perms": 31,
    		"Scheme": "world",
    		"ID": "anyone"
  	 }
	`)

}

func (zk *ZkGetAcl) Execute(runtime *Runtime) (interface{}, error) {
	log.Debugf("zkGetAcl.Execute %+v", runtime)
	acls, stat, err := runtime.Client.GetACL(runtime.ZkPath)
	if err != nil {
		return nil, err
	}
	log.Debug(stat)
	return acls, nil
}
