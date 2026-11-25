-- Restaura extensões, configuração de busca e índices que se perderam na
-- migração do banco do Supabase para a VPS (o dump não levou as extensões do
-- schema "extensions", e com elas a config portuguese_unaccent e os índices
-- GIN). Tudo idempotente: em bancos que já têm esses objetos não faz nada.

CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_ts_config WHERE cfgname = 'portuguese_unaccent'
    ) THEN
        CREATE TEXT SEARCH CONFIGURATION portuguese_unaccent (COPY = portuguese);
        ALTER TEXT SEARCH CONFIGURATION portuguese_unaccent
            ALTER MAPPING FOR hword, hword_part, word
            WITH unaccent, portuguese_stem;
    END IF;
END
$$;

-- 000029
CREATE INDEX IF NOT EXISTS idx_products_search
ON products
USING GIN (
    to_tsvector(
        'portuguese_unaccent',
        coalesce(name, '') || ' ' || coalesce(description, '')
    )
);
CREATE INDEX IF NOT EXISTS idx_products_barcode_trgm
ON products USING GIN (barcode gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_categories_name_trgm
ON product_categories USING GIN (name gin_trgm_ops);

-- 000031
CREATE INDEX IF NOT EXISTS idx_customers_search
ON customers
USING GIN (to_tsvector('portuguese_unaccent', coalesce(full_name, '')));
CREATE INDEX IF NOT EXISTS idx_customers_cpf_trgm
ON customers USING GIN (cpf gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_customers_rg_trgm
ON customers USING GIN (rg gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_customers_email_trgm
ON customers USING GIN (email gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_customers_whatsapp_trgm
ON customers USING GIN (whatsapp gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_customers_mobile_phone_trgm
ON customers USING GIN (mobile_phone gin_trgm_ops);

-- 000033
CREATE INDEX IF NOT EXISTS idx_customers_full_name_trgm
ON customers USING GIN (full_name gin_trgm_ops);

-- 000034
CREATE INDEX IF NOT EXISTS idx_bills_description_search
ON bills_payable
USING GIN (to_tsvector('portuguese_unaccent', coalesce(description, '')));
CREATE INDEX IF NOT EXISTS idx_vendors_name_trgm
ON vendors USING GIN (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_bill_categories_name_trgm
ON bill_categories USING GIN (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_payment_methods_name_trgm
ON payment_methods USING GIN (name gin_trgm_ops);
