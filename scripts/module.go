// Copyright © 2016 Kujtim Hoxha <kujtimii.h@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scripts

import (
	"encoding/json"
	"fmt"
	"github.com/Songmu/prompter"
	"github.com/kujtimiihoxha/plis/fs"
	"github.com/kujtimiihoxha/plis/generators"
	"github.com/kujtimiihoxha/plis/helpers"
	"github.com/kujtimiihoxha/plis/tpl"
	anko_core "github.com/mattn/anko/builtins"
	anko_encoding_json "github.com/mattn/anko/builtins/encoding/json"
	anko_errors "github.com/mattn/anko/builtins/errors"
	anko_flag "github.com/mattn/anko/builtins/flag"
	anko_fmt "github.com/mattn/anko/builtins/fmt"
	anko_io "github.com/mattn/anko/builtins/io"
	anko_io_ioutil "github.com/mattn/anko/builtins/io/ioutil"
	anko_math "github.com/mattn/anko/builtins/math"
	anko_math_rand "github.com/mattn/anko/builtins/math/rand"
	anko_net "github.com/mattn/anko/builtins/net"
	anko_net_http "github.com/mattn/anko/builtins/net/http"
	anko_net_url "github.com/mattn/anko/builtins/net/url"
	anko_os "github.com/mattn/anko/builtins/os"
	anko_os_exec "github.com/mattn/anko/builtins/os/exec"
	anko_os_signal "github.com/mattn/anko/builtins/os/signal"
	anko_path "github.com/mattn/anko/builtins/path"
	anko_path_filepath "github.com/mattn/anko/builtins/path/filepath"
	anko_regexp "github.com/mattn/anko/builtins/regexp"
	anko_runtime "github.com/mattn/anko/builtins/runtime"
	anko_sort "github.com/mattn/anko/builtins/sort"
	anko_strings "github.com/mattn/anko/builtins/strings"
	anko_time "github.com/mattn/anko/builtins/time"
	"github.com/mattn/anko/vm"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/alioygur/godash.v0"
	pongo2 "gopkg.in/flosch/pongo2.v3"
	"os"
	"path"
	"strconv"
	"strings"
	"errors"
)

var pkgs map[string]func(env *vm.Env) *vm.Env

