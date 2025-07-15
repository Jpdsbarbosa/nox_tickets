-- Não é possível reverter exatamente para o estado anterior
-- pois não sabemos qual era a capitalização original
-- Esta migração é uma correção de dados, então o rollback é apenas simbólico 

-- Remove as constraints
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS check_categoria_lowercase;
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS check_categoria_valid; 