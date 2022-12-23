package model

import "reflect"

var Topics = []string{
	reflect.TypeOf(PaymentKNEXT{}).Name(),
	reflect.TypeOf(InsuranceData{}).Name(),
}

type PaymentKNEXT struct {
	ID     string
	Status string
	Update int
}

type InsuranceData struct {
	ID         string
	IpCallBack string
	Status     string
	Update     int
}
