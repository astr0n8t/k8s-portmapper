package internal

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// ConfigStore defines a set of read-only methods for accessing the application
// configuration params as defined in one of the config files.
type ConfigStore interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
	UnmarshalKey(string, interface{}, ...viper.DecoderConfigOption) error
}

var defaultConfig *viper.Viper

// Config returns a default config providers
func Config() ConfigStore {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("cannot determine run status: %v", err)
	}

	// Check if we're running in a dev build
	dir := filepath.Dir(ex)
	if strings.Contains(dir, "go-build") {
		return readViperDevConfig("K8S_PORTMAPPER")
	}

	return readViperConfig("K8S_PORTMAPPER")
}

// LoadConfigProvider returns a configured viper instance
func LoadConfigProvider(appName string) ConfigStore {
	return readViperConfig(appName)
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("service_name", "")
	v.SetDefault("service_namespace", "")
	v.SetDefault("program_name_match", "disabled")
	v.SetDefault("program_name", "")
	v.SetDefault("interval", "5.0")
	v.SetDefault("mode", "production")
}

func setDevOverideDefaults(v *viper.Viper) {
	v.SetDefault("mode", "development")
	v.SetDefault("service_name", "test-service")
	v.SetDefault("service_namespace", "testing")
}

func readViperConfig(appName string) *viper.Viper {
	v := viper.New()

	setDefaults(v)

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.AddConfigPath(".")
	v.AddConfigPath("/etc/k8s-portmapper/")

	v.ReadInConfig()

	v.SetEnvPrefix(appName)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// workaround because viper does not treat env vars the same as other config
	for _, key := range v.AllKeys() {
		val := v.Get(key)
		v.Set(key, val)
	}

	return v
}

func readViperDevConfig(appName string) *viper.Viper {
	v := viper.New()

	setDefaults(v)
	setDevOverideDefaults(v)

	v.SetEnvPrefix(appName)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// workaround because viper does not treat env vars the same as other config
	for _, key := range v.AllKeys() {
		val := v.Get(key)
		v.Set(key, val)
	}

	return v
}
