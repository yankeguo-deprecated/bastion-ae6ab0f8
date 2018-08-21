package view // import "github.com/novakit/view"

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/novakit/binfs"
	"github.com/novakit/nova"
)

const (
	// ContentType header name for Content-Type
	ContentType = "Content-Type"

	// ContentLength header name for Content-Length
	ContentLength = "Content-Length"

	// FileExtension filename extension of template file
	FileExtension = ".html"

	// ContextKey key of View in nova.Context
	ContextKey = "_view"

	// I18nContextKey hard coded key of I18n in nova.Context
	I18nContextKey = "_i18n"
)

// FuncMapPlaceholder placeholder template.FuncMap for request-local functions
var FuncMapPlaceholder = template.FuncMap{
	// placeholder for i18n T func
	"T": func(key string, args ...string) string {
		return "[Error: i18n module not installed]"
	},
}

// Options options for view
type Options struct {
	// Directory of templates, default to "views"
	Directory string

	// BinFS using binfs
	BinFS bool
}

// View the view interface
type View struct {
	// Data data in html rendering
	Data map[string]interface{}

	ctx *nova.Context
	tpl *template.Template
}

type i18n interface {
	T(string, ...string) string
}

// TryUseI18n try to map template function T to I18n module T function
func (v *View) TryUseI18n(c *nova.Context) {
	// skip if no template
	if v.tpl == nil {
		return
	}
	// extract i18n module to package-local interface
	if i, ok := c.Values[I18nContextKey].(i18n); ok {
		v.tpl = v.tpl.Funcs(template.FuncMap{
			"T": func(key string, args ...string) string {
				return i.T(key, args...)
			},
		})
	}
}

// HTML render html template with status code 200
func (v *View) HTML(templateName string) {
	v.RenderHTML(http.StatusOK, templateName)
}

// RenderHTML render html template with custom status code
func (v *View) RenderHTML(statusCode int, templateName string) {
	if len(v.ctx.Res.Header().Get(ContentType)) == 0 {
		v.ctx.Res.Header().Set(ContentType, "text/html")
	}
	// create merged data
	data := make(map[string]interface{})
	for k, v := range v.ctx.Values {
		data[k] = v
	}
	for k, v := range v.Data {
		data[k] = v
	}
	// render template with specified name
	if err := v.tpl.ExecuteTemplate(v.ctx.Res, templateName, data); err != nil {
		panic(err)
	}
}

// DataAsJSON short-hand for v.JSON(v.Data)
func (v *View) DataAsJSON() {
	v.JSON(v.Data)
}

// RenderDataAsJSON short-hand for v.RenderJSON(statusCode, v.Data)
func (v *View) RenderDataAsJSON(statusCode int) {
	v.RenderJSON(statusCode, v.Data)
}

// JSON marshal object as JSON format with code 200
func (v *View) JSON(obj interface{}) {
	v.RenderJSON(http.StatusOK, obj)
}

// RenderJSON marshal object as JSON format with custom status code
func (v *View) RenderJSON(statusCode int, obj interface{}) {
	var err error
	var p []byte
	if p, err = json.Marshal(obj); err != nil {
		panic(err)
	}
	if len(v.ctx.Res.Header().Get(ContentType)) == 0 {
		v.ctx.Res.Header().Set(ContentType, "application/json")
	}
	v.ctx.Res.Header().Set(ContentLength, strconv.Itoa(len(p)))
	v.ctx.Res.WriteHeader(statusCode)
	v.ctx.Res.Write(p)
}

// Text write plain text with status code 200
func (v *View) Text(t string) {
	v.RenderText(http.StatusOK, t)
}

// RenderText with plain text with custom status code
func (v *View) RenderText(statusCode int, t string) {
	if len(v.ctx.Res.Header().Get(ContentType)) == 0 {
		v.ctx.Res.Header().Set(ContentType, "text/plain")
	}
	v.ctx.Res.Header().Set(ContentLength, strconv.Itoa(len(t)))
	v.ctx.Res.WriteHeader(statusCode)
	v.ctx.Res.Write([]byte(t))
}

