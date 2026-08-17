# Guía — Gestión de Usuarios

## 1. Propósito

Este documento explica dónde y cómo se registra, edita o inactiva un
usuario de acceso al sistema (la cuenta con email y contraseña que le
permite a un tercero iniciar sesión). Úsalo para responder preguntas de
tipo "¿dónde creo un usuario?", "¿cómo le doy acceso al sistema a un
colaborador?", "¿cómo cambio la contraseña de un usuario?", "¿cómo
inactivo a un usuario?".

## 2. Idea clave: un usuario siempre está vinculado a un tercero ya existente

Un usuario no es un registro independiente — siempre corresponde a un
**Tercero** que ya debe estar registrado de antes (ver el manual de
Terceros). Al crear un usuario, se busca y selecciona ese tercero; **no
se puede crear un tercero nuevo desde la pantalla de Usuarios**. Si la
persona a la que se le va a dar acceso todavía no está registrada como
tercero, primero hay que registrarla ahí.

## 3. Cómo llegar

1. Presiona **"Seleccionar módulo"** en la barra superior.
2. Elige el módulo **Configuración**.
3. Selecciona la función **Seguridad**.
4. Haz clic en el ítem **Usuarios**.

(Los 4 pasos son obligatorios y van en ese orden — Usuarios está dentro
de la función "Seguridad", igual que Roles, no directamente bajo
Configuración.)

## 4. Registrar un usuario nuevo

1. Presiona **"Seleccionar módulo"** en la barra superior.
2. Elige el módulo **Configuración**.
3. Selecciona la función **Seguridad**.
4. Haz clic en el ítem **Usuarios**.
5. Presiona el botón **"Nuevo Usuario"**.
6. Completa el formulario "Registro de Usuarios". Estos son
   **obligatorios**:
   - **Buscar Tercero**: escribe al menos 2 caracteres (documento o
     nombre) y selecciona el tercero de la lista de resultados. Solo se
     pueden vincular terceros ya registrados.
   - **Email** del usuario.
   - **Rol** que tendrá el usuario.
   - **Contraseña** y **Repita Contraseña**: mínimo 8 caracteres, y
     ambos campos deben coincidir.

   Estos son opcionales:
   - **Imagen Perfil**.
   - **Estado** (interruptor "Registro Activo" / "Registro Inactivo") —
     si no se toca, el comportamiento por defecto depende del sistema;
     si es importante que quede activo desde el inicio, revisa el
     interruptor antes de guardar.
7. Presiona **"Registrar"**.

## 5. Editar un usuario existente — cuidado con el tercero vinculado

1. Presiona **"Seleccionar módulo"** en la barra superior.
2. Elige el módulo **Configuración**.
3. Selecciona la función **Seguridad**.
4. Haz clic en el ítem **Usuarios**.
5. Ubica el usuario en el listado (ver la limitación de búsqueda en la
   sección 8) y presiona el ícono de **editar** (lápiz) en su fila.
6. En el formulario ya precargado, el **Email** y la **Contraseña**
   aparecen deshabilitados — no se pueden cambiar desde aquí (para
   contraseña, ver sección 7). Sí se pueden cambiar el **Rol**, la
   **Imagen Perfil** y el **Estado**.
7. Presiona **"Actualizar"**.

**Advertencia importante:** hoy, al guardar una edición de usuario, el
sistema tiene una falla que puede desvincular o vincular mal el tercero
asociado a ese usuario (guarda internamente un dato incorrecto en vez
del tercero real). Si vas a explicar este paso, no lo describas como
completamente seguro: si un usuario reporta que después de editarlo
"perdió" su vínculo con el tercero correcto, esto es una falla conocida
del sistema y no un error del usuario — sugiere contactar a soporte.

## 6. Cambiar la contraseña de un usuario — no es posible hoy

Aunque en el listado de usuarios existe un botón **"Cambiar
Contraseña"** (con ícono de llave), hoy ese botón hace exactamente lo
mismo que el botón "Editar": lleva al mismo formulario de edición, donde
el campo de contraseña aparece deshabilitado. **No hay ninguna forma
funcional de cambiar la contraseña de un usuario desde la interfaz
actualmente.** Si te preguntan cómo resetear o cambiar una contraseña,
dilo con honestidad y sugiere contactar a soporte.

## 7. Inactivar un usuario (no existe "eliminar")

No hay botón de eliminar en la tabla de usuarios. Lo que sí se puede
hacer es inactivarlo:

1. Entra a editar el usuario (ver sección 5).
2. Apaga el interruptor de **Estado**, dejándolo como "Registro
   Inactivo".
3. Presiona **"Actualizar"**.

(Recuerda la advertencia de la sección 5 sobre el tercero vinculado
antes de guardar cualquier edición, incluida esta.)

## 8. Limitación conocida — importante

La **búsqueda por texto** en la pantalla de Usuarios no funciona: por
una falla interna, usa por error la misma búsqueda de la pantalla de
Terceros, así que no devuelve usuarios reales al buscar por email o
nombre. Si te preguntan por qué no encuentran un usuario buscándolo,
sugiere revisar el listado inicial (los usuarios más recientes) que
aparece al entrar a la pantalla, y si tiene muchos usuarios y no lo
encuentra ahí, recomienda contactar a soporte.

## 9. Guía de tono para preguntas típicas de este tema

Estas notas son solo para que tú (el asistente) entiendas qué se espera
de cada tipo de pregunta — nunca repitas estas notas ni las palabras
"pregunta"/"respuesta" al usuario, solo respóndele directamente.

- Si preguntan cómo dar acceso al sistema a alguien (un colaborador, un
  cliente, etc.): explica que primero debe existir como Tercero, y luego
  sigue los pasos de la sección 4 para crear su usuario.
- Si preguntan cómo cambiar la contraseña de un usuario: sé directo en
  que hoy no es posible desde la interfaz (sección 6), sin describir el
  botón "Cambiar Contraseña" como si funcionara.
- Si preguntan cómo editar el email de un usuario: aclara que el email
  no se puede modificar una vez creado el usuario (queda deshabilitado
  en el formulario de edición).
- Si preguntan cómo inactivar o "eliminar" un usuario: explica que no
  existe eliminar, solo inactivar (sección 7), y guía con esos pasos.
- Si dicen que buscaron un usuario y no lo encuentran: explica la
  limitación de búsqueda (sección 8).

## 10. Reglas para el asistente sobre este tema

- No digas que se puede crear un tercero nuevo desde la pantalla de
  Usuarios — siempre hay que buscar uno ya existente (sección 2).
- No ofrezcas el botón "Cambiar Contraseña" como una forma real de
  cambiar la contraseña — hoy no funciona (sección 6).
- No prometas que la búsqueda por texto de Usuarios va a funcionar —
  advierte la limitación de la sección 8 si preguntan por qué no
  encuentran algo.
- No digas que se puede editar el email de un usuario ya creado.
- No ofrezcas "eliminar" un usuario como si fuera posible — solo se
  puede inactivar (sección 7).
- Si preguntan por editar un usuario, no lo presentes como una acción
  100% segura sin matices — menciona la advertencia sobre el tercero
  vinculado de la sección 5 si la pregunta se presta para ello (por
  ejemplo, si reportan un problema después de editar a alguien).
