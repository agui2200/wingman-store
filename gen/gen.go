package gen

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

var genTmpl *template.Template

//go:embed  templates/*.tmpl
var f embed.FS

func init() {

	var err error
	// 载入模板
	genTmpl, err = template.ParseFS(f, "templates/*.tmpl")
	if err != nil {
		log.Fatalf("load template  error : %+v", err)
	}

}

// 生成store相关目录的文件
func Generate(schemaPath, targetPath, packageName, configFile string) error {

	// 获取配置文件的绝对路径
	cpath, err := filepath.Abs(configFile)
	if err != nil {
		return err
	}
	tpath, err := filepath.Abs(targetPath)
	if err != nil {
		return err
	}
	cc, err := filepath.Rel(tpath, cpath)
	if err != nil {
		return err
	}
	// schemaPath 换成相对路径
	absSp, err := filepath.Abs(schemaPath)
	if err != nil {
		return err
	}
	sp, err := filepath.Rel(tpath, absSp)
	if err != nil {
		return err
	}
	d := struct {
		PackageName string
		SchemaPath  string
		ConfigFile  string
	}{
		PackageName: path.Base(packageName),
		SchemaPath:  sp,
		ConfigFile:  cc,
	}
	b, err := executeTemplate("generate.tmpl", d)
	if err != nil {
		return err
	}
	gfile := path.Join(targetPath, "generate.go")
	if _, err := os.Stat(gfile); errors.Is(err, os.ErrNotExist) {
		err = ioutil.WriteFile(gfile, b, 0755)
		if err != nil {
			return err
		}

	}

	b, err = executeTemplate("store.tmpl", struct {
		BasePackage string
		PackageName string
	}{
		BasePackage: packageName,
		PackageName: path.Base(packageName),
	})
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(targetPath, "store.go"), b, 0755)
	if err != nil {
		return err
	}
	return nil
}

func executeTemplate(name string, data interface{}) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	err := genTmpl.ExecuteTemplate(b, name, &data)
	if err != nil {
		return nil, err
	}
	bb, err := io.ReadAll(b)
	if err != nil {
		return nil, err
	}
	return bb, nil

}
