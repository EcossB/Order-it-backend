-- Estructura Inicial: Multi-Tenant para Order It!

CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL check (type in ( 'INDEPENDENT','PARK' )), -- 'INDEPENDENT' o 'PARK'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE kitchens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_kitchens_tenant_id ON kitchens(tenant_id);

CREATE TABLE tables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    table_number VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'FREE' check (status in('FREE','OCCUPIED','WAITING_FOOD' )), -- 'FREE', 'OCCUPIED', 'WAITING_FOOD'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_tables_tenant_id ON tables(tenant_id);

CREATE TABLE menu_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    kitchen_id UUID NOT NULL REFERENCES kitchens(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_menu_items_tenant_id ON menu_items(tenant_id);

CREATE TABLE users_role(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    description VARCHAR(50) NOT NULL CHECK (description IN ('WAITER', 'CHEF', 'ADMIN')),
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, description)
);
CREATE INDEX idx_users_role_tenant_id ON users_role(tenant_id);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    role UUID REFERENCES users_role(id),
    kitchen_id UUID REFERENCES kitchens(id) ON DELETE SET NULL, -- Solo requerido para 'CHEF'
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    password_hash VARCHAR(255),
    pin_hash VARCHAR(255), -- Ahora puede ser nulo para los admins que solo usan web
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, name), -- Un usuario es único solo dentro de su organización (Tenant)
    UNIQUE (tenant_id, email) -- El email también debe ser único por Tenant
);
CREATE INDEX idx_users_tenant_id ON users(tenant_id);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    table_id UUID REFERENCES tables(id) ON DELETE SET NULL,
    waiter_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(50) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'IN_PROGRESS', 'READY','DELIVERED')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_orders_tenant_status ON orders(tenant_id, status);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id UUID REFERENCES menu_items(id) ON DELETE CASCADE,
    kitchen_id UUID REFERENCES kitchens(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1,
    status VARCHAR(50) DEFAULT 'PENDING', -- 'PENDING', 'READY'
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_order_items_tenant_id ON order_items(tenant_id);

-- --------------------------------------------------------
-- ROW-LEVEL SECURITY (RLS)
-- --------------------------------------------------------

-- Habilitar RLS en todas las tablas orientadas al inquilino
ALTER TABLE kitchens ENABLE ROW LEVEL SECURITY;
ALTER TABLE kitchens FORCE ROW LEVEL SECURITY;

ALTER TABLE tables ENABLE ROW LEVEL SECURITY;
ALTER TABLE tables FORCE ROW LEVEL SECURITY;

ALTER TABLE menu_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE menu_items FORCE ROW LEVEL SECURITY;

ALTER TABLE users_role ENABLE ROW LEVEL SECURITY;
ALTER TABLE users_role FORCE ROW LEVEL SECURITY;

ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE users FORCE ROW LEVEL SECURITY;

ALTER TABLE orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE orders FORCE ROW LEVEL SECURITY;

ALTER TABLE order_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE order_items FORCE ROW LEVEL SECURITY;

-- Políticas de aislamiento por Tenant
-- Tu backend DEBE ejecutar `SET LOCAL app.current_tenant_id = 'uuid-del-tenant';` al inicio de cada petición.

CREATE POLICY tenant_isolation_kitchens ON kitchens
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_isolation_tables ON tables
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_isolation_menu_items ON menu_items
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_isolation_users_role ON users_role
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_isolation_users ON users
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_isolation_orders ON orders
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);

CREATE POLICY tenant_isolation_order_items ON order_items
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);
