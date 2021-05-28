package apollo

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func parseFlagWithSubCmd() {
	fs := parseApolloEnv()
	var index = -1
	for k, v :=range os.Args {
		if !strings.HasPrefix(v,"-") {//当前为子命令，跳过
			continue
		} else { //从该参数往后都是正常参数 开始解析flag
			index = k
			break
		}
	}
	if index != -1 {
		err := fs.Parse(os.Args[index:]) //除去子命令
		fmt.Printf("fs.Parse with subcmd error:%v\n", err)
	} else {
		err := fs.Parse(os.Args[index:]) //除去子命令
		fmt.Printf("fs.Parse without subcmd error:%v\n", err)
	}
}


//
// 支持带子命令的命令行解析
// return:
func parseApolloEnv() *flag.FlagSet{
	fs := flag.NewFlagSet("parseApolloEnv", flag.ContinueOnError)
	fs.StringVar(&env, envCmdKey, os.Getenv(envEnvKey), "running env")
	fs.StringVar(&apolloAppId, apolloAppIdCmdKey, os.Getenv(apolloAppIdEnvKey), "apollo AppId")
	fs.StringVar(&apolloServer, apolloServerCmdKey, os.Getenv(apolloServerEnvKey), "apollo Meta server")

	return fs
}