func Build(env *vm.Env, generator *generators.PlisGenerator, args []string) {
	anko_core.Import(env)
	pkgs = map[string]func(env *vm.Env) *vm.Env{
		"encoding/json": anko_encoding_json.Import,
		"errors":        anko_errors.Import,
		"flag":          anko_flag.Import,
		"fmt":           anko_fmt.Import,
		"io":            anko_io.Import,
		"io/ioutil":     anko_io_ioutil.Import,
		"math":          anko_math.Import,
		"math/rand":     anko_math_rand.Import,
		"net":           anko_net.Import,
		"net/http":      anko_net_http.Import,
		"net/url":       anko_net_url.Import,
		"os":            anko_os.Import,
		"os/exec":       anko_os_exec.Import,
		"os/signal":     anko_os_signal.Import,
		"path":          anko_path.Import,
		"path/filepath": anko_path_filepath.Import,
		"regexp":        anko_regexp.Import,
		"runtime":       anko_runtime.Import,
		"sort":          anko_sort.Import,
		"strings":       anko_strings.Import,
		"time":          anko_time.Import,
		"plis": func(env *vm.Env) *vm.Env {
			return plisModule(env, generator, args)
		},
	}

	env.Define("import", func(s string) interface{} {
		if loader, ok := pkgs[s]; ok {
			m := loader(env)
			return m
		} else if loader, ok := pkgs[generator.GetRootParent().Config.Name+"/"+s]; ok {
			//Search for the root function ex. angular2 don't use the current all packages will be in the root func.
			m := loader(env)
			return m
		}
		panic(fmt.Sprintf("package '%s' not found", s))
	})
	env.Define("register", func(m string, obj *vm.Env) {
		if _, ok := pkgs[m]; ok {
			panic(fmt.Sprintf("package '%s' already exists", m))
		} else if _, ok := pkgs[generator.GetRootParent().Config.Name+"/"+m]; ok {
			panic(fmt.Sprintf("package '%s' already exists", m))
		}
		//Search for the root function ex. angular2 don't use the current all packages will be in the root func.
		pkgs[generator.GetRootParent().Config.Name+"/"+m] = func(env *vm.Env) *vm.Env {
			return obj
		}
	})
}
func addUserConfig(plis *vm.Env, command string) {
	data, err := afero.ReadFile(fs.WorkingDirFs(), helpers.GeneratorUserConfigPath(command))
	if err != nil {
		return
	}
	config := &map[string]interface{}{}
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println(fmt.Sprintf("Could not read json from `%s`", helpers.GeneratorUserConfigPath(command)))
		os.Exit(-1)
	}
	plis.Define("UserConfig", *config)
}
func templateFunctions(plis *vm.Env, gen *generators.PlisGenerator) {
	fs.InitTemplatesDirFs(gen.GetRootParent().Config.Name)
	pongo := pongo2.NewSet(gen.Config.Name)
	pongo.SetBaseDirectory(helpers.GeneratorTemplatePath(gen.Config.Name))
	plis.Define("CopyTpl", func(t string, dest string, context map[string]interface{}) error {
		a, _ := afero.IsDir(fs.TemplatesDirFs(), t)
		if a {
			return errors.New("The path must be to a file, please use `CopyAll` to copy complete folders.")
		}
		pongo2.FromFile(t)
		filename := path.Base(t)
		t = strings.TrimPrefix(t, ".")
		t = strings.TrimPrefix(t, "./")
		t = strings.TrimPrefix(t, "/")
		data, err := afero.ReadFile(fs.TemplatesDirFs(), t)
		if err != nil {
			return err
		}
		temp, err := pongo.FromString(string(data))
		if err != nil {
			return err
		}
		addDefaultContextFuncs(context)
		res, err := temp.Execute(context)
		if err != nil {
			return err
		}
		destExt := path.Ext(dest)
		if destExt == "" {
			filename = strings.Replace(filename, ".tpl", "", -1)
		} else {
			filename = path.Base(dest)
			dest = path.Dir(dest)
		}

		return tpl.CopyTpl(res, filename, dest)
	})
	plis.Define("CopyAll", func(v string, dest string, context map[string]interface{}) error {
		v = strings.TrimPrefix(v, ".")
		v = strings.TrimPrefix(v, "./")
		v = strings.TrimPrefix(v, "/")
		v = strings.TrimSuffix(v, "/")
		dest = strings.TrimPrefix(dest, ".")
		dest = strings.TrimPrefix(dest, "./")
		dest = strings.TrimSuffix(dest, "/")
		a, err := afero.IsDir(fs.TemplatesDirFs(), v)
		if !a || err != nil {
			return errors.New("This template folder does not exist")
		}
		if path.Ext(dest) != "" {
			return errors.New("The destination path must be a folder")
		}
		tpls, err := getTemplates(gen.GetRootParent().Config.Name, v)
		if err != nil {
			return err
		}

		for _, t := range tpls {
			filename := path.Base(t)
			directory := path.Dir(t)
			data, err := afero.ReadFile(fs.TemplatesDirFs(), t)
			if err != nil {
				return err
			}
			temp, err := pongo.FromString(string(data))
			if err != nil {
				return err
			}
			addDefaultContextFuncs(context)
			res, err := temp.Execute(context)
			if err != nil {
				return err
			}
			if v != "" {
				dirParts := strings.Split(directory, "/")
				directory = ""
				for _, v := range dirParts[1:] {
					directory += v + "/"
				}
			}
			filename = strings.Replace(filename, ".tpl", "", -1)
			directory = strings.TrimSuffix(directory, ".")
			directory = strings.TrimSuffix(dest+"/"+directory, "/")
			err = tpl.CopyTpl(res, filename, directory)
			if err != nil {
				return err
			}
		}
		return nil

	})
	plis.Define("Mkdir", func(path string) error {
		return fs.WorkingDirFs().Mkdir(path, os.ModePerm)
	})
	plis.Define("MkdirAll", func(path string) error {
		return fs.WorkingDirFs().MkdirAll(path, os.ModePerm)
	})
}
func addDefaultContextFuncs(context map[string]interface{}) {
	//context["len"]= func(v interface{}) {
	//	len(v)
	//};
}
func getTemplates(generator string, p string) ([]string, error) {
	files := []string{}
	err := afero.Walk(fs.WorkingDirFs(), helpers.GeneratorTemplateFile(generator, p), func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			path = strings.Replace(path, "\\", "/", -1)
			files = append(files, strings.Replace(path, strings.Replace(helpers.GeneratorTemplatePath(generator), "\\", "/", -1), "", -1))
		}
		return nil
	})
	return files, err
}
func plisModule(env *vm.Env, gen *generators.PlisGenerator, args []string) *vm.Env {
	plis := env.NewPackage("plis")
	arguments := map[string]interface{}{}
	flags := map[string]interface{}{}
	if gen.Config.Arguments != nil && len(*gen.Config.Arguments) > 0 {
		for i, v := range args {
			arguments[(*gen.Config.Arguments)[i].Name] = argumentToType(v, (*gen.Config.Arguments)[i].Name, (*gen.Config.Arguments)[i].Type)
		}
	}
	if gen.Config.Flags != nil {
		for _, v := range *gen.Config.Flags {
			flags[v.Long] = flagToType(v.Long, gen.Cmd, v.Type)
		}
	}
	addPersistentFlags(gen, &flags)
	plis.Define("Args", arguments)
	plis.Define("Flags", flags)
	addUserConfig(plis, gen.GetRootParent().Config.Name)
	templateFunctions(plis, gen)
	plis.Define("Help", gen.Cmd.Help)
	addGoDashFuncs(plis)
	addPrompterFuncs(plis)
	helperFunctions(plis)
	jsonFuncs(plis)
	fsFuncs(plis)
	return plis
}

