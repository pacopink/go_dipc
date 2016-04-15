package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type RedisCfg struct {
	Addr  string `json:"address"`
	Auth  string `json:"auth"`
	DbNum int    `json:"db"`
}

type AgentCfg struct {
	ID   string `json:"agent"`
	Addr string `json:"address"`
}

type CommonCfg struct {
	Redis RedisCfg          `json:"redis"`
	Agent []AgentCfg        `json:"agent"`
	Extra map[string]string `json:"extra"`
}

func (cfg *RedisCfg) String() string {
	return fmt.Sprintf("Redis:\n\tAddr:%s\n\tAuth:%s\n\tDB:%d\n",
		cfg.Addr, cfg.Auth, cfg.DbNum)
}

func (cfg *AgentCfg) String() string {
	return fmt.Sprintf("Agnet:[%s],Address:[%s]\n", cfg.ID, cfg.Addr)
}

func (cfg *CommonCfg) String() string {
	ss := cfg.Redis.String()
	for i := 0; i < len(cfg.Agent); i++ {
		ss += cfg.Agent[i].String()
	}
	if cfg.Extra != nil {
		for k, v := range cfg.Extra {
			ss += fmt.Sprintf("Extra:[%s]=[%s]\n", k, v)
		}
	}
	return ss
}

func Store2JsonFile(v interface{}, filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("StoreJson2File fail to open file:", filename, "error:", err)
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	return encoder.Encode(v)

}

func Store2JsonFile_Indent(v interface{}, filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("StoreJson2File fail to open file:", filename, "error:", err)
		return err
	}
	defer f.Close()
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}
	f.Write(b)
	return nil
}

func ReadJsonFile2Struct(v interface{}, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	return decoder.Decode(v)
}

func main() {
	f_cfg := flag.String("c", "", "configure file")
	f_store := flag.Bool("s", false, "store mode to save a cfg file, by default false")
	flag.Parse()
	var cfg CommonCfg
	if *f_store {
		cfg.Redis.Addr = "139.122.10.181:6379"
		cfg.Redis.Auth = "ericsson"
		cfg.Redis.DbNum = 3
		if cfg.Extra == nil {
			cfg.Extra = make(map[string]string)
		}
		cfg.Extra["aaa"] = "bbbb"
		cfg.Extra["ccc"] = "dddd"
		cfg.Extra["fff"] = "kkkkkkkkkkkkk"
		agent := &AgentCfg{
			"Agent001",
			"139.122.10.181:60099",
		}
		cfg.Agent = append(cfg.Agent, *agent)
		agent = &AgentCfg{
			"Agent002",
			"139.122.10.182:60099",
		}
		cfg.Agent = append(cfg.Agent, *agent)
		agent = &AgentCfg{
			"Agent003",
			"139.122.10.183:60099",
		}
		cfg.Agent = append(cfg.Agent, *agent)
		fmt.Println("Agent Num:", len(cfg.Agent))
		err := Store2JsonFile_Indent(cfg, *f_cfg)
		if err != nil {
			fmt.Println("StoreJson2File error:", err)
		}
	} else {
		err := ReadJsonFile2Struct(&cfg, *f_cfg)
		if err != nil {
			fmt.Println("ReadJsonFile2Struct failed:", err)
		} else {
			fmt.Println(cfg)
			fmt.Println(cfg.String())
		}
	}
}
