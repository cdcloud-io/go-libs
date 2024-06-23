package config

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func Load() Config {
	clearTerminal()

	fmt.Println("🟩 STARTUP INFO: Loading configuration from file")
	startTime := time.Now()

	var config Config

	data, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		fmt.Println("🟥 STARTUP ERROR: Could not load config file from ./config/config.yaml. Exiting...")
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println("🟩 STARTUP INFO: Successfully read the config file")

	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("🟥 STARTUP ERROR: Could not unmarshal config data: %v", err)
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println("🟩 STARTUP INFO: Successfully unmarshaled the config data")

	// Replace placeholders with environment variables
	ReplaceEnvVars(&config)

	fmt.Printf("🟩 STARTUP INFO: configs loaded in: %v \n", time.Since(startTime))

	return config
}

func ReplaceEnvVars(config *Config) {
	v := reflect.ValueOf(config).Elem()
	replaceEnvVars(v)
}

func replaceEnvVars(v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if field.Kind() == reflect.Struct {
			replaceEnvVars(field)
		} else if field.Kind() == reflect.String {
			fieldValue := field.String()
			if strings.HasPrefix(fieldValue, "${") && strings.HasSuffix(fieldValue, "}") {

				envVarName := fieldValue[2 : len(fieldValue)-1]
				envVarValue := os.Getenv(envVarName)

				if envVarValue != "" {
					field.SetString(envVarValue)
				} else {
					fmt.Printf("🟥 STARTUP ERROR: Missing environment variable: %s. Exit 1\n\n", envVarName)
					os.Exit(1)
				}
			}
		}
	}
}

func clearTerminal() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear") // Systems
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls") // Windows
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("🟨 STARTUP WARN: clear terminal FAILED, unsupported OS")
	}
}

func CheckConfig(value string, errorFlag *bool) string {
	if value == "" {
		*errorFlag = true
		return "missing"
	}
	if value == "missing" {
		*errorFlag = true
		return "missing"
	}
	return value
}
