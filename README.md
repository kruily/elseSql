# elseSql
| 一个简单的sql语句定义器。

## 安装
```shell
go get github.com/jingxiu1016/elseSql
```
## 使用

```go
package main

import (
	"github.com/jingxiu/elseSql/core"
)

func main() {
	// 第一个参数未表名，第二个车参数为操作类型：insert，delete，update，select
	app,err := core.NewElseApp("user","insert")
	if err != nil && app == nil {
		fmt.Println(err.Error())
	}
}

```
## 增删改查
```go
app,_ := core.NewElseApp("table","insert")
app,_ := core.NewElseApp("table","delete")
app,_ := core.NewElseApp("table","update")
app,_ := core.NewElseApp("table","select")
```

## 解析字段
```go
app.ParseStruct(&struct{
	...
}{})
```
## 排除字段
```go
app.Emit("id,username")
```

## 设置查询条件
```go
app.And("id = 1234")
```

## 联表查询
```go
app.Join("auth","auth.user_id = user.id")
```