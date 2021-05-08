package gtools

import (
	"fmt"
	"log"
	"testing"
)

func TestSetByFields(t *testing.T) {

	type Respond struct {
		Server int64
		Info   struct {
			Int         int
			Int64       int64
			String      string
			SliceInt    []int
			SliceString []string
			Struct      struct {
				Int    int
				String string
			}
		}
	}

	resp := Respond{}
	_ = SetByFields(&resp, "Server", int64(1564451339282))
	if resp.Server != 1564451339282 {
		t.Errorf(`SetByFields Server no ok. [%v] `, resp.Server)
	} else {
		log.Println("SetByFields Server ok")
	}

	_ = SetByFields(&resp, "Info.Int", int(1))
	if resp.Info.Int != 1 {
		t.Errorf(`SetByFields Info.Int no ok. [%v] `, resp.Info.Int)
	} else {
		log.Println("SetByFields Info.Int ok")
	}

	_ = SetByFields(&resp, "Info.String", "test_string")
	if resp.Info.String != "test_string" {
		t.Errorf(`SetByFields Info.String no ok. [%v] `, resp.Info.String)
	} else {
		log.Println("SetByFields Info.String ok")
	}

	_ = SetByFields(&resp, "Info.Struct.String", "test_struct_string")
	if resp.Info.Struct.String != "test_struct_string" {
		t.Errorf(`SetByFields Info.Struct.String no ok. [%v] `, resp.Info.Struct.String)
	} else {
		log.Println("SetByFields Info.Struct.String ok")
	}

	_ = SetByFields(&resp, "Info.SliceString", "test_struct_string")
	if len(resp.Info.SliceString) != 1 {
		t.Errorf(`SetByFields Info.SliceString no ok. [%v] `, resp.Info.SliceString)
	} else {
		log.Println("SetByFields Info.SliceString  ok")
	}

	log.Println(fmt.Sprintf("resp:%v", resp))
}
