INSERT INTO customers (name, email) VALUES ($1, $2) RETURNING *;

SELECT * FROM customers ORDER BY id;

UPDATE customers SET name = $2, email = $3 WHERE id = $1;

DELETE FROM customers WHERE id = $1;

INSERT INTO appointments (customer_id, appointment_date, service_name) VALUES ($1, $2, $3) RETURNING *;

SELECT 
    c.id, c.name, c.email,
    COALESCE(json_agg(
        json_build_object(
            'id', a.id,
            'appointment_date', a.appointment_date,
            'service_name', a.service_name
        )
    ) FILTER (WHERE a.id IS NOT NULL), '[]')::json AS appointments
FROM customers c
LEFT JOIN appointments a ON c.id = a.customer_id
WHERE c.id = $1
GROUP BY c.id;