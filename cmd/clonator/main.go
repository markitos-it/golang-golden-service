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

var (
	entitySingular string
	entityPlural   string
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

		cpCmd := exec.Command("cp", "-r", currentDir, targetDir)
		if err := cpCmd.Run(); err != nil {
			fmt.Printf("❌ Error copying files: %v\n", err)
			os.Exit(1)
		}

		os.RemoveAll(filepath.Join(targetDir, ".git"))
		os.RemoveAll(filepath.Join(targetDir, "cmd", "clonator"))

		makefile := filepath.Join(targetDir, "Makefile")
		if content, err := os.ReadFile(makefile); err == nil {
			strContent := string(content)
			strContent = strings.Replace(strContent, " clonator", "", 1)
			clonatorHelp := "\t@echo \"  clonator - Start the interactive Clonator CLI to generate a new service\"\n\t@echo \"\"\n"
			strContent = strings.Replace(strContent, clonatorHelp, "", 1)
			clonatorTarget := `clonator:
	@if [ "$(FILE)" != "" ]; then \
		go run cmd/clonator/*.go --from-file=$(FILE); \
	else \
		go run cmd/clonator/*.go; \
	fi`
			strContent = strings.Replace(strContent, clonatorTarget, "", 1)
			os.WriteFile(makefile, []byte(strContent), 0644)
		} else {
			fmt.Printf("⚠️  Warning: Could not read or clean Makefile: %v\n", err)
		}

		os.Remove(filepath.Join(targetDir, "README.md"))

		pluralUpper := strings.ToUpper(entityPlural)
		pluralLower := strings.ToLower(entityPlural)
		pluralTitle := strings.ToUpper(pluralLower[:1]) + pluralLower[1:]

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

		fmt.Println("🧹 Formatting code...")
		fmtCmd := exec.Command("go", "fmt", "./...")
		fmtCmd.Dir = targetDir
		if err := fmtCmd.Run(); err != nil {
			fmt.Printf("❌ Error formatting code: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("🛠️ Tidying modules...")
		tidyCmd := exec.Command("go", "mod", "tidy")
		tidyCmd.Dir = targetDir
		if err := tidyCmd.Run(); err != nil {
			fmt.Printf("❌ Error tidying modules: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("⚙️ Generating Protobuf files...")
		protoCmd := exec.Command("make", "proto")
		protoCmd.Dir = targetDir
		if err := protoCmd.Run(); err != nil {
			fmt.Printf("❌ Error generating Protobuf files: %v\n", err)
			os.Exit(1)
		}

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

	cloneCmd.Run(cloneCmd, []string{serviceName})
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
