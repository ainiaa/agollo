package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/ainiaa/go-apollo"
	"github.com/shima-park/agollo"
)

type TestConf struct {
	A string `json:"a"`
	B string `json:"b"`
	C int    `json:"c"`
}

func getApolloUnmarshal(c *TestConf) (agollo.Agollo, error) {
	return apollo.StartAndUnmarshal(&c, apollo.BackFile("/tmp/test.log"))
}
func getApollo() (agollo.Agollo, error) {
	return apollo.Start(apollo.BackFile("/tmp/test.log"))
}

/**
 * 执行需要添加命令行参数 test abc def --env TE --apollo_server http://apollocfgv2.xxx.xxx --apollo_appid xxx-service --unmarshal 0
 */
func main() {
	var unmarshal string
	var c TestConf
	flag.StringVar(&unmarshal, "unmarshal", "1", "Unmarshal")
	if unmarshal == "1" {
		_, err := getApolloUnmarshal(&c)
		if err != nil {
			fmt.Printf("apollo.Start error:%s \n", err.Error())
		} else {
			fmt.Println("apollo.Start success")
			fmt.Printf("conf:%+v\n", c)
		}
		go func() {
			for {
				fmt.Printf("[%s]t:%+v\n", time.Now().Format("2006-01-02 15:04:05"), c)
				time.Sleep(2 * time.Second)
			}
		}()
	} else {
		_, err := getApollo()
		if err != nil {
			fmt.Printf("apollo.Start error:%s \n", err.Error())
		} else {
			fmt.Println("apollo.Start success")
		}
	}

	fmt.Printf("apollo.StartAndUnmarshal:%+v\n", c)

	time.Sleep(100 * time.Hour)
}
