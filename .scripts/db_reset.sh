#!/bin/bash

create_database() {
  local db_name=$1
  local collection_name=$2

  mongosh --eval "use $db_name; db.createCollection('$collection_name');"
}

insert_sample_data() {
  local db_name=$1
  local collection_name=$2
  local data=$3

  mongosh $db_name --eval "db.$collection_name.insert($data)"
}

create_database "test" "users"
create_database "test" "movies"

insert_sample_data "test" "users" '{"name": "John Doe", "email": "john.doe@example.com"}'
insert_sample_data "test" "movies" '{"title": "Inception", "year": 2010, "genre": "Sci-Fi"}'

echo "Databases and collections created successfully."