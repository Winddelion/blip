package checker

import (
	"fmt"
	"net/http"
	"time"
)

type StatusCode struct {
	UsrURL       string
	Resp         int16
	ResponseTime time.Duration
}

func GetStatusCode(url string) (StatusCode, error) {
	if len(url) == 0 {
		return StatusCode{}, fmt.Errorf("No urls provided")
	}

	start := time.Now()
	resp, err := http.Get(url)
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
