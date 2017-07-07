package common

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"
)

type HTTPResponseFunc func(w http.ResponseWriter, r *http.Request)

//ParseResponse - takes http.Response and pointer to interface,
// parses the data to according with the fields of the interface
func ParseResponse(resp *http.Response, rData interface{}) error {
	log.Println("ParseResponse rData type: ", reflect.TypeOf(rData))
	d := json.NewDecoder(resp.Body)
	err := d.Decode(rData)
	if err != nil {
		log.Println("ParseResponse json decode error: ", err)
		return err
	}
	return nil
}

//ParseJSONBody - takes io.Reader and pointer to interface,
// parses the data to according with the fields of the interface
func ParseJSONBody(body io.Reader, rData interface{}) error {
	log.Println("ParseResponse rData type: ", reflect.TypeOf(rData))
	d := json.NewDecoder(body)
	err := d.Decode(rData)
	if err != nil {
		log.Println("ParseResponse json decode error: ", err)
		return err
	}
	return nil
}
