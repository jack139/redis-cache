package helper

import (
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v3"
)


type serverYaml struct {
	HTTP2_LISTEN string `yaml:"Http2Listen"`
	REDIS_SERVER string `yaml:"RedisServer"`
	REDIS_PASSWD string `yaml:"RedisPasswd"`
	ORA_CONNECTION string `yaml:"OraConnString"`
	MSSQL_CONNECTION string `yaml:"MssqlConnString"`
	SSL_CERT_PATH string `yaml:"SSLCertPath"`
}


type configYaml struct{
	Server serverYaml `yaml:"Server"`
}


var (
	Settings = configYaml{}
)

func readSettings(yamlFilepath string){
	config, err := ioutil.ReadFile(yamlFilepath)
	if err != nil {
		log.Fatal("Read settings file FAIL: ", err)
	}

	yaml.Unmarshal(config, &Settings)

	log.Println("Settings loaded: ", yamlFilepath)
}

func InitSettings(yamlFilepath string){
	var err error
	readSettings(yamlFilepath)

	// 初始化redis连接, 
	err = redis_init()
	if err!=nil {
		log.Fatal("Redis connecting FAIL: ", err)
	}

	/*
	// 初始化Ora连接, 
	err = ora_init()
	if err!=nil {
		log.Fatal("Oracle connecting FAIL: ", err)
	}
	*/

	// 初始化Mssql连接, 
	err = mssql_init()
	if err!=nil {
		log.Fatal("MS-Sql connecting FAIL: ", err)
	}

}