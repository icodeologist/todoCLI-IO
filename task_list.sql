DROP TABLE IF EXISTS tasks;
CREATE TABLE tasks(
  id INT AUTO_INCREMENT NOT NULL,
  name VARCHAR(128) NOT NULL,
  status BOOLEAN,
  PRIMARY KEY(`id`)
);


INSERT INTO tasks 
  (name , status)
VALUES
  ("Code todoapp in go", 0),
  ("Write a para", 1);

  
