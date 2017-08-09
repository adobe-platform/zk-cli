package cli

import (
	"github.com/behance/go-logging/log"
	"flag"
	"io"
	"fmt"
)

type ZkGet struct {
	convertToString bool //AclList
	convertToHex    bool //AclList
}

func ( zk *ZkGet) FlagSet(flags *flag.FlagSet) *flag.FlagSet {
	flags.BoolVar(&zk.convertToString, "as-string", false, `convert the node value to a string`)
	flags.BoolVar(&zk.convertToHex, "as-hex", false, `convert the node value to a hex`)
	return flags
}

func (zk *ZkGet) Usage(writer io.Writer) {
	flags := flag.NewFlagSet("zk get", flag.ExitOnError)
	zk.FlagSet(flags)
	fmt.Fprintln(writer,"get")

	flags.SetOutput(writer)
	flags.PrintDefaults()
	fmt.Fprintln(writer,`
	Example:  ./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw2@172.20.0.2:2181/foo6 -debug get --as-string
	INFO[0000] result "something else"

	Example - no auth
	./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw2@172.20.0.2:2181/foo7 -debug get --as-string
	INFO[0000] result "wide open"
	  `)

}
func (zk *ZkGet) Parse(args []string) (exec CommandExec, err error) {
	log.Debugf(" %+v", args)

	flags := flag.NewFlagSet("zk get", flag.ExitOnError)
	zk.FlagSet(flags)
	if err = flags.Parse(args); err != nil {
		return nil, err
	}
	return zk, nil
}

func (zk *ZkGet) Execute(runtime *Runtime) (interface{}, error) {
	log.Debugf("zkGet.Execute %+v", runtime)
	data, stat, err := runtime.Client.Get(runtime.ZkPath)
	if err != nil {
		return nil, err
	}
	log.Debug(stat)
	if zk.convertToString {
		return string(data), nil
	} else if zk.convertToHex{
		return fmt.Sprintf("%X",data),nil
	}
	return data,nil
}

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
	fmt.Fprintln(writer,`
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
	fmt.Fprintln(writer,`
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
	stat, err = runtime.Client.SetACL(runtime.ZkPath, zk.AclList,stat.Version)
	if err != nil {
		return nil, err
	}
	log.Debugf("create %s -> %#v", runtime.ZkPath,stat)
	return stat, nil
}


