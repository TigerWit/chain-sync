package controllers

import (
	"github.com/astaxie/beego"
)

type Status struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (s *Status) Success() {
	s.Code = 200
	s.Msg = "success"
}

type MainController struct {
	beego.Controller
}
