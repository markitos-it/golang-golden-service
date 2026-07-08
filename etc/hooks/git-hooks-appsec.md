# Guía de Instalación y Uso de Git Hooks Nativos para Go

Esta guía detalla cómo instalar los binarios requeridos y configurar los interceptores nativos de Git (*Hooks*) en tu entorno local para automatizar las validaciones de calidad de código y seguridad (AppSec) en proyectos de Go.

---

## 1. Instalación de Herramientas (Aprovisionamiento Local)

Para garantizar la máxima compatibilidad con la versión del compilador instalada en tu sistema operativo, instalaremos y compilaremos las herramientas directamente a través del motor de Go.

Ejecuta los siguientes comandos en tu terminal:

### 1.1. Detección de Secretos (Gitleaks)

Instala el motor encargado de escanear el área de preparación en busca de API tokens, llaves privadas o credenciales expuestas:

```bash
go install github.com/zricethezav/gitleaks/v8@latest

```

### 1.2. Calidad de Código y Linters (Golangci-lint)

Descarga y compila la última versión del agregador de linters directamente con tu motor de Go para evitar conflictos con formatos de metadatos antiguos:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

```

### 1.3. Análisis Estático de Seguridad (GoSec)

Instala el analizador SAST que inspecciona el Árbol Sintáctico Abstracto (AST) de tu código en busca de fallos lógicos explotables:

```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest

```

### 1.4. Análisis de la Cadena de Suministro (Govulncheck)

Instala la herramienta oficial de Go para detectar librerías de terceros con vulnerabilidades conocidas que realmente afecten a tus grafos de ejecución:

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest

```

### 1.4. Configuración del PATH en el Sistema

Para que tu terminal localice los nuevos binarios que acabas de compilar, asegúrate de tener el directorio `bin` de Go indexado en tu variable de entorno. Añade esta línea al final de tu archivo `~/.bashrc` (o el archivo de configuración de tu shell):

```bash
export PATH=$(go env GOPATH)/bin:$PATH

```

Luego, recarga la configuración ejecutando `source ~/.bashrc`.


### 1.5. Autenticación de Auditoría Externa (Snyk)

El interceptor de fase profunda (pre-push) utiliza el motor de Snyk para complementar el análisis estático y de dependencias. Para que estos comandos se ejecuten de forma no interactiva sin interrumpir la terminal, es mandatorio exportar tu token de autenticación.

Obtén tu token desde el panel web de Snyk (Account Settings > Auth Token).

Añade la variable de entorno a tu archivo de configuración del shell (~/.bashrc o ~/.zshrc):

```bash
export SNYK_TOKEN="tu_token_secreto_aqui"
```

Recarga tu configuración actual ejecutando 

```bash
source ~/.bashrc.
```

Nota de Seguridad: Nunca dejes este token expuesto directamente en el script del hook dentro de un repositorio compartido. La lectura debe ser siempre delegada al entorno local del sistema operativo del desarrollador.



---

## 2. Configuración de Calidad (`.golangci.yml`)

Crea el archivo `.golangci.yml` en la raíz de tu proyecto para gobernar las reglas de calidad y formato. Esto optimiza el rendimiento del pre-commit y silencia advertencias redundantes:

```yaml
run:
  concurrency: 4
  timeout: 5m
  issues-exit-code: 1
  tests: false

output:
  formats:
    - format: colored-line-number
  print-issued-lines: true
  print-linter-name: true

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true
  goimports:
    local-prefixes: mdk-app

linters:
  disable-all: true
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - typecheck
    - unused
    - gocritic
    - gofmt
    - goimports

```

---

## 3. Implementación de los Scripts Locales

Git gestiona los hooks locales dentro del directorio oculto `.git/hooks/` de tu repositorio. Debes crear dos archivos sin extensión en esa ruta.

### 3.1. Interceptor de Fase Rápida (`.git/hooks/pre-commit`)

Este script se ejecuta al realizar un `git commit`. Valida secretos en memoria intermedia (*staging area*) y que el formato cumpla con las convenciones definidas.

```bash
#!/usr/bin/env bash

CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${CYAN}==> Iniciando validacion de calidad local...${NC}"

STATUS_GITLEAKS="${YELLOW}Omitido${NC}"
STATUS_GOLANGCI="${YELLOW}Omitido${NC}"
STATUS_GOFMT="${YELLOW}Omitido${NC}"
FAILURES=0

# 1. Intercepción de Secretos en Staging
if command -v gitleaks >/dev/null 2>&1; then
    echo -e "\n--> Escaneando credenciales en memoria intermedia..."
    if gitleaks protect -v --staged; then
        STATUS_GITLEAKS="${GREEN}PASADO${NC}"
    else
        STATUS_GITLEAKS="${RED}FALLADO${NC}"
        FAILURES=1
    fi
else
    echo -e "${YELLOW}Advertencia: Gitleaks no encontrado en el sistema.${NC}"
fi

# 2. Control de Calidad y Convenciones Go
if command -v golangci-lint >/dev/null 2>&1; then
    echo -e "\n--> Verificando calidad y formato de codigo..."
    if golangci-lint run; then
        STATUS_GOLANGCI="${GREEN}PASADO${NC}"
    else
        STATUS_GOLANGCI="${RED}FALLADO${NC}"
        FAILURES=1
    fi
else
    echo -e "${YELLOW}Advertencia: golangci-lint no encontrado en el sistema.${NC}"
fi

# 3. Formateo de Código Go
echo -e "\n--> Comprobando formato de código Go..."
GO_FILES=$(find . -name "*.go" -not -path "./vendor/*")
UNFORMATTED_FILES=$(gofmt -l $GO_FILES)

if [ -n "$UNFORMATTED_FILES" ]; then
    echo -e "${RED}Error: Los siguientes ficheros no están formateados:${NC}"
    echo "$UNFORMATTED_FILES"
    echo -e "${YELLOW}Ejecutando 'gofmt -w' para corregirlos...${NC}"
    gofmt -w $UNFORMATTED_FILES
    STATUS_GOFMT="${RED}FALLADO (Corregido)${NC}"
    FAILURES=1
else
    STATUS_GOFMT="${GREEN}PASADO${NC}"
fi

echo ""
echo -e "${CYAN}==============================================${NC}"
echo -e "${CYAN}           RESUMEN PRE-COMMIT                 ${NC}"
echo -e "${CYAN}==============================================${NC}"
echo -e "HERRAMIENTA             RESULTADO"
echo -e "----------------------------------------------"
echo -e "Gitleaks                ${STATUS_GITLEAKS}"
echo -e "Golangci-lint           ${STATUS_GOLANGCI}"
echo -e "Go Format               ${STATUS_GOFMT}"
echo -e "${CYAN}==============================================${NC}"

if [ $FAILURES -ne 0 ]; then
    echo -e "${RED}==> Bloqueando commit. Corrige los errores detallados arriba.${NC}"
    exit 1
else
    echo -e "${GREEN}==> Commit autorizado. Arbol de trabajo limpio.${NC}"
    exit 0
fi

```

