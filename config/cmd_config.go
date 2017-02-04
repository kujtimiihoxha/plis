package config

import (
	"github.com/asaskevich/govalidator"
	"github.com/kujtimiihoxha/plis/logger"
)

var (
	InputTypes  = []string{"string", "int", "float", "bool"}
	ScriptTypes = []string{"lua"}
)

type GeneratorConfig struct {
	Name            string          `json:"name" valid:"required"`
	Description     string          `json:"description" valid:"required"`
	LongDescription []string        `json:"long_description"`
	Aliases []string        	`json:"aliases"`
	Args            []GeneratorArgs `json:"args"`
	Flags           []GeneratorFlag `json:"flags"`
	SubCommands     []string        `json:"sub_commands"`
	ScriptType      string          `json:"script_type" valid:"scriptType,required"`
}
type GeneratorProjectConfig struct{
	GeneratorName string `json:"generator_name" valid:"required"`
	TestDir string `json:"test_dir" valid:"required"`
}
type GeneratorFlag struct {
	Name        string      `json:"name" valid:"required"`
	Description string      `json:"description" valid:"required"`
	Type        string      `json:"type" valid:"inputType"`
	Default     interface{} `json:"default"`
	Required    bool        `json:"required"`
	Persistent  bool        `json:"persistent"`
	Short       string      `json:"short" valid:"lenOne"`
}
type GeneratorArgs struct {
	Name        string `json:"name" valid:"required"`
	Description string `json:"description" valid:"required"`
	Type        string `json:"type" valid:"inputType"`
	Required    bool   `json:"required"`
}

func (c *GeneratorConfig) Validate() bool {
	result, err := govalidator.ValidateStruct(c)
	if govalidator.ErrorsByField(err)["Name"] != "" {
		logger.GetLogger().Warn("The name of the generator is required")
		return false
	}
	if govalidator.ErrorsByField(err)["Description"] != "" {
		logger.GetLogger().Warn("The description of the generator is required")
		return false
	}
	if govalidator.ErrorsByField(err)["ScriptType"] != "" {
		if c.ScriptType == "" {
			logger.GetLogger().Warn("The generator needs to specify the script type")
			return false
		}
		logger.GetLogger().Warnf("The script type `%s` is not suported , the suported types are `%s`", c.ScriptType, ScriptTypes)
		return false
	}
	return result
}
func (cf *GeneratorFlag) Validate() bool {
	result, err := govalidator.ValidateStruct(cf)
	if govalidator.ErrorsByField(err)["Name"] != "" {
		logger.GetLogger().Warn("The name of the flag is required")
		return false
	}
	if govalidator.ErrorsByField(err)["Description"] != "" {
		logger.GetLogger().Warn("The description of the flag is required")
		return false
	}

	if govalidator.ErrorsByField(err)["Type"] != "" {
		logger.GetLogger().Warnf("The flag type `%s` is not suported , the suported types are `%s`", cf.Type, InputTypes)
		logger.GetLogger().Warnf("The type of `%s` will be set to the default `string` type", cf.Name)
		cf.Type = "string"
		result = true
	}
	if !checkDefault(cf) {
		logger.GetLogger().Warn("The default value of the flag must match the type of the flag")
		return false
	}
	if govalidator.ErrorsByField(err)["Short"] != "" {
		logger.GetLogger().Warn("The shorthand flag can only be one character long")
		cf.Short = ""
		result = true
	}
	return result
}
func checkDefault(flag *GeneratorFlag) bool {
	switch flag.Type {
	case "string":
		if flag.Default == nil {
			flag.Default = ""
		}
		if _, ok := flag.Default.(string); ok {
			return true
		}
		return false
	case "int", "float":
		if flag.Default == nil {
			flag.Default = 0.0
		}
		if _, ok := flag.Default.(float64); ok {
			return true
		}
		return false
	case "bool":
		if flag.Default == nil {
			flag.Default = false
		}
		if _, ok := flag.Default.(bool); ok {
			return true
		}
		return false
	}
	return false
}
func (ca *GeneratorArgs) Validate() bool {
	result, err := govalidator.ValidateStruct(ca)
	if govalidator.ErrorsByField(err)["Name"] != "" {
		logger.GetLogger().Warn("The name of the argument is required")
		return false
	}
	if govalidator.ErrorsByField(err)["Description"] != "" {
		logger.GetLogger().Warn("The description of the argument is required")
		return false
	}
	if govalidator.ErrorsByField(err)["Type"] != "" {
		logger.GetLogger().Warnf("The argument type `%s` is not suported , the suported types are `%s`", ca.Type, InputTypes)
		logger.GetLogger().Warnf("The type of `%s` will be set to the default `string` type", ca.Name)
		ca.Type = "string"
		result = true
	}
	return result
}
