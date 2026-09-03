-- Estructura Inicial: Multi-Tenant para Order It!

CREATE TABLE tenants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'INDEPENDENT' o 'PARK'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE kitchens (
    id SERIAL PRIMARY KEY,
    tenant_id INT REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tables (
    id SERIAL PRIMARY KEY,
    tenant_id INT REFERENCES tenants(id) ON DELETE CASCADE,
    table_number VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'FREE', -- 'FREE', 'OCCUPIED', 'WAITING_FOOD'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE menu_items (
    id SERIAL PRIMARY KEY,
    kitchen_id INT REFERENCES kitchens(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users_role(
    id SERIAL PRIMARY KEY,
    description VARCHAR(50) NOT NULL, -- 'WAITER', 'CHEF', 'ADMIN'
    status int,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    tenant_id INT REFERENCES tenants(id) ON DELETE CASCADE,
    role INT REFERENCES users_role(id), -- 'WAITER', 'CHEF', 'ADMIN'
    kitchen_id INT REFERENCES kitchens(id) ON DELETE SET NULL, -- Solo requerido para 'CHEF'
    name VARCHAR(255) NOT NULL,
    pin_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, name) -- Un usuario es único solo dentro de su organización (Tenant)
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    tenant_id INT REFERENCES tenants(id) ON DELETE CASCADE,
    table_id INT REFERENCES tables(id) ON DELETE SET NULL,
    waiter_id INT REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(50) DEFAULT 'PENDING', -- 'PENDING', 'IN_PROGRESS', 'READY', 'DELIVERED'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
     id SERIAL PRIMARY KEY,
     order_id INT REFERENCES orders(id) ON DELETE CASCADE,
     menu_item_id INT REFERENCES menu_items(id) ON DELETE CASCADE,
     kitchen_id INT REFERENCES kitchens(id) ON DELETE CASCADE,
     quantity INT NOT NULL DEFAULT 1,
     status VARCHAR(50) DEFAULT 'PENDING', -- 'PENDING', 'READY'
     notes TEXT,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
