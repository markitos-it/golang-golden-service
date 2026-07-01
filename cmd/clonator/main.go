package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Field struct {
	Name         string
	Title        string
	Type         string
	EnumValues   string
	DefaultValue string
	Required     bool
	Validation   string
}

var (
	entitySingular string
	entityPlural   string
	customFields   []Field
	fromFile       string
)

var rootCmd = &cobra.Command{
	Use:   "clonator",
	Short: "🚀 Clonator - Golden Microservice Generator",
	Long:  "A cool CLI built with Cobra to clone and customize the markitos-it-service-golden project without relying on bash scripts.",
	Run: func(cmd *cobra.Command, args []string) {
		if fromFile != "" {
			processFromFile(fromFile)
			return
		}
		showMenu()
	},
}

func showMenu() {
	var choice string
	for {
		fmt.Println("\n--- 🚀 Clonator Menu ---")
		fmt.Println("1- New Service")
		fmt.Println("2- Help")
		fmt.Println("3- About")
		fmt.Println("4- Exit")
		fmt.Print("👉 Choose an option: ")

		fmt.Scanln(&choice)

		switch choice {
		case "1":
			NewServiceAction()
		case "2":
			HelpAction()
		case "3":
			AboutAction()
		case "4", "exit", "quit":
			fmt.Println("Bye! 👋")
			os.Exit(0)
		default:
			fmt.Println("❌ Invalid option, please try again.")
		}
	}
}

var cloneCmd = &cobra.Command{
	Use:   "clone [target_service_name]",
	Short: "Clone the current project to a new service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]

		if entitySingular == "" || entityPlural == "" {
			fmt.Println("❌ Error: You must specify --singular and --plural")
			cmd.Help()
			os.Exit(1)
		}

		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("❌ Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		parentDir := filepath.Dir(currentDir)
		targetDir := filepath.Join(parentDir, serviceName)

		fmt.Printf("✨ Cloning project into: %s\n", targetDir)
		if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
			fmt.Printf("❌ Directory %s already exists.\n", targetDir)
			os.Exit(1)
		}

		// Copiar recursivamente (puedes implementarlo 100% en Go luego si quieres eliminar la dependencia de "cp")
		cpCmd := exec.Command("cp", "-r", currentDir, targetDir)
		if err := cpCmd.Run(); err != nil {
			fmt.Printf("❌ Error copying files: %v\n", err)
			os.Exit(1)
		}

		// Limpiar basura del proyecto nuevo
		os.RemoveAll(filepath.Join(targetDir, ".git"))
		os.RemoveAll(filepath.Join(targetDir, "cmd", "clonator")) // Se borra a sí mismo
		os.Remove(filepath.Join(targetDir, "README.md"))

		// Variables de reemplazo
		pluralUpper := strings.ToUpper(entityPlural)
		pluralLower := strings.ToLower(entityPlural)
		pluralTitle := strings.ToUpper(pluralLower[:1]) + pluralLower[1:]

		if len(customFields) > 0 {
			fmt.Println("⚙️ Generating custom fields and value objects dynamically...")
			generateCustomFields(targetDir, customFields)
		}

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " 🔄 Applying customizations (names, imports in testsuite, etc)..."
		s.Start()
		startTime := time.Now()

		err = filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}

			if strings.Contains(path, ".git/") || strings.HasSuffix(path, ".png") {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			strContent := string(content)
			newContent := strings.ReplaceAll(strContent, "markitos-it-service-golden", serviceName)
			newContent = strings.ReplaceAll(newContent, "markitos-it-svc-golden", serviceName)
			newContent = strings.ReplaceAll(newContent, "GOLDEN", pluralUpper)
			newContent = strings.ReplaceAll(newContent, "Golden", pluralTitle)
			newContent = strings.ReplaceAll(newContent, "golden", pluralLower)

			if newContent != strContent {
				os.WriteFile(path, []byte(newContent), d.Type().Perm())
			}
			return nil
		})

		if elapsed := time.Since(startTime); elapsed < time.Second {
			time.Sleep(time.Second - elapsed)
		}
		s.Stop()

		if err != nil {
			fmt.Printf("❌ Error during replacements: %v\n", err)
			os.Exit(1)
		}

		if len(customFields) > 0 {
			fmt.Println("📦 Downloading Faker dependencies for tests...")
			getCmd := exec.Command("go", "get", "github.com/brianvoe/gofakeit/v6")
			getCmd.Dir = targetDir
			getCmd.Run()
		}

		fmt.Println("🧹 Formatting code...")
		fmtCmd := exec.Command("go", "fmt", "./...")
		fmtCmd.Dir = targetDir
		fmtCmd.Run()

		fmt.Println("🛠️ Tidying modules...")
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = targetDir
		tidyCmd.Run()

		fmt.Println("⚙️ Generating Protobuf files...")
		protoCmd := exec.Command("make", "proto")
		protoCmd.Dir = targetDir
		protoCmd.Run()

		fmt.Println("✅ Project successfully cloned and configured using Clonator!")
	},
}

