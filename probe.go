package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/alexec/kit/internal/types"
)

func probeLoop(ctx context.Context, stopEverything func(), name string, probe types.Probe, callback func(name string, ok bool, err error)) {
	defer handleCrash(stopEverything)
	initialDelay := probe.GetInitialDelay()
	period := probe.GetPeriod()
	time.Sleep(initialDelay)
	successes, failures := 0, 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var err error
			if tcp := probe.TCPSocket; tcp != nil {
				_, err = net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", tcp.Port.String()))
			} else if httpGet := probe.HTTPGet; httpGet != nil {
				err = func() error {
					resp, err := http.Get(httpGet.GetURL())
					if err != nil {
						return err
					}
					if resp.StatusCode >= 300 {
						return fmt.Errorf("%s", resp.Status)
					}
					return nil
				}()
			} else {
				panic(fmt.Errorf("probe not supported"))
			}

			if err == nil {
				failures = 0
				successes++
			} else {
				successes = 0
				failures++
			}

			if successes == probe.GetSuccessThreshold() {
				callback(name, true, nil)
			} else if failures == probe.GetFailureThreshold() {
				callback(name, false, err)
			}
			time.Sleep(period)
		}
	}
}
