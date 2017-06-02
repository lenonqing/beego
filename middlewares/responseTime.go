package middlewares

import (
	"fmt"
	"regexp"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var begin = func(ctx *context.Context) {
	ctx.Input.SetData("beginAt", time.Now())
}

var finish = func(ctx *context.Context) {
	re, _ := regexp.Compile("/$")
	path := re.ReplaceAllString(ctx.Input.URL(), "")

	duration := time.Duration(0)
	if beginAt, ok := ctx.Input.GetData("beginAt").(time.Time); ok {
		duration = time.Since(beginAt)
	}
	ms := float64(duration.Nanoseconds()) / 1000000
	beego.Debug(`Router: [`, path, `]`, ` Spend time: `, fmt.Sprintf(`%.5f`, ms), `ms`)
}

func init() {
	beego.InsertFilter(`*`, beego.BeforeRouter, begin)
	beego.InsertFilter(`*`, beego.FinishRouter, finish, false)
}
