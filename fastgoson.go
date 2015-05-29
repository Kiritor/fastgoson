package fastgoson

import (
	"bytes"
	"encoding/json"
    "errors"
	"io"
	"io/ioutil"
)

/**
   Value:
   json value,it can be any data type.
*/
type Value struct {
	data        interface{}
	exists      bool           //indicate nil or non-value
}

/**
   JSONObject
   inherits from Value,offers convenient access.
*/

type JSONObject struct {
	Value
	m           map[string]*Value     //key-value
	valid       bool
}

//return map
func (obj *JSONObject) Map() map[string]*Value{
	return obj.m
}

//create a value from an io.reader.

func NewValueByReader(reader io.Reader) (*Value,error){
	value:=new(Value)
	//decode jsonobject from io into Value.data
	decoder:=json.NewDecoder(reader)
	decoder.UseNumber()   //Use Number Type
	err:=decoder.Decode(&value.data)
	return value,err
}

//create a jsonvalue by bytes
func NewValueByBytes(b []byte)(*Value,error){
	reader :=bytes.NewReader(b)
	return NewValueByReader(reader)
}

//create an JSONObject by Reader
func NewJSONObjectByReader(reader io.Reader)(*JSONObject,error){
	return jsonObjectFromValue(NewValueByReader(reader))
}

//create an JSONObject by bytesONObjectByReader
func NewJSONObjectByBytes(b []byte)(*JSONObject,error){
	return jsonObjectFromValue(NewValueByBytes(b))
}

//Value to JSONObject
func jsonObjectFromValue(v *Value,err error)(*JSONObject,error){
	if err!=nil {
		return nil,err
	}
	o,err:=v.Object()
	if err!=nil {
		return nil,err
	}
	return o,err
}


//AttemptsByBytes to current value into an JSONObject
func (v *Value) Object()(*JSONObject,error) {
	var valid bool
	switch v.data.(type) {
		case map[string]interface{}:
		     valid = true
		break
	}

	if valid {
		obj :=new(JSONObject)
		obj.valid = valid
		m:=make(map[string]*Value)
		if valid {
			for key,element :=range v.data.(map[string]interface{}) {
				m[key] = &Value{element,true}
			}
		}
		obj.data = v.data     //set value's data
		obj.m = m
		return obj,nil
	}

	return nil,errors.New("Not JSONObject")
}

//Marshal into bytes.
//generate json into bytes.
func (v*Value) Marshal() ([]byte,error) {
	return json.Marshal(v.data)
}


//get value by key
func (v*Value) get(key string)(*Value,error) {
	obj,err:=v.Object()
	if err!=nil {
		child,ok:obj.Map()[key]
		if ok {
			return child,nil
		}else {
			return nil,errors.New("key not found!")
		}
	}
}
