DROP TABLE IF EXISTS tld CASCADE;
DROP SEQUENCE IF EXISTS tld_id CASCADE;
CREATE TABLE tld (
	id SERIAL,
	name VARCHAR(255),
	PRIMARY KEY (id),
	UNIQUE (name)
);

DROP TABLE IF EXISTS domain CASCADE;
CREATE TABLE domain (
	uuid UUID,
	name VARCHAR(255),
	tld INT NOT NULL references tld(id),
	PRIMARY KEY (uuid)
);

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
	uuid UUID,
	domain UUID NOT NULL references domain(uuid),
	name VARCHAR(255),
	args jsonb,
	record_type INT NOT NULL references record_type(id),
	parser_date DATE NOT NULL,
	parser INT DEFAULT 0,
	added TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (uuid)
);

DROP TABLE IF EXISTS whois CASCADE;
DROP SEQUENCE IF EXISTS whois_id CASCADE;
CREATE TABLE whois (
	id SERIAL,
	domain UUID,
	data jsonb,
	raw_whois jsonb,
	contacts jsonb,
	emails jsonb,
	added TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	uuid UUID,
	PRIMARY KEY (id)
);

DROP TABLE IF EXISTS parser CASCADE;
DROP SEQUENCE IF EXISTS parser_id CASCADE;
CREATE TABLE parser (
	id SERIAL,
	filename VARCHAR(150),
	started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	finished_at TIMESTAMP WITH TIME ZONE,
	parser_date DATE NOT NULL,
	tld INT NOT NULL references tld(id),
	logs jsonb,
	PRIMARY KEY (id),
	UNIQUE (filename)
);

DROP TABLE IF EXISTS interval CASCADE;
DROP SEQUENCE IF EXISTS interval_id CASCADE;
CREATE TABLE interval (
	id SERIAL,
	value VARCHAR(150),
	PRIMARY KEY (id),
	UNIQUE (value)
);

DROP TABLE IF EXISTS watcher CASCADE;
DROP SEQUENCE IF EXISTS watcher_id CASCADE;
CREATE TABLE watcher (
	id SERIAL,
	domain UUID,
	added TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	interval INT NOT NULL references interval(id),
	logs jsonb DEFAULT '[]'::jsonb,
	PRIMARY KEY (id),
	UNIQUE (domain)
);

CREATE OR REPLACE FUNCTION update_updated_column()
	RETURNS TRIGGER
	LANGUAGE plpgsql
	AS $$
	BEGIN
		IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
			NEW.updated = now();
			RETURN NEW;
		ELSE
			RETURN OLD;
		END IF;
	END;
	$$;


CREATE TRIGGER update_watcher_updated BEFORE UPDATE ON watcher
	FOR EACH ROW EXECUTE PROCEDURE update_updated_column();

CREATE OR REPLACE FUNCTION insert_interval(VARCHAR)
	RETURNS INT
	LANGUAGE plpgsql
	AS $$
	DECLARE
		interval_id INT;
	BEGIN
		INSERT INTO interval (value) VALUES ($1) RETURNING id INTO interval_id;
		RETURN interval_id;
	EXCEPTION WHEN UNIQUE_VIOLATION THEN
		SELECT id INTO interval_id FROM interval WHERE value = $1;
		RETURN interval_id;
	END;
	$$;

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
	       ') INHERITS (domain)';
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
		insert_sql := 'INSERT INTO domain__' || NEW.tld::text || 
			' (uuid, name, tld) VALUES (' || quote_literal(NEW.uuid) || ',' || quote_literal(NEW.name) || ',' || NEW.tld || ')';
		EXECUTE insert_sql;
		RETURN NULL;
	EXCEPTION WHEN UNIQUE_VIOLATION THEN
		RETURN NULL;
	END;
	$$;

CREATE TRIGGER insert_domain_in_to_partition BEFORE INSERT ON domain
	FOR EACH ROW EXECUTE PROCEDURE insert_domain();

CREATE OR REPLACE FUNCTION insert_whois()
	RETURNS TRIGGER
	LANGUAGE plpgsql
	AS $$
	BEGIN
		PERFORM 1 FROM whois WHERE uuid = NEW.uuid AND added > now() - interval '1 day';
		IF NOT FOUND THEN
			RETURN NEW;
		END IF;
		RETURN NULL;
	END;
	$$;

CREATE TRIGGER insert_whois_check BEFORE INSERT ON whois
	FOR EACH ROW EXECUTE PROCEDURE insert_whois();

-- Only use this to populate record_type table as it creates a
-- partitions the record table
CREATE OR REPLACE FUNCTION insert_record_type(VARCHAR)
	RETURNS INT
	LANGUAGE plpgsql
	AS $$
	DECLARE
		rt_id INT;
	BEGIN
		INSERT INTO record_type (name) VALUES ($1) RETURNING id INTO rt_id;
		RETURN rt_id;
	EXCEPTION WHEN UNIQUE_VIOLATION THEN
		SELECT id INTO rt_id FROM record_type WHERE name = $1;
		RETURN rt_id;
	END;
	$$;

-- Takes record_type.name and tld.id as input
CREATE OR REPLACE FUNCTION ensure_record_table(VARCHAR, INT)
	RETURNS INT
	LANGUAGE plpgsql
	AS $$
	DECLARE
		rt_id INT;
		create_sql TEXT;
	BEGIN
		SELECT insert_record_type($1) INTO rt_id;
		create_sql := 'CREATE TABLE record__' || rt_id::text || '_' || $2::text || ' ( ' ||
			'CHECK (record_type = ' || rt_id || ' ), ' ||
			'PRIMARY KEY (uuid)' ||
	       ') INHERITS (record)';
		EXECUTE create_sql;
		RETURN rt_id;
	EXCEPTION WHEN duplicate_table THEN
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
			' (uuid, domain, name, args, record_type, parser_date) VALUES (' || 
				quote_literal(NEW.uuid) || ',' ||
				quote_literal(NEW.domain) || ',' ||
				quote_literal(NEW.name) || ',' ||
				quote_literal(NEW.args) || ',' ||
				NEW.record_type || ',' ||
				quote_literal(NEW.parser_date) || ')';
		EXECUTE insert_sql;
		RETURN NULL;
	EXCEPTION WHEN UNIQUE_VIOLATION THEN
		RETURN NULL;
	END;
	$$;

CREATE TRIGGER insert_record_in_to_partition BEFORE INSERT ON record
	FOR EACH ROW EXECUTE PROCEDURE insert_record();