func init() {
	rootCmd.Flags().StringVarP(&fromFile, "from-file", "f", "", "Load configuration from a YAML file")
	cloneCmd.Flags().StringVarP(&entitySingular, "singular", "s", "", "Singular entity name (e.g. user)")
	cloneCmd.Flags().StringVarP(&entityPlural, "plural", "p", "", "Plural entity name (e.g. users)")
	rootCmd.AddCommand(cloneCmd)
}

type YamlConfig struct {
	Project     string `yaml:"project"`
	Entity      string `yaml:"entity"`
	Entities    string `yaml:"entities"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Fields      []struct {
		Name        string      `yaml:"name"`
		Type        string      `yaml:"type"`
		Description string      `yaml:"description"`
		Required    bool        `yaml:"required"`
		Default     interface{} `yaml:"default"`
		Validation  string      `yaml:"validation"`
	} `yaml:"fields"`
}

func processFromFile(filePath string) {
	fmt.Printf("📄 Reading configuration from %s...\n", filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("❌ Error reading file: %v\n", err)
		os.Exit(1)
	}

	var config YamlConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("❌ Error parsing YAML: %v\n", err)
		os.Exit(1)
	}

	serviceName := config.Project
	entitySingular = config.Entity
	entityPlural = config.Entities

	if serviceName == "" || entitySingular == "" || entityPlural == "" {
		fmt.Println("❌ Error: YAML must contain 'project', 'entity', and 'entities' fields.")
		os.Exit(1)
	}

	fmt.Println("--------------------------------------------------")
	fmt.Printf("🚀 Service Name    : %s\n", serviceName)
	fmt.Printf("👤 Entity Singular : %s\n", entitySingular)
	fmt.Printf("👥 Entity Plural   : %s\n", entityPlural)
	fmt.Println("--------------------------------------------------")

	customFields = []Field{}
	for _, f := range config.Fields {
		defVal := ""
		if f.Default != nil {
			if slice, ok := f.Default.([]interface{}); ok && len(slice) == 0 {
				defVal = "" // Si el default es un array vacío []
			} else {
				defVal = fmt.Sprintf("%v", f.Default)
			}
		}

		// Generar un Title bonito (e.g. "my_field" -> "My Field")
		words := strings.Split(f.Name, "_")
		for i, w := range words {
			if len(w) > 0 {
				words[i] = strings.ToUpper(w[:1]) + w[1:]
			}
		}
		title := strings.Join(words, " ")

		customFields = append(customFields, Field{
			Name:         f.Name,
			Title:        title,
			Type:         f.Type,
			DefaultValue: defVal,
			Required:     f.Required,
			Validation:   f.Validation,
		})
	}

	fmt.Printf("✅ Loaded %d custom fields from YAML.\n", len(customFields))
	cloneCmd.Run(cloneCmd, []string{serviceName})
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func generateCustomFields(targetDir string, fields []Field) {
	var structFields, constructorParams, voInstantiations, structAssignments []string
	var testGoldenDefaults, protoFields, mapperToProto, mapperToDomain []string
	var requestArgs, simpleStructFields, responseAssignments []string
	var testStructAssignments, requiredTestStructAssignments []string
	var sqlColumns []string
	var e2eCreateFields, e2eUpdateFields []string
	var foundDomainMapper, foundProtoMapper, foundSql bool
	protoIndex := 10

	for _, f := range fields {
		words := strings.Fields(f.Title)
		for i := range words {
			if len(words[i]) > 0 {
				words[i] = strings.ToUpper(words[i][:1]) + strings.ToLower(words[i][1:])
			}
		}
		titleNoSpace := strings.Join(words, "")
		varName := strings.ToLower(titleNoSpace[:1]) + titleNoSpace[1:]

		fieldType := f.Type
		if f.Type == "enum" {
			fieldType = "string"
		}

		bindingStr := ""
		if f.Required {
			bindingStr = "required"
		}

		tag := fmt.Sprintf("`json:\"%s\"", f.Name)
		if bindingStr != "" {
			tag += fmt.Sprintf(" binding:\"%s\"", bindingStr)
		}
		if f.DefaultValue != "" {
			tag += fmt.Sprintf(" default:\"%s\"", f.DefaultValue)
		}
		tag += "`"

		structFields = append(structFields, fmt.Sprintf("\t%s %s %s", titleNoSpace, fieldType, tag))
		simpleStructFields = append(simpleStructFields, fmt.Sprintf("\t%s %s `json:\"%s\"`", titleNoSpace, fieldType, f.Name))

		constructorParams = append(constructorParams, fmt.Sprintf("%s %s", varName, fieldType))
		requestArgs = append(requestArgs, "request."+titleNoSpace)
		responseAssignments = append(responseAssignments, fmt.Sprintf("\t\t%s: golden.%s,", titleNoSpace, titleNoSpace))

		voInstantiations = append(voInstantiations, fmt.Sprintf("\tsecure%s, err := types.NewGolden%s(%s)\n\tif err != nil {\n\t\treturn nil, err\n\t}", titleNoSpace, titleNoSpace, varName))
		structAssignments = append(structAssignments, fmt.Sprintf("\t\t%s: secure%s.Value(),", titleNoSpace, titleNoSpace))

		testDefVal := `gofakeit.Word()`
		if f.Type == "enum" && f.EnumValues != "" {
			vals := strings.Split(f.EnumValues, ",")
			if len(vals) > 0 {
				testDefVal = `"` + strings.TrimSpace(vals[0]) + `"`
			}
		} else if f.Validation != "" && fieldType == "string" {
			cleanRegex := strings.TrimSuffix(strings.TrimPrefix(f.Validation, "^"), "$")
			testDefVal = fmt.Sprintf("gofakeit.Regex(`%s`)", cleanRegex)
		} else if fieldType == "int" {
			testDefVal = "gofakeit.Number(1, 100)"
		} else if fieldType == "int64" {
			testDefVal = "int64(gofakeit.Number(1, 100))"
		} else if fieldType == "float64" {
			testDefVal = "gofakeit.Float64Range(1.0, 100.0)"
		} else if fieldType == "bool" {
			testDefVal = "gofakeit.Bool()"
		} else if fieldType == "time.Time" {
			testDefVal = "gofakeit.Date()"
		}
		testGoldenDefaults = append(testGoldenDefaults, testDefVal)
		testStructAssignments = append(testStructAssignments, fmt.Sprintf("\t\t%s: %s,", titleNoSpace, testDefVal))
		if f.Required {
			requiredTestStructAssignments = append(requiredTestStructAssignments, fmt.Sprintf("\t\t%s: %s,", titleNoSpace, testDefVal))
		}

		e2eValCreate := `"test_` + strings.ToLower(titleNoSpace) + `"`
		e2eValUpdate := `\"test_` + strings.ToLower(titleNoSpace) + `_mod\"`
		if f.Type == "enum" && f.EnumValues != "" {
			vals := strings.Split(f.EnumValues, ",")
			if len(vals) > 0 {
				e2eValCreate = `"` + strings.TrimSpace(vals[0]) + `"`
				e2eValUpdate = `\"` + strings.TrimSpace(vals[0]) + `\"`
			}
		} else if fieldType == "int" || fieldType == "int64" {
			e2eValCreate = "42"
			e2eValUpdate = "43"
		} else if fieldType == "float64" {
			e2eValCreate = "42.42"
			e2eValUpdate = "43.43"
		} else if fieldType == "bool" {
			e2eValCreate = "true"
			e2eValUpdate = "false"
		} else if fieldType == "time.Time" {
			e2eValCreate = `"2026-05-16T00:00:00Z"`
			e2eValUpdate = `\"2026-05-17T00:00:00Z\"`
		}
		e2eCreateFields = append(e2eCreateFields, fmt.Sprintf(",\n  \"%s\": %s", f.Name, e2eValCreate))
		e2eUpdateFields = append(e2eUpdateFields, fmt.Sprintf(",\n  \\\"%s\\\": %s", f.Name, e2eValUpdate))

		colType := "VARCHAR(255)"
		switch fieldType {
		case "int", "int64":
			colType = "INTEGER"
		case "float64":
			colType = "DECIMAL"
		case "bool":
			colType = "BOOLEAN"
		case "time.Time":
			colType = "TIMESTAMP"
		}
		reqSql := ""
		if f.Required {
			reqSql = " NOT NULL"
		}
		sqlColumns = append(sqlColumns, fmt.Sprintf("    %s %s%s,", f.Name, colType, reqSql))

		protoType := "string"
		switch fieldType {
		case "int":
			protoType = "int32"
		case "int64":
			protoType = "int64"
		case "float64":
			protoType = "double"
		case "bool":
			protoType = "bool"
		}
		protoFields = append(protoFields, fmt.Sprintf("\t%s %s = %d;", protoType, f.Name, protoIndex))
		protoIndex++

		mapperToProto = append(mapperToProto, fmt.Sprintf("\t\t%s: entity.%s,", titleNoSpace, titleNoSpace))
		if f.Type == "time.Time" {
			mapperToDomain = append(mapperToDomain, fmt.Sprintf("\t\t// TODO: Convert req.%s manually", titleNoSpace))
		} else {
			mapperToDomain = append(mapperToDomain, fmt.Sprintf("\t\t%s: req.%s,", titleNoSpace, titleNoSpace))
		}

		generateTypeFile(targetDir, f, titleNoSpace, fieldType)
	}

	modelGoPath := filepath.Join(targetDir, "internal", "domain", "model", "model.go")
	if content, err := os.ReadFile(modelGoPath); err == nil {
		strContent := string(content)
		replacement := strings.Join(structFields, "\n") + "\n\t/* ___CUSTOM_FIELDS___*/"
		strContent = strings.Replace(strContent, "/* ___CUSTOM_FIELDS___*/", replacement, 1)

		if len(constructorParams) > 0 {
			paramsStr := ", " + strings.Join(constructorParams, ", ")
			strContent = strings.Replace(strContent, `baseDir string)`, `baseDir string`+paramsStr+`)`, 1)

			voStr := strings.Join(voInstantiations, "\n\n") + "\n\n\treturn &Golden{"
			strContent = strings.Replace(strContent, `return &Golden{`, voStr, 1)

			assignStr := "UpdatedAt: time.Now(),\n" + strings.Join(structAssignments, "\n")
			strContent = strings.Replace(strContent, `UpdatedAt: time.Now(),`, assignStr, 1)
		}
		os.WriteFile(modelGoPath, []byte(strContent), 0644)
	}

	protoFiles, _ := filepath.Glob(filepath.Join(targetDir, "internal", "infrastructure", "proto", "*.proto"))
	for _, pFile := range protoFiles {
		if pContent, err := os.ReadFile(pFile); err == nil && strings.Contains(string(pContent), "/* ___CUSTOM_FIELDS___*/") {
			pReplacement := strings.Join(protoFields, "\n") + "\n\t/* ___CUSTOM_FIELDS___*/"
			newPContent := strings.ReplaceAll(string(pContent), "/* ___CUSTOM_FIELDS___*/", pReplacement)
			os.WriteFile(pFile, []byte(newPContent), 0644)
		}
	}

	filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || (!strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".sql") && !strings.HasSuffix(path, ".sh")) {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		strContent := string(content)
		modified := false

		if len(structFields) > 0 && strings.Contains(strContent, "/* ___CUSTOM_FIELDS___*/") && filepath.Base(path) != "model.go" {
			rep := strings.Join(structFields, "\n") + "\n\t/* ___CUSTOM_FIELDS___*/"
			strContent = strings.ReplaceAll(strContent, "/* ___CUSTOM_FIELDS___*/", rep)
			modified = true
		}

		if len(simpleStructFields) > 0 && strings.Contains(strContent, "/* ___CUSTOM_STRUCT_FIELDS___*/") {
			rep := strings.Join(simpleStructFields, "\n") + "\n\t/* ___CUSTOM_STRUCT_FIELDS___*/"
			strContent = strings.ReplaceAll(strContent, "/* ___CUSTOM_STRUCT_FIELDS___*/", rep)
			modified = true
		}

		if len(responseAssignments) > 0 && strings.Contains(strContent, "/* ___CUSTOM_RESPONSE_FIELDS___*/") {
			rep := strings.Join(responseAssignments, "\n") + "\n\t\t/* ___CUSTOM_RESPONSE_FIELDS___*/"
			strContent = strings.ReplaceAll(strContent, "/* ___CUSTOM_RESPONSE_FIELDS___*/", rep)
			modified = true
		}

		if len(mapperToProto) > 0 && strings.Contains(strContent, "/* ___CUSTOM_FIELDS_TO_PROTO___*/") {
			rep := strings.Join(mapperToProto, "\n") + "\n\t\t/* ___CUSTOM_FIELDS_TO_PROTO___*/"
			strContent = strings.ReplaceAll(strContent, "/* ___CUSTOM_FIELDS_TO_PROTO___*/", rep)
			modified = true
			foundProtoMapper = true
		}
		if len(mapperToDomain) > 0 && strings.Contains(strContent, "/* ___CUSTOM_FIELDS_TO_DOMAIN___*/") {
			rep := strings.Join(mapperToDomain, "\n") + "\n\t\t/* ___CUSTOM_FIELDS_TO_DOMAIN___*/"
			strContent = strings.ReplaceAll(strContent, "/* ___CUSTOM_FIELDS_TO_DOMAIN___*/", rep)
			modified = true
			foundDomainMapper = true
		}

		// Agresivo: Reemplazar TODAS las llamadas de assert por require en tests para evitar panics absolutos
		if strings.HasSuffix(path, "_test.go") && strings.Contains(strContent, "assert.") {
			strContent = strings.ReplaceAll(strContent, "assert.", "require.")
			strContent = strings.ReplaceAll(strContent, "\"github.com/stretchr/testify/assert\"", "\"github.com/stretchr/testify/require\"")
			modified = true
		}

		if len(testGoldenDefaults) > 0 && (strings.HasSuffix(path, "_test.go") || strings.HasSuffix(path, "mother.go")) {
			if strings.Contains(strContent, "NewGolden(") {
				defStr := strings.Join(testGoldenDefaults, ", ")
				strContent = injectDefaultsToNewGolden(strContent, defStr)
				modified = true
			}
			if len(testStructAssignments) > 0 && strings.Contains(strContent, "/* ___CUSTOM_TEST_FIELDS___*/") {
				rep := strings.Join(testStructAssignments, "\n") + "\n\t\t/* ___CUSTOM_TEST_FIELDS___*/"
				strContent = strings.ReplaceAll(strContent, "/* ___CUSTOM_TEST_FIELDS___*/", rep)
				modified = true
			}

			if modified && strings.Contains(strContent, "gofakeit.") && !strings.Contains(strContent, "github.com/brianvoe/gofakeit/v6") {
				if strings.Contains(strContent, "import (") {
					strContent = strings.Replace(strContent, "import (", "import (\n\t\"github.com/brianvoe/gofakeit/v6\"", 1)
				} else {
					strContent = strings.Replace(strContent, "\n", "\n\nimport \"github.com/brianvoe/gofakeit/v6\"\n", 1)
				}
			}
		}

		if len(sqlColumns) > 0 && strings.HasSuffix(path, ".sql") {
			if strings.Contains(strContent, "-- ___CUSTOM_SQL_FIELDS___") {
				rep := strings.Join(sqlColumns, "\n") + "\n    -- ___CUSTOM_SQL_FIELDS___"
				strContent = strings.ReplaceAll(strContent, "-- ___CUSTOM_SQL_FIELDS___", rep)
				modified = true
				foundSql = true
			}
		}

		if len(e2eCreateFields) > 0 && strings.HasSuffix(path, "test_e2e_grpc.sh") {
			repCreate := strings.Join(e2eCreateFields, "")
			strContent = strings.ReplaceAll(strContent, "__CUSTOM_E2E_FIELDS_CREATE__", repCreate)

			repUpdate := strings.Join(e2eUpdateFields, "")
			strContent = strings.ReplaceAll(strContent, "__CUSTOM_E2E_FIELDS_UPDATE__", repUpdate)
			modified = true
		}

		if len(requestArgs) > 0 && (strings.HasSuffix(path, "create.go") || strings.HasSuffix(path, "update.go") || strings.HasSuffix(path, "handler.go")) {
			if strings.Contains(strContent, "NewGolden(") {
				argsToInject := requestArgs
				if strings.Contains(strContent, "req.Id") || strings.Contains(strContent, "req.Name") {
					var reqArgs []string
					for _, arg := range requestArgs {
						reqArgs = append(reqArgs, strings.Replace(arg, "request.", "req.", 1))
					}
					argsToInject = reqArgs
				} else if strings.Contains(strContent, "cmd.Id") || strings.Contains(strContent, "cmd.Name") {
					var cmdArgs []string
					for _, arg := range requestArgs {
						cmdArgs = append(cmdArgs, strings.Replace(arg, "request.", "cmd.", 1))
					}
					argsToInject = cmdArgs
				}

				defStr := strings.Join(argsToInject, ", ")
				strContent = injectDefaultsToNewGolden(strContent, defStr)
				modified = true
			}
		}

		if modified {
			os.WriteFile(path, []byte(strContent), 0644)
		}
		return nil
	})

	// Alertas inteligentes si el usuario olvidó las marcas en el Golden template
	if len(fields) > 0 {
		if !foundDomainMapper {
			fmt.Println("\n⚠️  WARNING: '/* ___CUSTOM_FIELDS_TO_DOMAIN___*/' not found in any file!")
			fmt.Println("   Your gRPC fields will NOT be mapped to the Domain, causing 'bad request' errors.")
		}
		if !foundProtoMapper {
			fmt.Println("⚠️  WARNING: '/* ___CUSTOM_FIELDS_TO_PROTO___*/' not found in any file!")
			fmt.Println("   Your Domain fields will NOT be mapped to the gRPC response, resulting in missing data.")
		}
		if !foundSql {
			fmt.Println("⚠️  WARNING: '-- ___CUSTOM_SQL_FIELDS___' not found in any .sql file!")
			fmt.Println("   Your database will lack the new columns, causing GORM insert errors.")
		}
	}
}

