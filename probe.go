package main

import (
	"context"
	"fmt"
	"github.com/alexec/kit/internal/types"
	"net"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/runtime"
)

func probeLoop(ctx context.Context, name string, probe types.Probe, callback func(name string, ok bool, err error)) {
	defer runtime.HandleCrash()
	initialDelay := probe.GetInitialDelay()
	period := probe.GetPeriod()
	time.Sleep(initialDelay)
	successes, failures := 0, 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if tcp := probe.TCPSocket; tcp != nil {
				_, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", tcp.Port.IntVal))
				callback(name, err == nil, err)
			} else if httpGet := probe.HTTPGet; httpGet != nil {
				resp, err := http.Get(httpGet.GetURL())
				ok := err == nil && resp.StatusCode < 300
				if ok {
					successes++
					failures = 0
				} else {
					successes = 0
					failures++
				}
				successThreshold := probe.GetSuccessThreshold()
				failureThreshold := probe.GetFailureThreshold()
				if successes == successThreshold {
					callback(name, ok, nil)
					successes = 0
				} else if failures == failureThreshold {
					if err != nil {
						callback(name, ok, err)
					} else {
						callback(name, ok, fmt.Errorf("%s", resp.Status))
					}
					failures = 0
				}
			} else {
				callback(name, false, fmt.Errorf("probe not supported"))
			}
			time.Sleep(period)
		}
	}
}
