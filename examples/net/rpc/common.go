package main

const (
	CmdCalc int32 = iota + 1
)

const (
	ServerVersion = 1.0
	ClientVersion = 1.0
	ServerAddr    = ":51000"
)

type CalcBody struct {
	Num1 int32 `json:"num1"`
	Num2 int32 `json:"num2"`
}
