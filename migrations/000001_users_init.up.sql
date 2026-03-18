CREATE TABLE IF NOT EXISTS public.users (
	id UUID PRIMARY KEY,
	username char(6) NOT NULL UNIQUE,
	password varchar(255) NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	last_updated TIMESTAMP DEFAULT NOW()
);