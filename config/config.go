// Copyright 2020 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

// Metric contains values that define a metric
type Metric struct {
	Name   string
	Path   string
	Labels map[string]string
	Type   MetricType
	Help   string
	Values map[string]string
}

type MetricType string

const (
	ValueScrape  MetricType = "value" // default
	ObjectScrape MetricType = "object"
)

// Config contains metrics and headers defining a configuration
type Config struct {
	Headers map[string]string
	Metrics []Metric
	Global  GlobalConfig
}

type GlobalConfig struct {
	TimeoutSeconds float64
}

func (metric *Metric) LabelNames() []string {
	var labelNames []string
	for name := range metric.Labels {
		labelNames = append(labelNames, name)
	}
	return labelNames
}

func LoadConfig(configPath string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Errorf("Failed to load config: %s, Error: %s", configPath, err)
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Errorf("Failed to parse YAML: %s", err)
		return config, err
	}

	// Complete Defaults
	for i := 0; i < len(config.Metrics); i++ {
		if config.Metrics[i].Type == "" {
			config.Metrics[i].Type = ValueScrape
			//config.Metrics[i].Type = "DefaultScrapeType"
		}
		if config.Metrics[i].Help == "" {
			config.Metrics[i].Help = config.Metrics[i].Name
		}
	}
	if config.Global.TimeoutSeconds == 0 {
		config.Global.TimeoutSeconds = 10
	}

	return config, nil
}
