package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	slack = flag.String("slack", "", "Slack channel to send notifications")
)

type eventT struct {
	Message        string `json:"message"`
	Reason         string `json:"reason"`
	Type           string `json:"type"`
	InvolvedObject struct {
		Kind string `json:"kind"`
		Name string `json:"name"`
		Uid  string `json:"uid"`
	} `json:"involvedObject"`
	Metadata struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Source struct {
		Component string `json:"component"`
		Host      string `json:"host"`
	} `json:"source"`
}

func handleEvent(e eventT) {
	if e.Reason != "OOMKilling" {
		return
	}

	slog.Info("OOMKilled event:", "v", e)

	if *slack == "" {
		return
	}

	var sm strings.Builder
	fmt.Fprintf(&sm, "%v", e.Message)
	node := e.InvolvedObject.Name
	if node == "" {
		node = e.Source.Host
	}

	if node != "" {
		fmt.Fprintf(&sm, ", from [%v]", node)
	}

	p := map[string]interface{}{
		"attachments": []map[string]interface{}{
			{
				"color":  "danger",
				"title":  "oomkill-watch | notify",
				"text":   sm.String(),
				"footer": "oomkill-watch",
				"ts":     time.Now().Unix(),
			},
		},
	}

	bp, _ := json.Marshal(p)
	req, err := http.NewRequest("POST", *slack, bytes.NewReader(bp))
	if err != nil {
		slog.Error("NewRequest failed:", "err", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Do failed:", "err", err)
		return
	}

	resp.Body.Close()
}

func main() {
	flag.Parse()
	var tries int
	exit := make(chan struct{})
	done := make(chan struct{}, 1)
	var w sync.WaitGroup
	w.Add(1)
	go func() {
		defer w.Done()
		s := make(chan os.Signal)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
		sig := fmt.Sprintf("%s", <-s)
		slog.Info("Signal received, cleaning up:", "sig", sig)
		exit <- struct{}{}
	}()

	for {
		tries++
		var ix int32
		args := []string{
			"get",
			"events",
			"-w",
			"--watch-only",
			"--field-selector",
			"type=Warning",
			"-o",
			"json",
		}

		cmd := exec.Command("kubectl", args...)
		outpipe, err := cmd.StdoutPipe()
		if err != nil {
			slog.Error("StdoutPipe failed:", "err", err)
			return
		}

		errpipe, err := cmd.StderrPipe()
		if err != nil {
			slog.Error("StdoutPipe failed:", "err", err)
			return
		}

		slog.Info("Starting kubectl:", "tries", tries)
		err = cmd.Start()
		if err != nil {
			slog.Error("Start failed:", "err", err)
			return
		}

		ch0 := make(chan struct{}, 1)
		ch1 := make(chan struct{}, 1)
		ch2 := make(chan struct{}, 1)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-exit:
					slog.Info("Cleaning up child processes...")
					err := cmd.Process.Signal(syscall.SIGTERM)
					if err != nil {
						slog.Info("Failed to terminate process, force kill...")
						cmd.Process.Signal(syscall.SIGKILL)
					}

					atomic.StoreInt32(&ix, 1)
					done <- struct{}{}
					return
				case <-ch0:
					return
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			outscan := bufio.NewScanner(outpipe)
			var j strings.Builder
		loop1:
			for {
				select {
				case <-ch1:
					break loop1
				default:
					chk := outscan.Scan()
					if !chk {
						break
					}

					stxt := outscan.Text()
					var e eventT
					fmt.Fprintf(&j, "%v", stxt)
					err := json.Unmarshal([]byte(j.String()), &e)
					if err == nil {
						slog.Info("[event]", "v", e)
						handleEvent(e)
						j.Reset()
					}
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			errscan := bufio.NewScanner(errpipe)
		loop2:
			for {
				select {
				case <-ch2:
					break loop2
				default:
					chk := errscan.Scan()
					if !chk {
						break
					}

					stxt := errscan.Text()
					slog.Info("[stderr]", "v", stxt)
				}
			}
		}()

		err = cmd.Wait()
		slog.Info("[exit]", "err", err)
		ch0 <- struct{}{} // exit term
		ch1 <- struct{}{} // exit stdout
		ch2 <- struct{}{} // exit stderr
		wg.Wait()

		// See if this is final exit:
		if atomic.LoadInt32(&ix) > 0 {
			break
		}
	}

	w.Wait()
	<-done
}
