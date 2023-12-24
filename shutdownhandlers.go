package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscalls that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		log.Println("shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Printf("cleaning up: %s", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Printf("%s: clean up failed: %s", innerKey, err.Error())
					return
				}

				log.Printf("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}

func revertBkp() {

	paths := getVTOLDir()

	if paths == nil {
		log.Panic(fmt.Errorf("no VTOL VR Directories found. Is game installed?"))
	}

	for _, x := range paths {
		files, err := os.ReadDir(x + "RadioMusic\\")
		if err != nil {
			log.Panic(err)
		}

		for _, y := range files {
			checkRemoval(x, y.Name())
		}
		for _, y := range files {
			checkRestore(x, y.Name())
		}

	}
}

func checkRestore(path string, name string) {
	if strings.Contains(name, ".bkp") {
		fmt.Println("Restoring file: " + name + ".")
		splitStr := strings.Split(name, ".bkp")
		if splitStr[0] == name {
			return
		}
		os.Rename(path+"RadioMusic\\"+name, path+"RadioMusic\\"+splitStr[0])
	}
}

func checkRemoval(path string, name string) {
	if validRemoval(name) {
		fmt.Println("Cleaning up blank mp3 files.")
		os.Remove(path + "RadioMusic\\" + name)
	}
}

func validRemoval(name string) bool {
	return name == "0.mp3" || name == "1.mp3" || name == "2.mp3"
}