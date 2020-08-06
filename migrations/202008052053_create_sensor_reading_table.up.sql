create table sensor_reading
(
    id int auto_increment
        primary key,
    sensor_id int null,
    created_date timestamp not null,
    value float default 0 not null,
    constraint sensor_reading_sensor_id_created_date_uindex
        unique (sensor_id, created_date),
    constraint sensor_reading_sensor_id_fk
        foreign key (sensor_id) references sensor (id)
            on delete cascade
);