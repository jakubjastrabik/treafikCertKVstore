package consul

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

// PutToKV insret key to storage
func PutToKV(Key, Value string) {
	kv := consul()

	fmt.Println(Key)
	fmt.Println(Value)

	// PUT a new KV pair
	p := &api.KVPair{Key: Key, Value: []byte(Value)}
	_, err := kv.Put(p, nil)
	if err != nil {
		log.Println(err)
	}
}

// DelFromKV insret key to storage
func DelFromKV(key string) {
	kv := consul()
	// Delete a KV pair
	_, err := kv.Delete(key, nil)
	if err != nil {
		log.Println(err)
	}
}

// GetFromKV insret key to storage
func GetFromKV(key string) string {
	kv := consul()
	// Lookup the pair
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		log.Println(err)
		return ""
	}

	if pair == nil {
		return ""
	}

	return string(pair.Value)

}

func consul() *api.KV {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Println(err)
	}
	kv := client.KV()

	return kv
}
