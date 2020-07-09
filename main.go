package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jakubjastrabik/Ha-cert-manager-for-traefik/consul"
)

var (
	members               = flag.String("members", "", "comma seperated list of members")
	httpPort              = flag.String("httpPort", "7900", "Port to be use for connection")
	httpAddress           = flag.String("httpAddress", "0.0.0.0", "Address to be use for connection")
	traefikCertLocalStore = flag.String("localStore", "/etc/traefik/acme.json", "path with file name where are stored certificates")
	consulKey             = flag.String("consulKey", "traefik/acme.json", "Consul key for storage certificates")
)

func init() {
	flag.Parse()
}

func checkFileChange() {

	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("ERROR", err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				done <- true
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					time.Sleep(time.Millisecond * 100)
					content, err := ioutil.ReadFile(*traefikCertLocalStore)
					if err != nil {
						log.Fatal(err)
					}
					consul.PutToKV(*consulKey, string(content))

					s := strings.Split(*members, ",")
					for i := range s {
						resp, err := http.Get("http://" + s[i] + ":" + *httpPort + "/update")
						if err != nil {
							log.Println(err)
						}
						defer resp.Body.Close()
					}

					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
				done <- false
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add(*traefikCertLocalStore); err != nil {
		log.Println("ERROR", err)
	}
	<-done
}

func saveData() {
	d1 := []byte(consul.GetFromKV(*consulKey))

	err := ioutil.WriteFile(*traefikCertLocalStore, d1, 0644)
	if err != nil {
		log.Println(err)
	}
	return
}

func httpServer() {
	log.Printf("Start web servers on address %s:%s", *httpAddress, *httpPort)
	http.HandleFunc("/update", handleUpdate)
	http.ListenAndServe(":"+*httpPort, nil)
}

func main() {
	// Get data from consul after start app
	saveData()
	go httpServer()
	for {
		checkFileChange()
	}

}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	saveData()
}
