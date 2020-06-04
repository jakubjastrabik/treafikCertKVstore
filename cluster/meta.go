package cluster

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/memberlist"
)

//MyMetaData comment
type MyMetaData struct {
	Location string `json:"location"`
	Zone     string `json:"zone"`
	ShardID  int    `json:"shard-id"`
	Weight   int    `json:"weight"`
}

//MyDelegate comment
type MyDelegate struct {
	meta MyMetaData
}

//NodeMeta comment
func (d *MyDelegate) NodeMeta(limit int) []byte {
	return d.meta.Bytes()
}

//GetBroadcasts comment
func (d *MyDelegate) GetBroadcasts(overhead, limit int) [][]byte {
	// not use, noop
	return nil
}

//LocalState comment
func (d *MyDelegate) LocalState(join bool) []byte {
	// not use, noop
	return []byte("")
}

//MergeRemoteState comment
func (d *MyDelegate) MergeRemoteState(buf []byte, join bool) {
	// not use
}

//NotifyMsg comment
func (d *MyDelegate) NotifyMsg(msg []byte) {
	// not use
}

//Bytes comment
func (m MyMetaData) Bytes() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return []byte("")
	}
	return data
}

//ParseMyMetaData comment
func ParseMyMetaData(data []byte) (MyMetaData, bool) {
	meta := MyMetaData{}
	if err := json.Unmarshal(data, &meta); err != nil {
		return meta, false
	}
	return meta, true
}

//MyEventDelegate comment
type MyEventDelegate struct {
}

//NotifyJoin comment
func (d *MyEventDelegate) NotifyJoin(node *memberlist.Node) {
	hostPort := fmt.Sprintf("%s:%d", node.Addr.To4().String(), node.Port)
	log.Printf("New node join: %s", hostPort)

	NodeState = clusterStatus(node, "add")
}

//NotifyLeave comment
func (d *MyEventDelegate) NotifyLeave(node *memberlist.Node) {
	hostPort := fmt.Sprintf("%s:%d", node.Addr.To4().String(), node.Port)
	log.Printf("Node leave cluster: %s", hostPort)

	NodeState = clusterStatus(node, "del")
}

//NotifyUpdate comment
func (d *MyEventDelegate) NotifyUpdate(node *memberlist.Node) {
	hostPort := fmt.Sprintf("%s:%d", node.Addr.To4().String(), node.Port)
	log.Printf("Node will be Upadated: %s", hostPort)

	NodeState = clusterStatus(node, "edit")
}
