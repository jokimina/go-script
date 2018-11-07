package main

import (
	"github.com/henrylee2cn/pholcus/exec"
	//_ "github.com/henrylee2cn/pholcus_lib" // 此为公开维护的spider规则库
	_ "github.com/jokimina/go-script/pkg/spider_lib" // 同样你也可以自由添加自己的规则库
)

//cmd options:
// -_ui=cmd -a_mode=0 -c_spider=0 -a_outtype=csv -a_thread=20 -a_dockercap=5000 -a_pause=300 -a_proxyminute=0 -a_keyins="<pholcus><golang>" -a_limit=10 -a_success=true -a_failure=true
func main() {
	// 设置运行时默认操作界面，并开始运行
	// 运行软件前，可设置 -a_ui 参数为"web"、"gui"或"cmd"，指定本次运行的操作界面
	// 其中"gui"仅支持Windows系统
	exec.DefaultRun("cmd")
}
