package main

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runprobe/prober"
	"sync"
)

var (
	srcRegions = map[string]string{
		"arn1": "aws eu-north-1",
		"bru1": "aws eu-central-1",
		"cdg1": "aws eu-west-3",
		"lhr1": "aws eu-west-2",
		"dub1": "aws eu-west-1",
		"sfo1": "aws us-west-1",
		"pdx1": "aws us-west-2",
		"cle1": "aws us-east-2",
		"iad1": "aws us-east-1",
		"gru1": "aws sa-east-1",
		"hnd1": "aws ap-northeast-1",
		"icn1": "aws ap-northeast-2",
		"sin1": "aws ap-southeast-1",
		"syd1": "aws ap-southeast-2",
		"bom1": "aws ap-south-1",
	}

	dstRegions = map[string]string{
		"us-central1":     "gcp us-central1",
		"us-east1":        "gcp us-east1",
		"asia-northeast1": "gcp asia-northeast1",
		"europe-west1":    "gcp europe-west1",
	}
)

func dstURL(region string) string {
	return "https://example-server-" + region + "-dpyb4duzqq-ew.a.run.app"
}

func srcURL(region, targetURL string) string {
	v := &url.Values{}
	v.Set("url", targetURL)
	return "https://prober-" + region + ".ahmetalpbalkan.now.sh/api/probe?" + v.Encode()
}

func main() {
	if err := warmup(); err != nil {
		log.Fatal(err)
	}

	src := []string{"arn1", "bru1", "cdg1"}
	_ = src

	table := tablewriter.NewWriter(os.Stdout)
	header := []string{""}
	for dstLbl := range dstRegions {
		header = append(header, dstLbl)
	}
	table.SetAutoFormatHeaders(false)
	table.SetAutoMergeCells(true)
	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetHeader(header)
	table.SetHeaderLine(false)

	for s, sL := range srcRegions {
		//if i > 2 {
		//	break
		//}
		row := []string{sL}
		for d, _ := range dstRegions {
			v, err := measure(s, d)
			if err != nil {
				log.Printf("[%s->%s] failed: %v", s, d, err)
				row = append(row, "probing_err")
				continue
			}

			if v.Error != nil {
				log.Printf("[%s->%s] measurement error: (kind=%v) %v", s, d, v.Error.Kind, v.Error.Message)
				row = append(row, "measure_err")
				continue
			}
			if !v.Success {
				row = append(row, "not_success")
				continue
			}

			kind := ""
			if v.Data.RequestState == "cold" {
				kind = "*"
			}
			row = append(row, fmt.Sprintf("%dms%s", v.Data.TotalMs, kind))
		}
		table.Append(row)
	}
	table.Render()
}

func warmup() error {
	var errOut error
	var wg sync.WaitGroup
	for dst := range dstRegions {
		wg.Add(1)
		go func(r string) {
			resp, err := http.Get(dstURL(r))
			if err != nil {
				errOut = fmt.Errorf("error warming up %s: %v", r, err)
			}
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				errOut = fmt.Errorf("error warming up %s, status code=http %d", r, resp.StatusCode)
			}
		}(dst)
	}
	wg.Done()
	return errOut
}

func measure(from, to string) (*prober.ProbeResp, error) {
	target := dstURL(to)
	resp, err := http.Get(srcURL(from, target))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var v prober.ProbeResp
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
