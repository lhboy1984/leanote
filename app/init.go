package app

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lhboy1984/leanote/app/controllers"
	"github.com/lhboy1984/leanote/app/controllers/admin"
	"github.com/lhboy1984/leanote/app/controllers/api"
	"github.com/lhboy1984/leanote/app/controllers/member"
	"github.com/lhboy1984/leanote/app/db"
	. "github.com/lhboy1984/leanote/app/lea"
	_ "github.com/lhboy1984/leanote/app/lea/binder"
	"github.com/lhboy1984/leanote/app/lea/i18n"
	"github.com/lhboy1984/leanote/app/lea/route"
	"github.com/lhboy1984/leanote/app/service"
	"github.com/revel/revel"
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter, // Recover from panics and display an error page instead.
		route.RouterFilter,
		// revel.RouterFilter,            // Use the routing table to select the right Action
		// AuthFilter,						// Invoke the action.
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.

		// 使用SessionFilter标准版从cookie中得到sessionID, 然后通过MssessionFilter从Memcache中得到
		// session, 之后MSessionFilter将session只存sessionID然后返回给SessionFilter返回到web
		// session.SessionFilter,         // leanote session
		// session.MSessionFilter,         // leanote memcache session

		revel.FlashFilter,      // Restore and write the flash cookie.
		revel.ValidationFilter, // Restore kept validation errors and save new ones from cookie.
		// revel.I18nFilter,        // Resolve the requested language
		i18n.I18nFilter,         // Resolve the requested language by leanote
		revel.InterceptorFilter, // Run interceptors around the action.
		revel.CompressFilter,    // Compress the result.
		revel.ActionInvoker,     // Invoke the action.
	}

	revel.TemplateFuncs["raw"] = func(str string) template.HTML {
		return template.HTML(str)
	}
	revel.TemplateFuncs["trim"] = func(str string) string {
		str = strings.Trim(str, " ")
		str = strings.Trim(str, " ")

		str = strings.Trim(str, "\n")
		str = strings.Trim(str, "&nbsp;")

		// 以下两个空格不一样
		str = strings.Trim(str, " ")
		str = strings.Trim(str, " ")
		return str
	}
	revel.TemplateFuncs["add"] = func(i int) string {
		i = i + 1
		return fmt.Sprintf("%v", i)
	}
	revel.TemplateFuncs["sub"] = func(i int) int {
		i = i - 1
		return i
	}
	// 增加或减少
	revel.TemplateFuncs["incr"] = func(n, i int) int {
		n = n + i
		return n
	}
	revel.TemplateFuncs["join"] = func(arr []string) template.HTML {
		if arr == nil {
			return template.HTML("")
		}
		return template.HTML(strings.Join(arr, ","))
	}
	revel.TemplateFuncs["concat"] = func(s1, s2 string) template.HTML {
		return template.HTML(s1 + s2)
	}
	revel.TemplateFuncs["concatStr"] = func(strs ...string) string {
		str := ""
		for _, s := range strs {
			str += s
		}
		return str
	}
	revel.TemplateFuncs["decodeUrlValue"] = func(i string) string {
		v, _ := url.ParseQuery("a=" + i)
		return v.Get("a")
	}
	revel.TemplateFuncs["json"] = func(i interface{}) string {
		b, _ := json.Marshal(i)
		return string(b)
	}
	revel.TemplateFuncs["jsonJs"] = func(i interface{}) template.JS {
		b, _ := json.Marshal(i)
		return template.JS(string(b))
	}
	revel.TemplateFuncs["datetime"] = func(t time.Time) template.HTML {
		return template.HTML(t.Format("2006-01-02 15:04:05"))
	}
	revel.TemplateFuncs["dateFormat"] = func(t time.Time, format string) template.HTML {
		return template.HTML(t.Format(format))
	}
	revel.TemplateFuncs["unixDatetime"] = func(unixSec string) template.HTML {
		sec, _ := strconv.Atoi(unixSec)
		t := time.Unix(int64(sec), 0)
		return template.HTML(t.Format("2006-01-02 15:04:05"))
	}

	// interface是否有该字段
	revel.TemplateFuncs["has"] = func(i interface{}, key string) bool {
		t := reflect.TypeOf(i)
		_, ok := t.FieldByName(key)
		return ok
	}

	// tags
	// 2014/12/30 标签添加链接
	revel.TemplateFuncs["blogTags"] = func(renderArgs map[string]interface{}, tags []string) template.HTML {
		if tags == nil || len(tags) == 0 {
			return ""
		}
		locale, _ := renderArgs[revel.CurrentLocaleRenderArg].(string)
		tagStr := ""
		lenTags := len(tags)

		tagPostUrl, _ := renderArgs["tagPostsUrl"].(string)

		for i, tag := range tags {
			str := revel.Message(locale, tag)
			var classes = "label"
			if strings.HasPrefix(str, "???") {
				str = tag
			}
			if InArray([]string{"red", "blue", "yellow", "green"}, tag) {
				classes += " label-" + tag
			} else {
				classes += " label-default"
			}

			classes += " label-post"
			var url = tagPostUrl + "/" + url.QueryEscape(tag)
			tagStr += "<a class=\"" + classes + "\" href=\"" + url + "\">" + str + "</a>"
			if i != lenTags-1 {
				tagStr += " "
			}
		}
		return template.HTML(tagStr)
	}

	revel.TemplateFuncs["blogTagsForExport"] = func(renderArgs map[string]interface{}, tags []string) template.HTML {
		if tags == nil || len(tags) == 0 {
			return ""
		}
		tagStr := ""
		lenTags := len(tags)

		for i, tag := range tags {
			str := tag
			var classes = "label"
			if InArray([]string{"red", "blue", "yellow", "green"}, tag) {
				classes += " label-" + tag
			} else {
				classes += " label-default"
			}

			classes += " label-post"
			tagStr += "<span class=\"" + classes + "\" >" + str + "</span>"
			if i != lenTags-1 {
				tagStr += " "
			}
		}
		return template.HTML(tagStr)
	}

	revel.TemplateFuncs["msg"] = func(renderArgs map[string]interface{}, message string, args ...interface{}) template.HTML {
		str, ok := renderArgs[revel.CurrentLocaleRenderArg].(string)
		if !ok {
			return ""
		}
		return template.HTML(i18n.Message(str, message, args...))
	}

	// 不用revel的msg
	revel.TemplateFuncs["leaMsg"] = func(renderArgs map[string]interface{}, key string) template.HTML {
		locale, _ := renderArgs[revel.CurrentLocaleRenderArg].(string)
		str := i18n.Message(locale, key)
		if strings.HasPrefix(str, "???") {
			str = key
		}
		return template.HTML(str)
	}

	// lea++
	revel.TemplateFuncs["blogTagsLea"] = func(renderArgs map[string]interface{}, tags []string, typeStr string) template.HTML {
		if tags == nil || len(tags) == 0 {
			return ""
		}
		locale, _ := renderArgs[revel.CurrentLocaleRenderArg].(string)
		tagStr := ""
		lenTags := len(tags)

		tagPostUrl := "http://lea.leanote.com/"
		if typeStr == "recommend" {
			tagPostUrl += "?tag="
		} else if typeStr == "latest" {
			tagPostUrl += "latest?tag="
		} else {
			tagPostUrl += "subscription?tag="
		}

		for i, tag := range tags {
			str := revel.Message(locale, tag)
			var classes = "label"
			if strings.HasPrefix(str, "???") {
				str = tag
			}
			if InArray([]string{"red", "blue", "yellow", "green"}, tag) {
				classes += " label-" + tag
			} else {
				classes += " label-default"
			}
			classes += " label-post"
			var url = tagPostUrl + url.QueryEscape(tag)
			tagStr += "<a class=\"" + classes + "\" href=\"" + url + "\">" + str + "</a>"
			if i != lenTags-1 {
				tagStr += " "
			}
		}
		return template.HTML(tagStr)
	}

	revel.TemplateFuncs["li"] = func(a string) string {
		return ""
	}
	// str连接
	revel.TemplateFuncs["urlConcat"] = func(url string, v ...interface{}) string {
		html := ""
		for i := 0; i < len(v); i = i + 2 {
			item := v[i]
			if i+1 == len(v) {
				break
			}
			value := v[i+1]
			if item != nil && value != nil {
				keyStr, _ := item.(string)
				valueStr, err := value.(string)
				if !err {
					valueInt, _ := value.(int)
					valueStr = strconv.Itoa(valueInt)
				}
				if keyStr != "" && valueStr != "" {
					s := keyStr + "=" + valueStr
					if html != "" {
						html += "&" + s
					} else {
						html += s
					}
				}
			}
		}

		if html != "" {
			if strings.Index(url, "?") >= 0 {
				return url + "&" + html
			} else {
				return url + "?" + html
			}
		}
		return url
	}

	revel.TemplateFuncs["urlCond"] = func(url string, sorterI, keyords interface{}) template.HTML {
		return ""
	}

	// http://stackoverflow.com/questions/14226416/go-lang-templates-always-quotes-a-string-and-removes-comments
	revel.TemplateFuncs["rawMsg"] = func(renderArgs map[string]interface{}, message string, args ...interface{}) template.JS {
		str, ok := renderArgs[revel.CurrentLocaleRenderArg].(string)
		if !ok {
			return ""
		}
		return template.JS(revel.Message(str, message, args...))
	}

	// 为后台管理sorter th使用
	// 必须要返回HTMLAttr, 返回html, golang 会执行安全检查返回ZgotmplZ
	// sorterI 可能是nil, 所以用interfalce{}来接收
	/*
		data-url="/adminUser/index"
		data-sorter="email"
		class="th-sortable {{if eq .sorter "email-up"}}th-sort-up{{else}}{{if eq .sorter "email-down"}}th-sort-down{{end}}{{end}}"
	*/
	revel.TemplateFuncs["sorterTh"] = func(url, sorterField string, sorterI interface{}) template.HTMLAttr {
		sorter := ""
		if sorterI != nil {
			sorter, _ = sorterI.(string)
		}
		html := "data-url=\"" + url + "\" data-sorter=\"" + sorterField + "\""
		html += " class=\"th-sortable "
		if sorter == sorterField+"-up" {
			html += "th-sort-up\""
		} else if sorter == sorterField+"-down" {
			html += "th-sort-down"
		}
		html += "\""
		return template.HTMLAttr(html)
	}

	// pagination
	revel.TemplateFuncs["page"] = func(urlBase string, page, pageSize, count int) template.HTML {
		if count == 0 {
			return ""
		}
		totalPage := int(math.Ceil(float64(count) / float64(pageSize)))

		preClass := ""
		prePage := page - 1
		if prePage == 0 {
			prePage = 1
		}
		nextClass := ""
		nextPage := page + 1
		var preUrl, nextUrl string

		preUrl = urlBase + "?page=" + strconv.Itoa(prePage)
		nextUrl = urlBase + "?page=" + strconv.Itoa(nextPage)

		// 没有上一页了
		if page == 1 {
			preClass = "disabled"
			preUrl = "#"
		}
		// 没有下一页了
		if totalPage <= page {
			nextClass = "disabled"
			nextUrl = "#"
		}
		return template.HTML("<li class='" + preClass + "'><a href='" + preUrl + "'>Previous</a></li> <li  class='" + nextClass + "'><a href='" + nextUrl + "'>Next</a></li>")
	}

	revel.TemplateFuncs["N"] = func(start, end int) (stream chan int) {
		stream = make(chan int)
		go func() {
			for i := start; i <= end; i++ {
				stream <- i
			}
			close(stream)
		}()
		return
	}

	// init Email
	revel.OnAppStart(func() {
		// 数据库
		db.Init()
		// email配置
		InitEmail()
		InitVd()
		// memcache.InitMemcache() // session服务
		// 其它service
		service.InitService()
		controllers.InitService()
		admin.InitService()
		member.InitService()
		service.ConfigS.InitGlobalConfigs()
		api.InitService()
	})
}
