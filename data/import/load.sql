-- part
CREATE TABLE IF NOT EXISTS parts (id INTEGER PRIMARY KEY, kind TEXT, name TEXT);
-- kit
CREATE TABLE IF NOT EXISTS kits (id INTEGER PRIMARY KEY, name TEXT, schematic TEXT, diagram TEXT);
-- links
CREATE TABLE IF NOT EXISTS links (id INTEGER PRIMARY KEY, url TEXT);
-- kit part associations
CREATE TABLE IF NOT EXISTS kitassociations (id INTEGER PRIMARY KEY, partId UNSIGNED BIG INT NOT NULL, kitId UNSIGNED BIG INT NOT NULL, quantity UNSIGNED BIG INT NOT NULL);
-- kit links
CREATE TABLE IF NOT EXISTS kitlinks (id INTEGER PRIMARY KEY, kitId UNSIGNED BIG INT NOT NULL, linkId UNSIGNED BIG INT NOT NULL);
-- part links
CREATE TABLE IF NOT EXISTS partlinks (id INTEGER PRIMARY KEY, partId UNSIGNED BIG INT NOT NULL, linkId UNSIGNED BIG INT NOT NULL);

-- Parts
insert into parts(id, kind, name) values(1, "Resistor", "2.2M");
insert into parts(id, kind, name) values(2, "Resistor", "2k");
insert into parts(id, kind, name) values(3, "Resistor", "330k");
insert into parts(id, kind, name) values(4, "Resistor", "100k");
insert into parts(id, kind, name) values(5, "Resistor", "470r");
insert into parts(id, kind, name) values(6, "Resistor", "470k");
insert into parts(id, kind, name) values(7, "Resistor", "10k");
insert into parts(id, kind, name) values(8, "Resistor", "47k");
insert into parts(id, kind, name) values(9, "Resistor", "1k");
insert into parts(id, kind, name) values(10, "Resistor", "47k");
insert into parts(id, kind, name) values(11, "Resistor", "4k7");
insert into parts(id, kind, name) values(12, "Capacitor", "47nf");
insert into parts(id, kind, name) values(13, "Capacitor", "1.5nf");
insert into parts(id, kind, name) values(14, "Capacitor", "4.7uf");
insert into parts(id, kind, name) values(15, "Capacitor", "47pf");
insert into parts(id, kind, name) values(16, "Capacitor", "10nf");
insert into parts(id, kind, name) values(17, "Capacitor", "2.2nf");
insert into parts(id, kind, name) values(18, "Capacitor", "6.8nf");
insert into parts(id, kind, name) values(19, "IC","TL072");
insert into parts(id, kind, name) values(20, "IC", "CD4066");
insert into parts(id, kind, name) values(21, "Transistor", "2N3904");
insert into parts(id, kind, name) values(22, "Diode", "1N4148");
insert into parts(id, kind, name) values(23, "Diode", "1N5817");
insert into parts(id, kind, name) values(24, "Potentiometer", "B1k");
insert into parts(id, kind, name) values(25, "Potentiometer", "B100k");
insert into parts(id, kind, name) values(26, "Potentiometer", "A10k");
insert into parts(id, kind, name) values(27, "Switch", "2P4T Rotary");
insert into parts(id, kind, name) values(28, "Trim Potemtiometer", "1k 3362P");

