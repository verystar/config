package config

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Config struct {
	ConfigMap map[string]map[string]string
}

func Load(path string, paths ...string) (*Config, error) {
	conf := new(Config)
	conf.ConfigMap = make(map[string]map[string]string)

	if err := conf.parseDataSource(path); err != nil {
		return nil, err
	}

	for _, v := range paths {
		if err := conf.parseDataSource(v); err != nil {
			return nil, err
		}
	}
	return conf, nil
}

func (c *Config) parseDataSource(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var section string
	r := bufio.NewReader(f)
	m := make(map[string]string)

	for {
		byts, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		s := strings.TrimSpace(string(byts))
		if s == "" {
			continue
		}
		if strings.HasPrefix(s, "#") || strings.HasPrefix(s, ";") {
			continue
		}

		if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
			i := strings.Index(s, "[")
			j := strings.LastIndex(s, "]")
			section = strings.TrimSpace(s[i+1 : j])
			m = make(map[string]string)
			continue
		}

		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		key := strings.TrimSpace(s[:index])
		if key == "" {
			continue
		}

		value := strings.TrimSpace(s[index+1:])
		if value == "" {
			continue
		}

		pos := strings.Index(value, "\t#")
		if pos > -1 {
			value = value[0:pos]
		}

		pos = strings.Index(value, " #")
		if pos > -1 {
			value = value[0:pos]
		}

		pos = strings.Index(value, " ;")
		if pos > -1 {
			value = value[0:pos]
		}

		pos = strings.Index(value, "\t//")
		if pos > -1 {
			value = value[0:pos]
		}

		pos = strings.Index(value, " //")
		if pos > -1 {
			value = value[0:pos]
		}

		value = strings.TrimSpace(value)

		if _, ok := c.ConfigMap[section]; ok {
			c.ConfigMap[section][key] = value
		} else {
			m[key] = value
			c.ConfigMap[section] = m
		}
	}
	return nil
}

func (c *Config) Read(section, key string) string {
	v, ok := c.ConfigMap[section][key]
	if !ok {
		return ""
	}
	return v
}
