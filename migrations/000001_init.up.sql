CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE transports (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    can_be_rented BOOLEAN NOT NULL DEFAULT TRUE,
    transport_type TEXT NOT NULL CHECK (transport_type IN ('Car', 'Bike', 'Scooter')),
    model TEXT NOT NULL,
    color TEXT NOT NULL,
    identifier TEXT NOT NULL UNIQUE,
    description TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    minute_price BIGINT,
    day_price BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rents (
    id BIGSERIAL PRIMARY KEY,
    transport_id BIGINT NOT NULL REFERENCES transports(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    time_start TIMESTAMPTZ NOT NULL,
    time_end TIMESTAMPTZ,
    price_of_unit BIGINT NOT NULL,
    price_type TEXT NOT NULL CHECK (price_type IN ('Minutes', 'Days')),
    final_price BIGINT,
);
