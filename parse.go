package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)
const NOTELINE = 1;
const CONTENTLINE = 2;

type Config struct {
	path string
	cache map[string]map[string]interface{};
	filePtr *os.File
	error error
	group string
	noteBySlash bool
}

func New (path string,note_by_slash bool) *Config {
	var config = Config{path :path,noteBySlash:note_by_slash,group:"mysql"};
	config.cache = make(map[string]map[string]interface{})
	return &config;
}

func  (this *Config )  Load () bool {
	ok := this.open();
	defer this.filePtr.Close()
	if ok {
		if ok := this.parse();ok {
			return true;
		}
		return false;
	}

	return false;
}

func (this *Config) SetDefaultGroup (str string) *Config {
	this.group = strings.ToLower(str);
	return this;
}

//打开文件系统
func (this *Config ) open () bool {
	this.filePtr,this.error = os.Open(this.path)
	if this.error != nil {
		return false;
	}

	return true;
}


//获取key
func (this *Config) Get (key string) interface{} {
	this.error = nil;
	group,key := this.parseKey(key)
	if group == "*" {
		return nil;
	}

	if group_val,ok := this.GetGroup(group); ok {
		if val,exists := group_val[key];exists {
			return val;
		}
		this.error = errors.New(fmt.Sprintf("未找到group为【%s】 key为【%s】的值",group,key))
		return nil;
	}

	this.error = errors.New(fmt.Sprintf("没有名称为【%s】的group ",group))
	return nil;
}

//获取group
func (this *Config) GetGroup (key string) (map[string]interface{},bool) {
	this.error = nil;
	if val,exists := this.cache[key];exists {
		return val,true;
	}
	this.error = errors.New("未找到该group")
	return nil,false;
}

//解析key
func (this *Config) parseKey (key string) (string,string) {
	tmp := strings.Split(key,".")
	tmp_len := len(tmp)
	var group string

	if tmp_len < 1 {
		group = "*"
		key = "*"
		this.error = errors.New("请输入正确的key")
	} else if tmp_len == 1 {
		group = this.group
		key = strings.ToLower(tmp[0]);
	} else {
		group = strings.ToLower(tmp[0]);
		key = strings.ToLower(tmp[1]);
	}

	return group,key;
}

//解析config文件
func (this *Config) parse () bool {
	var config_str []byte;
	config_str, this.error = ioutil.ReadAll(this.filePtr)
	if this.error != nil {
		return false;
	}

	ok := this.read(config_str)
	if ok {
		return true;
	}
	return false;
}

//读取文件
func (this *Config) read(config_str []byte) bool {

	all_bytes := bytes.FieldsFunc(config_str,this.fieldFunc)
	lens := len(all_bytes)
	var group string = "global";

	for i := 0; i < lens;i++ {
		//清空注释
		tmp,line_type := this.clearNote(all_bytes[i])
		if line_type == NOTELINE {
			continue;
		}

		tmp_len := len(all_bytes[i]);
		if tmp[0] == '[' && tmp[tmp_len - 1] == ']' && tmp_len > 2 {
			group = string(all_bytes[i][1:tmp_len-1]);
			group = strings.ToLower(strings.Trim(group," "))
			if this.cache[group] == nil {
				this.cache[group] = make(map[string]interface{})
			}
		} else {
			key,val,ok := this.parseValue(tmp);
			if (!ok) {
				return false;
			}

			this.cache[group][key] = val;
		}

	}

	return true;
}

//分割文件
func (this *Config) fieldFunc (r rune) bool {
	if r == '\n' || r == '\r' {
		return true;
	}
	return false;
}
//解析值
func (this *Config) parseValue (tmps []byte) (string,interface{},bool) {
	tmp := bytes.Split(tmps,[]byte{'='});
	if len(tmp) != 2 {
		this.error = errors.New("有=号的键值对不完整")
		return "",nil,false;
	}

	//键
	key := strings.ToLower(strings.Trim(string(tmp[0])," "))

	//值
	values := bytes.Split(tmp[1],[]byte{','})
	if len(values) == 1 {
		return key, strings.Trim(string(values[0])," "),true;
	}

	len_values := len(values);
	var res []string =make ([]string,len_values);
	for i := 0; i < len_values; i++ {
		res[i] = strings.Trim(string(values[i])," ")
	}

	return key,res,true;
}

func (this *Config) clearNote (arr_byte []byte) ([]byte,int) {
	arr_byte = bytes.Trim(arr_byte," ")

	if arr_byte[0] == '/' && this.noteBySlash{
		return nil,NOTELINE
	}

	if  arr_byte[0] == '#'  {
		return nil,NOTELINE
	}
	if this.noteBySlash{
		arr_byte = bytes.Split(arr_byte,[]byte("//"))[0]
	}

	arr_byte = bytes.Split(arr_byte,[]byte("#"))[0]
	return arr_byte,CONTENTLINE
}

func (this *Config) LastErr () error {
	return this.error;
}