package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

const nginxTemplate = `
geo $googlebot {
    default 0;
    {{- range .Prefixes }}
    {{- if .Ipv4Prefix }}
    {{ .Ipv4Prefix }} 1;
    {{- end }}
    {{- if .Ipv6Prefix }}
    {{ .Ipv6Prefix }} 1;
    {{- end }}
    {{- end }}
}
`

type BotConfig struct {
	File     string            `yaml:"file"`
	Template string            `yaml:"template"`
	Bots     map[string]BotURL `yaml:",inline"`
}

type BotURL struct {
	URL string `yaml:"url"`
}

type IPRange struct {
	CreationTime string `json:"creationTime"`
	Prefixes     []struct {
		Ipv4Prefix string `json:"ipv4Prefix,omitempty"`
		Ipv6Prefix string `json:"ipv6Prefix,omitempty"`
	} `json:"prefixes"`
}

func (ipr *IPRange) Merge(other *IPRange) {
	ipr.Prefixes = append(ipr.Prefixes, other.Prefixes...)
}

func parseYAMLConfig(filePath string) (*BotConfig, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config BotConfig
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func fetchJSON(url string) (*IPRange, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ipRange IPRange
	err = json.Unmarshal(body, &ipRange)
	if err != nil {
		return nil, err
	}

	return &ipRange, nil
}

func generateNginxConfig(ipRange *IPRange, config BotConfig) error {
	tmpl, err := template.New("nginx").Parse(config.Template)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ipRange); err != nil {
		return err
	}

	return os.WriteFile(config.File, buf.Bytes(), 0644)
}

func main() {
	configFilePath := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	// Load YAML configuration
	config, err := parseYAMLConfig(*configFilePath)
	if err != nil {
		log.Fatalf("Error parsing YAML config: %v", err)
	}

	combinedIPRange := &IPRange{}

	for botName, bot := range config.Bots {
		ipRange, err := fetchJSON(bot.URL)
		if err != nil {
			log.Fatalf("Error fetching JSON for %s: %v", botName, err)
		}
		combinedIPRange.Merge(ipRange)
	}

	err = generateNginxConfig(combinedIPRange, *config)
	if err != nil {
		log.Fatalf("Error generating Nginx config: %v", err)
	}
}
