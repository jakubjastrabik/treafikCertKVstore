package main

import (
	"flag"
	"time"

	"github.com/jakubjastrabik/Ha-cert-manager-for-traefik/cluster"
)

var (
	members  = flag.String("members", "", "comma seperated list of members")
	nodePort = flag.Int("nodePort", 7900, "Port to be use for connection")
	location = flag.String("nodeLocation", "", "Node location")
	zone     = flag.String("nodeZone", "", "Node zone name")
	nodeID   = flag.Int("nodeID", 0, "Node ID")
	weight   = flag.Int("nodeWeight", 0, "Node weight")
)

func init() {
	flag.Parse()
}

func main() {

	clusterConf := cluster.Cluster{
		Members:  *members,
		NodePort: *nodePort,
		MyData: []cluster.MetaData{
			{
				Location: *location,
				Zone:     *zone,
				ShardID:  *nodeID,
				Weight:   *weight,
			},
		},
	}

	cluster.InitCluster(clusterConf)

	for i := 3; i > 0; i-- {
		time.Sleep(10 * time.Second)
	}
}
