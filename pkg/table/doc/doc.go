package doc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"framework/pkg/cap/database/mysql"
	"framework/pkg/table/data"
	cap "framework/pkg/table/proto"
	"framework/pkg/table/registry"
)

// ServeHTTP ...
func ServeHTTP(addr string, db *mysql.DB, tr *registry.TableRegistry) {
	http.ListenAndServe(addr, &docHander{reg: tr, db: db})
}

type docHander struct {
	reg *registry.TableRegistry
	db  *mysql.DB
}

var BaseURL = "/tabledoc"

/***************************************************************************
/ - root
/options/*
/meta/*
/meta/xxxx/list
```
****************************************************************************/

func (dh *docHander) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	uri = strings.Replace(uri, BaseURL, "", 1)
	if uri == "/" {
		dh.serveRoot(w)
		return
	}
	segments := strings.Split(uri, "/")
	if len(segments) == 3 {
		segments = segments[1:]
		if segments[0] == "options" {
			dh.serveOption(segments[1], w)
			return
		} else if segments[0] == "meta" {
			dh.serveMeta(segments[1], w)
			return
		}
	} else if len(segments) == 4 {
		if segments[1] == "meta" && strings.HasPrefix(segments[3], "list") {
			u, err := url.Parse(uri)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			tableID := segments[2]

			page := 0
			pageSize := 30
			page, err = strconv.Atoi(u.Query().Get("page"))
			if err != nil {
				page = 0
			}
			pageSize, err = strconv.Atoi(u.Query().Get("pageSize"))
			if err != nil {
				pageSize = 30
			}
			dh.serveData(tableID, page, pageSize, w)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

func (dh *docHander) serveRoot(w http.ResponseWriter) {
	w.Write([]byte("<html>"))
	w.Write([]byte(fmt.Sprintf(`<h1 style="%s">Options</h1>`, titleStyle)))
	w.Write([]byte(fmt.Sprintf(`<table border="1" style="%s">`, tableStyle)))
	w.Write([]byte(`<tr style="background-color: #4CAF50; color:white; font-weight:bold;">`))
	w.Write([]byte(`<td style="font-weight:bold; padding:8px;">NO</td>`))
	w.Write([]byte(`<td style="font-weight:bold; padding:8px;">ID</td>`))
	w.Write([]byte(`<td style="font-weight:bold; padding:8px;">Detail</td>`))
	w.Write([]byte("</tr>"))
	idx := 0
	dh.reg.OptionReg.Store().Range(func(key, _ interface{}) bool {
		idx++
		if idx%2 == 1 {
			w.Write([]byte(`<tr style="background-color: #f2f2f2;">`))
		} else {
			w.Write([]byte(`<tr>`))
		}
		w.Write([]byte(fmt.Sprintf(`<td style="padding:8px;">%d</td>`, idx)))
		w.Write([]byte(fmt.Sprintf(`<td style="padding:8px;">%s</td>`, key)))
		w.Write([]byte(fmt.Sprintf(`<td style="padding:8px;"><a href="%s/options/%s">Link</a></td>`, BaseURL, key)))
		w.Write([]byte(`</tr>`))
		return true
	})
	w.Write([]byte(fmt.Sprintf(`</table>`)))

	w.Write([]byte(fmt.Sprintf(`<h1 style="%s">Metadata</h1>`, titleStyle)))
	w.Write([]byte(fmt.Sprintf(`<table border="1" style="%s">`, tableStyle)))
	w.Write([]byte(`<tr style="background-color: #4CAF50; color:white; font-weight:bold;">`))
	w.Write([]byte(`<td style="font-weight:bold; padding:8px;">NO</td>`))
	w.Write([]byte(`<td style="font-weight:bold; padding:8px;">ID</td>`))
	w.Write([]byte(`<td style="font-weight:bold; padding:8px;">Name</td>`))
	w.Write([]byte(`<td style="font-weight:bold; padding:8px;">Detail</td>`))
	w.Write([]byte("</tr>"))
	idx = 0
	dh.reg.TableMetaReg.Store().Range(func(key, value interface{}) bool {
		idx++
		if idx%2 == 1 {
			w.Write([]byte(`<tr style="background-color: #f2f2f2;">`))
		} else {
			w.Write([]byte(`<tr>`))
		}
		w.Write([]byte(fmt.Sprintf(`<td style="padding:8px;">%d</td>`, idx)))
		w.Write([]byte(fmt.Sprintf(`<td style="padding:8px;">%s</td>`, value.(registry.TableMetaData).ID())))
		w.Write([]byte(fmt.Sprintf(`<td style="padding:8px;">%s</td>`, value.(registry.TableMetaData).Name())))
		w.Write([]byte(fmt.Sprintf(`<td style="padding:8px;"><a href="%s/meta/%s">Link</a></td>`, BaseURL, key)))
		w.Write([]byte(`</tr>`))
		return true
	})
	w.Write([]byte(fmt.Sprintf(`</table>`)))
	w.Write([]byte("</html>"))
}

func (dh *docHander) serveOption(optID string, w http.ResponseWriter) {
	opts, err := dh.reg.OptionReg.GetOptions(optID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("<html>"))
	w.Write([]byte(fmt.Sprintf(`<h1 style="%s">%s</h1>`, titleStyle, optID)))
	w.Write([]byte(fmt.Sprintf(`<table border="1" style="%s">`, tableStyle)))
	w.Write([]byte(`<tr style="background-color: #4CAF50; color:white; font-weight:bold;">`))
	w.Write([]byte(`<td style="font-weight:bold; padding:8px;">ID</td>`))
	w.Write([]byte(`<td style="font-weight:bold; padding:8px;">Name</td>`))
	w.Write([]byte("</tr>"))
	idx := 0
	for _, opt := range opts {
		idx++
		if idx%2 == 1 {
			w.Write([]byte(`<tr style="background-color: #f2f2f2;">`))
		} else {
			w.Write([]byte(`<tr>`))
		}
		w.Write([]byte(fmt.Sprintf(`<td style="padding:8px;">%d</td>`, opt.Id)))
		w.Write([]byte(fmt.Sprintf(`<td style="padding:8px;">%s</td>`, opt.Name)))
		w.Write([]byte(`</tr>`))
	}
	w.Write([]byte("</html>"))
}

func (dh *docHander) serveMeta(tableID string, w http.ResponseWriter) {
	tmd, err := dh.reg.TableMetaReg.Find(tableID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("<html>"))
	w.Write([]byte(fmt.Sprintf(`<h1 style="%s">%s</h1>`, titleStyle, tableID)))
	w.Write([]byte(fmt.Sprintf(`<a href="%s/meta/%s/list?page=0&pageSize=30" style="%s">View Data</a>`, BaseURL, tableID, buttonStyle)))
	w.Write([]byte(`<pre style="word-wrap: break-word; white-space: pre-wrap; font-size: 15px; display: block; word-wrap: break-word;
	font-family: monospace; background: #f4f4f4; border: 1px solid #ddd; margin-bottom: 1.6em; overflow: auto; padding: 1em 1.5em;">`))
	tmd.Print(w)
	w.Write([]byte("</pre>"))
	w.Write([]byte("</html>"))
}

var buttonStyle = `background-color: #4CAF50; text-decoration: none; color:white; 
padding: 10px 24px; display: inline-block; font-family: 'Trebuchet MS', Arial, Helvetica, sans-serif; margin-right:10px`

var tableStyle = `border-collapse:collapse; width:90%; font-size: 16px; font-family: 'Trebuchet MS', Arial, Helvetica, sans-serif;
margin-left:20px; margin-right:20px; margin-top:20px; margin-bottom:20px;`

var (
	titleStyle = `border-collapse:collapse; width:90%; font-size: 32px; font-family: 'Trebuchet MS', Arial, Helvetica, sans-serif;`
	spanStyle  = `margin-right:10px;width:90%; font-size: 16px; font-family: 'Trebuchet MS', Arial, Helvetica, sans-serif;`
)

func (dh *docHander) serveData(tableID string, page, pageSize int, w http.ResponseWriter) {
	tmd, err := dh.reg.TableMetaReg.Find(tableID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	ss, err := dh.db.NewSession()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer ss.Close(err)
	tpl := tmd.DefaultTpl(context.Background())
	rsp, err := data.GlobalManager().FindRows(context.Background(), ss, tpl,
		&cap.PageParam{Page: int32(page), PageSize: int32(pageSize)}, &cap.OrderParam{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("<html>"))
	w.Write([]byte(fmt.Sprintf(`<h1 style="%s">%s Data View</h1>`, titleStyle, tableID)))
	w.Write([]byte(fmt.Sprintf(`<table border="1" style="%s">`, tableStyle)))
	w.Write([]byte(`<tr style="background-color: #4CAF50; color:white; font-weight:bold;">`))
	for _, o := range tpl.Body.Output.VisibleColumns {
		w.Write([]byte(`<td style="font-weight:bold; padding:8px;">`))
		col, err := tmd.Columns().Find(o.ColumnId)
		if err == nil {
		}
		w.Write([]byte(col.Name))
		w.Write([]byte("</td>"))
	}
	w.Write([]byte("</tr>"))
	for i, r := range rsp.Rows {
		if i%2 == 1 {
			w.Write([]byte(`<tr style="background-color: #f2f2f2;">`))
		} else {
			w.Write([]byte(`<tr>`))
		}
		for _, c := range r.Cells {
			w.Write([]byte(`<td style="padding:8px;">`))
			values := c.Values
			if len(values) == 0 {
				values = append(values, c.Value)
			} else {
				if len(values) > 1 {
					fmt.Println(values)
				}
			}
			for i, v := range values {
				if v.Href == "" {
					w.Write([]byte(`<a>`))
				} else {
					w.Write([]byte(fmt.Sprintf(`<a href=%s>`, v.Href)))
				}
				if s, ok := v.V.(*cap.Value_VString); ok {
					w.Write([]byte(s.VString))
				} else if o, ok := v.V.(*cap.Value_VOption); ok {
					w.Write([]byte(o.VOption.Name))
				} else if o, ok := v.V.(*cap.Value_VDouble); ok {
					w.Write([]byte(fmt.Sprintf("%g", o.VDouble)))
				} else if o, ok := v.V.(*cap.Value_VInt); ok {
					w.Write([]byte(fmt.Sprintf("%d", o.VInt)))
				} else if o, ok := v.V.(*cap.Value_VTime); ok {
					w.Write([]byte(o.VTime))
				} else if o, ok := v.V.(*cap.Value_VDate); ok {
					w.Write([]byte(o.VDate))
				} else {
					w.Write([]byte(fmt.Sprintf("%v", v.V)))
				}
				w.Write([]byte(`</a>`))
				if i != len(values)-1 {
					w.Write([]byte(`<br/>`))
				}
			}
			w.Write([]byte("</td>"))

		}
		w.Write([]byte("</tr>"))
	}
	w.Write([]byte("</table>"))
	w.Write([]byte("<br/>"))
	w.Write([]byte(fmt.Sprintf(`<span style="%s"> %d Results | Page %d / %d | %d Items / Page</span>`, spanStyle,
		rsp.PageInfo.TotalResults, rsp.PageInfo.GetCurrentPage()+1, rsp.PageInfo.GetTotalPages(), rsp.PageInfo.PageSize)))
	if page > 0 {
		w.Write([]byte(fmt.Sprintf(`<a href="%s/meta/%s/list?page=%d&pageSize=%d" style="%s">< Pre</a>`, BaseURL, tableID, page-1, pageSize, buttonStyle)))
	}
	if rsp.PageInfo.CurrentPage < rsp.PageInfo.TotalPages-1 {
		w.Write([]byte(fmt.Sprintf(`<a href="%s/meta/%s/list?page=%d&pageSize=%d" style="%s">Next ></a>`, BaseURL, tableID, page+1, pageSize, buttonStyle)))
	}
	w.Write([]byte("</html>"))
}
