package mpsqsd

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

// SqsdPlugin plugin for JSON
type SqsdPlugin struct {
	Tempfile           string
	URL                string
	Prefix             string
	InsecureSkipVerify bool
}

type sqsdStats struct {
	IsWorking      bool    `json:"is_working"`
	TotalHandled   float64 `json:"total_handled"`
	TotalSucceeded float64 `json:"total_succeeded"`
	TotalFailed    float64 `json:"total_failed"`
	MaxWorker      float64 `json:"max_worker"`
	BusyWorker     float64 `json:"busy_worker"`
	IdleWorker     float64 `json:"idle_worker"`
}

func (ss sqsdStats) ToMackerelMetrics() map[string]float64 {
	isWorking := float64(0)
	if ss.IsWorking {
		isWorking = 1
	}

	return map[string]float64{
		"is_working": isWorking,
		"handled":    ss.TotalHandled,
		"succeeded":  ss.TotalSucceeded,
		"failed":     ss.TotalFailed,
		"max":        ss.MaxWorker,
		"busy":       ss.BusyWorker,
		"idle":       ss.IdleWorker,
	}
}

func (p *SqsdPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(p.MetricKeyPrefix())
	return map[string]mp.Graphs{
		"is_working": {
			Label: labelPrefix + " IsWorking",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "is_working", Label: "IsWorking"},
			},
		},
		"workers": {
			Label: labelPrefix + " Workers",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "max", Label: "Max"},
				{Name: "busy", Label: "Busy"},
				{Name: "idle", Label: "Idle"},
			},
		},
		"jobs": {
			Label: labelPrefix + " Jobs",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "handled", Label: "Handled", Diff: true},
				{Name: "succeeded", Label: "Succeeded", Diff: true},
				{Name: "failed", Label: "Failed", Diff: true},
			},
		},
	}
}

// FetchMetrics interface for mackerel-plugin
func (p *SqsdPlugin) FetchMetrics() (map[string]float64, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: p.InsecureSkipVerify},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Get(p.URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var stats sqsdStats
	if err := json.NewDecoder(response.Body).Decode(&stats); err != nil {
		return nil, err
	}

	return stats.ToMackerelMetrics(), nil
}

func (p *SqsdPlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "sqsd"
	}
	return p.Prefix
}

// Do do doo
func Do() {
	url := flag.String("url", "", "URL to get a JSON")
	prefix := flag.String("prefix", "", "Prefix for metric names")
	insecure := flag.Bool("insecure", false, "Skip certificate verifications")
	tempfile := flag.String("tempfile", ``, "Temp file name")
	flag.Parse()

	if *url == "" {
		fmt.Println("-url is mandatory")
		os.Exit(1)
	}

	plugin := &SqsdPlugin{
		URL:                *url,
		Prefix:             *prefix,
		InsecureSkipVerify: *insecure,
	}

	helper := mp.NewMackerelPlugin(plugin)
	helper.Tempfile = *tempfile
	helper.Run()
}
