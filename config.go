package viper_x

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func ReadConf(configName string, obj interface{}) {
	// read flag
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(fmt.Sprintf("Fatal error while bind pflags: %s\n", err))
	}

	// read env
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// read config file
	if configName != "" {
		viper.SetConfigName(configName)
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/")
		viper.AddConfigPath(".")
		err = viper.ReadInConfig()
		if err != nil {
			panic(fmt.Sprintf("Fatal error while reading config file: %s\n", err))
		}

		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			// TODO when config file changes, what you expect to do in here
		})
	}

	unmarshal(obj)
}

func unmarshal(obj interface{}) {
	typeObj := reflect.TypeOf(obj)
	if typeObj.Kind() != reflect.Ptr {
		panic("obj must be ptr")
	}
	valueObj := reflect.ValueOf(obj)
	parse("", typeObj, valueObj)
}

func parse(prefix string, typeObj reflect.Type, valueObj reflect.Value) {
	if typeObj.Kind() == reflect.Ptr {
		typeObj = typeObj.Elem()
		valueObj = valueObj.Elem()
	}

	switch typeObj.Kind() {
	case reflect.Struct:
		for i := 0; i < typeObj.NumField(); i++ {
			field := typeObj.Field(i)
			parse(getName(prefix, field), field.Type, valueObj.Field(i))
		}
	case reflect.Bool:
		valueObj.SetBool(viper.GetBool(trimPrefixDot(prefix)))

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		valueObj.SetInt(viper.GetInt64(trimPrefixDot(prefix)))

	case reflect.Float32, reflect.Float64:
		valueObj.SetFloat(viper.GetFloat64(trimPrefixDot(prefix)))

	case reflect.String:
		valueObj.SetString(viper.GetString(trimPrefixDot(prefix)))

	default:
		panic("unsupported type")
	}
}

func getName(prefix string, field reflect.StructField) string {
	if field.Tag.Get("viper") != "" {
		prefix = prefix + "." + field.Tag.Get("viper")
	} else {
		prefix = prefix + "." + field.Name
	}
	return prefix
}

func trimPrefixDot(prefix string) string {
	if strings.HasPrefix(prefix, ".") {
		prefix = strings.TrimPrefix(prefix, ".")
	}
	return prefix
}
