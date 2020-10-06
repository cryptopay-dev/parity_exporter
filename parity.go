package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type ParityLastBlock struct {
	url    string
	client *http.Client
	desc   *prometheus.Desc
}

type parityLastBlockRequest struct {
	Version string `json:"jsonrpc"`
	Method  string `json:"method"`
	ID      int    `json:"id"`
}

type parityLastBlockResponse struct {
	Result string                        `json:"result"`
	Error  *parityLastBlockErrorResponse `json:"error"`
}

type parityLastBlockErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (p *parityLastBlockErrorResponse) Error() string {
	if p.Message == "" {
		return fmt.Sprintf("json-rpc error %d", p.Code)
	}
	return p.Message
}

func NewParityLastBlock(url string) *ParityLastBlock {
	return &ParityLastBlock{url, http.DefaultClient, prometheus.NewDesc(
		"parity_last_block", "Parity last block", nil, nil)}
}

func (p *ParityLastBlock) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.desc
}

func (p *ParityLastBlock) Collect(ch chan<- prometheus.Metric) {
	lastBlockNumber, err := p.getLastBlockNumber()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(p.desc, err)
		return
	}

	ch <- prometheus.MustNewConstMetric(p.desc, prometheus.GaugeValue, float64(lastBlockNumber))
}

func (p *ParityLastBlock) getLastBlockNumber() (uint64, error) {
	var res *http.Response

	reqBody, err := json.Marshal(&parityLastBlockRequest{
		Version: "2.0",
		Method:  "eth_blockNumber",
		ID:      int(time.Now().UnixNano()),
	})
	if err != nil {
		return 0, err
	}

	buf := bytes.NewBuffer(reqBody)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	req, err := http.NewRequest(http.MethodPost, p.url, nil)
	if err != nil {
		return 0, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Body = ioutil.NopCloser(buf)
	req.ContentLength = int64(buf.Len())

	res, err = p.client.Do(req)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return 0, fmt.Errorf("incorrect response code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	resp := &parityLastBlockResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return 0, err
	}

	if resp.Error != nil {
		return 0, resp.Error
	}

	num, err := strconv.ParseUint(resp.Result, 0, 64)
	if err != nil {
		return 0, err
	}

	return num, nil
}
