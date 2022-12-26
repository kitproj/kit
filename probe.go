package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/alexec/kit/internal/types"
)

func probeLoop(ctx context.Context, stop func(), probe types.Probe, callback func(ok bool, err error)) {
	defer handleCrash(stop)
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
				_, err = net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", tcp.Port.IntVal))
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
				callback(true, nil)
			} else if failures == probe.GetFailureThreshold() {
				callback(false, err)
			}
			time.Sleep(period)
		}
	}
}
