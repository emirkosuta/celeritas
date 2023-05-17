CREATE TABLE $MIGRATIONNAME$ (
  id int(11) NOT NULL AUTO_INCREMENT,
  title varchar(255) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  PRIMARY KEY (id)
);

-- add auto update of updated_at
-- If you already have this trigger
-- you can delete the next line
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON $MIGRATIONNAME$
FOR EACH ROW
SET NEW.updated_at = NOW();