CREATE TABLE roles (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password TEXT NOT NULL,
                       role_id INT NOT NULL REFERENCES roles(id),
                       last_access TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE role_rights (
                             id SERIAL PRIMARY KEY,
                             role_id INT REFERENCES roles(id) ON DELETE CASCADE,
                             section TEXT NOT NULL,
                             route TEXT NOT NULL,
                             r_create BOOLEAN DEFAULT FALSE,
                             r_read BOOLEAN DEFAULT FALSE,
                             r_update BOOLEAN DEFAULT FALSE,
                             r_delete BOOLEAN DEFAULT FALSE
);

-- Insert Roles
INSERT INTO roles (id, name) VALUES
                                 (1, 'Admin'),
                                 (2, 'Editor'),
                                 (3, 'Viewer');

-- Insert users
INSERT INTO users (email, password, id)
VALUES
    ('admin@gmail.com', 'adminadmin', 1),
    ('rachman@gmail.com', 'adminadmin', 2),
    ('user@gmail.com', 'user123', 3);

-- Insert Role Rights
INSERT INTO role_rights (role_id, section, route, r_create, r_read, r_update, r_delete)
VALUES
    (1, 'be', '/users/user', TRUE, TRUE, TRUE, TRUE),  -- Admin: Full access
    (2, 'be', '/users/user', TRUE, TRUE, FALSE, FALSE), -- Editor: Can read and create, but not update or delete
    (3, 'be', '/users/user', FALSE, TRUE, FALSE, FALSE); -- Can only read
