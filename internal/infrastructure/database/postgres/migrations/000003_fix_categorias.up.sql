-- Primeiro, corrigir os dados existentes
UPDATE tickets SET categoria = LOWER(categoria);

-- Adicionar uma constraint para garantir que categoria seja sempre minúscula
ALTER TABLE tickets ADD CONSTRAINT check_categoria_lowercase CHECK (categoria = LOWER(categoria));

-- Adicionar uma constraint para garantir que categoria seja válida
ALTER TABLE tickets ADD CONSTRAINT check_categoria_valid CHECK (
    categoria IN (
        'financeiro',
        'comercial',
        'compliance',
        'contratos',
        'gestores',
        'meds',
        'onboarding',
        'operacional',
        'reclamacoes',
        'ti',
        'trading'
    )
);

-- Corrige categorias específicas que podem estar erradas
UPDATE tickets 
SET categoria = 'ti' 
WHERE LOWER(categoria) = 'ti';

UPDATE tickets 
SET categoria = 'financeiro' 
WHERE LOWER(categoria) = 'financeiro';

UPDATE tickets 
SET categoria = 'comercial' 
WHERE LOWER(categoria) = 'comercial';

UPDATE tickets 
SET categoria = 'compliance' 
WHERE LOWER(categoria) = 'compliance';

UPDATE tickets 
SET categoria = 'contratos' 
WHERE LOWER(categoria) = 'contratos';

UPDATE tickets 
SET categoria = 'gestores' 
WHERE LOWER(categoria) = 'gestores';

UPDATE tickets 
SET categoria = 'meds' 
WHERE LOWER(categoria) = 'meds';

UPDATE tickets 
SET categoria = 'onboarding' 
WHERE LOWER(categoria) = 'onboarding';

UPDATE tickets 
SET categoria = 'operacional' 
WHERE LOWER(categoria) = 'operacional';

UPDATE tickets 
SET categoria = 'reclamacoes' 
WHERE LOWER(categoria) = 'reclamacoes';

UPDATE tickets 
SET categoria = 'trading' 
WHERE LOWER(categoria) = 'trading'; 