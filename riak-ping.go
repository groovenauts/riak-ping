package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"time"

	"github.com/tpjg/goriakpbc"
)

// Logger is used to write to syslog.
var Logger *log.Logger

func main() {
	// err is used to store the error.
	var err error
	// argv
	var flagAddr = flag.String("i", "127.0.0.1:8087", "Riak Server IP Address and Port")
	var flagProtocol = flag.String("p", "TCP", "Protocol (TCP/PB)")
	var flagVersion = flag.Bool("v", false, "Print Version Info")
	flag.Parse()

	if *flagVersion {
		fmt.Printf("riak-ping version: %s\n", Version)
		os.Exit(0)
	}

	// write setting to syslog
	Logger, err = SetLogger()
	if err != nil {
		panic(err)
	}

	// create riak client
	myRiak := SetRiakClient(*flagAddr)

	// check riak connection
	err = CheckConnect(myRiak, *flagProtocol)
	if err != nil {
		os.Exit(1)
	}

	return
}

// SetLogger define NewLogger
func SetLogger() (*log.Logger, error) {
	// Logger set (syslog & crit)
	l, err := syslog.NewLogger(syslog.LOG_CRIT|syslog.LOG_SYSLOG, 0)
	return l, err
}

// SetRiakClient define NewClient(riak)
func SetRiakClient(a string) *riak.Client {
	c := riak.NewClient(a)
	return c
}

// CheckConnect wait for the result of limits of the connection
func CheckConnect(c *riak.Client, t string) error {
	// use goroutine, define channel to use to store result
	resultCh := make(chan error, 1)
	// use goroutine(asynchronous), because the func set a timeout
	go ConnectRiak(c, t, resultCh)

	// exit conditions: func ConnectRiak is finish or 10sec elasped
	select {
	case err := <-resultCh: // func ConnectRiak is finish
		if err != nil {
			// NG: display the screen and write syslog
			fmt.Printf("%s Connect Failure.\n", t)
			if Logger != nil {
				Logger.Printf("%s Connect Failure.\n", t)
			}
		} else {
			// OK: display the screen only
			fmt.Printf("%s Connect Success.\n", t)
		}
		return err
	case <-time.After(time.Second * 10): // 10sec elasped
		// display the screen and write syslog
		fmt.Printf("%s Connect Timed Out.\n", t)
		if Logger != nil {
			Logger.Printf("%s Connect Timed Out.\n", t)
		}
		// store the error status
		err := errors.New("timed out")
		return err
	}
}

// ConnectRiak TCP(L4) or PB(L7) connection check to riak
func ConnectRiak(c *riak.Client, t string, rch chan error) {
	var err error
	defer c.Close()
	// exec riak connection check for each protocol
	switch {
	case t == "TCP":
		// TCP: set TCP(L4) connection
		err = c.Connect()
	case t == "PB":
		// PB: set PB(L7) connection
		err = c.Ping()
	default:
		// others: error
		err = errors.New("unknown protocol")
	}
	// store the error
	if err != nil {
		rch <- err
	} else {
		// if err=nil then, rch <- err is locked and wait the routine
		rch <- nil
	}
}
