-- +goose Up
-- +goose StatementBegin

INSERT INTO "permission" (name, description) VALUES
("CanLogin", "Permission to login"),
("CanaccessAdminAPI", "Permission to access to admin APIs"),
("CanaccessSuperAdminAPI", "Permission to access to super admin APIs");


INSERT INTO "role" (name, description) VALUES
('SUPER ADMIN', 'Super admin role'),
('ADMIN', 'Admin role'),
('USER', 'User role'),
('GUEST', 'Guest role');

WITH temp_role_permissions (role_name, permission_name) AS (
    VALUES
    ('SUPER ADMIN', 'CanLogin'),
    ('SUPER ADMIN', 'CanaccessSuperAdminAPI'),
    ('ADMIN', 'CanLogin'),
    ('ADMIN', 'CanaccessAdminAPI'),
    ('USER', 'CanLogin')
)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM temp_role_permissions trp
JOIN roles r ON r.name = trp.role_name
JOIN permissions p ON p.name = trp.permission_name;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
