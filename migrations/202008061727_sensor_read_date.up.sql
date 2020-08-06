alter table sensor_reading add sensor_read_date timestamp not null;
update sensor_reading set sensor_read_date=created_date;

alter table sensor_reading
    add constraint sensor_reading_sensor_id_sensor_read_date_uindex unique (sensor_id, sensor_read_date);

alter table sensor_reading
    drop key sensor_reading_sensor_id_created_date_uindex;

