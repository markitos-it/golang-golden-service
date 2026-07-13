package configuration

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

//[.'.]:> 📦 SISTEMA DE CONFIGURACIÓN
//[.'.]:> ========================
//[.'.]:> Este módulo maneja la configuración de nuestra aplicación siguiendo este flujo:
//[.'.]:>
//[.'.]:> 1️⃣ ARCHIVO CONFIG: Busca primero un archivo config.yaml en el directorio especificado
//[.'.]:>    Si lo encuentra, carga todas sus variables como configuración base
//[.'.]:>
//[.'.]:> 2️⃣ VARIABLES DE ENTORNO: Después de cargar el archivo (o si no existe):
//[.'.]:>    - Comprueba si existen variables de entorno con los mismos nombres
//[.'.]:>    - Las variables de entorno tienen PRIORIDAD y sobrescriben los valores del archivo
//[.'.]:>
//[.'.]:> 3️⃣ VALORES POR DEFECTO: Como último recurso, si algún valor sigue vacío,
//[.'.]:>    se aplican valores predeterminados para garantizar que la app pueda funcionar

// [.'.]:> 🧩 Estructura que contiene toda la configuración de la aplicación
// [.'.]:> Cada campo se mapea a una variable de entorno o valor en config.yaml del mismo nombre
type GoldenConfiguration struct {
	DatabaseDsn       string `mapstructure:"DATABASE_DSN"`
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	BaseDir           string `mapstructure:"GOLDEN_UPLOADS_BASEDIR"`
}

// [.'.]:> 🔄 Función principal que carga toda la configuración
// [.'.]:> Recibe la ruta donde buscar el archivo config.yaml y devuelve la configuración completa
// [.'.]:> Si hay algún error durante la carga, lo devuelve para que la aplicación pueda manejarlo
func LoadConfiguration(configFilesPath string) (config GoldenConfiguration, err error) {
	viper.New()
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(configFilesPath)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.BindEnv("DATABASE_DSN")
	viper.BindEnv("GRPC_SERVER_ADDRESS")
	viper.BindEnv("GOLDEN_UPLOADS_BASEDIR")
	viper.AutomaticEnv()

	if err := loadConfigFile(); err != nil {
		return config, err
	}

	overrideWithEnvVars()

	err = viper.Unmarshal(&config)
	if err == nil {
		fmt.Println("['.']:> ✨ Configuración cargada correctamente ✨")
		fmt.Println("['.']:> ----------------------------------------")
		fmt.Printf("['.']:> 🚀 gRPC Server: %s\n", config.GRPCServerAddress)
		fmt.Printf("['.']:> 📁 Base Directory: %s\n", config.BaseDir)
		fmt.Println("['.']:> ----------------------------------------")
		applyFallbackEnvVars(&config)
	}

	return config, err
}

// [.'.]:> 📄 Intenta cargar el archivo de configuración config.yaml
// [.'.]:> Si el archivo no existe, lo maneja elegantemente y permite continuar
// [.'.]:> usando solo variables de entorno
func loadConfigFile() error {
	// Intenta leer config.yaml
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
		fmt.Println("['.']:> 📋 Archivo 'config.yaml' no encontrado, buscando otras fuentes.")
	} else {
		fmt.Println("['.']:> 📋 Archivo 'config.yaml' cargado correctamente.")
	}

	// Intenta fusionar con .env si existe
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	if err := viper.MergeInConfig(); err == nil {
		fmt.Println("['.']:> 📋 Archivo '.env' cargado y fusionado correctamente.")
	}

	return nil
}

// [.'.]:> 🔀 Sobrescribe valores del archivo con variables de entorno
// [.'.]:> Esta es la clave para que las variables de entorno tengan prioridad
// [.'.]:> sobre el archivo de configuración
func overrideWithEnvVars() {
	dsnEnv := os.Getenv("DATABASE_DSN")
	if dsnEnv != "" && viper.GetString("DATABASE_DSN") == "" {
		viper.Set("DATABASE_DSN", dsnEnv)
	}

	grpcEnv := os.Getenv("GRPC_SERVER_ADDRESS")
	if grpcEnv != "" && viper.GetString("GRPC_SERVER_ADDRESS") == "" {
		viper.Set("GRPC_SERVER_ADDRESS", grpcEnv)
	}

	baseDirEnv := os.Getenv("GOLDEN_UPLOADS_BASEDIR")
	if baseDirEnv != "" && viper.GetString("GOLDEN_UPLOADS_BASEDIR") == "" {
		viper.Set("GOLDEN_UPLOADS_BASEDIR", baseDirEnv)
	}
}

// [.'.]:> 🔒 Aplica valores de respaldo directamente desde variables de entorno
// [.'.]:> como última red de seguridad para los campos que aún estén vacíos
// [.'.]:> después de procesar el archivo y las variables a través de viper
func applyFallbackEnvVars(config *GoldenConfiguration) {
	if config.DatabaseDsn == "" {
		config.DatabaseDsn = os.Getenv("DATABASE_DSN")
	}
	if config.GRPCServerAddress == "" {
		config.GRPCServerAddress = os.Getenv("GRPC_SERVER_ADDRESS")
	}
	if config.BaseDir == "" {
		config.BaseDir = os.Getenv("GOLDEN_UPLOADS_BASEDIR")
		if config.BaseDir == "" {
			config.BaseDir = "/tmp/goldens"
		}
	}
}
