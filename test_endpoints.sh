#!/bin/bash

# Script para probar los endpoints de Wishlist
# Asegúrate de que los servidores estén corriendo:
# - AuthGo en http://localhost:3000
# - CatalogGo en http://localhost:3002
# - CartGo en http://localhost:3003
# - WishlistGo en http://localhost:3005

BASE_URL="http://localhost:3005"
AUTH_URL="http://localhost:3000"

echo "Probando endpoints de Wishlist..."
echo ""
echo "Servicios requeridos:"
echo "   - AuthGo: $AUTH_URL"
echo "   - CatalogGo: http://localhost:3002"
echo "   - CartGo: http://localhost:3003"
echo "   - WishlistGo: $BASE_URL"
echo ""

# Verificar que AuthGo esté corriendo
echo "Verificando que AuthGo este corriendo..."
if ! curl -s -f "$AUTH_URL/health" > /dev/null 2>&1 && ! curl -s -f "$AUTH_URL" > /dev/null 2>&1; then
  echo "WARNING: AuthGo no parece estar corriendo en $AUTH_URL"
  echo "   Intentando continuar de todas formas..."
fi
echo ""

# Intentar crear usuario si no existe (ignorar errores si ya existe)
echo "Creando usuario de prueba (si no existe)..."
SIGNUP_RESPONSE=$(curl -s -X POST "$AUTH_URL/users/signup" \
  -H "Content-Type: application/json" \
  -d '{"login": "testuser", "password": "test123", "name": "Test User"}' 2>/dev/null)
# Ignorar si ya existe
echo ""

# Obtener token JWT
echo "Obteniendo token JWT..."
TOKEN_RESPONSE=$(curl -s -X POST "$AUTH_URL/users/signin" \
  -H "Content-Type: application/json" \
  -d '{"login": "testuser", "password": "test123"}' 2>/dev/null)

if [ $? -ne 0 ] || [ -z "$TOKEN_RESPONSE" ]; then
  echo "ERROR: No se pudo obtener el token."
  echo "   Respuesta del servidor:"
  echo "$TOKEN_RESPONSE"
  echo ""
  echo "   Asegurate de que:"
  echo "   1. AuthGo este corriendo en http://localhost:3000"
  echo "   2. El usuario 'testuser' exista (se intento crear automaticamente)"
  exit 1
fi

TOKEN=$(echo $TOKEN_RESPONSE | jq -r '.token' 2>/dev/null)

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "ERROR: No se pudo extraer el token de la respuesta:"
  echo "$TOKEN_RESPONSE"
  exit 1
fi

echo "Token obtenido: ${TOKEN:0:20}..."
echo ""

# Obtener o crear artículos de prueba
echo "Preparando articulos de prueba..."
ART1=""
ART2=""

# Intentar obtener artículos existentes
ARTICLES_JSON=$(curl -s -X GET "http://localhost:3002/articles" \
  -H "Authorization: Bearer $TOKEN" 2>/dev/null)

if [ $? -eq 0 ] && [ -n "$ARTICLES_JSON" ] && [ "$ARTICLES_JSON" != "null" ]; then
  ART1=$(echo "$ARTICLES_JSON" | jq -r 'if type == "array" then .[0]._id // empty else ._id // empty end' 2>/dev/null)
  ART2=$(echo "$ARTICLES_JSON" | jq -r 'if type == "array" then .[1]._id // empty else empty end' 2>/dev/null)
fi

# Si no hay artículos, crear algunos de prueba
if [ -z "$ART1" ] || [ "$ART1" = "null" ] || [ "$ART1" = "" ]; then
  echo "Creando articulos de prueba en CatalogGo..."
  
  # Crear artículo 1
  ART1_RESPONSE=$(curl -s -X POST "http://localhost:3002/articles" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Producto Test 1",
      "description": "Descripción del producto test 1",
      "image": "https://example.com/image1.jpg",
      "price": 99.99,
      "stock": 10
    }' 2>/dev/null)
  
  ART1=$(echo "$ART1_RESPONSE" | jq -r '._id // empty' 2>/dev/null)
  sleep 1
fi

if [ -z "$ART2" ] || [ "$ART2" = "null" ] || [ "$ART2" = "" ]; then
  # Crear artículo 2
  ART2_RESPONSE=$(curl -s -X POST "http://localhost:3002/articles" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Producto Test 2",
      "description": "Descripción del producto test 2",
      "image": "https://example.com/image2.jpg",
      "price": 149.99,
      "stock": 5
    }' 2>/dev/null)
  
  ART2=$(echo "$ART2_RESPONSE" | jq -r '._id // empty' 2>/dev/null)
  sleep 1
