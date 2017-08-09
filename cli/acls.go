package cli

import (
	"flag"
	"github.com/samuel/go-zookeeper/zk"
	"fmt"
	"strings"
	"github.com/behance/go-logging/log"
	"crypto/sha1"
	"encoding/base64"
)
// ArtifactList - thin type providing Flags Value interface implementation for Metronome artifacts
type AclList  []zk.ACL

// String - Value interface implementation
func (list *AclList) String() string {
	return fmt.Sprintf("%v", *list)
}

func GetDigest(user string, password string) (string, error) {
	userPass := []byte(fmt.Sprintf("%s:%s", user, password))
	h := sha1.New()
	if n, err := h.Write(userPass); err != nil || n != len(userPass) {
		panic("SHA1 failed")
	}
	digest:=base64.StdEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s:%s", user, digest), nil
}

// Set - Value interface implemention
func (list *AclList) Set(value string) (err error) {
	var arty zk.ACL

	for _, pairs := range strings.Split(strings.TrimSpace(value), " ") {
		log.Debugf("pairs : %+v", pairs)
		kv := strings.SplitN(strings.TrimSpace(pairs), "=", 2)
		log.Debugf("kv=%+v", kv)
		switch item := strings.TrimSpace(kv[0]); item{
		case "perms":
			if kv[1] == "all" {
				arty.Perms = zk.PermAll
				continue
			}
			for _, perm := range strings.Split(strings.TrimSpace(kv[1]), "") {
				log.Debugf("perm : %+v", pairs)
				switch perm {
				case "r": arty.Perms |= zk.PermRead
				case "w": arty.Perms |= zk.PermWrite
				case "c": arty.Perms |= zk.PermCreate
				case "d": arty.Perms |= zk.PermDelete
				case "a": arty.Perms |= zk.PermAdmin
				default:
					return fmt.Errorf("Unknown perm >%s<", perm)
				}
			}

		case "id":
			arty.ID = kv[1]

		case "scheme":
			switch kv[1] {
			case "world", "auth", "digest", "ip","host":
				arty.Scheme = kv[1]
			default:
				return fmt.Errorf("Unknown scheme %s", kv[1])

			}
		}

	}
	if arty.Scheme == "digest" {
		up := strings.Split(arty.ID, ":")
		if len(up) > 1 {
			rez:=zk.DigestACL(arty.Perms,up[0], up[1])
			arty = rez[0]
		}
	}

	if arty.Perms == 0 {
		return fmt.Errorf("Unknown permissions")
	}
	log.Debugf("permissions: %#v",arty)
	*list = append(*list, arty)
	return nil
}


func ( zz *AclList) FlagSet(flags *flag.FlagSet) *flag.FlagSet {
	flags.Var(zz, "acl",
		`perms=(all|(r|w|c|d|a)" scheme=(world|digest|auth|ip)  id=(anyone|id).
			if scheme is digest and id has ':' then id value is assumed to be cleartext user:password which digested with sha1`)

	return flags
}
