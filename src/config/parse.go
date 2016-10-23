package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

const (
	//参数默认名称
	configParamName = "config"

	//默认配置文件夹名
	configFileName = "config.ini"

	//提示语
	usageText = "config file name"
)

func Init() {

	var configFileName = flag.String(configParamName, configFileName, usageText)

	flag.Parse();


	var config config

	if _, err := toml.DecodeFile(*configFileName, &config); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println(config.PathAlias)
	//
	//for _, name := range []string{"alpha", "beta"} {
	//	s := config[name]
	//	fmt.Printf("Server: %s (ip: %s) in %s created on %s\n",
	//		name, s.IP, s.Config.Location,
	//		s.Config.Created.Format("2006-01-02"))
	//	fmt.Printf("Ports: %v\n", s.Config.Ports)
	//}

}