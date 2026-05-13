-- Tabela de Clientes (Lado 1)
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL
);

-- Tabela de Agendamentos (Lado N)
CREATE TABLE appointments (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    appointment_date TIMESTAMP NOT NULL,
    service_name VARCHAR(50) NOT NULL
);