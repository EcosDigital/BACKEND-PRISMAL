# Guía — Gestión de Empresas

## 1. Propósito

Este documento explica dónde y cómo se edita la ficha de una empresa
(razón social, NIT, datos tributarios, logo, e integración del asistente
de IA). Úsalo para responder preguntas de tipo "¿dónde edito los datos de
mi empresa?", "¿dónde cambio el NIT?", "¿dónde configuro el Agente?".

## 2. Cómo llegar

1. Presiona **"Seleccionar módulo"** en la barra superior.
2. Elige el módulo **Configuración**.
3. Selecciona la función **Gestionar Empresa**.
4. Haz clic en el ítem **Empresas**.

(Los 4 pasos son obligatorios y van en ese orden — Empresas está dentro
de la función "Gestionar Empresa", igual que Sedes, no directamente bajo
Configuración.)

## 3. Editar los datos de una empresa

1. Presiona **"Seleccionar módulo"** en la barra superior.
2. Elige el módulo **Configuración**.
3. Selecciona la función **Gestionar Empresa**.
4. Haz clic en el ítem **Empresas**.
5. Ubica la empresa en el listado inicial (ver la limitación de búsqueda
   en la sección 6) y presiona el ícono de **editar**.
6. El formulario tiene 3 pestañas:
   - **Información General**: razón social, NIT y dígito de verificación
     (DV), tipo de empresa, dirección, ubicación, teléfono, email,
     página web, logo, si usa sedes y el estado (activo/inactivo).
   - **Información Tributaria**: naturaleza, actividad económica, tipo de
     contribuyente, régimen tributario, representante legal, y el código
     de licencia (este último normalmente se asigna solo al momento del
     registro inicial de la empresa — no debería necesitar tocarse
     manualmente).
   - **Integración IA**: el endpoint y el token del servicio de IA que
     atiende al Agente Prismar de esta empresa (ambos opcionales).
7. Guarda los cambios.

**Cuidado:** si cambias el NIT de la empresa, el sistema no revisa si ya
existe otra empresa con ese mismo NIT — asegúrate de no repetir el NIT de
otra empresa ya registrada.

## 4. "Usa sedes" y "N° de sedes" son solo informativos

Estos 2 campos del formulario **no crean ni administran sedes reales** —
son solo un dato indicativo. Para crear o editar una sede de verdad, hay
que ir a la pantalla de Sedes (módulo Configuración, función Gestionar
Empresa, ítem Sedes), no aquí.

## 5. Inactivar una empresa (no existe "eliminar")

Igual que con Sedes y Terceros: no hay forma de eliminar una empresa de
forma definitiva. Lo que se puede hacer es entrar a editarla y apagar el
interruptor de **Estado** en la pestaña "Información General", dejándola
como inactiva.

## 6. Limitación conocida — importante

La **búsqueda por texto** en la pantalla de Empresas no funciona: no
filtra por nombre ni por NIT como debería. Si te preguntan por qué la
búsqueda no encuentra una empresa, sugiere revisar el listado inicial que
aparece al entrar a la pantalla, y si tiene muchas empresas registradas y
no la encuentra ahí, recomienda contactar a soporte.

## 7. Guía de tono para preguntas típicas de este tema

Estas notas son solo para que tú (el asistente) entiendas qué se espera
de cada tipo de pregunta — nunca repitas estas notas ni las palabras
"pregunta"/"respuesta" al usuario, solo respóndele directamente.

- Si preguntan dónde editar los datos de la empresa: guía con los pasos
  de la sección 3.
- Si preguntan dónde configurar el Agente Prismar: aclara que se hace
  desde la misma ficha de la empresa, pestaña "Integración IA".
- Si preguntan cómo crear una sede desde esta pantalla: aclara que aquí
  no se crean sedes de verdad (sección 4), y redirige al ítem Sedes.
- Si preguntan cómo eliminar una empresa: aclara que no es posible, y
  explica cómo inactivarla en su lugar (sección 5).

## 8. Reglas para el asistente sobre este tema

- No prometas que la búsqueda por texto de Empresas va a funcionar —
  indica la limitación de la sección 6.
- No digas que marcar "Usa sedes" o poner un número ahí crea sedes de
  verdad — siempre aclara que eso se hace desde la pantalla de Sedes.
- No ofrezcas "eliminar" una empresa — solo se puede inactivar.
- Si preguntan por cambiar el NIT, recuerda la advertencia de la sección
  3 sobre no repetir el de otra empresa.
