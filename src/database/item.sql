USE 'test'; -- pick the database
CREATE TABLE IF NOT EXISTS item (
       -- NOT NULL means it can never be nothing (this is due to it being our primary key)
       -- AUTO_INCREMENT means for a new entry in the table,
       --                we automatically give it id of previous +1
       Id int unsigned NOT NULL AUTO_INCREMENT,
       Title text,
       Description text,
       Price int,
       PRIMARY KEY (Id)
);


