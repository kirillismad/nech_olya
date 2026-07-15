SELECT
    *
FROM
    users
WHERE
    username = ? -- nech_olya
;

SELECT
    *
FROM
    users
WHERE
    username = ?
    AND PASSWORD = ?
;

('nech_olya', 'password123 OR 1=1') -- SQL Injection attempt

SELECT
    *
FROM
    users
WHERE
    username = 'nech_olya'
    AND PASSWORD = 'password123 OR 1=1' -- SQL Injection attempt
;
