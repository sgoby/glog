package gfmt

import (
	"testing"
	"encoding/json"
)

func Test_sprintf(t *testing.T){
	m := make(map[string]int)
	data,_ := json.Marshal(t)
	t.Log(string(data))
	//m["abc"] = 554321
	//print("asdfa")
	str := Sprintf("number = %j5656",m)
	t.Log(str)
}
