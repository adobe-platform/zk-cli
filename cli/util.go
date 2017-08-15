package cli

import (
	"fmt"
	"github.com/behance/go-logging/log"
	"github.com/samuel/go-zookeeper/zk"
)

func GetChildren(path string, client *zk.Conn) (tree []string, err error) {

	queue := make([]string, 0)
	tree = make([]string, 0)
	queue = append(queue, path)
	tree = append(tree, path)
	for {
		node := func() string {
			if len(queue) > 0 {
				retval := queue[0]
				queue = queue[1:]
				return retval
			}
			return ""
		}()
		if node == "" {
			break
		}
		children, _, err := client.Children(node)
		log.Debugf("%s children: %v", node, children)
		if err != nil {
			log.Errorf("children %s error %v", node, err)
			continue
			//return nil,err
		}
		for _, child := range children {
			childPath := func() string {
				if node != "/" {
					return fmt.Sprintf("%s/%s", node, child)
				}
				return fmt.Sprintf("/%s", child)
			}()
			log.Debugf("push %s", childPath)
			tree = append(tree, childPath)
			queue = append(queue, childPath)
		}

	}
	log.Debugf("%s children %v", path, tree)
	return tree, nil
}
func Delete(path string, client *zk.Conn) (err error) {
	err = client.Delete(path, -1)
	if err != nil {
		return fmt.Errorf("Can delete %s %v", path, err)
	}

	return nil
}
