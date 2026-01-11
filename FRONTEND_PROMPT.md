# Prompt para Desarrollo del Frontend - Phoenix Alliance Dashboard

## Contexto del Proyecto

Necesito crear un dashboard web moderno y funcional para **Phoenix Alliance**, una aplicación de seguimiento de entrenamiento físico. El backend ya está completamente desarrollado en Go y está listo para recibir requests. El frontend debe ser una Single Page Application (SPA) que permita a los usuarios gestionar sus ejercicios, entrenamientos y sets de entrenamiento.

## Base URL del Backend

```
http://localhost:8080
```

## Autenticación

El backend usa **JWT (JSON Web Tokens)** para autenticación. Todas las rutas protegidas requieren el header:

```
Authorization: Bearer <token>
```

### Endpoints de Autenticación

#### POST /signup
Registra un nuevo usuario.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "email": "user@example.com",
  "created_at": "2024-01-01T00:00:00Z"
}
```

**Errores:**
- `400 Bad Request`: Email y password son requeridos, password mínimo 8 caracteres
- `409 Conflict`: Usuario con ese email ya existe

#### POST /login
Autentica un usuario y retorna un JWT token.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

**Errores:**
- `400 Bad Request`: Email y password son requeridos
- `401 Unauthorized`: Credenciales inválidas

## Endpoints Protegidos (Requieren JWT)

Todas las siguientes rutas requieren el header `Authorization: Bearer <token>`

### Ejercicios

#### GET /exercises
Obtiene todos los ejercicios del usuario autenticado.

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "name": "Bench Press",
    "created_at": "2024-01-01T00:00:00Z"
  },
  {
    "id": 2,
    "user_id": 1,
    "name": "Squat",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

#### POST /exercises
Crea un nuevo ejercicio.

**Request:**
```json
{
  "name": "Deadlift"
}
```

**Response (201 Created):**
```json
{
  "id": 3,
  "user_id": 1,
  "name": "Deadlift",
  "created_at": "2024-01-01T00:00:00Z"
}
```

**Errores:**
- `400 Bad Request`: Nombre del ejercicio es requerido

#### GET /exercises/{id}/history
Obtiene el historial completo de sets para un ejercicio específico con métricas agregadas.

**Response (200 OK):**
```json
{
  "exercise_id": 1,
  "exercise_name": "Bench Press",
  "sets": [
    {
      "id": 1,
      "workout_id": 1,
      "exercise_id": 1,
      "weight": 60.0,
      "reps": 10,
      "rest_seconds": 120,
      "rpe": 7,
      "notes": "Felt good",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "metrics": {
    "total_sets": 15,
    "total_volume": 9000.0,
    "max_weight": 100.0,
    "max_reps": 12,
    "average_weight": 75.5,
    "average_reps": 8.2,
    "average_rest": 120.0,
    "average_rpe": 7.5,
    "first_recorded_at": "2024-01-01T00:00:00Z",
    "last_recorded_at": "2024-01-15T00:00:00Z"
  }
}
```

#### GET /exercises/{id}/progress?range=week|month|year
Obtiene datos de progreso para un ejercicio en un rango de tiempo específico.

**Query Parameters:**
- `range`: `week`, `month`, o `year` (default: `month`)

**Response (200 OK):**
```json
{
  "exercise_id": 1,
  "exercise_name": "Bench Press",
  "range": "month",
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-01-31T00:00:00Z",
  "data_points": [
    {
      "date": "2024-01-01T00:00:00Z",
      "total_volume": 600.0,
      "max_weight": 60.0,
      "total_sets": 3,
      "average_rpe": 7.0
    }
  ],
  "summary": {
    "total_sets": 30,
    "total_volume": 18000.0,
    "max_weight": 100.0,
    "max_reps": 12,
    "average_weight": 75.5,
    "average_reps": 8.2
  }
}
```

### Entrenamientos (Workouts)

#### GET /workouts
Obtiene todos los entrenamientos del usuario autenticado.

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "name": "Push Day",
    "created_at": "2024-01-15T00:00:00Z"
  }
]
```

#### POST /workouts
Crea un nuevo entrenamiento.

**Request:**
```json
{
  "name": "Push Day"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "user_id": 1,
  "name": "Push Day",
  "created_at": "2024-01-15T00:00:00Z"
}
```

#### GET /workouts/{id}
Obtiene un entrenamiento específico.

**Response (200 OK):**
```json
{
  "id": 1,
  "user_id": 1,
  "name": "Push Day",
  "created_at": "2024-01-15T00:00:00Z"
}
```

**Errores:**
- `404 Not Found`: Entrenamiento no encontrado

#### GET /workouts/{id}/sets
Obtiene todos los sets de un entrenamiento específico.

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "workout_id": 1,
    "exercise_id": 1,
    "weight": 60.0,
    "reps": 10,
    "rest_seconds": 120,
    "rpe": 7,
    "notes": "Felt good",
    "created_at": "2024-01-15T00:00:00Z"
  }
]
```

#### POST /workouts/{id}/sets
Crea un nuevo set para un entrenamiento.

**Request:**
```json
{
  "exercise_id": 1,
  "weight": 60.0,
  "reps": 10,
  "rest_seconds": 120,
  "rpe": 7,
  "notes": "Felt good"
}
```

**Campos:**
- `exercise_id` (required): ID del ejercicio
- `weight` (required, min: 0): Peso en kg
- `reps` (required, min: 1): Número de repeticiones
- `rest_seconds` (optional, min: 0): Segundos de descanso
- `rpe` (optional, 1-10): Rate of Perceived Exertion
- `notes` (optional): Notas adicionales

**Response (201 Created):**
```json
{
  "id": 1,
  "workout_id": 1,
  "exercise_id": 1,
  "weight": 60.0,
  "reps": 10,
  "rest_seconds": 120,
  "rpe": 7,
  "notes": "Felt good",
  "created_at": "2024-01-15T00:00:00Z"
}
```

**Errores:**
- `400 Bad Request`: Datos inválidos (weight negativo, reps < 1, RPE fuera de rango)
- `404 Not Found`: Entrenamiento o ejercicio no encontrado

### Health Check

#### GET /health
Verifica que el servidor esté funcionando.

**Response (200 OK):**
```
OK
```

## Modelos de Datos

### User
```typescript
interface User {
  id: number;
  email: string;
  created_at: string; // ISO 8601 date string
}
```

### Exercise
```typescript
interface Exercise {
  id: number;
  user_id: number;
  name: string;
  created_at: string;
}
```

### Workout
```typescript
interface Workout {
  id: number;
  user_id: number;
  name: string;
  created_at: string;
}
```

### Set
```typescript
interface Set {
  id: number;
  workout_id: number;
  exercise_id: number;
  weight: number;
  reps: number;
  rest_seconds?: number;
  rpe?: number; // 1-10
  notes?: string;
  created_at: string;
}
```

### ExerciseHistory
```typescript
interface ExerciseHistory {
  exercise_id: number;
  exercise_name: string;
  sets: Set[];
  metrics?: ExerciseMetrics;
}

