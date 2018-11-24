package migrations1

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
	"io/ioutil"
	"strings"
)

type Driver struct {
	Logger *zap.Logger
	Client *tarantool.Connection
}

func (d Driver) Run() error {
	migrations := strings.Builder{}
	box := packr.New("Migrations", "./tnt_stored_procedures")
	box.Walk(func(name string, file packr.File) error {
		if s, err := ioutil.ReadAll(file); err != nil {
			return err
		} else {
			d.Logger.Info("got file", zap.String("name", name))
			migrations.Write(s)
			return nil
		}
	})
	_, err := d.Client.Eval(migrations.String(), []interface{}{})
	return err
}
