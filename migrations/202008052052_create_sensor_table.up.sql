create table sensor
(
    id int auto_increment
        primary key,
    logger_id int not null,
    name varchar(255) not null,
    type varchar(255) not null,
    unit varchar(255) null,
    constraint sensor_logger_id_fk
        foreign key (logger_id) references logger (id)
            on delete cascade
);