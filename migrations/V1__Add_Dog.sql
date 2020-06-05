CREATE TYPE dog_breed AS ENUM ('Chihuahua');

create table dog (
    id uuid PRIMARY KEY not null,
    name text not null CHECK (name <> ''),
    breed dog_breed not null,
    birth_date timestamp with time zone not null,
    owner_id text not null CHECK (owner_id <> '')
);
