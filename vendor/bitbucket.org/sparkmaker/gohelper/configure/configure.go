package configure

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var defaultMode = "default"

// Config keep the configuration in the buffer
type Config struct {
	data                   map[string]interface{}
	path                   string
	mode                   map[string]string
	currentMode            string
	strictModes            []string
	allowEmptyInStrictMode []string
}

func root() (dir string) {
	defer recover()

	rootDir, _ := os.Getwd()
	dir = rootDir

	gopath := os.Getenv("GOPATH")
	if len(gopath) == 0 {
		return
	}

	prefix := strings.Replace(dir, gopath, "", -1)
	spliter := "/"
	if strings.ContainsAny(prefix, "\\") {
		spliter = "\\"
	}
	splited := strings.Split(prefix, spliter)[0:3]
	projectDir := strings.Join(splited, spliter)
	dir = gopath + projectDir
	return
}

func hasPrefix(path string, prefix ...string) bool {
	for _, p := range prefix {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func contains(source string, in []string) bool {
	for _, v := range in {
		if source == v {
			return true
		}
	}
	return false
}

func floattostr(fv float64) string {
	return strconv.FormatFloat(fv, 'f', 0, 64)
}

func keyToEnv(key string) string {
	return strings.Replace(strings.ToUpper(key), ".", "_", -1)
}

func flatten(m map[string]interface{}) map[string]interface{} {
	o := make(map[string]interface{})
	for k, v := range m {
		switch child := v.(type) {
		case map[string]interface{}:
			nm := flatten(child)
			for nk, nv := range nm {
				o[k+"."+nk] = nv
			}
		default:
			o[k] = v
		}
	}
	return o
}

func unique(values []string) []string {
	temp := make(map[string]int)
	for _, value := range values {
		temp[value] = 0
	}
	result := []string{}
	for value := range temp {
		result = append(result, value)
	}
	return result
}

func compareMapKey(source map[string]interface{}, destination map[string]interface{}) bool {
	for key := range source {
		if _, ok := destination[key]; !ok {
			return false
		}
	}
	return true
}

// New configuration and initialize default configuration file path
// Adding default mode (environment variable `RUN_MODE` no set) to config path /conf/application.json
func New() *Config {
	c := &Config{mode: map[string]string{defaultMode: root() + "/conf/application.json"}}
	c.reload()
	return c
}

func (c *Config) reload() {
	c.filesVerify()
	c.assignPath()
	c.assignData()
	c.checkOverrideAllVarialbes()
}

func (c *Config) filesVerify() {
	paths := []string{}
	for _, path := range c.Mode() {
		paths = append(paths, filepath.FromSlash(path))
	}
	paths = unique(paths)
	if len(paths) < 2 {
		return
	}

	allModeData := []map[string]interface{}{}
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			panic(fmt.Errorf("invalid config file path '%v'", path))
		}
		m := make(map[string]interface{})
		readErr := json.NewDecoder(file).Decode(&m)
		file.Close()
		if readErr != nil {
			panic(readErr)
		}
		allModeData = append(allModeData, flatten(m))
	}

	source := allModeData[0]
	for index, destination := range allModeData {
		if index == 0 {
			continue
		}
		if !compareMapKey(source, destination) {
			panic(fmt.Errorf("config files have some keys not match"))
		}
	}
}

func (c *Config) checkOverrideAllVarialbes() {
	if !contains(c.currentMode, c.strictModes) {
		return
	}
	eks := []string{}
	for variable := range c.data {
		env := keyToEnv(variable)
		if len(strings.TrimSpace(os.Getenv(env))) == 0 && !contains(env, c.allowEmptyInStrictMode) {
			eks = append(eks, env)
		}
	}
	if len(eks) > 0 {
		panic(errors.New(fmt.Sprintf("these variables need to override [%v]", strings.Join(eks, ", "))))
	}
}

func (c *Config) assignPath() {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("RUN_MODE")))
	if len(mode) == 0 {
		mode = defaultMode
	}
	c.currentMode = mode
	fp, ok := c.mode[mode]
	if !ok {
		fp, ok := c.mode[defaultMode]
		if !ok {
			panic(fmt.Sprintf("run mode '%v' config file not found", mode))
		}
		c.path = fp
		return
	}
	c.path = fp
}

func (c *Config) assignData() {
	file, err := os.Open(c.path)
	if err == nil {
		c.assignDataFromFile(file)
		return
	}
	c.data = make(map[string]interface{})
}

func (c *Config) assignDataFromFile(f *os.File) {
	m := make(map[string]interface{})
	readErr := json.NewDecoder(f).Decode(&m)
	if readErr != nil {
		panic(readErr)
	}
	c.data = flatten(m)
}

// Add path of file in run mode
// read directory tear down from your root project
func (c *Config) Add(mode string, path string) {
	if !hasPrefix(path, "#", "~", "../", "./", "/") {
		path = fmt.Sprintf("%v/%v", root(), path)
	}
	c.mode[mode] = path
	for k, v := range c.mode {
		c.mode[k] = filepath.FromSlash(v)
	}
	c.reload()
}

// Del using for delete configuration path from run mode
func (c *Config) Del(mode string) {
	delete(c.mode, mode)
	c.reload()
}

// Mode for getting all of mode in the configuration
func (c *Config) Mode() map[string]string {
	return c.mode
}

// StrictMode will accept only override all variables, allow empty fields if set it via AllowEmptyInStrictMode
func (c *Config) StrictMode(modes []string, allowEmptyKeys []string) {
	strictModes := []string{}
	for _, v := range modes {
		strictModes = append(strictModes, strings.ToLower(strings.TrimSpace(v)))
	}
	c.strictModes = strictModes

	allowEmptyInStrictMode := []string{}
	for _, key := range allowEmptyKeys {
		allowEmptyInStrictMode = append(allowEmptyInStrictMode, keyToEnv(key))
	}
	c.allowEmptyInStrictMode = allowEmptyInStrictMode
	c.reload()
}

// GetString using for getting string type of a value in the configuration file or environment variable
func (c *Config) GetString(key string) string {
	if value := os.Getenv(keyToEnv(key)); len(value) != 0 {
		return value
	}
	if c.data[key] == nil {
		return ""
	}
	return c.data[key].(string)
}

// GetStrings using for getting list of string type of a value in the configuration file or environment variable
func (c *Config) GetStrings(key string) []string {
	var str string
	var slice []string
	if value := os.Getenv(keyToEnv(key)); len(value) != 0 {
		str = value
	} else if c.data[key] == nil {
		return slice
	} else {
		str = c.data[key].(string)
	}
	for _, v := range strings.Split(str, ",") {
		slice = append(slice, strings.TrimSpace(v))
	}
	return slice
}

// GetInt using for getting integer type of a value in the configuration file or environment variable
func (c *Config) GetInt(key string) int {
	value := os.Getenv(keyToEnv(key))
	if len(value) == 0 {
		value = floattostr(c.data[key].(float64))
	}
	i, _ := strconv.Atoi(value)
	return i
}

// GetBool using for getting boolean type of a value in the configuration file or environment variable
func (c *Config) GetBool(key string) bool {
	result := false
	if value, ok := c.data[key].(bool); ok {
		result = value
	}
	if value, ok := c.data[key].(string); ok {
		b, err := strconv.ParseBool(value)
		if err == nil {
			result = b
		}
	}
	if value := os.Getenv(keyToEnv(key)); len(value) != 0 {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return result
		}
		return b
	}
	return result
}
