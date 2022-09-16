package core

import (
	"fmt"
	"reflect"
	"strings"
)

//IsEmpty 校验对象是否未空值
func IsEmpty(obj interface{}) bool {
	if obj == nil {
		return true
	}
	if obj.(string) == "" || len(obj.(string)) <= 0 {
		return true
	}
	return false
}

// NewElseApp 定义一个ElseSql应用
func NewElseApp(table, operate string) (*Application, error) {
	if IsEmpty(table) {
		return nil, NoTableSpecified
	}
	if IsEmpty(operate) {
		return nil, ErrorOperate
	}
	app := &Application{
		Table:    table,
		Selector: "*",
	}
	switch operate {
	case INSERT:
		app.Operate = "insert into @TABLE (@SELECTOR) values (@VALUES)"
	case DELETE:
		app.Operate = "delete from @TABLE where 1 = 1 @CONDITION"
	case UPDATE:
		app.Operate = "update @TABLE set @COLUMN where 1 = 1 @CONDITION"
	case SELECT:
		app.Operate = "select @SELECTOR from @TABLE @JOIN where 1 = 1 @CONDITION"
	default:
		return nil, ErrorOperate
	}
	return app, nil
}

// ParseStruct 从结构体中解析属性到selector
func (a *Application) ParseStruct(st interface{}) error {
	if a == nil {
		return ErrorNilApp
	}
	out := make([]string, 0)
	v := reflect.ValueOf(st)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// we only accept structs
	if v.Kind() != reflect.Struct {
		//panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
		return ErrorParseStruct
	}
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		tagv := fi.Tag.Get(dbTag)
		switch tagv {
		case "-":
			continue
		case "":
			out = append(out, fmt.Sprintf("`%s`", fi.Name))
		default:
			if strings.Contains(tagv, ",") {
				tagv = strings.TrimSpace(strings.Split(tagv, ",")[0])
			}
			if len(tagv) == 0 {
				tagv = fi.Name
			}
			out = append(out, fmt.Sprintf("`%s`", tagv))
		}
	}
	a.Selector = strings.Join(out, ",")
	return nil
}

// Emit 把exclude从selector中移除
func (a *Application) Emit(ex string) error {
	if a == nil {
		return ErrorNilApp
	}
	a.Exclude = ex
	arr := strings.Split(ex, ",")
	a.Selector = strings.Join(Remove(strings.Split(a.Selector, ","), arr), ",")
	return nil
}

// Remove removes given strs from strings.
func Remove(strings []string, strs []string) []string {
	out := append([]string(nil), strings...)
	for _, str := range strs {
		var n int
		for _, v := range out {
			if v != str && v != fmt.Sprintf("`%s`", str) {
				out[n] = v
				n++
			}
		}
		out = out[:n]
	}
	return out
}

// Result 得到sql
func (a *Application) Result() (string, error) {
	if a == nil {
		return "", ErrorNilApp
	}
	ope := strings.Split(a.Operate, " ")[0]
	a.Operate = strings.Replace(a.Operate, "@TABLE", a.Table, 1)
	switch ope {
	case INSERT:
		//a.Operate = "insert into @TABLE (@SELECTOR) values (@VALUES)"
		arr := strings.Split(a.Selector, ",")
		if len(arr) <= 0 {
			return "", ErrorSelector
		}
		strs := make([]string, 0)
		for i := 0; i < len(arr); i++ {
			strs = append(strs, "?")
		}
		a.Operate = strings.Replace(a.Operate, "@SELECTOR", a.Selector, 1)
		a.Operate = strings.Replace(a.Operate, "@VALUES", strings.Join(strs, ","), 1)
	case DELETE:
		//a.Operate = "delete from @TABLE where 1 = 1 @CONDITION"
		a.Operate = strings.Replace(a.Operate, "@CONDITION,", "", 1)
	case UPDATE:
		//a.Operate = "update @TABLE set @COLUMN where 1 = 1 @CONDITION"
		arr := strings.Split(a.Selector, ",")
		if len(arr) <= 0 {
			return "", ErrorSelector
		}
		a.Operate = strings.Replace(a.Operate, "@COLUMN", strings.Join(arr, "=?,")+"=?", 1)
	case SELECT:
		//a.Operate = "select @SELECTOR from @TABLE @JOIN where 1 = 1 @CONDITION"
		a.Operate = strings.Replace(a.Operate, "@SELECTOR", a.Selector, 1)
	}
	a.Operate = strings.Replace(a.Operate, "@CONDITION", "", 1)
	a.Operate = strings.Replace(a.Operate, "@JOIN", "", 1)
	return a.Operate, nil
}

// And 条件增加 示例：And(" a.id = 1")
func (a *Application) And(s string) *Application {
	if a == nil {
		return nil
	}
	index := strings.Index(a.Operate, "@CONDITION")
	if index >= 0 {
		str1 := a.Operate[:index]
		str2 := a.Operate[index:]
		a.Operate = strings.Join([]string{str1, "and", s, str2}, " ")
	}
	return a
}

/*
Join 连表
示例：Join("inner join","a.id = b.id")

	inner join b on a.id = b.id
	left join b on a.id = b.id
	right join b on a.id = b.id
	full join b on a.id = b.id
*/
func (a *Application) Join(join, on string) *Application {
	index := strings.Index(a.Operate, "@JOIN")
	if index >= 0 {
		str1 := a.Operate[:index]
		str2 := a.Operate[index:]
		a.Operate = strings.Join([]string{str1, join, "on", on, str2}, " ")
	}
	return a
}
