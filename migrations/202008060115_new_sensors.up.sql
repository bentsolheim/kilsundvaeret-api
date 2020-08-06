insert into logger (name, description) values ('met', 'Data fra Meteorologisk Institutt');
insert into sensor (logger_id, name, type, unit) values ((select id from logger where name='met'), 'air-temperature', 'temperature', 'C');
insert into sensor (logger_id, name, type, unit) values ((select id from logger where name='met'), 'relative-humidity', 'humidity', '%');
insert into sensor (logger_id, name, type, unit) values ((select id from logger where name='met'), 'wind-from-direction', 'wind-direction', 'degrees');
insert into sensor (logger_id, name, type, unit) values ((select id from logger where name='met'), 'wind-speed', 'wind_speed', 'm/s');
insert into sensor (logger_id, name, type, unit) values ((select id from logger where name='met'), 'air-pressure-at-sea-level', 'air-pressure', 'hPa');