func fsFuncs(plis *vm.Env) {
	f := map[string]interface{}{}
	f["ReadFile"] = func(file interface{}) (string, error) {
		val, ok := file.(string)
		if !ok {
			fmt.Println("The file path to ReadFile must be a string")
			os.Exit(-1)
		}
		bt, err := afero.ReadFile(fs.WorkingDirFs(), val)
		return string(bt), err
	}
	f["ReadDir"] = func(file interface{}) ([]os.FileInfo, error) {
		val, ok := file.(string)
		if !ok {
			fmt.Println("The file path to ReadDir must be a string")
			os.Exit(-1)
		}
		return afero.ReadDir(fs.WorkingDirFs(), val)
	}
	f["WriteFile"] = func(file interface{}, data interface{}) error {
		val, ok := file.(string)
		if !ok {
			fmt.Println("The file path to WriteFile must be a string")
			os.Exit(-1)
		}
		dt, ok := data.(string)
		if !ok {
			fmt.Println("The data to WriteFile must be a string")
			os.Exit(-1)
		}
		return afero.WriteFile(fs.WorkingDirFs(), val, []byte(dt), os.ModePerm)
	}
	plis.Define("Fs", f)
}

func jsonFuncs(plis *vm.Env) {
	js := map[string]interface{}{}
	js["Indent"] = json.Indent
	js["Marshal"] = func(data interface{}) string {
		dt, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		return string(dt)
	}
	js["Unmarshal"] = func(data interface{}) (map[string]interface{}, error) {
		dt, ok := data.(string)
		if !ok {
			fmt.Println("The data to Unmarshal must be a string")
			os.Exit(-1)
		}
		resp := map[string]interface{}{}
		err := json.Unmarshal([]byte(dt), &resp)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		return resp, err
	}
	js["MarshalIndent"] = func(data interface{}) string {
		dt, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		return string(dt)
	}
	plis.Define("Json", js)
}
func helperFunctions(plis *vm.Env) {
	hlp := map[string]interface{}{}
	hlp["BasePath"] = helpers.BasePath
	hlp["GeneratorsPath"] = helpers.GeneratorsPath
	hlp["RootGeneratorPath"] = func(gen interface{}) string {
		g, ok := gen.(string)
		if !ok {
			fmt.Println("The generator parameter to RootGeneratorPath must be a string")
			os.Exit(-1)
		}
		return helpers.RootGeneratorPath(g)
	}
	hlp["RootGeneratorConfig"] = func(gen interface{}) string {
		g, ok := gen.(string)
		if !ok {
			fmt.Println("The generator parameter to RootGeneratorConfig must be a string")
			os.Exit(-1)
		}
		return helpers.RootGeneratorConfig(g)
	}
	hlp["RootGeneratorScript"] = func(gen interface{}) string {
		g, ok := gen.(string)
		if !ok {
			fmt.Println("The generator parameter to RootGeneratorScript must be a string")
			os.Exit(-1)
		}
		return helpers.RootGeneratorScript(g)
	}
	hlp["ChildGeneratorConfigPath"] = func(gen interface{}) string {
		g, ok := gen.(string)
		if !ok {
			fmt.Println("The generator parameter to ChildGeneratorConfigPath must be a string")
			os.Exit(-1)
		}
		return helpers.ChildGeneratorConfigPath(g)
	}
	hlp["GeneratorScriptsPath"] = func(gen interface{}) string {
		g, ok := gen.(string)
		if !ok {
			fmt.Println("The generator parameter to GeneratorScriptsPath must be a string")
			os.Exit(-1)
		}
		return helpers.GeneratorScriptsPath(g)
	}
	hlp["GeneratorTemplatesPath"] = func(gen interface{}) string {
		g, ok := gen.(string)
		if !ok {
			fmt.Println("The generator parameter to GeneratorTemplatesPath must be a string")
			os.Exit(-1)
		}
		return helpers.GeneratorTemplatesPath(g)
	}
	hlp["ChildGeneratorConfig"] = func(gen interface{}, child interface{}) string {
		g, ok := gen.(string)
		if !ok {
			fmt.Println("The generator parameter to ChildGeneratorConfig must be a string")
			os.Exit(-1)
		}
		c, ok := child.(string)
		if !ok {
			fmt.Println("The child parameter to ChildGeneratorConfig must be a string")
			os.Exit(-1)
		}
		return helpers.ChildGeneratorConfig(g, c)
	}
	hlp["ChildGeneratorScript"] = func(gen interface{}, child interface{}) string {
		g, ok := gen.(string)
		if !ok {
			fmt.Println("The generator parameter to ChildGeneratorScript must be a string")
			os.Exit(-1)
		}
		c, ok := child.(string)
		if !ok {
			fmt.Println("The child parameter to ChildGeneratorScript must be a string")
			os.Exit(-1)
		}
		return helpers.ChildGeneratorScript(g, c)
	}
	hlp["GeneratorUserConfigPath"] = func(gen interface{}) string {
		g, ok := gen.(string)
		if !ok {
			fmt.Println("The generator parameter to GeneratorUserConfigPath must be a string")
			os.Exit(-1)
		}
		return helpers.GeneratorUserConfigPath(g)
	}
	hlp["GeneratorModulesPath"] = func(gen interface{}) string {
		g, ok := gen.(string)
		if !ok {
			fmt.Println("The generator parameter to GeneratorModulesPath must be a string")
			os.Exit(-1)
		}
		return helpers.GeneratorModulesPath(g)
	}
	hlp["GeneratorTemplatePath"] = func(gen interface{}) string {
		g, ok := gen.(string)
		if !ok {
			fmt.Println("The generator parameter to GeneratorTemplatePath must be a string")
			os.Exit(-1)
		}
		return helpers.GeneratorTemplatePath(g)
	}
	hlp["GeneratorTemplateFile"] = func(gen interface{}, template interface{}) string {
		g, ok := gen.(string)
		if !ok {
			fmt.Println("The generator parameter to GeneratorTemplateFile must be a string")
			os.Exit(-1)
		}
		c, ok := template.(string)
		if !ok {
			fmt.Println("The template parameter to GeneratorTemplateFile must be a string")
			os.Exit(-1)
		}
		return helpers.GeneratorTemplateFile(g, c)
	}
	hlp["GeneratorModulesFile"] = func(gen interface{}, module interface{}) string {
		g, ok := gen.(string)
		if !ok {
			fmt.Println("The generator parameter to GeneratorModulesFile must be a string")
			os.Exit(-1)
		}
		c, ok := module.(string)
		if !ok {
			fmt.Println("The module parameter to GeneratorModulesFile must be a string")
			os.Exit(-1)
		}
		return helpers.GeneratorModulesFile(g, c)
	}
	plis.Define("Helpers", hlp)
}
func addGoDashFuncs(plis *vm.Env) {
	str := make(map[string]interface{})
	str["IsASCII"] = godash.IsASCII
	str["IsAlpha"] = godash.IsAlpha
	str["IsAlphanumeric"] = godash.IsAlphanumeric
	str["IsBase64"] = godash.IsBase64
	str["IsByteLength"] = godash.IsByteLength
	str["IsCreditCard"] = godash.IsCreditCard
	str["IsDNSName"] = godash.IsDNSName
	str["IsDataURI"] = godash.IsDataURI
	str["IsDialString"] = godash.IsDialString
	str["IsDivisibleBy"] = godash.IsDivisibleBy
	str["IsEmail"] = godash.IsEmail
	str["IsFilePath"] = godash.IsFilePath
	str["IsFloat"] = godash.IsFloat
	str["IsFullWidth"] = godash.IsFullWidth
	str["IsHalfWidth"] = godash.IsHalfWidth
	str["IsHexadecimal"] = godash.IsHexadecimal
	str["IsHexcolor"] = godash.IsHexcolor
	str["IsIP"] = godash.IsIP
	str["IsIPv4"] = godash.IsIPv4
	str["IsIPv6"] = godash.IsIPv6
	str["IsISBN"] = godash.IsISBN
	str["IsISBN10"] = godash.IsISBN10
	str["IsISBN13"] = godash.IsISBN13
	str["IsISO3166Alpha2"] = godash.IsISO3166Alpha2
	str["IsISO3166Alpha3"] = godash.IsISO3166Alpha3
	str["IsInRange"] = godash.IsInRange
	str["IsInt"] = godash.IsInt
	str["IsJSON"] = godash.IsJSON
	str["IsLatitude"] = godash.IsLatitude
	str["IsLongitude"] = godash.IsLongitude
	str["IsLowerCase"] = godash.IsLowerCase
	str["IsMAC"] = godash.IsMAC
	str["IsMatches"] = godash.IsMatches
	str["IsMongoID"] = godash.IsMongoID
	str["IsMultibyte"] = godash.IsMultibyte
	str["IsNatural"] = godash.IsNatural
	str["IsNegative"] = godash.IsNegative
	str["IsNonNegative"] = godash.IsNonNegative
	str["IsNonPositive"] = godash.IsNonPositive
	str["IsNull"] = godash.IsNull
	str["IsNumeric"] = godash.IsNumeric
	str["IsPort"] = godash.IsPort
	str["IsPositive"] = godash.IsPositive
	str["IsPrintableASCII"] = godash.IsPrintableASCII
	str["IsRGBcolor"] = godash.IsRGBcolor
	str["IsRequestURI"] = godash.IsRequestURI
	str["IsSSN"] = godash.IsSSN
	str["IsSemver"] = godash.IsSemver
	str["IsStringLength"] = godash.IsStringLength
	str["IsStringMatches"] = godash.IsStringMatches
	str["IsURL"] = godash.IsURL
	str["IsUTFDigit"] = godash.IsUTFDigit
	str["IsUTFLetter"] = godash.IsUTFLetter
	str["IsUTFLetterNumeric"] = godash.IsUTFLetterNumeric
	str["IsUTFNumeric"] = godash.IsUTFNumeric
	str["IsUUID"] = godash.IsUUID
	str["IsUUIDv3"] = godash.IsUUIDv3
	str["IsUUIDv4"] = godash.IsUUIDv4
	str["IsUUIDv5"] = godash.IsUUIDv5
	str["IsUpperCase"] = godash.IsUpperCase
	str["IsVariableWidth"] = godash.IsVariableWidth
	str["IsWhole"] = godash.IsWhole
	str["ToCamelCase"] = godash.ToCamelCase
	str["ToString"] = godash.ToString
	str["ToBoolean"] = godash.ToBoolean
	str["ToSnakeCase"] = godash.ToSnakeCase
	str["ToFloat"] = godash.ToFloat
	str["ToInt"] = godash.ToInt
	str["ToJSON"] = godash.ToJSON
	str["ToKebabCase"] = func(t string) string {
		return strings.Replace(godash.ToSnakeCase(t), "_", "-", -1)
	}
	str["ToLowerFirst"] = func(t string) string {
		str := strings.ToLower(string(t[0]))
		additional := string(t[1:])

		return str + additional
	}
	str["ToStringArray"] = func(t string, s string) []string {
		if t == "" {
			return make([]string, 0)
		}
		str := strings.Split(t, s)
		return str
	}
	plis.Define("Strings", str)
}
func addPrompterFuncs(plis *vm.Env) {
	pr := make(map[string]interface{})
	pr["Prompt"] = prompter.Prompt
	pr["Choose"] = prompter.Choose
	pr["Password"] = prompter.Password
	pr["YN"] = prompter.YN
	pr["YesNo"] = prompter.YesNo
	plis.Define("Prompter", pr)
}
func addPersistentFlags(gen *generators.PlisGenerator, flags *map[string]interface{}) {
	current := gen.Parent
	for current.Config.Name != "plis" {
		if current.Config.Flags == nil {
			current = current.Parent
			continue
		}
		for _, v := range *current.Config.Flags {
			if v.Persistent {
				(*flags)[v.Long] = flagToType(v.Long, gen.Cmd, v.Type)
			}
		}
		current = current.Parent
	}

}
func flagToType(name string, cmd *cobra.Command, tp string) interface{} {
	switch tp {
	case "string":
		v, _ := cmd.Flags().GetString(name)
		return v
	case "int":
		v, _ := cmd.Flags().GetInt64(name)
		return v
	case "float":
		v, _ := cmd.Flags().GetFloat64(name)
		return v

	case "bool":
		v, _ := cmd.Flags().GetBool(name)
		return v
	default:
		return ""
	}
}
func argumentToType(arg string, name string, tp string) interface{} {
	switch tp {
	case "string":
		return arg
	case "int":
		number, err := strconv.ParseInt(arg, 10, 0)
		if err != nil {
			fmt.Println(fmt.Sprintf("The argument '%s' must be an INT type", name))
			os.Exit(-1)
		}
		return number
	case "float":
		floatNumber, err := strconv.ParseFloat(arg, 0)
		if err != nil {
			fmt.Println(fmt.Sprintf("The argument '%s' must be an Float type", name))
			os.Exit(-1)
		}
		return floatNumber
	case "bool":
		boolValue, err := strconv.ParseBool(arg)
		if err != nil {
			fmt.Println(fmt.Sprintf("The argument '%s' must be true/false", name))
			os.Exit(-1)
		}
		return boolValue
	default:
		return arg
	}
}
