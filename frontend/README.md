# Frontend - Guia de ejecucion

## Requisitos

- Node.js
- npm
- Docker Desktop

## Variables de entorno

Crear o verificar el archivo frontend/.env con estos valores:

```dotenv
VITE_API_BASE_URL=/api
VITE_PROXY_TARGET=http://localhost:80
```

## Ejecucion local del frontend

Desde la raiz del repositorio:

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

Tomar las credenciales desde .env.admin:

- DEFAULT_ADMIN_EMAIL
- DEFAULT_ADMIN_PASSWORD

## Proxy de Vite y CORS

El frontend envía peticiones a rutas relativas /api. El servidor de desarrollo de Vite intercepta esas rutas y las redirige a VITE_PROXY_TARGET, evitando errores de CORS en desarrollo.

## Prueba rapida de login (curl)

```bash
curl -X POST http://localhost:80/api/auth/sign-in -H "Content-Type: application/json" -d '{"email":"admin@seneprojects.com","password":"seneprojects123"}'
```

## Problemas comunes

- Failed to fetch:
  Verificar que docker compose este arriba y que frontend/.env use VITE_API_BASE_URL=/api.
- Cambios en .env no reflejados:
  Reiniciar npm run dev o reconstruir con Docker.
- Base de datos con estado viejo:
  Ejecutar docker compose down -v.

## Alcance actual del frontend

- Login por rol.
- Panel de administracion.
- Panel de profesor.
- Panel de monitor/asistente.
- Creacion de usuarios, periodos, cursos/proyectos y vinculaciones.
- Registro y gestion basica de tareas.
- Generacion y descarga de reportes PDF semanales.

## Deuda tecnica conocida

- Reportes con IA externa real aun no implementados.
- Adjuntos de tareas aun no implementados.
- Persistencia de archivos en Cloud Storage pendiente para despliegue GCP.
- Frontend minimo orientado a demostracion, no diseno final de produccion.
- Filtros avanzados y busquedas aun no implementados.
- Validaciones frontend completas pendientes; actualmente se muestran errores reales del backend.
- Validaciones de monitores y demás
