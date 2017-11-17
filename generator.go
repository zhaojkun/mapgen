package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func generator(st Struct) string {
	rs := "// AUTOGENERATED FILE BY github.com/zhaojkun/mapgen\n\npackage main\n"
	rs = rs + generateStruct(st)
	rs = rs + generateInsert(st)
	rs = rs + generateDel(st)
	return rs
}

func generateType(fields []Field) string {
	if len(fields) <= 1 {
		return ""
	}
	res := fields[len(fields)-1].Type
	for i := len(fields) - 2; i >= 0; i-- {
		res = fmt.Sprintf("map[%s]%s", fields[i].Type, res)
	}
	return res
}

func generateStruct(st Struct) string {
	return fmt.Sprintf(`
type %sMapperImpl struct{
   vals %s
}
`, st.Name, generateType(st.Fields))
}

var insertTpl = `
{{.NextFieldName}}s,ok:={{.LastMapper}}[d.{{.FieldName}}]
if !ok{
    {{.NextFieldName}}s = make({{.LeftFields}}) 
    {{.LastMapper}}[d.{{.FieldName}}] = {{.NextFieldName}}s
}
`

func generateInsert(st Struct) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(`func(impl *%sMapperImpl)insert(d %s){`, st.Name, st.Name))
	t := template.Must(template.New("letter").Parse(insertTpl))
	lastMaper := "impl.vals"
	for i := 0; i < len(st.Fields)-2; i++ {
		data := struct {
			NextFieldName string
			LastMapper    string
			FieldName     string
			LeftFields    string
		}{
			NextFieldName: st.Fields[i+1].Name,
			LastMapper:    lastMaper,
			FieldName:     st.Fields[i].Name,
			LeftFields:    generateType(st.Fields[i+1:]),
		}
		t.Execute(&buf, data)
		lastMaper = data.NextFieldName + "s"
	}
	fmt.Fprintf(&buf, "%s[d.%s] = d.%s\n", lastMaper, st.Fields[len(st.Fields)-2].Name, st.Fields[len(st.Fields)-1].Name)
	buf.Write([]byte("}\n"))
	return buf.String()
}

var delTpl = `
{{.NextFieldName}}s,ok:={{.LastMapper}}[d.{{.FieldName}}]
if !ok{
     return
}
defer func(){
   if(len({{.NextFieldName}}s)==0){
         delete({{.LastMapper}},d.{{.FieldName}})
    }
}()
`

func generateDel(st Struct) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(`func(impl *%sMapperImpl)del(d %s){`, st.Name, st.Name))
	t := template.Must(template.New("letter").Parse(delTpl))
	lastMaper := "impl.vals"
	for i := 0; i < len(st.Fields)-2; i++ {
		data := struct {
			NextFieldName string
			LastMapper    string
			FieldName     string
			LeftFields    string
		}{
			NextFieldName: st.Fields[i+1].Name,
			LastMapper:    lastMaper,
			FieldName:     st.Fields[i].Name,
			LeftFields:    generateType(st.Fields[i+1:]),
		}
		t.Execute(&buf, data)
		lastMaper = data.NextFieldName + "s"
	}
	fmt.Fprintf(&buf, "delete(%s,d.%s)", lastMaper, st.Fields[len(st.Fields)-2].Name)
	buf.Write([]byte("}\n"))
	return buf.String()
}