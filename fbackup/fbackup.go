package fbackup

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

/*
	1. File Name
	2. Nuber of the file version
*/

type Backup struct {
	fileName     string
	versionCount int
}

// NewBackup create struct for backup rotate
func NewBackup(fileName string, versionCount int) *Backup {
	bs := &Backup{
		fileName:     fileName,
		versionCount: versionCount,
	}
	return bs
}

func (bs *Backup) BackupRotate() {
	for i := bs.versionCount; i > 0; i-- {
		/*
			--> Pokus o otvorenie suboru
			Ak sa subor n-1 otvori tak sa zrotuje na n
			Ak je (n-1) == 0, tak sme na koreni, to sa presunie na 1
			n-1 n
			fp  fd
			2   3
			1	2
			0	1
		*/

		// First need generated n-1 file name
		var fp, fd string

		if i-1 == 0 {
			fp = "/etc/traefik/acme.json"
			fd = "/etc/traefik/acme.json" + "-" + strconv.Itoa(i) + ".back"
		} else {
			fp = "/etc/traefik/acme.json" + "-" + strconv.Itoa(i-1) + ".back"
			fd = "/etc/traefik/acme.json" + "-" + strconv.Itoa(i) + ".back"
		}

		content, err := ioutil.ReadFile(fp)
		if err != nil {
			// File does not exist, skip this number
			continue
		}

		// fmt.Println("zo suboru:", fp)
		// fmt.Println("do suboru:", fd)

		// If file exist, rotating
		err = ioutil.WriteFile(fd, content, 0644)
		if err != nil {
			// fileWriteError.Inc()
			// Logg.LoggWrite("ERROR", "Error Write first Backup File ", err)
			fmt.Println("Error Write first Backup File")
		}
		fmt.Println(content)
	}
}

/*
	1. Open file
	2. Rotate OLD Backup
	3. Create New Latest Backup
*/

// fn := *traefikCertLocalStore + "-1.back"
// fn2 := *traefikCertLocalStore + "-2.back"

// fmt.Println(fn2)

// input, err := ioutil.ReadFile(fn)
// if err != nil {
// 	fileReadError.Inc()
// 	Logg.LoggWrite("ERROR", "Error read first Backup file", err)
// 	return
// }

// err = ioutil.WriteFile(fn2, input, 0644)
// if err != nil {
// 	fileWriteError.Inc()
// 	Logg.LoggWrite("ERROR", "Error Write second Backup File ", err)
// 	return
// }
