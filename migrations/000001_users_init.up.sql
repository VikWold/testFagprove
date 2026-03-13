CREATE TABLE IF NOT EXISTS public.users (
	id UUID PRIMARY KEY,
	username char(6) NOT NULL UNIQUE,
	password varchar(255) NOT NULL,
	createdAt TIMESTAMP,
	lastUpdated TIMESTAMP DEFAULT NOW()
);