interface ExerciseMetrics {
  total_sets: number;
  total_volume: number;
  max_weight: number;
  max_reps: number;
  average_weight: number;
  average_reps: number;
  average_rest?: number;
  average_rpe?: number;
  first_recorded_at?: string;
  last_recorded_at?: string;
}
```

### ExerciseProgress
```typescript
interface ExerciseProgress {
  exercise_id: number;
  exercise_name: string;
  range: string; // "week" | "month" | "year"
  start_date: string;
  end_date: string;
  data_points: ProgressDataPoint[];
  summary?: ExerciseMetrics;
}

interface ProgressDataPoint {
  date: string;
  total_volume: number;
  max_weight: number;
  total_sets: number;
  average_rpe?: number;
}
```

## Requisitos del Dashboard

### 1. Autenticación y Registro
- ✅ Página de login con email y password
- ✅ Página de registro con validación de email y password (mínimo 8 caracteres)
- ✅ Manejo de errores de autenticación
- ✅ Almacenamiento seguro del JWT token
- ✅ Redirección automática si el usuario no está autenticado
- ✅ Logout funcional

### 2. Dashboard Principal
- Vista general con estadísticas rápidas:
  - Total de ejercicios
  - Total de entrenamientos
  - Último entrenamiento
  - Progreso reciente

### 3. Gestión de Ejercicios
- Lista de ejercicios del usuario
- Crear nuevo ejercicio (modal o página)
- Ver historial de un ejercicio con:
  - Lista de todos los sets realizados
  - Métricas agregadas (volumen total, peso máximo, promedio, etc.)
- Ver progreso de un ejercicio con gráficos:
  - Gráfico de volumen total por fecha
  - Gráfico de peso máximo por fecha
  - Gráfico de RPE promedio por fecha
  - Selector de rango de tiempo (semana, mes, año)

### 4. Gestión de Entrenamientos
- Lista de entrenamientos ordenados por fecha de creación (más recientes primero)
- Crear nuevo entrenamiento (ingresar nombre)
- Ver detalles de un entrenamiento:
  - Nombre del entrenamiento
  - Fecha de creación
  - Lista de sets agrupados por ejercicio
  - Posibilidad de agregar nuevos sets
  - Editar/eliminar sets (si el backend lo permite en el futuro)

### 5. Crear Sets
- Formulario para crear un set dentro de un entrenamiento:
  - Selector de ejercicio (dropdown con ejercicios del usuario)
  - Input de peso (kg)
  - Input de repeticiones
  - Input opcional de descanso (segundos)
  - Input opcional de RPE (1-10, slider o input numérico)
  - Textarea opcional para notas
  - Validación de campos

### 6. Visualizaciones y Gráficos
- Gráficos interactivos para progreso de ejercicios:
  - Línea de tiempo con volumen total
  - Línea de tiempo con peso máximo
  - Comparación de métricas entre períodos
- Implementar visualizaciones usando la librería de gráficos de tu elección

### 7. Diseño y UX
- Diseño moderno y responsive (mobile-first)
- Navegación intuitiva con sidebar o navbar
- Feedback visual para acciones (loading states, success/error messages)
- Confirmaciones para acciones destructivas
- Formateo de fechas legible
- Formateo de números (pesos, repeticiones)

### 8. Manejo de Errores
- Manejo centralizado de errores HTTP
- Mensajes de error amigables para el usuario
- Manejo de token expirado (redirigir a login)
- Manejo de errores de red

## Consideraciones Técnicas

### CORS
El backend tiene CORS habilitado, pero asegúrate de configurar correctamente las peticiones desde el frontend.

### Manejo de Tokens
- Almacenar el JWT token en `localStorage` o `sessionStorage`
- Incluir el token en todas las peticiones protegidas:
  ```
  Authorization: Bearer <token>
  ```
- Verificar si el token está expirado antes de hacer peticiones
- Implementar refresh token si es necesario (actualmente el backend no tiene refresh tokens)

### Formateo de Fechas
El backend retorna fechas en formato ISO 8601. Formatea las fechas de manera legible usando la librería de manejo de fechas de tu elección.

### Validación de Formularios
- Validar email con regex o librería de validación
- Validar password mínimo 8 caracteres
- Validar peso >= 0
- Validar reps >= 1
- Validar RPE entre 1 y 10

## Estructura de Carpetas Sugerida

```
src/
├── components/
│   ├── common/
│   │   ├── Button.tsx
│   │   ├── Input.tsx
│   │   ├── Modal.tsx
│   │   └── Loading.tsx
│   ├── exercises/
│   │   ├── ExerciseList.tsx
│   │   ├── ExerciseCard.tsx
│   │   ├── ExerciseHistory.tsx
│   │   └── ExerciseProgress.tsx
│   ├── workouts/
│   │   ├── WorkoutList.tsx
│   │   ├── WorkoutCard.tsx
│   │   ├── WorkoutDetail.tsx
│   │   └── SetForm.tsx
│   └── charts/
│       └── ProgressChart.tsx
├── pages/
│   ├── Login.tsx
│   ├── Signup.tsx
│   ├── Dashboard.tsx
│   ├── Exercises.tsx
│   ├── ExerciseDetail.tsx
│   ├── Workouts.tsx
│   └── WorkoutDetail.tsx
├── services/
│   ├── api.ts
│   ├── auth.ts
│   ├── exercises.ts
│   ├── workouts.ts
│   └── sets.ts
├── hooks/
│   ├── useAuth.ts
│   ├── useExercises.ts
│   └── useWorkouts.ts
├── types/
│   └── index.ts
├── utils/
│   ├── formatDate.ts
│   ├── formatNumber.ts
│   └── validation.ts
└── App.tsx
```

## Ejemplo de Implementación de API Service

```typescript
// services/api.ts
// Ejemplo genérico - adapta según tu HTTP client

