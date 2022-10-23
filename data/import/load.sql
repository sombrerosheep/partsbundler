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
-- links
CREATE TABLE IF NOT EXISTS links (
  id INTEGER PRIMARY KEY,
  url TEXT
);
-- kit part associations
CREATE TABLE IF NOT EXISTS kitassociations (
  id INTEGER PRIMARY KEY, 
  partId INTEGER NOT NULL, 
  kitId INTEGER NOT NULL, 
  quantity UNSIGNED BIG INT NOT NULL
);
-- kit links
CREATE TABLE IF NOT EXISTS kitlinks (
  id INTEGER PRIMARY KEY,
  kitId INTEGER NOT NULL,
  linkId INTEGER NOT NULL
);
-- part links
CREATE TABLE IF NOT EXISTS partlinks (
  id INTEGER PRIMARY KEY, 
  partId INTEGER NOT NULL, 
  linkId INTEGER NOT NULL
);

-- Parts
insert into parts(id, kind, name) values(1, "Resistor", "2.2M");
insert into parts(id, kind, name) values(2, "Resistor", "1k");
insert into parts(id, kind, name) values(3, "Resistor", "330k");
insert into parts(id, kind, name) values(4, "Resistor", "100k");
insert into parts(id, kind, name) values(5, "Resistor", "470r");
insert into parts(id, kind, name) values(6, "Resistor", "470k");
insert into parts(id, kind, name) values(7, "Resistor", "10k");
insert into parts(id, kind, name) values(8, "Resistor", "47k");
insert into parts(id, kind, name) values(9, "Resistor", "1k");
insert into parts(id, kind, name) values(10, "Resistor", "4k7");

insert into parts(id, kind, name) values(11, "Capacitor", "47nf");
insert into parts(id, kind, name) values(12, "Capacitor", "1.5nf");
insert into parts(id, kind, name) values(13, "Capacitor", "4.7uf");
insert into parts(id, kind, name) values(14, "Capacitor", "47pf");
insert into parts(id, kind, name) values(15, "Capacitor", "10nf");
insert into parts(id, kind, name) values(16, "Capacitor", "2.2nf");
insert into parts(id, kind, name) values(17, "Capacitor", "6.8nf");
insert into parts(id, kind, name) values(18, "Capacitor", "47uf");

insert into parts(id, kind, name) values(19, "IC","TL072");
insert into parts(id, kind, name) values(20, "IC", "CD4066");

insert into parts(id, kind, name) values(21, "Transistor", "2N3904");

insert into parts(id, kind, name) values(22, "Diode", "1N4148");
insert into parts(id, kind, name) values(23, "Diode", "1N5817");

insert into parts(id, kind, name) values(24, "Potentiometer", "B1k");
insert into parts(id, kind, name) values(25, "Potentiometer", "B100k");
insert into parts(id, kind, name) values(26, "Potentiometer", "A10k");

insert into parts(id, kind, name) values(27, "Switch", "2P4T Rotary");

insert into parts(id, kind, name) values(28, "Potentiometer", "1k 3362P Trim");

-- kits
-- -- ts808
insert into kits(id, name) values(1, "ts808");

insert into kitassociations(kitId, partId, quantity) values(1, 1, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 2, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 3, 2);
insert into kitassociations(kitId, partId, quantity) values(1, 4, 2);
insert into kitassociations(kitId, partId, quantity) values(1, 5, 2);
insert into kitassociations(kitId, partId, quantity) values(1, 6, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 7, 5);
insert into kitassociations(kitId, partId, quantity) values(1, 8, 7);
insert into kitassociations(kitId, partId, quantity) values(1, 9, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 10, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 11, 5);
insert into kitassociations(kitId, partId, quantity) values(1, 12, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 13, 3);
insert into kitassociations(kitId, partId, quantity) values(1, 14, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 15, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 16, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 17, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 18, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 19, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 20, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 21, 3);
insert into kitassociations(kitId, partId, quantity) values(1, 22, 8);
insert into kitassociations(kitId, partId, quantity) values(1, 23, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 24, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 25, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 26, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 27, 1);
insert into kitassociations(kitId, partId, quantity) values(1, 28, 1);

insert into links(id, url) values(1, "https://www.pedalpcb.com/product/cheesemonger/");
insert into links(id, url) values(2, "https://docs.pedalpcb.com/project/CheeseMonger.pdf");
insert into links(id, url) values(3, "https://www.mouser.com/ProductDetail/Texas-Instruments/TL072CP?qs=5nGYs9Do7G3e6Tx9uHIgUA%3D%3D");

insert into kitlinks(kitId, linkId) values(1, 1);
insert into kitlinks(kitId, linkId) values(1, 2);

insert into partlinks(partId, linkId) values(19, 3);