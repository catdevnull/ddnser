package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type refreshChan chan string

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Must provide config file path as first argument")
	}
	config, err := LoadConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	refr := make(refreshChan)
	errch := make(chan error)
	go poller(config.Every, refr)
	go ipPoller(refr, errch)
	for {
		select {
		case reason := <-refr:
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
			defer cancel()
			log.Printf("Running because of %s", reason)
			var wg sync.WaitGroup
			wg.Add(len(config.Domains))
			for _, d := range config.Domains {
				go func(d Domain, wg *sync.WaitGroup) {
					record, err := d.NameServer.SetRecord(ctx, d.Name, config.Ip)
					if err != nil {
						log.Printf("ERROR: [%s] %s", d.Name, err.Error())
						// TODO: Reportar errores
						return
					}
					log.Printf("info: [%s] Set to %s", d.Name, record)
					wg.Done()

				}(d, &wg)
			}
		case err := <-errch:
			log.Fatal(err)
		}
	}
}

func poller(every int, refr refreshChan) {
	for {
		time.Sleep(time.Duration(every) * time.Minute)
		refr <- "poll"
	}
}

func ipPoller(refr refreshChan, errch chan error) {
	var lastWatched string
	for {
		cmd := exec.Command("ip", "address")
		out, err := cmd.Output()
		if err != nil {
			errch <- err
			return
		}
		var addrs string
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			prefix := "    inet "
			if strings.Index(line, prefix) == 0 {
				last := strings.Index(line[len(prefix):], "/")
				ip := line[len(prefix) : len(prefix)+last]
				addrs = addrs + ip + "/"
			}
		}
		if addrs != lastWatched {
			refr <- "ip changed"
			lastWatched = addrs
		}
		time.Sleep(time.Duration(2) * time.Second)
	}
}
