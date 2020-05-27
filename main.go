package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/pborman/uuid"
)

var (
	members = flag.String("members", "", "comma seperated list of members")
)

func init() {
	flag.Parse()
}

func start(m *memberlist.Memberlist) error {
	if len(*members) > 0 {
		parts := strings.Split(*members, ",")
		_, err := m.Join(parts)
		if err != nil {
			return err
		}
	}

	node := m.LocalNode()
	fmt.Printf("Local member %s:%d\n", node.Addr, node.Port)

	// Ask for members of the cluster
	for _, member := range m.Members() {
		fmt.Printf("Member: %s %s\n", member.Name, member.Addr)
	}

	return nil
}

func main() {
	hostname, _ := os.Hostname()

	c := memberlist.DefaultLocalConfig()
	c.BindPort = 7902
	c.AdvertisePort = 7902
	c.Name = hostname + "-" + uuid.NewUUID().String()

	m, err := memberlist.Create(c)
	if err != nil {
		fmt.Println(err)
	}

	for {
		if err := start(m); err != nil {
			fmt.Println(err)
		}

		log.Printf("wait 10 sec")
		time.Sleep(10 * time.Second)
	}

}
