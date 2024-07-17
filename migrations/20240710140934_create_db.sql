-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cities (  id serial primary key,
                                     name varchar,
                                     country varchar,
                                     latitude float,
                                     longitude float
);
CREATE INDEX IF NOT EXISTS idx_cities_name ON cities(name);
CREATE INDEX IF NOT EXISTS idx_country ON cities(country);

CREATE INDEX IF NOT EXISTS idx_cities_name ON cities(name);
CREATE INDEX IF NOT EXISTS idx_date ON cities(country);

CREATE TABLE IF NOT EXISTS weather ( id serial primary key,
                                     city_name varchar,
                                     temperature float,
                                     date date,
                                     city_id int,
                                     full_forecast jsonb,
                                     foreign key (city_id) references cities (id) ON DELETE CASCADE

);
CREATE INDEX IF NOT EXISTS idx_weather_date ON weather(date);

CREATE TABLE IF NOT EXISTS users ( id serial primary key,
                                     login varchar unique,
                                     password varchar
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop table cities, weather
-- +goose StatementEnd
