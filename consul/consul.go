package consul

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/hashicorp/consul/api"
)

var (
	serviceName = "_consul_"
	conPut      = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert" + serviceName + "put",
		Help: "Total count of PUT requests stored into KV storage",
	})
	conPutErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert" + serviceName + "put_error",
		Help: "Total count of errored PUT requests not stored into KV storage",
	})
	conGet = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert" + serviceName + "get",
		Help: "Total count of GET requests stored into KV storage",
	})
	conWhiteGet = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert" + serviceName + "empty_get",
		Help: "Total count of GET non exist KV in KV storage",
	})
	conGetErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert" + serviceName + "get_error",
		Help: "Total count of errored GET requests not stored into KV storage",
	})
	conDel = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert" + serviceName + "del",
		Help: "Total count of DELETE requests stored into KV storage",
	})
	conDelErr = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert" + serviceName + "del_error",
		Help: "Total count of errored DELETE requests not stored into KV storage",
	})
	conAPI = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert" + serviceName + "api",
		Help: "Total consul api clietn request",
	})
	conAPIError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert" + serviceName + "api_error",
		Help: "Total consul api faild clietn request",
	})
)

// PutToKV insret key to storage
func PutToKV(Key, Value string) {
	kv := consul()
	// PUT a new KV pair
	p := &api.KVPair{Key: Key, Value: []byte(Value)}
	_, err := kv.Put(p, nil)
	if err != nil {
		conPutErr.Inc()
		log.Println(err)
	} else {
		conPut.Inc()
	}
}

// DelFromKV insret key to storage
func DelFromKV(key string) {
	kv := consul()
	// Delete a KV pair
	_, err := kv.Delete(key, nil)
	if err != nil {
		conDelErr.Inc()
		log.Println(err)
	} else {
		conDel.Inc()
	}
}

// GetFromKV insret key to storage
func GetFromKV(key string) string {
	kv := consul()
	// Lookup the pair
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		conGetErr.Inc()
		log.Println(err)
		return ""
	}

	if pair == nil {
		conWhiteGet.Inc()
		return ""
	}

	conGet.Inc()
	return string(pair.Value)

}

func consul() *api.KV {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		conAPIError.Inc()
		log.Println(err)
	}
	kv := client.KV()
	conAPI.Inc()

	return kv
}