const API_BASE_URL = 'http://localhost:8080';

// Función helper para hacer requests con autenticación
async function apiRequest(endpoint: string, options: RequestInit = {}) {
  const token = localStorage.getItem('token');
  
  const headers = {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
    ...options.headers,
  };

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  // Manejar token expirado
  if (response.status === 401) {
    localStorage.removeItem('token');
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }

  return response;
}

export default apiRequest;
```

## Prioridades de Desarrollo

1. **Fase 1 - Autenticación y Estructura Base**
   - Setup del proyecto
   - Login y Signup
   - Routing básico
   - Manejo de autenticación

2. **Fase 2 - Gestión de Ejercicios**
   - Lista de ejercicios
   - Crear ejercicio
   - Ver historial de ejercicio

3. **Fase 3 - Gestión de Entrenamientos**
   - Lista de entrenamientos
   - Crear entrenamiento
   - Ver detalles de entrenamiento
   - Agregar sets

4. **Fase 4 - Visualizaciones**
   - Gráficos de progreso
   - Dashboard con estadísticas

5. **Fase 5 - Mejoras y Pulido**
   - Mejoras de UX
   - Optimizaciones
   - Testing

## Notas Adicionales

- El backend retorna errores en formato: `{ "error": "mensaje de error" }`
- Todos los IDs son números (int64)
- Las fechas están en formato ISO 8601
- El backend valida que los recursos pertenezcan al usuario autenticado
- No hay endpoints de actualización o eliminación aún (solo creación y lectura)

## Testing

Considera implementar tests para:
- Componentes críticos
- Hooks personalizados
- Servicios de API
- Validaciones de formularios

---

**¡Comienza a desarrollar el dashboard siguiendo estos requisitos y la estructura del backend proporcionada!**


