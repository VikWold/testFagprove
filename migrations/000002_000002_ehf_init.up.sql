CREATE TABLE IF NOT EXISTS public.ehf (
	id UUID PRIMARY KEY,
	fileName varchar(255),
	customerId int,
	supplierId int,
	invoiceNo varchar(20),
	buyerReference char(6),
	issueDate date,
	dueDate date,
	currency char(3),
	amount NUMERIC(15, 2)
);