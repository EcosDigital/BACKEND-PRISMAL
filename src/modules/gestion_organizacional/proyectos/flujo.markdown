ref_estado_proyecto      ref_prioridad_proyecto
       │                         │
       └──────┐    ┌─────────────┘
              ▼    ▼
         cfg_proyectos  ◄──────────────────────────────┐
              │                                         │
              ├──► cfg_stakeholders                     │
              │        (soft delete)                    │
              │                                         │
              └──► bit_comentarios ◄── ref_tipo_evidencia
                       │
                       └──► bit_evidencias
                                (ON DELETE CASCADE)