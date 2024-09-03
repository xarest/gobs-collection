-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "user" (
    id  UUID DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    updatedAt TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT user_pk PRIMARY KEY (id),
);

-- Create permission table
CREATE TABLE IF NOT EXISTS "permission" (
    id  UUID DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT NULL,
    updatedAt TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT permission_pk PRIMARY KEY (id),
    CONSTRAINT permission_uk_name UNIQUE (name)
);

-- Create role table
CREATE TABLE IF NOT EXISTS "role" (
    id  UUID DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT NULL,
    updatedAt TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT role_pk PRIMARY KEY (id),
    CONSTRAINT role_uk_name UNIQUE (name)
);

-- Create user_role table
CREATE TABLE IF NOT EXISTS "user_role" (
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    updatedAt TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT user_role_pk PRIMARY KEY (user_id, role_id),
    CONSTRAINT user_role_fk_user_id FOREIGN KEY (user_id) REFERENCES user(id),
    CONSTRAINT user_role_fk_role_id FOREIGN KEY (role_id) REFERENCES role(id)
);

-- Create role_permission table
CREATE TABLE IF NOT EXISTS "role_permission" (
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    updatedAt TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT role_permission_pk PRIMARY KEY (role_id, permission_id),
    CONSTRAINT role_permission_fk_role_id FOREIGN KEY (role_id) REFERENCES role(id),
    CONSTRAINT role_permission_fk_permission_id FOREIGN KEY (permission_id) REFERENCES permission(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd
