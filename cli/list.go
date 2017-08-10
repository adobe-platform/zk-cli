package cli

import (
	"flag"
	"fmt"
	"github.com/behance/go-logging/log"
	"github.com/samuel/go-zookeeper/zk"
	"io"
	"strings"
)

type ZkLs struct {
	recursive bool //AclList
	acls      bool
}

func (zk *ZkLs) FlagSet(flags *flag.FlagSet) *flag.FlagSet {
	flags.BoolVar(&zk.recursive, "recursive", false, `recursively retrieve children`)
	flags.BoolVar(&zk.acls, "with-acls", false, `with acls`)
	return flags
}

func (zk *ZkLs) Usage(writer io.Writer) {
	flags := flag.NewFlagSet("zk ls", flag.ExitOnError)
	zk.FlagSet(flags)
	fmt.Fprintln(writer, "ls")

	flags.SetOutput(writer)
	flags.PrintDefaults()
	fmt.Fprintln(writer, `
	Example:  ./zk-cli-linux-amd64 --zk-hosts zk://digest:user:pw2@172.20.0.2:2181/foo6 -debug ls  --recursive
	  `)

}
func (zk *ZkLs) Parse(args []string) (exec CommandExec, err error) {
	log.Debugf(" %+v", args)

	flags := flag.NewFlagSet("zk ls", flag.ExitOnError)
	zk.FlagSet(flags)
	if err = flags.Parse(args); err != nil {
		return nil, err
	}
	return zk, nil
}

func (zz *ZkLs) Execute(runtime *Runtime) (interface{}, error) {
	log.Debugf("zkLs.Execute %+v", runtime)
	children, err := func() ([]string, error) {
		if zz.recursive {
			return GetChildren(runtime.ZkPath, runtime.Client)
		}
		data, _, err2 := runtime.Client.Children(runtime.ZkPath)
		if err2 != nil {
			return nil, err2
		}
		return data, nil
	}()
	if err != nil {
		return nil, err
	}
	if zz.acls {
		vchildren := make([]string, len(children))
		for ii, child := range children {
			acls, _, err := runtime.Client.GetACL(child)
			if err != nil {
				vchildren[ii] = fmt.Sprintf("%s -- unknown %s", child, err)
				continue
			}
			aclSet := make([]string, len(acls))
			perms := ""
			for mm, acl := range acls {
				log.Debugf("Considering %v", acls)
				for kk := 1; kk <= zk.PermAll; kk = kk << 1 {
					switch val := acl.Perms & int32(kk); val {
					case zk.PermAdmin:
						perms += "a"
					case zk.PermRead:
						perms += "r"
					case zk.PermWrite:
						perms += "w"
					case zk.PermDelete:
						perms += "d"
					case zk.PermCreate:
						perms += "c"
					}
				}

				aclSet[mm] = fmt.Sprintf("%s:%s:%s:", acl.Scheme, acl.ID, perms)
			}
			vchildren[ii] = fmt.Sprintf("%s -- %s", child, strings.Join(aclSet, ","))
		}
		return GenericResult{Success: true,
			Message: fmt.Sprintf("children of %s ", runtime.ZkPath),
			Data:    vchildren,
		}, nil
	}

	return GenericResult{Success: true,
		Message: fmt.Sprintf("children of %s ", runtime.ZkPath),
		Data:    children,
	}, nil
}
