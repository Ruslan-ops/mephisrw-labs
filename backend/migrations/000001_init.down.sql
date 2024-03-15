DROP TABLE users_done;

DROP TABLE bank_variance;

ALTER TABLE users_done
    DROP CONSTRAINT constraint_line_unique;

ALTER TABLE users_done
    DROP CONSTRAINT check_percentage_range