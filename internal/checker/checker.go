package checker

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type StatusCode struct {
	UsrURL       string
	Resp         int16
	ResponseTime time.Duration
}

func GetStatusCode(url string) (StatusCode, error) {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, "tcp4", addr)
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	if len(url) == 0 {
		return StatusCode{}, fmt.Errorf("No urls provided")
	}

	start := time.Now()
	resp, err := client.Get(url)
	elapsed := time.Since(start)

	if err != nil {
		return StatusCode{}, fmt.Errorf("%v\n", err.Error())
	}
	defer resp.Body.Close()

	return StatusCode{
		UsrURL:       url,
		Resp:         int16(resp.StatusCode),
		ResponseTime: elapsed,
	}, nil
}
