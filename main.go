package main

import (
	"fmt"
	"github.com/jingxiu1016/elseSql/core"
)

func main() {
	app, err := core.NewElseApp("maile_user", "update")
	if err != nil && app == nil {
		fmt.Println(err.Error())
	}
	err = app.ParseStruct(&struct {
		Id       int64  `db:"id"`
		Username string `db:"username"`
		Password string `db:"password"`
	}{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	app.Emit("`password`,id")
	app.And("id = 1234").And("username = 'test'")
	fmt.Printf("%#v\n", app)
	query, _ := app.Result()
	fmt.Printf("result: %#v\n", query)
}
