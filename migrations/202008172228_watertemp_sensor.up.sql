insert into sensor (logger_id, name, type, unit) values ((select id from logger where name='bua'), 'water-temperature', 'water-temperature', 'C');
