package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

func init() {
	if err := NewViperConfig(); err != nil {
		panic(err)
	}
}

var once sync.Once

func NewViperConfig() (err error) {
	once.Do(func() {
		err = newViperConfig()
	})
	return
}

func newViperConfig() error {
	relpath, err := getRelativePathFromCaller()
	if err != nil {
		return err
	}

	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(relpath)
	viper.EnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	_ = viper.BindEnv("stripe-key", "STRIPE_KEY", "endpoint-stripe-secret", "END_STRIPE_SECRET")
	return viper.ReadInConfig()
}

func getRelativePathFromCaller() (relativePath string, err error) {
	callerPwd, err := os.Getwd()
	if err != nil {
		return
	}

	_, here, _, _ := runtime.Caller(0)
	relativePath, err = filepath.Rel(callerPwd, filepath.Dir(here))
	fmt.Printf("caller from %s, here , %s, relpath: %s", callerPwd, here, relativePath)
	return
}
