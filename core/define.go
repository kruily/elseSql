package core

import "errors"

const (
	dbTag  = "db"
	INSERT = "insert"
	DELETE = "delete"
	UPDATE = "update"
	SELECT = "select"
)

var NoTableSpecified = errors.New("未指定表名")
var ErrorSelector = errors.New("错误的选择器")
var ErrorParseStruct = errors.New("非可解析对象")
var ErrorOperate = errors.New("非支持的数据库操作形式：insert，delete，update，select")
var ErrorNilApp = errors.New("空的elseSql应用")

//var Error

type Application struct {
	Table    string `json:"table"`                // 表名
	Selector string `json:"selector" default:"*"` // 选择器
	Exclude  string `json:"exclude"`              // 排除器
	Operate  string `json:"operate"`              // 操作值：select create update delete insert
}
