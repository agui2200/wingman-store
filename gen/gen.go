package gen

import (
	"html/template"
	"io/ioutil"
	"log"
	"path"
)

var genTmpl *template.Template
var storeTmpl *template.Template

func init() {
	var err error
	// 载入模板
	genTmpl, err = template.ParseFiles("./template/generate.tmpl")
	if err != nil {
		log.Fatalf("load template generate.tmpl error : %+v", err)
	}
	storeTmpl, err = template.ParseFiles("./template/store.tmpl")
	if err != nil {
		log.Fatalf("load template store.yaml error : %+v", err)
	}

}

// 生成store相关目录的文件
func Generate(schemaPath, targetPath, pageName string) error {
	err := ioutil.WriteFile(path.Join(targetPath, "generate.go"), []byte(genTmpl.Name()), 0644)
	if err != nil {
		return err
	}
	return nil
}
