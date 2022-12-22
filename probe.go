package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
)

func probeLoop(ctx context.Context, name string, probe corev1.Probe, callback func(name string, ok bool, err error)) {
	defer runtime.HandleCrash()
	initialDelay := time.Duration(probe.InitialDelaySeconds) * time.Second
	period := time.Duration(probe.PeriodSeconds) * time.Second
	if period == 0 {
		period = 10 * time.Second
	}
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
				proto := strings.ToLower(string(httpGet.Scheme))
				if proto == "" {
					proto = "http"
				}
				resp, err := http.Get(fmt.Sprintf("%s://localhost:%v%s", proto, httpGet.Port.IntValue(), httpGet.Path))
				ok := err == nil && resp.StatusCode < 300
				if ok {
					successes++
					failures = 0
				} else {
					successes = 0
					failures++
				}
				successThreshold := int(probe.SuccessThreshold)
				if successThreshold == 0 {
					successThreshold = 1
				}
				failureThreshold := int(probe.FailureThreshold)
				if failureThreshold == 0 {
					failureThreshold = 1
				}
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
