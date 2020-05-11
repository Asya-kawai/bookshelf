CREATE TABLE IF NOT EXISTS default.books (
  id MEDIUMINT NOT NULL AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  author VARCHAR(255),
  published_at VARCHAR(255),
  image_url TEXT,
  description TEXT,
  PRIMARY KEY (id)
);
