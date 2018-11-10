package controllers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
)

var loggerKey = "logger"

func LogBody(ctx *gin.Context) {
	logger := ctx.Value(loggerKey).(*zap.Logger)
	if ctx.Request.Method == http.MethodPost || ctx.Request.Method == http.MethodPut {

		buf, _ := ioutil.ReadAll(ctx.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		ctx.Request.Body = rdr2
		logger.Debug("", zap.String("body", readBody(rdr1)))
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
