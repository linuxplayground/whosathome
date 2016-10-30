package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "regexp"
    "strconv"
    "time"
)

const (
    ARPREGEX="[a-f0-9]{2}:[a-f0-9]{2}:[a-f0-9]{2}:[a-f0-9]{2}:[a-f0-9]{2}:[a-f0-9]{2}"
    INTERVAL="30s"
    DECAY=5
)

type Person struct {
    name    string
    online  bool
    stamp   time.Time
}

var people = map[string]Person{
	"d0:e1:40:30:e0:7b": {name: "Ben", online: false, stamp: time.Now()},
	"48:db:50:97:4e:62": {name: "David", online: false, stamp: time.Now()},
	"e8:50:8b:39:aa:49": { name: "Caroline", online: false, stamp: time.Now()},
	"1c:cb:99:c6:93:48": { name: "Luca", online: false, stamp: time.Now()},
}

var arplog = map[string]Person{}

func die(msg string, code int) {
	log.Fatalln(msg)
	os.Exit(code)
}

func getTable() string {
    data, err := exec.Command("arp-scan", "-l").Output()
    if err != nil {
        log.Println("Failed to run arp-scan")
        return ""
    }
    return string(data)
}

func out(m string, entry Person) {
	fmt.Printf("{\"timestamp\":\"%s\", \"mac\":\"%s\", \"name\":\"%s\", \"online\":%s}\n", 
                   entry.stamp.String(), 
                   m, 
                   entry.name, 
                   strconv.FormatBool(entry.online)
        )
}

func main() {

    interval, err := time.ParseDuration(INTERVAL)
    if err != nil {
        die("Failed to parse time interval", 2)
    }
    var arpregexp = regexp.MustCompile(ARPREGEX)

    for {

        matches := arpregexp.FindAllStringSubmatch(getTable(), -1)
        //Make a map out of matches
        var matchmap = map[string]bool{}
        for i:=0; i<len(matches); i++ {
            matchmap[matches[i][0]] = true
        }

        for i:=0; i<len(matches); i++ {
            mac := matches[i][0]
            //Have we found a person?
            if person,present := people[mac]; present {
                //Is the person not in the arplog?
                if _,logged := arplog[mac]; !logged {
                        person.online = true
                        person.stamp = time.Now()
                        arplog[mac] = person
                        fmt.Printf("Added %s to arplog.\n", person.name)
                        out(mac, person)
                }
				//Is the person already in the arplog
				if _,logged := arplog[mac]; logged {
					person.stamp = time.Now()
					person.online = true
					if !arplog[mac].online {
						fmt.Printf("Setting %s to online\n", person.name)
						out(mac, person)
					}
					arplog[mac] = person
				}

            }
        }

        for m := range people {
            if _, present := matchmap[m]; !present {
                if entry, logged := arplog[m]; logged {
                    if entry.online {
                        duration := time.Since(entry.stamp)
                        if duration.Minutes() > float64(DECAY) {
                                fmt.Printf("Marking %s as offline due to no arp response for more than %d minutes.\n", entry.name, DECAY)
                                entry.online = false
                                entry.stamp = time.Now()
                                arplog[m] = entry
                                out(m, entry)
                        }
                    }
                }
            }
        }

        time.Sleep(interval)
    }
}
