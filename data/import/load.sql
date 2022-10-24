-- part
CREATE TABLE IF NOT EXISTS parts (
  id INTEGER PRIMARY KEY, 
  kind TEXT, 
  name TEXT
);
-- kit
CREATE TABLE IF NOT EXISTS kits (
  id INTEGER PRIMARY KEY, 
  name TEXT, 
  schematic TEXT, 
  diagram TEXT
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
  link TEXT
);
-- part links
CREATE TABLE IF NOT EXISTS partlinks (
  id INTEGER PRIMARY KEY, 
  partId INTEGER NOT NULL, 
  link TEXT
);

-- Parts
insert into parts(id, kind, name)
  values
    (1, "Resistor", "2.2M"),
    (2, "Resistor", "1k"),
    (3, "Resistor", "330k"),
    (4, "Resistor", "100k"),
    (5, "Resistor", "470r"),
    (6, "Resistor", "470k"),
    (7, "Resistor", "10k"),
    (8, "Resistor", "47k"),
    (9, "Resistor", "1k"),
    (10, "Resistor", "4k7"),
    (11, "Capacitor", "47nf"),
    (12, "Capacitor", "1.5nf"),
    (13, "Capacitor", "4.7uf"),
    (14, "Capacitor", "47pf"),
    (15, "Capacitor", "10nf"),
    (16, "Capacitor", "2.2nf"),
    (17, "Capacitor", "6.8nf"),
    (18, "Capacitor", "47uf"),
    (19, "IC","TL072"),
    (20, "IC", "CD4066"),
    (21, "Transistor", "2N3904"),
    (22, "Diode", "1N4148"),
    (23, "Diode", "1N5817"),
    (24, "Potentiometer", "B1k"),
    (25, "Potentiometer", "B100k"),
    (26, "Potentiometer", "A10k"),
    (27, "Switch", "2P4T Rotary"),
    (28, "Potentiometer", "1k 3362P Trim");

-- kits
-- -- ts808
insert into kits(id, name) values(1, "ts808");

insert into kitparts(kitId, partId, quantity)
  values
    (1, 1, 1),
    (1, 2, 1),
    (1, 3, 2),
    (1, 4, 2),
    (1, 5, 2),
    (1, 6, 1),
    (1, 7, 5),
    (1, 8, 7),
    (1, 9, 1),
    (1, 10, 1),
    (1, 11, 5),
    (1, 12, 1),
    (1, 13, 3),
    (1, 14, 1),
    (1, 15, 1),
    (1, 16, 1),
    (1, 17, 1),
    (1, 18, 1),
    (1, 19, 1),
    (1, 20, 1),
    (1, 21, 3),
    (1, 22, 8),
    (1, 23, 1),
    (1, 24, 1),
    (1, 25, 1),
    (1, 26, 1),
    (1, 27, 1),
    (1, 28, 1);

insert into kitlinks(kitId, link)
  values
    (1, "https://www.pedalpcb.com/product/cheesemonger/"),
    (2, "https://docs.pedalpcb.com/project/CheeseMonger.pdf");

insert into partlinks(partId, link) values(19, "https://www.mouser.com/ProductDetail/Texas-Instruments/TL072CP?qs=5nGYs9Do7G3e6Tx9uHIgUA%3D%3D");