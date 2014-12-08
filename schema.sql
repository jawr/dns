DROP TABLE IF EXISTS tld CASCADE;
DROP SEQUENCE IF EXISTS tld_id CASCADE;
CREATE TABLE tld (
	id SERIAL,
	name VARCHAR(20),
	PRIMARY KEY (id),
	UNIQUE (name)
);

DROP TABLE IF EXISTS domain_base CASCADE;
DROP SEQUENCE IF EXISTS domain_base_id CASCADE;
CREATE TABLE domain_base (
	id SERIAL,
	name VARCHAR(255),
	tld INT NOT NULL references tld(id),
	PRIMARY KEY (id)
);
CREATE OR REPLACE VIEW domain AS SELECT * FROM domain_base;

DROP TABLE IF EXISTS record_type CASCADE;
DROP SEQUENCE IF EXISTS record_type_id CASCADE;
CREATE TABLE record_type (
	id SERIAL,
	name VARCHAR(10),
	PRIMARY KEY (id),
	UNIQUE (name)
);

DROP TABLE IF EXISTS record CASCADE;
DROP SEQUENCE IF EXISTS record_id CASCADE;
CREATE TABLE record (
	id SERIAL,
	domain INT NOT NULL references domain_base(id),
	name VARCHAR(255),
	args json,
	record_type INT NOT NULL references record_type(id),
	hash VARCHAR(32),
	added TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	active BOOLEAN DEFAULT TRUE,
	history json,
	PRIMARY KEY (id)
);

-- Only use this to populate tld table as it creates a
-- partitions the domain table
CREATE OR REPLACE FUNCTION insert_tld(VARCHAR)
	RETURNS INT 
	LANGUAGE plpgsql
	AS $$
	DECLARE
		tld_id INT;
		create_sql TEXT;
	BEGIN
		INSERT INTO tld (name) VALUES ($1) RETURNING id INTO tld_id;
		create_sql := 'CREATE TABLE domain__' || tld_id::text || ' ( ' ||
			'CHECK (tld = ' || tld_id || ' ), ' ||
			'UNIQUE (name)' ||
	       ') INHERITS (domain_base)';
		EXECUTE create_sql;
		RETURN tld_id;
	EXCEPTION WHEN UNIQUE_VIOLATION THEN
		SELECT id INTO tld_id FROM tld WHERE name = $1;
		RETURN tld_id;
	END;
	$$;

CREATE OR REPLACE FUNCTION insert_domain() 
	RETURNS TRIGGER
	LANGUAGE plpgsql
       	AS $$
	DECLARE
		insert_sql TEXT;
	BEGIN
		IF NEW.id IS NULL THEN
			NEW.id := nextval('domain_base_id_seq');
		END IF;
		insert_sql := 'INSERT INTO domain__' || NEW.tld::text || 
			' (name, tld) VALUES (' || quote_literal(NEW.name) || ',' || NEW.tld || ')';
		EXECUTE insert_sql;
		RETURN NEW;
	EXCEPTION WHEN UNIQUE_VIOLATION THEN
		RETURN NEW;
	END;
	$$;

CREATE TRIGGER insert_domain_in_to_partition INSTEAD OF INSERT ON domain
	FOR EACH ROW EXECUTE PROCEDURE insert_domain();

-- Only use this to populate record_type table as it creates a
-- partitions the record table
CREATE OR REPLACE FUNCTION insert_record_type(VARCHAR)
	RETURNS INT
	LANGUAGE plpgsql
	AS $$
	DECLARE
		rt_id INT;
		create_sql TEXT;
	BEGIN
		INSERT INTO record_type (name) VALUES ($1) RETURNING id INTO rt_id;
		create_sql := 'CREATE TABLE record__' || rt_id::text || ' ( ' ||
			'CHECK (record_type = ' || rt_id || ' ), ' ||
			'UNIQUE (hash)' ||
	       ') INHERITS (record)';
		EXECUTE create_sql;
		RETURN rt_id;
	EXCEPTION WHEN UNIQUE_VIOLATION THEN
		SELECT id INTO rt_id FROM record_type WHERE name = $1;
		RETURN rt_id;
	END;
	$$;

CREATE OR REPLACE FUNCTION insert_record() 
	RETURNS TRIGGER
	LANGUAGE plpgsql
       	AS $$
	DECLARE
		insert_sql TEXT;
	BEGIN
		insert_sql := 'INSERT INTO record__' || NEW.record_type::text || 
			' (domain, name, args, record_type, hash) VALUES (' || 
				NEW.domin || ',' ||
				quote_literal(NEW.name) || ',' ||
				quote_literal(NEW.args) || ',' ||
				NEW.record_type || ',' ||
				quote_literal(NEW.hash) || ')';
		EXECUTE insert_sql;
		RETURN NULL;
	EXCEPTION WHEN UNIQUE_VIOLATION THEN
		RETURN NULL;
	END;
	$$;

CREATE TRIGGER insert_record_in_to_partition BEFORE INSERT ON record
	FOR EACH ROW EXECUTE PROCEDURE insert_record();

