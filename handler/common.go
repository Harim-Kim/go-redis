package handler

import "github.com/jinzhu/gorm"

type BaseHandler struct {
	CoreDB		*gorm.DB
}

func (h BaseHandler) WriteLog(code int, err error) {
	return
}

func (h BaseHandler) UpdateLog(code int, err error) {
	return
}

func (h BaseHandler) ReadLog () {
	return
}

func (h BaseHandler) DeleteLog () {
	return
}