func generateTypeFile(targetDir string, f Field, titleNoSpace string, fieldType string) {
	typesDir := filepath.Join(targetDir, "internal", "domain", "types")
	os.MkdirAll(typesDir, 0755)
	fileName := filepath.Join(typesDir, f.Name+".go")
	typeName := "Golden" + titleNoSpace

	validationCode := "return true"
	if f.Type == "enum" && f.EnumValues != "" {
		vals := strings.Split(f.EnumValues, ",")
		var conditions []string
		for _, v := range vals {
			conditions = append(conditions, fmt.Sprintf(`value == "%s"`, strings.TrimSpace(v)))
		}
		validationCode = fmt.Sprintf("return %s", strings.Join(conditions, " || "))
	} else if f.Validation != "" && fieldType == "string" {
		if f.Required {
			validationCode = fmt.Sprintf("pattern := `%s`\n\tmatched, err := regexp.MatchString(pattern, value)\n\tif err != nil {\n\t\treturn false\n\t}\n\n\treturn matched", f.Validation)
		} else {
			validationCode = fmt.Sprintf("if value == \"\" {\n\t\treturn true\n\t}\n\tpattern := `%s`\n\tmatched, err := regexp.MatchString(pattern, value)\n\tif err != nil {\n\t\treturn false\n\t}\n\n\treturn matched", f.Validation)
		}
	}

	imports := "\"markitos-it-svc-golden/internal/domain/shared\""
	if fieldType == "time.Time" {
		imports += "\n\t\"time\""
	}
	if f.Validation != "" && fieldType == "string" && f.Type != "enum" {
		imports += "\n\t\"regexp\""
	}

	content := fmt.Sprintf("package types\n\nimport (\n\t%s\n)\n\ntype %s struct {\n\tvalue %s\n}\n\nfunc New%s(value %s) (*%s, error) {\n\tif isValid%s(value) {\n\t\treturn &%s{value}, nil\n\t}\n\n\treturn nil, shared.ErrGoldenBadRequest\n}\n\nfunc isValid%s(value %s) bool {\n\t%s\n}\n\nfunc (b *%s) Value() %s {\n\treturn b.value\n}\n", imports, typeName, fieldType, typeName, fieldType, typeName, typeName, typeName, typeName, fieldType, validationCode, typeName, fieldType)

	os.WriteFile(fileName, []byte(content), 0644)
}

