package j2s

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

type ResponseError struct {
	ResCode string
	Msg     string
}

const (
	SUCCESS = "200"
	ERR1    = "LMS0001"
	ERR2    = "LMS0002"
	ERR3    = "LMS0003"
	ERR4    = "LMS0004"
	ERR5    = "LMS0005"
	ERR6    = "LMS0006"
	ERR7    = "LMS0007"
	ERR8    = "LMS0008"
	ERR9    = "LMS0009"
	ERR10   = "LMS0010"
	ERR11   = "LMS0011"
	ERR12   = "LMS0012"
	ERR13   = "LMS0013"
	ERR14   = "LMS0014"
	ERR15   = "LMS0015"
	ERR16   = "LMS0016"
	ERR17   = "LMS0017"
	ERR18   = "LMS0018"
	ERR19   = "LMS0019"
	ERR20   = "LMS0020"
	ERR21   = "LMS0021"
)

//Maping error
var ResCodeDict = map[string]string{
	"200":     "OK",
	"LMS0003": "Get data fail!",
	"LMS0004": "Convert data fail!",
	"LMS0005": "Overwrite data fail!",
	"LMS0006": "Convert Json fail!",
	"LMS0007": "Insert data fail!",
	"LMS0008": "No data valid!",
	"LMS0009": "Transaction type does not exist",
	"LMS0010": "Statistics type does not exist",
	"LMS0011": "Token Amount not enough to transaction",
	"LMS0012": "Token Invalid",
	"LMS0013": "FromMerchant and ToMerchant are duplicate",
	"LMS0014": "Redeemed",
	"LMS0015": "Activated",
	"LMS0016": "Expired",
	"LMS0017": "Not Activated",
	"LMS0018": "Not Expired",
	"LMS0019": "Invalid Voucher WalletAddress",
	"LMS0020": "Voucher is return value for expired",
	"LMS0021": "HTP error insert data",
}

//ParseArgsToStruct Parse args(string) to Struct
func ParseArgsToStruct(args string, reqObj interface{}) ResponseError {
	var jsonObj map[string]interface{}
	err := json.Unmarshal([]byte(args), &jsonObj)

	if err != nil {
		resErr := ResponseError{
			ResCode: ERR6,
			Msg:     ResCodeDict[ERR6] + " : " + err.Error()}
		return resErr
	}
	stringToDateTimeHook := func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t == reflect.TypeOf(time.Time{}) && f == reflect.TypeOf("") {
			return time.Parse(time.RFC3339, data.(string))
		}

		return data, nil
	}
	decoderConfig := mapstructure.DecoderConfig{
		DecodeHook: stringToDateTimeHook,
		Result:     &reqObj,
	}
	decoder, err := mapstructure.NewDecoder(&decoderConfig)

	errDecode := decoder.Decode(jsonObj)
	if errDecode != nil {
		resErr := ResponseError{
			ResCode: ERR4,
			Msg:     errDecode.Error()}
		return resErr
	}
	mapStructs := structs.Map(reqObj)
	fieldNames := MissingFieldInMap(mapStructs)
	if len(fieldNames) != 0 {
		resErr := ResponseError{
			ResCode: ERR8,
			Msg:     "Missing field in JSON: " + strings.Join(fieldNames, ", ")}
		return resErr
	}

	resErr := ResponseError{
		ResCode: SUCCESS,
		Msg:     ResCodeDict[SUCCESS]}
	return resErr
}

//MissingFieldInMap find field missing in JSON to parse
func MissingFieldInMap(mapStruct map[string]interface{}) []string {
	var fieldNames []string
	for key, value := range mapStruct {
		field := reflect.TypeOf(value)
		// fmt.Printf("key[%s] value[%s]---- kind[%s]\n ", key, value, field.Kind())
		switch field.Kind() {
		case reflect.Float64:
			if value.(float64) == 0 && value.(float64) != -1 {
				fieldNames = append(fieldNames, key+"("+field.Kind().String()+")")
			}
			break
		case reflect.Float32:
			if value.(float32) == 0 && value.(float32) != -1 {
				fieldNames = append(fieldNames, key+"("+field.Kind().String()+")")
			}
			break
		case reflect.Int:
			if value.(int) == 0 && value.(int) != -1 {
				fieldNames = append(fieldNames, key+"("+field.Kind().String()+")")
			}
			break
		case reflect.Int8:
			if value.(int8) == 0 && value.(int8) != -1 {
				fieldNames = append(fieldNames, key+"("+field.Kind().String()+")")
			}
			break
		case reflect.Int16:
			if value.(int16) == 0 && value.(int16) != -1 {
				fieldNames = append(fieldNames, key+"("+field.Kind().String()+")")
			}
			break
		case reflect.Int32:
			if value.(int32) == 0 && value.(int32) != -1 {
				fieldNames = append(fieldNames, key+"("+field.Kind().String()+")")
			}
			break
		case reflect.Int64:
			if value.(int64) == 0 && value.(int64) != -1 {
				fieldNames = append(fieldNames, key+"("+field.Kind().String()+")")
			}
			break
		case reflect.String:
			if len(value.(string)) == 0 && value != "null" {
				fieldNames = append(fieldNames, key+"("+field.Kind().String()+")")
			}
			break
		case reflect.Slice:
			if reflect.ValueOf(value).Len() == 0 {
				fieldNames = append(fieldNames, key+"[...]")
			} else if field.String() == "[]interface {}" {
				arrayMap := value.([]interface{})
				for _, mapObj := range arrayMap {
					fieldChilds := MissingFieldInMap(mapObj.(map[string]interface{}))
					if len(fieldChilds) != 0 {
						fieldNames = append(fieldNames, key+"[")
						fieldNames = append(fieldNames, fieldChilds...)
						fieldNames = append(fieldNames, "]")
					}
				}
			}
			break
		case reflect.Struct:
			defaultTime := time.Time{}
			if field.String() == "time.Time" && value == defaultTime {
				fieldNames = append(fieldNames, key)
			} else {
				mapStructs := structs.Map(value)
				fieldChilds := MissingFieldInMap(mapStructs)
				if len(fieldChilds) != 0 {
					fieldNames = append(fieldNames, key+"{")
					fieldNames = append(fieldNames, fieldChilds...)
					fieldNames = append(fieldNames, "}")
				}
			}
			break
		}
	}
	return fieldNames
}
