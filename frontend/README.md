# Frontend - Guia de Ejecucion

## Requisitos

- Node.js
- npm
- Docker Desktop

## Variables de entorno

Crear o verificar el archivo `frontend/.env` con los siguientes valores:

```dotenv
VITE_API_BASE_URL=/api
VITE_PROXY_TARGET=http://localhost:80
```

## Ejecucion local del frontend

Desde la carpeta raiz del proyecto:

```bash
cd frontend
npm install
npm run dev
```

## Ejecucion completa con Docker (desde la raiz)

```bash
docker compose down
docker compose up --build
```

## URLs de acceso

- Frontend: http://localhost:5173
- Backend por nginx: http://localhost:80/api

## Credenciales admin

Tomar las credenciales desde el archivo `.env.admin`:

- DEFAULT_ADMIN_EMAIL
- DEFAULT_ADMIN_PASSWORD

## Proxy de Vite y CORS

El frontend realiza peticiones a rutas relativas `/api`. Vite intercepta esas rutas y las redirige a `VITE_PROXY_TARGET` (nginx/backend), lo que evita errores de CORS en desarrollo.

## Prueba rapida de login (curl)

```bash
curl -X POST http://localhost:80/api/auth/sign-in -H "Content-Type: application/json" -d '{"email":"admin@seneprojects.com","password":"seneprojects123"}'
```

## Problemas comunes

- Failed to fetch:
   Revisar que `docker compose` este arriba y que `frontend/.env` use `VITE_API_BASE_URL=/api`.
- Cambios de `.env` no reflejados:
   Reiniciar `npm run dev` o hacer rebuild con Docker.
- Base de datos con estado viejo:
   Ejecutar `docker compose down -v`.
