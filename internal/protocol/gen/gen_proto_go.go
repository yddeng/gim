package main

import (
	"fmt"
	"initialthree/protocol/cs/proto_def"
	"os"
	"strings"
)

var message_template string = `syntax = "proto2";
package message;
option go_package = "initialthree/protocol/cs/message";

message %s_toS {}

message %s_toC {}`

func gen_proto(out_path string) {

	fmt.Printf("gen_proto message ............\n")

	for _, v := range proto_def.CS_message {
		filename := fmt.Sprintf("%s/%s.proto", out_path, v.Name)
		//检查文件是否存在，如果存在跳过不存在创建
		f, err := os.Open(filename)
		if nil != err && os.IsNotExist(err) {
			f, err = os.Create(filename)
			if nil == err {
				var content string
				content = fmt.Sprintf(message_template, v.Name, v.Name)
				_, err = f.WriteString(content)

				if nil != err {
					fmt.Printf("------ error -------- %s Write error:%s\n", v.Name, err.Error())
				}

				f.Close()

			} else {
				fmt.Printf("------ error --------%s Create error:%s\n", v.Name, err.Error())
			}
		} else if nil != f {
			fmt.Printf("%s.proto exist skip\n", v.Name)
			f.Close()
		}
	}

}

var register_template string = `
package cs
import (
	"initialthree/codec/pb"
	"initialthree/protocol/cs/message"
)

func init() {
	//toS
%s
	//toC
%s
}
`

//产生协议注册文件
func gen_register(out_path string) {

	f, err := os.OpenFile(out_path, os.O_RDWR, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(out_path)
			if err != nil {
				fmt.Printf("------ error -------- create %s failed:%s", out_path, err.Error())
				return
			}
		} else {
			fmt.Printf("------ error -------- open %s failed:%s", out_path, err.Error())
			return
		}
	}

	err = os.Truncate(out_path, 0)

	if err != nil {
		fmt.Printf("------ error -------- Truncate %s failed:%s", out_path, err.Error())
		return
	}

	toS_str := ""
	toC_str := ""

	nameMap := map[string]bool{}
	idMap := map[int]bool{}

	for _, v := range proto_def.CS_message {

		if ok, _ := nameMap[v.Name]; ok {
			panic("duplicate message:" + v.Name)
		}

		if ok, _ := idMap[v.MessageID]; ok {
			panic(fmt.Sprintf("duplicate messageID: %d", v.MessageID))
		}

		nameMap[v.Name] = true
		idMap[v.MessageID] = true

		toS_str = toS_str + fmt.Sprintf(`	pb.Register("cs",&message.%sToS{},%d)`, strings.Title(v.Name), v.MessageID) + "\n"
		toC_str = toC_str + fmt.Sprintf(`	pb.Register("sc",&message.%sToC{},%d)`, strings.Title(v.Name), v.MessageID) + "\n"
	}

	content := fmt.Sprintf(register_template, toS_str, toC_str)

	_, err = f.WriteString(content)

	//fmt.Printf(content)

	if nil != err {
		fmt.Printf("------ error -------- %s Write error:%s\n", out_path, err.Error())
	} else {
		fmt.Printf("%s Write ok\n", out_path)
	}

	f.Close()

}

func main() {
	os.MkdirAll("../message", os.ModePerm)
	gen_proto("../proto/message")
	gen_register("../register.go")
	fmt.Printf("cs gen_proto_go ok!\n")
}
