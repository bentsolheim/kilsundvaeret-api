insert into logger (name, description) values ('bua', 'Vanntemperaturlogger i bua på bryggen');
insert into sensor (logger_id, name, type, unit) values (1, 'inne-temp', 'temperature', 'C');
insert into sensor (logger_id, name, type, unit) values (1, 'inne-humidity', 'humidity', '%');