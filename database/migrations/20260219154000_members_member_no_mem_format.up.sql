-- Normalize all member_no values to role-based format using member UUID suffix
-- Admin: ADM-XXXXXXXX, Customer: MEM-XXXXXXXX
UPDATE members
SET member_no =
	CASE
		WHEN LOWER(COALESCE(role::text, '')) = 'admin'
			THEN 'ADM-' || UPPER(RIGHT(REPLACE(id::text, '-', ''), 8))
		ELSE 'MEM-' || UPPER(RIGHT(REPLACE(id::text, '-', ''), 8))
	END;