func injectDefaultsToNewGolden(content string, defaults string) string {
	searchStr := "NewGolden("
	var result strings.Builder
	currIdx := 0
	for {
		idx := strings.Index(content[currIdx:], searchStr)
		if idx == -1 {
			result.WriteString(content[currIdx:])
			break
		}
		startIdx := currIdx + idx
		openCount := 0
		matchIdx := -1
		for i := startIdx + len(searchStr) - 1; i < len(content); i++ {
			if content[i] == '(' {
				openCount++
			} else if content[i] == ')' {
				openCount--
				if openCount == 0 {
					matchIdx = i
					break
				}
			}
		}
		if matchIdx != -1 {
			lastNonSpaceIdx := matchIdx - 1
			for lastNonSpaceIdx > startIdx && (content[lastNonSpaceIdx] == ' ' || content[lastNonSpaceIdx] == '\t' || content[lastNonSpaceIdx] == '\n') {
				lastNonSpaceIdx--
			}
			hasComma := content[lastNonSpaceIdx] == ','
			result.WriteString(content[currIdx : lastNonSpaceIdx+1])
			if hasComma {
				result.WriteString(" " + defaults + ",")
			} else {
				result.WriteString(", " + defaults)
			}
			result.WriteString(content[lastNonSpaceIdx+1 : matchIdx])
			currIdx = matchIdx
		} else {
			result.WriteString(content[currIdx : startIdx+len(searchStr)])
			currIdx = startIdx + len(searchStr)
		}
	}
	return result.String()
}