// Binary write bytes with status code 200
func (v *View) Binary(p []byte) {
	v.RenderBinary(http.StatusOK, p)
}

// RenderBinary write bytes with custom status code
func (v *View) RenderBinary(statusCode int, p []byte) {
	if len(v.ctx.Res.Header().Get(ContentType)) == 0 {
		v.ctx.Res.Header().Set(ContentType, "application/octet-stream")
	}
	v.ctx.Res.Header().Set(ContentLength, strconv.Itoa(len(p)))
	v.ctx.Res.WriteHeader(statusCode)
	v.ctx.Res.Write(p)
}

func sanitizeOptions(opts ...Options) (opt Options) {
	if len(opts) > 0 {
		opt = opts[0]
	}
	if len(opt.Directory) == 0 {
		opt.Directory = "views"
	}
	return
}

// LoadTemplate load template from options
func LoadTemplate(opt Options) *template.Template {
	var err error
	// create the main template
	tpl := template.New("__MAIN__").Funcs(FuncMapPlaceholder)
	// create http.FileSystem
	var fs http.FileSystem
	if opt.BinFS {
		var n *binfs.Node
		n = binfs.Find(strings.Split(opt.Directory, "/")...)
		if n == nil {
			return nil
		}
		fs = n.FileSystem()
	} else {
		fs = http.Dir(opt.Directory)
	}
	// walk filesystem and update template
	if err = walkTemplateFiles(fs, "", tpl); err != nil {
		return nil
	}
	return tpl
}

func walkTemplateFiles(fs http.FileSystem, n string, tpl *template.Template) (err error) {
	// open file
	var f http.File
	if f, err = fs.Open(n); err != nil {
		return
	}
	defer f.Close()
	// stat file
	var fi os.FileInfo
	if fi, err = f.Stat(); err != nil {
		return
	}
	// must be called with a directory
	if !fi.IsDir() {
		err = fmt.Errorf("%s is not a directory", fi.Name())
		return
	}
	// read directory
	var fis []os.FileInfo
	if fis, err = f.Readdir(-1); err != nil {
		return
	}
	// iterate directories
	for _, fi1 := range fis {
		n1 := n + "/" + fi1.Name()
		if fi1.IsDir() {
			// recursive in
			if err = walkTemplateFiles(fs, n1, tpl); err != nil {
				return
			}
		} else {
			if path.Ext(n1) != FileExtension {
				continue
			}
			// template name
			tn := n1[:len(n1)-len(FileExtension)]
			if strings.HasPrefix(tn, "/") {
				tn = tn[1:]
			}
			// open file, read and parse
			var f1 http.File
			if f1, err = fs.Open(n1); err != nil {
				break
			}
			var buf []byte
			if err == nil {
				buf, err = ioutil.ReadAll(f1)
			}
			if err == nil {
				_, err = tpl.New(tn).Parse(string(buf))
			}
			// remember to close file
			f1.Close()
			// return if err
			if err != nil {
				break
			}
		}
	}
	return
}

// Handler create a nova.HandlerFunc
func Handler(opts ...Options) nova.HandlerFunc {
	opt := sanitizeOptions(opts...)
	tpl := LoadTemplate(opt)
	return func(c *nova.Context) error {
		// reload templates if is development and not using BinFS
		if c.Env.IsDevelopment() && !opt.BinFS {
			tpl = LoadTemplate(opt)
		}
		// build view
		v := &View{Data: map[string]interface{}{}, ctx: c, tpl: tpl}
		// try use request-local i18n instance
		v.TryUseI18n(c)
		// inject
		c.Values[ContextKey] = v
		// invoke next
		c.Next()
		return nil
	}
}

// Extract extract View from nova.Context
func Extract(c *nova.Context) (v *View) {
	v, _ = c.Values[ContextKey].(*View)
	return
}
