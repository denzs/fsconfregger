// Copyright 2020 Sebastian Denz
// GNU General Public License v3.0+
// see COPYING or https://www.gnu.org/licenses/gpl-3.0.txt)

package main

import (
	"flag"
	"fmt"
	"github.com/fiorix/go-eventsocket/eventsocket"
	"log"
	"os/exec"
	"strconv"
)

var ESHost string
var ESPort int
var ESPW string
var ScriptPath string

func init() {
	flag.StringVar(&ESHost, "eshost", "localhost", "FreeSWITCH Event Socket Host (default 'localhost')")
	flag.IntVar(&ESPort, "esport", 8021, "FreeSWITCH Event Socket Port (default 8021)")
	flag.StringVar(&ESPW, "espw", "ClueCon", "FreeSWITCH Event Socket Password (default 'ClueCon')")
	flag.StringVar(&ScriptPath, "script", "./sofia-generator.sh", "Path to XML Generator (default ./sofia-generator.sh)")
}

func create_confreg(cConf chan string, es *eventsocket.Connection) {
	for {
		select {
		case conf := <-cConf:
			log.Printf("%s: executing '%s %s'", conf, ScriptPath, conf)
			cmd := exec.Command(ScriptPath, conf)
			err := cmd.Run()
			if err != nil {
				log.Printf("%s finished with error: %v", conf, ScriptPath, err)
			}
			log.Printf("%s: executing 'sofia profile external rescan'\n", conf)
			es.Send("api sofia profile external rescan")
		}
	}
}

func destroy_confreg(dConf chan string, es *eventsocket.Connection) {
	for {
		select {
		case conf := <-dConf:
			log.Printf("%s: executing '%s %s del'", conf, ScriptPath, conf)
			cmd := exec.Command(ScriptPath, conf, "del")
			err := cmd.Run()
			if err != nil {
				log.Printf("%s '%s %s del' finished with error: %v", conf, ScriptPath, conf, err)
			}
			log.Printf("%s: executing 'api sofia profile external killgw conf_%s'\n", conf, conf)
			es.Send(fmt.Sprintf("api sofia profile external killgw conf_%s", conf))
		}
	}
}

func main() {
	flag.Parse()

	log.Printf("starting up..")
	log.Printf("connecting to: %s", ESHost+":"+strconv.Itoa(ESPort))
	eventsocket, err := eventsocket.Dial(ESHost+":"+strconv.Itoa(ESPort), ESPW)
	if err != nil {
		log.Fatal(err)
	}

	cChan := make(chan string)
	dChan := make(chan string)

	go create_confreg(cChan, eventsocket)
	go destroy_confreg(dChan, eventsocket)

	eventsocket.Send("events json ALL")

	for {
		event, err := eventsocket.ReadEvent()
		if err != nil { // might occur when gateway is deleted and did not exist before
			log.Printf("Error on Eventsocket: %s", err)
		} else {
			// event.PrettyPrint()
			if event.Get("Action") == "conference-create" {
				cChan <- event.Get("Conference-Name")
			}
			if event.Get("Action") == "conference-destroy" {
				dChan <- event.Get("Conference-Name")
			}
		}
	}
	eventsocket.Close()
	log.Printf("shutting down..")
}
