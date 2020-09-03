package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/jakubjastrabik/treafikCertKVstore/consul"
	"github.com/jakubjastrabik/treafikCertKVstore/logging"

	"github.com/fsnotify/fsnotify"
)

var (
	// Logg global log file
	Logg                  *logging.Logging
	members               = flag.String("members", "", "comma seperated list of members")
	httpPort              = flag.String("httpPort", "7900", "Port to be use for connection")
	httpAddress           = flag.String("httpAddress", "0.0.0.0", "Address to be use for connection")
	traefikCertLocalStore = flag.String("localStore", "/etc/traefik/acme.json", "path with file name where are stored certificates")
	consulKey             = flag.String("consulKey", "traefik/acme.json", "Consul key for storage certificates")
	path                  = flag.String("logFilePath", "/var/log/hacert.log", "Logi file path with name")
	logLevel              = flag.String("logLevel", "DEBUG", "DEBUG, WARN, INFO, ERROR")
	appName               = flag.String("appName", "traefikCertKVStore", "Aplication tag in log")

	watchError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert_watcher_error",
		Help: "Total count of ERROR start watcher",
	})
	fileReadError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert_file_read_error",
		Help: "Total count of ERROR reads cert file",
	})
	httpError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert_http_update_error",
		Help: "Total count of faild update notify cluster",
	})
	fileChange = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert_file_change",
		Help: "Total count of file changes",
	})
	watcherFileError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert_file_watcher_error",
		Help: "Total count of ERROR watching file",
	})
	fileWrite = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert_file_write",
		Help: "Total count of writing string to cert file",
	})
	fileWriteError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert_file_write_error",
		Help: "Total count of ERROR writing string to cert file",
	})
	traefikReload = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert_traefik_reload",
		Help: "Total count Traefik reload",
	})
	traefikReloadError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert_traefik_reload_error",
		Help: "Total count Traefik error realod",
	})
	certUpdate = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert_http_request_update",
		Help: "Total count of request to update cert",
	})
	httpLisError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert_http_listen_error",
		Help: "Total count of http listen ERROR",
	})
	httpUpdate = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert_http_accepted_update",
		Help: "Total count of accepted update from cluster",
	})
)

func init() {
	flag.Parse()
}

func checkFileChange() {

	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		Logg.LoggWrite("ERROR", "Error create new wathcer", err)
		watchError.Inc()
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
						fileReadError.Inc()
						Logg.LoggWrite("ERROR", "Error read file ", err)
					}
					consul.PutToKV(*consulKey, string(content))

					s := strings.Split(*members, ",")
					for i := range s {
						resp, err := http.Get("http://" + s[i] + ":" + *httpPort + "/update")
						httpUpdate.Inc()
						if err != nil {
							httpError.Inc()
							Logg.LoggWrite("ERROR", "Error handle http request ", err)
						}
						defer resp.Body.Close()
					}
					if err != nil {
						fileChange.Inc()
						Logg.LoggWrite("INFO", "File was change ", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				if err != nil {
					watchError.Inc()
					Logg.LoggWrite("ERROR", "Unable start watcher", err)
				}
				done <- false
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add(*traefikCertLocalStore); err != nil {
		watcherFileError.Inc()
		Logg.LoggWrite("ERROR", "Unable add watching file", err)
	}
	<-done
}

func saveData() {
	d1 := []byte(consul.GetFromKV(*consulKey))

	err := ioutil.WriteFile(*traefikCertLocalStore, d1, 0644)
	if err != nil {
		fileWriteError.Inc()
		Logg.LoggWrite("ERROR", "Unaible write data to file", err)
	} else {
		fileWrite.Inc()
		Logg.LoggWrite("DEBUG", "Write data to acme.json", err)
	}
	cmd := "systemctl reload traefik.service"
	_, err = exec.Command("bash", "-c", cmd).CombinedOutput()

	if err != nil {
		traefikReloadError.Inc()
		Logg.LoggWrite("ERROR", "Systemd faild reload traefik service", err)
	} else {
		traefikReload.Inc()
		Logg.LoggWrite("DEBUG", "Systemd reload traefik service", err)
	}
	return
}

func httpServer() {
	l := "Start web servers on address " + *httpAddress + ":" + *httpPort
	Logg.LoggWrite("INFO", l, nil)
	http.HandleFunc("/update", handleUpdate)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":"+*httpPort, nil)
	if err != nil {
		httpLisError.Inc()
		Logg.LoggWrite("ERROR", "Unable start listening", err)
	} else {
		Logg.LoggWrite("DEBUG", "Web server listening", err)
	}

}

func main() {
	// Get data from consul after start app
	Logg = logging.NewLogging(*path, *logLevel, *appName)
	Logg.LoggWrite("DEBUG", "Wait after start", nil)
	time.Sleep(1 * time.Second)

	saveData()
	go httpServer()
	for {
		checkFileChange()
	}
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	certUpdate.Inc()
	saveData()
}
