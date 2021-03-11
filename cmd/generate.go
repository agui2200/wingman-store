package cmd

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/agui2200/wingman-store/config"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"path"
	"strings"
)

func GenerateCommand() *cobra.Command {
	var cfile string
	var Generate = &cobra.Command{
		Use: "generate <path>",
		Example: examples(
			"store generate",
		),
		Run: func(cmd *cobra.Command, args []string) {

			// 检查文件是否存在
			_, err := os.Stat(cfile)
			if err != nil {
				log.Fatalf("config %s error: %+v", cfile, err)
			}
			if err := config.LoadConfig(cfile); err != nil {
				log.Fatalf("config file load error: %v", err)
			}
			var featureLists []gen.Feature
			if config.C.FeatureEntQL {
				featureLists = append(featureLists, gen.FeatureEntQL)
			}
			if config.C.FeaturePrivacy {
				featureLists = append(featureLists, gen.FeaturePrivacy)
			}
			if config.C.FeatureSchemaConfig {
				featureLists = append(featureLists, gen.FeatureSchemaConfig)
			}
			if config.C.FeatureSnapshot {
				featureLists = append(featureLists, gen.FeatureSnapshot)
			}
			baseDir, _ := path.Split(cfile)
			genTarget := "./" + path.Join(baseDir, config.C.SchemaPackage)

			// 拼接path
			pkgs, err := packages.Load(&packages.Config{
				Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo,
			}, genTarget)
			if err != nil {
				log.Fatalf("load schema package error: %+v", err)
			}
			if len(pkgs) == 0 {
				log.Fatalf("unknown %s package ", config.C.SchemaPackage)
			}

			log.Printf("load schema package: %s", genTarget)

			// 解析出应该输出的package
			pi := strings.Split(pkgs[0].String(), "/")
			npack := path.Join(path.Join(pi[:len(pi)-1]...), config.C.TargetPackage)
			err = entc.Generate(genTarget, &gen.Config{
				Features: featureLists,
				Header:   "// Code generated by wingman-store, DO NOT EDIT.",
				Package:  npack,
				Target:   "./" + path.Join(baseDir, config.C.TargetPackage),
			})
			if err != nil {
				log.Fatalf("running ent codegen: %v", err)
			}

			log.Printf("generate package %s done", npack)
		},
	}
	Generate.Flags().StringVarP(&cfile, "config", "c", defaultConfigFile, "load config file")
	return Generate
}
