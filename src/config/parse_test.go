package config

import (
	//"fmt"
	//"unsafe"
	//"strconv"
	//"strings"
	"fmt"
	"testing"
	"reflect"
)

func TestReadLinkText(s *testing.T) {
	type VehicleInfo struct {
		// ID         bson.ObjectId `bson:"_id,omitempty"`
		VehicleId string `bson:"编号"`
		Date      string `bson:"日期"`
		Type      string `bson:"类型"`
		Brand     string `bson:"型号"`
		Color     string `bson:"颜色"`
	}
	vinfo := VehicleInfo{
		VehicleId: "123456",
		Date:      "20140101",
		Type:      "Truck",
		Brand:     "Ford",
		Color:     "White",
	}

	fmt.Println(toMap(vinfo))

	vt := reflect.TypeOf(vinfo)
	vv := reflect.ValueOf(vinfo)
	for i := 0; i < vt.NumField(); i++ {
		f := vt.Field(i)
		chKey := f.Tag.Get("bson")
		fmt.Printf("%q.%q => %q, ", f.Name, chKey, vv.FieldByName(f.Name).String())
	}
}

func toMap(any interface{}) map[string]string {
	vt := reflect.TypeOf(any)
	vv := reflect.ValueOf(any)

	ret := make(map[string]string)

	for i := 0; i < vt.NumField(); i++ {
		f := vt.Field(i)
		ret[f.Name] = vv.FieldByName(f.Name).String()
	}
	return ret
}