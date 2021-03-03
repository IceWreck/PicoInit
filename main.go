package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strings"
	"sync"
	"syscall"

	"github.com/IceWreck/PicoInit/service"
	"github.com/postfinance/single"
)

func main() {
	// handle sigterm
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("User sent sigterm, quitting services and exitting.")
		os.Exit(0)
	}()

	// waitgroup to keep track of running goroutines
	var wg sync.WaitGroup

	// lockfile mech to ensure that only one instance of PicoInit is running at a time

	// create a new lockfile
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	one, err := single.New(".picoinit", single.WithLockPath(user.HomeDir))
	if err != nil {
		log.Fatal(err)
	}
	// lock and defer unlocking
	if err := one.Lock(); err != nil {
		if err == single.ErrAlreadyRunning {
			log.Fatal("An instance of PicoInit is already running. Please kill it and then start PicoInit.")
		} else {
			log.Fatal(err)
		}
	}

	// run the main application

	log.Println("Starting PicoInit ...")

	// start all services
	for _, item := range service.Config {
		sv := item

		// add 1 to waitgroup
		wg.Add(1)
		go func() {
			// tell wg that we're done when everything ends
			defer wg.Done()

			// everything is in a for loop because we need to restart services if they end
			// service is actually stopped and for loop is broken according to user's
			// restart policy mentioned in picoinit_config.json
			for {
				log.Println("Starting", sv.Name)

				// open the log file
				logFile, err := os.OpenFile(fmt.Sprintf("./logs/%s.log", sv.Name), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					log.Println("please create a folder called 'logs' in the same dir as the PicoInit executable.")
					log.Panic("error opening log file:", err)
				}
				defer logFile.Close()

				// split executable file and additional arguments
				baseExec := strings.Split(sv.Command, " ")[0]
				args := strings.Split(sv.Command, " ")[1:]
				log.Println("EXEC", baseExec, "ARGS", args)
				cmd := exec.Command(baseExec, args...)
				if sv.WorkingDir != "" {
					cmd.Dir = sv.WorkingDir
				}
				// set stdout and stderr to logfile
				cmd.Stdout = logFile
				cmd.Stderr = logFile

				// execute the command
				err = cmd.Run()
				if err != nil {
					log.Println("Command", sv.Name, "finished with error", err)
				} else {
					// did not end on error so don't restart if policy says so
					if sv.Restart == "on_error" {
						log.Println("Command", sv.Name, "finished. As per policy, it will not be restarted.")
						break
					}
				}
				// don't restart if policy says never
				if sv.Restart == "never" {
					log.Println("Command", sv.Name, "finished. As per policy, it will not be restarted.")
					break
				}

				// if its at this point, then it should restart
				log.Println("Command", sv.Name, "finished. Restarting...")
			}

			// recover from panic if any
			recover()
		}()
	}

	// keep everything running until killed
	wg.Wait()
	log.Println("All services ended as per policy. Quitting PicoInit.")
}
