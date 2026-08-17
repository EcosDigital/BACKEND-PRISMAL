# Guía — Gestión de Roles

## 1. Propósito

Este documento explica dónde se crean, editan e inactivan los roles de
usuario del sistema (por ejemplo "Administrador", "Operativo",
"Financiero"). Úsalo para responder preguntas de tipo "¿dónde creo un
rol nuevo?", "¿cómo edito un rol?", "¿cómo elimino/inactivo un rol?".

## 2. Cómo llegar

1. Presiona **"Seleccionar módulo"** en la barra superior.
2. Elige el módulo **Configuración**.
3. Selecciona la función **Seguridad**.
4. Haz clic en el ítem **Roles**.

(Los 4 pasos son obligatorios y van en ese orden — Roles está dentro de
la función "Seguridad", no directamente bajo Configuración.)

## 3. Registrar un rol nuevo — limitación crítica, hoy no guarda

1. Presiona **"Seleccionar módulo"** en la barra superior.
2. Elige el módulo **Configuración**.
3. Selecciona la función **Seguridad**.
4. Haz clic en el ítem **Roles**.
5. Completa el formulario que aparece en la parte superior de la
   pantalla. Estos son los campos, todos **obligatorios**:
   - **Nombre** del rol.
   - **Tipo**: se elige de una lista (por ejemplo Administrativo,
     Operativo, Gerencial, Financiero, Soporte).
   - **Intentos Login**: número máximo de intentos de inicio de sesión
     permitidos para los usuarios con este rol.

   Estos son opcionales (casillas de verificación):
   - **Multisession**: permite que un usuario con este rol inicie
     sesión simultáneamente desde varios dispositivos.
   - **Externo**: marca el rol como perteneciente a un usuario externo.
6. Presiona el botón **"Registrar"**.

**Importante — hoy esto no funciona:** al presionar "Registrar" para un
rol nuevo, el sistema no guarda nada; el formulario no envía la
información al servidor. Si un usuario dice que llenó el formulario, le
dio a "Registrar" y el rol nunca aparece en la lista, esto no es un
error de su parte — es una falla conocida del sistema. No hay hoy una
forma de crear un rol nuevo desde la pantalla de Roles. Sé honesto sobre
esto si te preguntan; no sugieras reintentar el mismo botón como si
fuera a funcionar distinto.

## 4. Editar un rol existente — limitación crítica, hoy no está disponible

En la tabla de roles (debajo del formulario, en la misma pantalla del
paso 2), cada fila tiene un ícono de **editar** (lápiz). Sin embargo,
hoy ese ícono **no abre ningún formulario de edición ni lleva a ninguna
otra pantalla** — no tiene ningún efecto visible para el usuario.

Si te preguntan cómo editar un rol (por ejemplo, cambiar su nombre, su
tipo, o los intentos de login permitidos), dilo con honestidad: hoy esa
función no está disponible desde la interfaz, y sugiere contactar a
soporte para corregir un rol existente.

## 5. Inactivar o eliminar un rol — no existe ningún mecanismo

A diferencia de otras pantallas del sistema (como Empresas, Sedes o
Terceros), un rol **no tiene ni un botón de eliminar ni un interruptor
de estado activo/inactivo** — ni en el formulario, ni en la tabla. Un
rol que ya está registrado en el sistema no se puede inactivar ni
eliminar por ningún camino disponible hoy en la interfaz.

Si te preguntan cómo dar de baja un rol, aclara que hoy no existe esa
opción, y sugiere contactar a soporte.

## 6. Limitaciones conocidas — importante

Hoy, la pantalla de Roles tiene su función de creación completamente
inoperante y no ofrece edición ni inactivación:

- **Crear**: el botón "Registrar" no envía la solicitud de creación al
  servidor — no se guarda el rol (sección 3).
- **Editar**: el ícono de editar de la tabla no realiza ninguna acción
  (sección 4).
- **Inactivar/eliminar**: no existe esa opción en ningún lugar de la
  pantalla (sección 5).

En resumen: hoy, desde la interfaz, no es posible crear, editar ni
inactivar un rol. Si un usuario necesita cualquiera de estas tres
acciones, la respuesta honesta es que debe contactar a soporte.

## 7. Guía de tono para preguntas típicas de este tema

Estas notas son solo para que tú (el asistente) entiendas qué se espera
de cada tipo de pregunta — nunca repitas estas notas ni las palabras
"pregunta"/"respuesta" al usuario, solo respóndele directamente.

- Si preguntan dónde se crea un rol nuevo: da los pasos de navegación y
  del formulario (sección 3), pero sé claro y directo en que hoy el
  botón "Registrar" no guarda el rol — no lo presentes como un simple
  detalle menor, es la respuesta principal a la pregunta.
- Si preguntan cómo editar un rol: aclara de inmediato que hoy esa
  función no está disponible (sección 4), sin describir un flujo que no
  existe.
- Si preguntan cómo eliminar o inactivar un rol: aclara que no existe
  esa opción hoy (sección 5).
- Si preguntan por qué no ven reflejado el rol que intentaron crear:
  explica la limitación de la sección 3 en vez de sugerir que revisen su
  conexión o vuelvan a intentarlo.

## 8. Reglas para el asistente sobre este tema

- No digas que "Registrar" crea el rol — hoy no lo hace (sección 3).
- No describas el ícono de editar de la tabla de roles como si abriera
  un formulario de edición — hoy no hace nada (sección 4).
- No ofrezcas "eliminar" ni "inactivar" un rol como si fuera posible —
  hoy no existe esa opción en absoluto (sección 5).
- No inventes que existe una confirmación o mensaje de éxito al crear un
  rol — no se debe prometer un resultado que el sistema no produce.
