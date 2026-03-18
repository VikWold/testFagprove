CREATE TABLE IF NOT EXISTS public.ehf (
	ehf_id UUID PRIMARY KEY,
	file_Name varchar(255),
	customer_Id int,
	supplier_Id int,
	invoice_No varchar(20),
	buyer_Reference char(6),
	issue_Date date,
	due_Date date,
	currency char(3),
	amount NUMERIC(15, 2)
);