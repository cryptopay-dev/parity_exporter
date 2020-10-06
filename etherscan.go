package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type EtherscanLastBlock struct {
	url, key string
	client   *http.Client
	desc     *prometheus.Desc
}

type etherscanLastBlockResponse struct {
	Result string `json:"result"`
}

func NewEtherscanLastBlock(url, key string) *EtherscanLastBlock {
	return &EtherscanLastBlock{url, key, http.DefaultClient, prometheus.NewDesc(
		"etherscan_last_block", "Etherscan last block", nil, nil)}
}

func (e *EtherscanLastBlock) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.desc
}

func (e *EtherscanLastBlock) Collect(ch chan<- prometheus.Metric) {
	lastBlockNumber, err := e.getLastBlockNumber()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(e.desc, err)
		return
	}

	ch <- prometheus.MustNewConstMetric(e.desc, prometheus.GaugeValue, float64(lastBlockNumber))
}

func (e *EtherscanLastBlock) getLastBlockNumber() (uint64, error) {
	resp, err := e.client.Get(fmt.Sprintf("%s/api?module=proxy&action=eth_blockNumber&apikey=%s", e.url, e.key))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	obj := new(etherscanLastBlockResponse)
	err = json.Unmarshal(b, obj)
	if err != nil {
		return 0, err
	}

	num, err := strconv.ParseUint(obj.Result, 0, 64)
	if err != nil {
		return 0, err
	}

	return num, nil
}
