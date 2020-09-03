package logging

import (
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	serviceName = "_logging_"
	logErrTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hacert" + serviceName + "error",
		Help: "Total count of ERROR logs",
	})
)

// Logging structure contain basic setup parameters
type Logging struct {
	path    string
	level   string
	file    *os.File
	appName string
}

// NewLogging init logging
func NewLogging(p, logLevel, appName string) *Logging {
	l := &Logging{
		path:    p,
		file:    nil,
		level:   logLevel,
		appName: appName,
	}
	createLogFile(l)
	return l
}

func createLogFile(l *Logging) error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		log.Print("error opening file: ", err)
		return err
	}
	l.file = f

	log.SetOutput(l.file)
	log.Println("This is a haCert log file")
	return nil
}

// LoggWrite virte log to file
func (l *Logging) LoggWrite(p, s string, err error) {
	logAllow := 4
	switch l.level {
	case "DEBUG":
		logAllow = 1
	case "WARN":
		logAllow = 2
	case "INFO":
		logAllow = 3
	case "ERROR":
		logAllow = 4
	}

	switch p {
	case "DEBUG":
		if logAllow == 1 {
			if err == nil {
				log.Println("[" + l.appName + "-DEBUG] " + s)
			} else {
				log.Println("["+l.appName+"-DEBUG] "+s, err)
			}
		}
	case "WARN":
		if logAllow <= 2 {
			if err == nil {
				log.Println("[" + l.appName + "-WARN] " + s)
			} else {
				log.Println("["+l.appName+"-WARN] "+s, err)
			}
		}
	case "INFO":
		if logAllow <= 3 {
			if err == nil {
				log.Println("[" + l.appName + "-INFO] " + s)
			} else {
				log.Println("["+l.appName+"-INFO] "+s, err)
			}
		}
	case "ERROR":
		if logAllow <= 4 {
			if err == nil {
				logErrTotal.Inc()
				log.Println("[" + l.appName + "-ERROR] " + s)
			} else {
				logErrTotal.Inc()
				log.Println("["+l.appName+"-ERROR] "+s, err)
			}
		}
	}
}
