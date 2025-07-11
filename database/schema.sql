create table if not exists sets (
    name varchar(50) unique not null,
    performed_at timestamp not null default (now()),
    weight decimal not null default 0,
    unit varchar(4) not null,
    reps smallint not null,
    set_count smallint not null,
    notes text,
    split_day varchar(20),
    program varchar(50),
    tags text
)
