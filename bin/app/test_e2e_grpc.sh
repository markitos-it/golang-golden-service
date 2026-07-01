#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'
ENVIRONMENT_FILE="bin/shared/environment.sh"
source $ENVIRONMENT_FILE


SERVER="localhost:30000"
SERVICE="golden.Goldenservice"

echo "🚀 Iniciando Test E2E gRPC para $SERVICE en $SERVER..."

UPLOADS_DIR=${GOLDEN_UPLOADS_BASEDIR:-./uploads}
mkdir -p $UPLOADS_DIR

# Garantizar que el directorio por defecto exista en caso de que el server no tenga las variables de entorno
mkdir -p /tmp/goldens

# Check dependencies
if ! command -v grpcurl &> /dev/null; then
    echo "❌ grpcurl no está instalado. Ejecuta 'make support-install-grpc-tools' primero."
    exit 1
fi

if ! command -v jq &> /dev/null; then
    echo "❌ jq no está instalado. Por favor instálalo (ej. sudo apt install jq / brew install jq)."
    exit 1
fi

echo "--------------------------------------------------"
echo "1️⃣  Creando un nuevo Golden..."
CREATE_PAYLOAD='{
  "name": "GoldenTest",
  "content": "GoldenContent",
  "poster_data": ""__CUSTOM_E2E_FIELDS_CREATE__
}'
PLACEHOLDER_CREATE="__CUSTOM_E2E_""FIELDS_CREATE__"
echo "🔍 Placeholders encontrados: $PLACEHOLDER_CREATE"
CREATE_PAYLOAD="${CREATE_PAYLOAD//$PLACEHOLDER_CREATE/}"
echo "🔍 Payload enviado: $CREATE_PAYLOAD"
CREATE_RESP=$(grpcurl -plaintext -d "$CREATE_PAYLOAD" $SERVER $SERVICE/CreateGolden)

echo "$CREATE_RESP"
GOLDEN_ID=$(echo "$CREATE_RESP" | jq -r '.id')

if [ -z "$GOLDEN_ID" ] || [ "$GOLDEN_ID" == "null" ]; then
    echo "❌ Error: No se pudo extraer el ID de la respuesta."
    exit 1
fi
echo "✅ Creado exitosamente con ID: $GOLDEN_ID"
echo "--------------------------------------------------"

echo "2️⃣  Obteniendo Golden por ID..."
grpcurl -plaintext -d "{\"id\": \"$GOLDEN_ID\"}" $SERVER $SERVICE/GetGolden
echo "✅ Get exitoso."
echo "--------------------------------------------------"

echo "3️⃣  Actualizando Golden..."
UPDATE_PAYLOAD="{
  \"id\": \"$GOLDEN_ID\",
  \"name\": \"GoldenTestUpdated\",
  \"content\": \"GoldenContentUpdated\",
  \"poster_data\": \"\"__CUSTOM_E2E_FIELDS_UPDATE__
}"
PLACEHOLDER_UPDATE="__CUSTOM_E2E_""FIELDS_UPDATE__"
UPDATE_PAYLOAD="${UPDATE_PAYLOAD//$PLACEHOLDER_UPDATE/}"

echo "🔍 Payload enviado: $UPDATE_PAYLOAD"
grpcurl -plaintext -d "$UPDATE_PAYLOAD" $SERVER $SERVICE/UpdateGolden
echo "✅ Update exitoso."
echo "--------------------------------------------------"

echo "4️⃣  Listando Goldens..."
grpcurl -plaintext $SERVER $SERVICE/ListGoldens
echo "✅ List exitoso."
echo "--------------------------------------------------"

echo "5️⃣  Borrando Golden..."
grpcurl -plaintext -d "{\"id\": \"$GOLDEN_ID\"}" $SERVER $SERVICE/DeleteGolden
echo "✅ Delete exitoso."
echo "--------------------------------------------------"

echo "6️⃣  Verificando que fue borrado..."
set +e
GET_DELETED_RESP=$(grpcurl -plaintext -d "{\"id\": \"$GOLDEN_ID\"}" $SERVER $SERVICE/GetGolden 2>&1)
set -e

if echo "$GET_DELETED_RESP" | grep -i -q "not found\|notfound"; then
    echo "✅ Verificación exitosa: El registro ya no existe."
else
    echo "❌ Advertencia: El registro podría seguir existiendo o ocurrió un error inesperado:"
    echo "$GET_DELETED_RESP"
    exit 1
fi

echo "--------------------------------------------------"
echo "🎉 ¡Todos los tests E2E pasaron exitosamente!"