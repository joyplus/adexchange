package main

import (
	"github.com/astaxie/beego"
	"github.com/golang/protobuf/proto"
)

func TestProtobuf() {

	test := &Test{
		Label: proto.String("hello"),
		Type:  proto.Int32(17),
		Optionalgroup: &Test_OptionalGroup{
			RequiredField: proto.String("good bye"),
		},
	}
	data, err := proto.Marshal(test)
	if err != nil {
		beego.Error("marshaling error: ", err)
	}
	newTest := &Test{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		beego.Error("unmarshaling error: ", err)
	}
	// Now test and newTest contain the same data.
	if test.GetLabel() != newTest.GetLabel() {
		beego.Error("data mismatch %q != %q", test.GetLabel(), newTest.GetLabel())
	}

	beego.Info("Unmarshalled to: %+v", newTest)

}
