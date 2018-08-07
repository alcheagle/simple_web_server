package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net/http"
	"os"

	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	filemutex "github.com/alexflint/go-filemutex"
	"github.com/urfave/cli"
)

var (
	counter         = uint32(0)
	counterFilename string
	listenAddress   string
	// mutex_filename   = counterFilename + ".lock"

	mutex *filemutex.FileMutex
)

func init() {
	log.SetFormatter(&prefixed.TextFormatter{})
	// log.SetFormatter(&log.JSONFormatter{})
	filenameHook := filename.NewHook()
	filenameHook.Field = "filename" // Customize source field name
	log.AddHook(filenameHook)
}

func handler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()

	f, err := os.OpenFile(counterFilename, os.O_RDWR|os.O_CREATE, 0660)
	rwter := bufio.NewReadWriter(
		bufio.NewReader(f),
		bufio.NewWriter(f),
	)

	if err != nil {
		log.Error(err)
	}

	bs := make([]byte, 4)

	// rwter.Read(bs)
	f.Read(bs)
	global_counter := binary.LittleEndian.Uint32(bs)

	defer f.Close()

	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	counter++
	global_counter++

	binary.LittleEndian.PutUint32(bs, global_counter)
	// _, err = rwter.Write(bs)
	f.WriteAt(bs, 0)
	rwter.Flush()
	mutex.Unlock()
	if err != nil {
		log.Error(err)
	}

	log.Infof("request count: %d", counter)
	log.Infof("global request count: %d", global_counter)
	fmt.Fprintf(w, "Hi there, I am %s\n", name)
	fmt.Fprintf(w, "request count: %d\n", counter)
	fmt.Fprintf(w, "global request count: %d\n", global_counter)

	// time.Sleep(5 * time.Second)
}

func main() {
	app := cli.NewApp()
	app.Name = "raptor to influx"
	app.Usage = "InfluxDB integrating script for Raptorbox"
	app.Flags = defaultFlags
	app.Before = func(c *cli.Context) error {
		var err error
		mutex, err = filemutex.New(counterFilename + ".lock")

		if err != nil {
			log.Fatal("Directory did not exist or file could not created")
		}
		return nil
	}

	app.Action = func(c *cli.Context) error {
		http.HandleFunc("/", handler)
		log.Println("Starting server...")
		log.Fatal(http.ListenAndServe(listenAddress, nil))

		return nil
	}

	app.Run(os.Args)
}

var defaultFlags = []cli.Flag{
	cli.StringFlag{
		Name:        "listen-address, la",
		Value:       "0.0.0.0:8080",
		Usage:       "Listen address for the web application",
		EnvVar:      "LISTEN_ADDRESS",
		Destination: &listenAddress,
	},
	cli.StringFlag{
		Name:        "counter-file, cf",
		Value:       "counter.bin",
		Usage:       "file where to store the counters",
		EnvVar:      "COUNTER_FILE",
		Destination: &counterFilename,
	},
}
