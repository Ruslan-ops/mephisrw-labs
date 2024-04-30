DROP TABLE users_done;

DROP TABLE bank_variance_1b_lab;

ALTER TABLE users_done
    DROP CONSTRAINT constraint_line_unique;

ALTER TABLE users_done
    DROP CONSTRAINT check_percentage_range