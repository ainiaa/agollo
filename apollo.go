package apollo

import (
	"bytes"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/shima-park/agollo"
)

var ag agollo.Agollo
var errCh <-chan *agollo.LongPollerError
var watchCh <-chan *agollo.ApolloResponse
var env, apolloAppId, apolloServer string

const (
	apolloAppIdCmdKey  = "apollo_appid"
	apolloAppIdEnvKey  = "apollo_appid"
	apolloServerCmdKey = "apollo_server"
	apolloServerEnvKey = "apollo_server"
	envCmdKey          = "env"
	envEnvKey          = "env"
)

//
// convert apollo config to a struct
// param: c
// param: cs
func convertConf(c interface{}, cs interface{}) (err error) {
	var buff bytes.Buffer
	buff.WriteString("{")
	l := len(c.(agollo.Configurations))
	if l > 0 {
		i := 0
		for key, val := range c.(agollo.Configurations) {
			buff.WriteString(fmt.Sprintf("\"%s\":%s", key, val))
			i++
			if l != i {
				buff.WriteString(",")
			}
		}
	}
	buff.WriteString("}")

	if cs != nil {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		err = json.Unmarshal(buff.Bytes(), &cs)
	}
	return
}

func GetApollo() agollo.Agollo {
	return ag
}

func GetApolloEnv() (e string, appId string, server string) {
	if env == "" || apolloAppId == "" || apolloServer == "" {
		parseFlagWithSubCmd()
	}
	return env, apolloAppId, apolloServer
}

//
// get apollo server config   cmd args > environment variable
// param: ac
// param: opts
func loadApolloConf(ac *Conf, opts Options) {
	parseFlagWithSubCmd()
	if env != "" {
		ac.Cluster = env
	}
	if apolloServer != "" {
		ac.Endpoint = apolloServer
	}
	if apolloAppId != "" {
		ac.AppId = apolloAppId
	}

	ac.BackupFile = opts.BackFile
	ac.DefaultNamespace = opts.DefaultNamespace
	ac.LongPollerInterval = opts.LongPollerInterval
}

//
// 启动goroutine去轮训apollo通知接口
// param: c
// return:
func doStart(c Conf) agollo.Agollo {
	var err error

	// apollo默认的配置文件是properties格式的
	ag, err = agollo.New(c.Endpoint, c.AppId,
		agollo.Cluster(c.Cluster),
		agollo.DefaultNamespace(c.DefaultNamespace),
		agollo.AutoFetchOnCacheMiss(),
		agollo.BackupFile(c.BackupFile),
		agollo.FailTolerantOnBackupExists(),
		agollo.LongPollerInterval(time.Duration(c.LongPollerInterval)*time.Second),
	)
	if err != nil {
		fmt.Printf("go-apollo.New error:%s \n", err.Error())
		return ag
	}

	// 如果想监听并同步服务器配置变化，启动apollo长轮训
	// 返回一个期间发生错误的error channel,按照需要去处理
	errCh = ag.Start()

	return ag
}


// Start 启动goroutine去轮训apollo通知接口
// param: opts
// return:
func Start(opts ...Option) (agollo.Agollo, error) {
	opts = append(opts, WithoutConvertStruct())
	return loadStruct(nil, opts...)
}

// StartAndUnmarshal 启动goroutine去轮训apollo通知接口 序列化到对应的结构体
func StartAndUnmarshal(cs interface{}, opts ...Option) (agollo.Agollo, error) {
	return loadStruct(cs, opts...)
}

// 启动goroutine去轮训apollo通知接口 并且将 配置文件序列化为指定结构体
// param: cs
// param: opts
// return:
func loadStruct(cs interface{}, opts ...Option) (agollo.Agollo, error) {
	var c Conf
	options := newOption(opts...)
	loadApolloConf(&c, options)
	doStart(c)

	if options.ConvertStruct {
		doWatch(&cs)
		//获取当前配置
		v := ag.GetNameSpace(c.DefaultNamespace)
		return nil, convertConf(v, &cs)
	}
	return ag, nil
}

//
// 监控 有变动的时候 更新结构体
// param: cs
func doWatch(cs interface{}) {
	var err error
	// 监听apollo配置更改事件
	// 返回namespace和其变化前后的配置,以及可能出现的error
	watchCh = ag.Watch()

	go func() {
		for {
			select {
			case err := <-errCh:
				fmt.Println("Error:", err)
			case resp := <-watchCh:
				err = convertConf(resp.NewValue, &cs) //配置项有变化的时候 更新配置
				if err != nil {
					fmt.Printf("Watch Apollo error:%s\n", err.Error())
				} else {
					fmt.Printf("newValue:%+v\n", resp.NewValue)
				}
			}
		}
	}()
}
