package viperx

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	gConfigPath = pflag.String("config", "config.yml", "yaml config file path")
)

// 加载配置
func Load(config interface{}) {
	if reflect.TypeOf(config).Kind() != reflect.Ptr {
		panic("config is't ptr type")
	}

	// 从结构体中解析 tag ，自动注册命令行并绑定参数
	var tag = "yaml"
	NewCmdParser().SetCmdTag(tag).Parse(config)

	// 自动绑定环境变量
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	viper.SetConfigFile(*gConfigPath)
	err := viper.ReadInConfig()
	if err != nil {
		if os.IsNotExist(err) {
			println("[WARN] config file not exist: " + *gConfigPath)
		} else {
			panic("Fatal error while reading config file: " + err.Error())
		}
	}

	err = viper.Unmarshal(config, func(c *mapstructure.DecoderConfig) {
		c.TagName = tag
	})
	if err != nil {
		panic("Fatal error while unmarshal config")
	}
}

type CmdParser struct {
	cmdTag        string
	defaultValTag string
	usageTag      string
}

func NewCmdParser() *CmdParser {
	return &CmdParser{
		cmdTag:        "yaml",
		defaultValTag: "default",
		usageTag:      "usage",
	}
}

func (p *CmdParser) SetCmdTag(tag string) *CmdParser {
	p.cmdTag = tag
	return p
}

func (p *CmdParser) Parse(obj interface{}) {
	typeObj := reflect.TypeOf(obj)
	if typeObj.Kind() == reflect.Ptr || typeObj.Kind() == reflect.Interface {
		typeObj = typeObj.Elem()
	}
	if typeObj.Kind() != reflect.Struct {
		panic("obj is't struct type")
	}

	p.parse("", typeObj)

	pflag.Parse()
}

func (p *CmdParser) parse(prefix string, typeObj reflect.Type) {
	for i := 0; i < typeObj.NumField(); i++ {
		field := typeObj.Field(i)
		switch field.Type.Kind() {
		case reflect.Bool, reflect.String,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:

			p.regPFlag(prefix, field)

		case reflect.Struct:
			p.parse(p.fullTag(prefix, field.Tag.Get(p.cmdTag)), field.Type)

		default:
			panic("unsupported type")
		}
	}
}

func (p *CmdParser) regPFlag(prefix string, field reflect.StructField) {
	if field.Tag.Get(p.cmdTag) == "" {
		return
	}

	key := p.fullTag(prefix, field.Tag.Get(p.cmdTag))
	cmd := p.keyFunc(key)
	val := field.Tag.Get(p.defaultValTag)
	usage := field.Tag.Get(p.usageTag)

	var err error
	switch field.Type.Kind() {
	case reflect.Bool:
		var v bool
		if val != "" {
			v, err = strconv.ParseBool(val)
			if err != nil {
				panic(err)
			}
		}
		pflag.Bool(cmd, v, usage)

	case reflect.Int:
		var v int
		if val != "" {
			v, err = strconv.Atoi(val)
			if err != nil {
				panic(err)
			}
		}
		pflag.Int(cmd, v, usage)

	case reflect.Int8:
		var v int8
		if val != "" {
			tmp, err := strconv.ParseInt(val, 10, 8)
			if err != nil {
				panic(err)
			}
			v = int8(tmp)
		}
		pflag.Int8(cmd, v, usage)

	case reflect.Int16:
		var v int16
		if val != "" {
			tmp, err := strconv.ParseInt(val, 10, 16)
			if err != nil {
				panic(err)
			}
			v = int16(tmp)
		}
		pflag.Int16(cmd, v, usage)

	case reflect.Int32:
		var v int32
		if val != "" {
			tmp, err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				panic(err)
			}
			v = int32(tmp)
		}
		pflag.Int32(cmd, v, usage)

	case reflect.Int64:
		var v int64
		if val != "" {
			v, err = strconv.ParseInt(val, 10, 64)
			if err != nil {
				panic(err)
			}
		}
		pflag.Int64(cmd, v, usage)

	case reflect.Uint:
		var v uint
		if val != "" {
			tmp, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				panic(err)
			}
			v = uint(tmp)
		}
		pflag.Uint(cmd, v, usage)

	case reflect.Uint8:
		var v uint8
		if val != "" {
			tmp, err := strconv.ParseUint(val, 10, 8)
			if err != nil {
				panic(err)
			}
			v = uint8(tmp)
		}
		pflag.Uint8(cmd, v, usage)

	case reflect.Uint16:
		var v uint16
		if val != "" {
			tmp, err := strconv.ParseUint(val, 10, 16)
			if err != nil {
				panic(err)
			}
			v = uint16(tmp)
		}
		pflag.Uint16(cmd, v, usage)

	case reflect.Uint32:
		var v uint32
		if val != "" {
			tmp, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				panic(err)
			}
			v = uint32(tmp)
		}
		pflag.Uint32(cmd, v, usage)

	case reflect.Uint64:
		var v uint64
		if val != "" {
			v, err = strconv.ParseUint(val, 10, 64)
			if err != nil {
				panic(err)
			}
		}
		pflag.Uint64(cmd, v, usage)

	case reflect.Float32:
		var v float32
		if val != "" {
			tmp, err := strconv.ParseFloat(val, 32)
			if err != nil {
				panic(err)
			}
			v = float32(tmp)
		}
		pflag.Float32(cmd, v, usage)

	case reflect.Float64:
		var v float64
		if val != "" {
			v, err = strconv.ParseFloat(val, 64)
			if err != nil {
				panic(err)
			}
		}
		pflag.Float64(cmd, v, usage)

	case reflect.String:
		pflag.String(cmd, val, usage)

	default:
		panic("unsupported type")
	}

	err = viper.BindPFlag(key, pflag.Lookup(cmd))
	if err != nil {
		panic(err)
	}
}

func (p *CmdParser) fullTag(prefix, tag string) string {
	if prefix == "" {
		return tag
	}
	return fmt.Sprintf("%s.%s", prefix, tag)
}

func (p *CmdParser) keyFunc(key string) string {
	key = strings.Replace(key, ".", "-", -1)
	key = strings.Replace(key, "_", "-", -1)
	return key
}
