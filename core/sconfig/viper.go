// @Title
// @Description
// @Author  Wangwengang  2023/12/10 12:05
// @Update  Wangwengang  2023/12/10 12:05
package sconfig

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/wwengg/threego/core/internal"
)

// Viper //
// 优先级: 命令行 > 环境变量 > 默认值

func Viper(rawValue any, path ...string) *viper.Viper {
	var config string

	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse()
		if config == "" { // 判断命令行参数是否为空
			config = os.Getenv(internal.ConfigEnv)
			fmt.Printf("您正在使用%s环境变量,config的路径为%s\n", internal.ConfigEnv, config)
			if config == "" {
				panic(fmt.Errorf("请传入config.yaml的路径，可通过-c 或者 SIMPLE_CONFIG 环境变量 \n"))
			}
		} else { // 命令行参数不为空 将值赋值于config
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%s\n", config)
		}
	} else { // 函数传递的可变参数的第一个值赋值于config
		config = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%s\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(rawValue); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(rawValue); err != nil {
		panic(err)
	}

	// root 适配性 根据root位置去找到对应迁移位置,保证root路径有效
	//rawVal.RootPath, _ = filepath.Abs("..")
	//
	//// gateway basepath
	Show(rawValue)

	return v
}

// Show Zinx Config Info
func Show(rawValue any) {
	objVal := reflect.ValueOf(rawValue).Elem()
	objType := reflect.TypeOf(rawValue)

	fmt.Println("===== Global Config =====")
	for i := 0; i < objVal.NumField(); i++ {
		field := objVal.Field(i)
		typeField := objType.Field(i)

		fmt.Printf("%s: %v\n", typeField.Name, field.Interface())
	}
	fmt.Println("==============================")
}