fi

if [ -z "$ART1" ] || [ "$ART1" = "null" ] || [ "$ART1" = "" ]; then
  echo "ERROR: No se pudieron obtener o crear articulos del catalogo"
  echo "   Asegurate de que CatalogGo este corriendo en http://localhost:3002"
  exit 1
fi

echo "Articulos a usar: ART1=$ART1, ART2=$ART2"
echo ""
echo "---"
echo ""

# 1. GET /v1/wishlist (debería retornar wishlist vacía o existente)
echo "1. GET /v1/wishlist"
curl -X GET "$BASE_URL/v1/wishlist" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  2>/dev/null | jq . 2>/dev/null || echo "Error o respuesta no JSON"
echo ""
echo "---"
echo ""

# 2. POST /v1/wishlist/article (agregar artículo)
echo "2. POST /v1/wishlist/article"
curl -X POST "$BASE_URL/v1/wishlist/article" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"article_id\": \"$ART1\", \"notes\": \"Regalo para mamá\"}" \
  -w "\nHTTP Status: %{http_code}\n" \
  2>/dev/null | jq . 2>/dev/null || echo "Error o respuesta no JSON"
echo ""
echo "---"
echo ""

# 3. GET /v1/wishlist (debería mostrar el artículo agregado con info del catálogo)
echo "3. GET /v1/wishlist (despues de agregar)"
curl -X GET "$BASE_URL/v1/wishlist" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  2>/dev/null | jq . 2>/dev/null || echo "Error o respuesta no JSON"
echo ""
echo "---"
echo ""

# 4. POST /v1/wishlist/article (agregar otro artículo)
echo "4. POST /v1/wishlist/article (segundo articulo)"
curl -X POST "$BASE_URL/v1/wishlist/article" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"article_id\": \"$ART2\"}" \
  -w "\nHTTP Status: %{http_code}\n" \
  2>/dev/null | jq . 2>/dev/null || echo "Error o respuesta no JSON"
echo ""
echo "---"
echo ""

# 5. POST /v1/wishlist/article (intentar agregar duplicado - debería retornar 409)
echo "5. POST /v1/wishlist/article (intentar duplicado - esperado 409)"
curl -X POST "$BASE_URL/v1/wishlist/article" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"article_id\": \"$ART1\"}" \
  -w "\nHTTP Status: %{http_code}\n" \
  2>/dev/null | jq . 2>/dev/null || echo "Error o respuesta no JSON"
echo ""
echo "---"
echo ""

# 6. GET /v1/wishlist (debería mostrar 2 artículos con info enriquecida)
echo "6. GET /v1/wishlist (con 2 articulos)"
curl -X GET "$BASE_URL/v1/wishlist" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  2>/dev/null | jq . 2>/dev/null || echo "Error o respuesta no JSON"
echo ""
echo "---"
echo ""

# 7. DELETE /v1/wishlist/article/:article_id (eliminar artículo)
echo "7. DELETE /v1/wishlist/article/$ART1"
curl -X DELETE "$BASE_URL/v1/wishlist/article/$ART1" \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nHTTP Status: %{http_code}\n" \
  2>/dev/null
echo ""
echo "---"
echo ""

# 8. GET /v1/wishlist (debería mostrar solo 1 artículo)
echo "8. GET /v1/wishlist (despues de eliminar)"
curl -X GET "$BASE_URL/v1/wishlist" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  2>/dev/null | jq . 2>/dev/null || echo "Error o respuesta no JSON"
echo ""
echo "---"
echo ""

# 9. POST /v1/wishlist/articles/:article_id/cart (mover a carrito)
echo "9. POST /v1/wishlist/articles/$ART2/cart"
curl -X POST "$BASE_URL/v1/wishlist/articles/$ART2/cart" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  2>/dev/null | jq . 2>/dev/null || echo "Error o respuesta no JSON"
echo ""
echo "---"
echo ""

# 10. GET /v1/wishlist (debería estar vacía)
echo "10. GET /v1/wishlist (despues de mover a carrito)"
curl -X GET "$BASE_URL/v1/wishlist" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  2>/dev/null | jq . 2>/dev/null || echo "Error o respuesta no JSON"
echo ""
echo "Pruebas completadas"

