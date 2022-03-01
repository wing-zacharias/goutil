package util

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type ConfigSet struct{}

type Prop struct {
	LineNumber int    `json:"line"`
	KeyStr     string `json:"key"`
	ValueStr   string `json:"value"`
}

func NewConfigSet() *ConfigSet {
	return &ConfigSet{}
}

func (s *ConfigSet) GetFilePerm(filePath string) os.FileMode {
	f, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	fm := f.Mode()
	return fm.Perm()
}

func (s *ConfigSet) GetProp(configFile string, key string, sep string) (*Prop, error) {
	prop := &Prop{}
	matchKey, err := regexp.Compile(fmt.Sprintf("^\\s*%s\\s*%s.*?", key, sep))
	if err != nil {
		return nil, err
	}
	f, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(f), "\n")
	cmt := regexp.MustCompile(`^#|\[.*?`)
	count := 0
	for index, line := range lines {
		if (!cmt.MatchString(line)) && matchKey.MatchString(line) {
			var res []string
			prop.LineNumber = index
			if sep == "" || regexp.MustCompile("\\s+").MatchString(sep) || regexp.MustCompile("\t").MatchString(sep) {
				res = strings.Fields(line)
				prop.KeyStr = res[0]
				prop.KeyStr = res[1]
			} else {
				res = strings.SplitN(line, sep, 2)
				prop.KeyStr = strings.TrimSpace(res[0])
				prop.ValueStr = strings.TrimSpace(res[1])
			}
			count++
		}
	}
	if count == 1 {
		return prop, nil
	}
	return nil, errors.Errorf("%s get prop %s error!", configFile, key)
}

func (s *ConfigSet) SetProp(configFile string, key string, value string, sep string) error {
	var result string
	f, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	lines := strings.Split(string(f), "\n")
	matchKey, err := regexp.Compile(fmt.Sprintf("^\\s*%s\\s*%s.*?", key, sep))
	if err != nil {
		return err
	}
	for index, line := range lines {
		if matchKey.MatchString(line) {
			result += fmt.Sprintf("%s%s%s", key, sep, value)
		} else {
			if index == (len(lines) - 1) {
				result += line
			} else {
				result += line + "\n"
			}
		}
	}
	if err := ioutil.WriteFile(configFile, []byte(result), s.GetFilePerm(configFile)); err != nil {
		return err
	}
	return nil
}

func (s *ConfigSet) AddProp(configFile string, key string, value string, sep string) error {
	var result string
	prop, _ := s.GetProp(configFile, key, sep)
	if prop == nil {
		f, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		lines := strings.Split(string(f), "\n")
		for _, line := range lines {
			result += line + "\n"
		}
		result += fmt.Sprintf("%s%s%s", key, sep, value)
		if err := ioutil.WriteFile(configFile, []byte(result), s.GetFilePerm(configFile)); err != nil {
			return err
		}
	}
	return nil
}

func (s *ConfigSet) DelProp(configFile string, key string, sep string) error {
	var result string
	lineRegex := fmt.Sprintf("^\\s*%s\\s*%s.*?", key, sep)
	re, err := regexp.Compile(lineRegex)
	if err != nil {
		return err
	}
	prop, _ := s.GetProp(configFile, key, sep)
	if prop != nil {
		f, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		lines := strings.Split(string(f), "\n")
		for index, line := range lines {
			if !re.MatchString(line) {
				if index == (len(lines) - 1) {
					result += line
				} else {
					result += line + "\n"
				}
			}
		}
		if err := ioutil.WriteFile(configFile, []byte(result), s.GetFilePerm(configFile)); err != nil {
			return err
		}
	}
	return nil
}

func (s *ConfigSet) GetAllProp(configFile string, sep string) ([]Prop, error) {
	var resProps []Prop
	f, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(f), "\n")
	cmt := regexp.MustCompile(`^#|\[.*?`)
	for index, line := range lines {
		if !cmt.MatchString(line) {
			var res []string
			if sep == "" || regexp.MustCompile("\\s+").MatchString(sep) || regexp.MustCompile("\t").MatchString(sep) {
				res = strings.Fields(line)
			} else {
				res = strings.Split(line, sep)
			}
			if len(res) == 2 {
				prop := Prop{
					LineNumber: index,
					KeyStr:     strings.TrimSpace(res[0]),
					ValueStr:   strings.TrimSpace(res[1]),
				}
				resProps = append(resProps, prop)
			}
		}
	}
	return resProps, nil
}
