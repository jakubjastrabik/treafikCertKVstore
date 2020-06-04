package cluster

import (
	"log"
	"os"
	"strings"

	"github.com/hashicorp/memberlist"
)

var (
	// NodeState tag node as master/backup
	NodeState string
	// MembersL tag node as master/backup
	MembersL    memberlist.Memberlist
	hostname, _ = os.Hostname()
	globCluster *Cluster
	a           []*memberlist.Node
)

// MetaData information
type MetaData struct {
	Location string
	Zone     string
	ShardID  int
	Weight   int
}

// Cluster contain all variables
type Cluster struct {
	NodePort int
	MyData   []MetaData
	Members  string
}

type nodes struct {
	node *memberlist.Node
}

func start(m *memberlist.Memberlist, cluster Cluster) error {
	if len(cluster.Members) > 0 {
		parts := strings.Split(cluster.Members, ",")
		_, err := m.Join(parts)
		if err != nil {
			return err
		}
	}

	node := m.LocalNode()
	log.Printf("Local member %s:%d\n", node.Addr, node.Port)

	return nil
}

// InitCluster function create cluster
// Handle cluster Comunication
// Etc. Notify, metadata ...
func InitCluster(cluster Cluster) {
	globCluster = &cluster
	d := new(MyDelegate)
	d.meta = MyMetaData(cluster.MyData[0])

	config := memberlist.DefaultLocalConfig()
	config.BindPort = cluster.NodePort
	config.AdvertisePort = cluster.NodePort
	config.Name = hostname
	config.Delegate = d
	config.Events = new(MyEventDelegate)

	MembersL, err := memberlist.Create(config)

	if err != nil {
		log.Println(err)
	}

	if err := start(MembersL, cluster); err != nil {
		log.Println(err)
	}

}

// clusterStatus return master node
// Return Node information
func clusterStatus(node *memberlist.Node, option string) (state string) {
	switch option {
	case "add":
		state = clusterAdd(node)
	case "del":
		state = clusterDel(node)
	case "edit":
		state = clusterEdir(node)
	}

	log.Println("Node has status: " + state)

	return state
}

func clusterAdd(node *memberlist.Node) (state string) {
	a = append(a, node)

	return checkCluster(node)
}
func clusterDel(node *memberlist.Node) (state string) {

	for i := range a {
		if a[i].String() == node.String() {

			a[len(a)-1], a[i] = a[i], a[len(a)-1]
			a = a[:len(a)-1]
		}
	}

	return checkCluster(node)
}

func clusterEdir(node *memberlist.Node) (state string) {

	for i := range a {
		if a[i].String() == node.String() {
			oldData, _ := ParseMyMetaData(a[i].Meta)
			freshData, _ := ParseMyMetaData(node.Meta)

			if oldData.Weight != freshData.Weight {
				a[i].Meta = node.Meta
			}
		}
	}

	return checkCluster(node)
}

func checkCluster(node *memberlist.Node) (state string) {
	masterNode := 0
	maxWeig := 0

	for i := range a {
		data, _ := ParseMyMetaData(a[i].Meta)

		if data.Weight > maxWeig {
			maxWeig = data.Weight
			masterNode = data.ShardID
		}
	}

	if masterNode != globCluster.MyData[0].ShardID {
		state = "backup"
	} else {
		state = "master"
	}

	return state
}
