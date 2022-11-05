-- part
CREATE TABLE IF NOT EXISTS parts (
  id INTEGER PRIMARY KEY, 
  kind TEXT NOT NULL,
  name TEXT NOT NULL
);
-- kit
CREATE TABLE IF NOT EXISTS kits (
  id INTEGER PRIMARY KEY, 
  name TEXT NOT NULL,
  schematic TEXT DEFAULT "" NOT NULL, 
  diagram TEXT DEFAULT "" NOT NULL
);
-- kit part associations
CREATE TABLE IF NOT EXISTS kitparts (
  id INTEGER PRIMARY KEY, 
  partId INTEGER NOT NULL, 
  kitId INTEGER NOT NULL, 
  quantity UNSIGNED BIG INT NOT NULL
);
-- kit links
CREATE TABLE IF NOT EXISTS kitlinks (
  id INTEGER PRIMARY KEY,
  kitId INTEGER NOT NULL,
  link TEXT NOT NULL
);
-- part links
CREATE TABLE IF NOT EXISTS partlinks (
  id INTEGER PRIMARY KEY, 
  partId INTEGER NOT NULL, 
  link TEXT NOT NULL
);