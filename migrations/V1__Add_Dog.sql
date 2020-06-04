CREATE TYPE dog_breed AS ENUM ('Chihuahua');

create table dog (
    id uuid PRIMARY KEY not null,
    name text not null,
    breed dog_breed not null,
    birth_date timestamp with time zone not null
);
