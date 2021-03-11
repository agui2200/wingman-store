package cmd

import (
	"bytes"
	"errors"
	"github.com/agui2200/wingman-store/config"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
	"unicode"
)

// schema template for the "init" command.
var tmpl = template.Must(template.New("schema").
	Parse(`package {{ .PackName }}

import "entgo.io/ent"

// {{ .Schema }} holds the schema definition for the {{ .Schema }} entity.
type {{ .Schema }} struct {
	ent.Schema
}

// Fields of the {{ .Schema }}.
func ({{ .Schema }}) Fields() []ent.Field {
	return nil
}

// Edges of the {{ .Schema }}.
func ({{ .Schema }}) Edges() []ent.Edge {
	return nil
}
`))

func InitCommand() *cobra.Command {
	var cfile string
	var pname = "store"
	c := &cobra.Command{
		Use:   "init [schemas]",
		Short: "initialize an environment with zero or more schemas",
		Example: examples(
			"store init Example",
		),
		Args: func(c *cobra.Command, names []string) error {
			if err := cobra.MinimumNArgs(1)(c, names); err != nil {
				return err
			}
			for _, name := range names {
				if !unicode.IsUpper(rune(name[0])) {
					return errors.New("schema names must begin with uppercase")
				}
			}
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			// 判断文件是否存在
			_, err := os.Stat(cfile)
			if err != nil {
				log.Fatalf("load %s file error: %+v", cfile, err)
			}

			// 读取配置
			err = config.LoadConfig(cfile)
			if err != nil {
				log.Fatalf("load config error: %+v", err)
			}
			baseDir, _ := path.Split(cfile)
			// 根据配置文件的路径创建目录
			dir := path.Join(baseDir, config.C.SchemaPackage)
			log.Printf("create target %s dir", dir)
			err = createDir(dir)
			if err != nil {
				log.Fatalf("create target dir error: %+v", err)
			}
			dirinfo := strings.Split(dir, "/")
			if len(dirinfo) > 0 {
				pname = dirinfo[len(dirinfo)-1:][0]
			}
			for _, schema := range args {
				b := bytes.NewBuffer(nil)
				// 根据他要的schema初始化文件进去
				err := tmpl.Execute(b, struct {
					PackName string
					Schema   string
				}{
					Schema:   schema,
					PackName: pname,
				})
				if err != nil {
					log.Fatalf("parse %s schema error: %+v", schema, err)
				}
				// 写入对应文件
				fn := strings.ToLower(schema) + ".go"
				fn = path.Join(dir, fn)
				bb, _ := io.ReadAll(b)
				err = os.WriteFile(fn, bb, 0755)
				if err != nil {
					log.Fatalf("write %s schema error: %+v", fn, err)
				}

			}

		},
	}

	c.Flags().StringVarP(&cfile, "config", "c", defaultConfigFile, "load config file")
	return c
}