### 3.2. Interceptor de Fase Profunda (`.git/hooks/pre-push`)

Este script se ejecuta únicamente cuando intentas enviar tus commits al repositorio remoto (`git push`). Aplica análisis estático avanzado (SAST) y análisis de dependencias (SCA) filtrando solo las amenazas de severidad alta.

```bash
#!/usr/bin/env bash

CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${CYAN}==> Iniciando auditoria de seguridad profunda...${NC}"

STATUS_GOVULNCHECK="${YELLOW}Omitido${NC}"
STATUS_SNYK_CODE="${YELLOW}Omitido${NC}"
STATUS_SNYK_SCA="${YELLOW}Omitido${NC}"
FAILURES=0

# 2. Análisis SCA Nativo con Contexto de Ejecucion
if command -v govulncheck >/dev/null 2>&1; then
    echo -e "\n--> Consultando Go Vulnerability Database..."
    if GOVULN_OUT=$(govulncheck ./... 2>&1); then
        STATUS_GOVULNCHECK="${GREEN}PASADO${NC}"
    else
        STATUS_GOVULNCHECK="${RED}FALLADO${NC}"
        FAILURES=1
        echo -e "${RED}$GOVULN_OUT${NC}"
    fi
fi

# 3. Análisis SAST Semantico Avanzado
if command -v snyk >/dev/null 2>&1; then
    echo -e "\n--> Ejecutando Snyk Code Engine..."
    if SNYK_CODE_OUT=$(snyk code test --severity-threshold=high 2>&1); then
        STATUS_SNYK_CODE="${GREEN}PASADO${NC}"
    else
        STATUS_SNYK_CODE="${RED}FALLADO${NC}"
        FAILURES=1
        echo -e "${RED}$SNYK_CODE_OUT${NC}"
    fi
fi

# 4. Análisis SCA Avanzado Multilenguaje
if command -v snyk >/dev/null 2>&1; then
    echo -e "\n--> Ejecutando Snyk Open Source..."
    if SNYK_SCA_OUT=$(snyk test --severity-threshold=high 2>&1); then
        STATUS_SNYK_SCA="${GREEN}PASADO${NC}"
    else
        STATUS_SNYK_SCA="${RED}FALLADO${NC}"
        FAILURES=1
        echo -e "${RED}$SNYK_SCA_OUT${NC}"
    fi
fi

echo ""
echo -e "${CYAN}==============================================${NC}"
echo -e "${CYAN}           RESUMEN PRE-PUSH                   ${NC}"
echo -e "${CYAN}==============================================${NC}"
echo -e "HERRAMIENTA             RESULTADO"
echo -e "----------------------------------------------"
echo -e "Govulncheck             ${STATUS_GOVULNCHECK}"
echo -e "Snyk Code               ${STATUS_SNYK_CODE}"
echo -e "Snyk Open Source        ${STATUS_SNYK_SCA}"
echo -e "${CYAN}==============================================${NC}"

if [ $FAILURES -ne 0 ]; then
    echo -e "${RED}==> Amenazas criticas detectadas. Empuje al remoto bloqueado.${NC}"
    exit 1
else
    echo -e "${GREEN}==> Auditoria perfecta. Transmision autorizada.${NC}"
    exit 0
fi

```

---

## 4. Inicialización y Activación

Para poner los hooks en funcionamiento dentro de tu repositorio local, ejecuta los siguientes comandos desde la raíz de tu proyecto:

### Paso 1: Otorgar permisos de ejecución

Los scripts de Git deben tener los privilegios adecuados en el sistema operativo para ser invocados durante el ciclo de vida:

```bash
chmod +x .git/hooks/pre-commit .git/hooks/pre-push

```

### Paso 2: Limpieza preliminar de la caché

Asegúrate de purgar los datos antiguos remanentes de compilaciones previas para forzar un análisis limpio desde cero:

```bash
go clean -cache
golangci-lint cache clean

```

### Paso 3: Verificación práctica

Realiza una pequeña modificación en tu código fuente, añade los cambios al área de preparación mediante `git add .` e intenta hacer una confirmación con `git commit` para comprobar el correcto funcionamiento y la renderización visual del panel tabular.