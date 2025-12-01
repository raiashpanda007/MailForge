CREATE TABLE IF NOT EXISTS users(
  id UUID PRIMARY KEY,
  name VARCHAR(200) NOT NULL ,
  password TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()

); 


CREATE TABLE IF NOT EXISTS apikeys(
  id UUID PRIMARY KEY,
  organization VARCHAR(200) UNIQUE,
  apikey TEXT UNIQUE NOT NULL,
  user_id UUID NOT NULL ,
  email_app_password TEXT ,
  FOREIGN KEY (user_id) REFERENCES users(id),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()

);
