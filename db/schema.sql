
\connect temba_latest

CREATE TABLE fcapp_orgunits (
    id SERIAL NOT NULL PRIMARY KEY,
    uid VARCHAR(11) NOT NULL,
    name VARCHAR(230) NOT NULL,
    code VARCHAR(50) NOT NULL DEFAULT '',
    shortname VARCHAR(50) DEFAULT '',
    parentid BIGINT,
    path VARCHAR(255),
    hierarchylevel INTEGER,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION get_ou_parent(_id INT) RETURNS INTEGER AS
$$
    DECLARE
    parentuid TEXT;
    parentid INT := 0;
    BEGIN
        SELECT split_part(path, '/', hierarchylevel) INTO parentuid
            FROM fcapp_orgunits WHERE id = _id;
        IF FOUND THEN
            SELECT id INTO parentid FROM fcapp_orgunits WHERE uid = parentuid;
            RETURN parentid;
        END IF;

        RETURN 0;
    END;
$$
LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION count_ou_children(_id INT) RETURNS INT
AS
$$
    DECLARE
        ret INT := 0;
    BEGIN
        SELECT count(*) INTO ret FROM fcapp_orgunits WHERE parentid = _id;
        RETURN ret;
    END
$$
LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION add_node(_parent INT, _name TEXT) RETURNS BOOLEAN AS
$$
    DECLARE
    parent_path TEXT;
    _level INT;
    _uid TEXT;
    BEGIN
        SELECT path, hierarchylevel + 1, gen_code() INTO parent_path, _level, _uid WHERE id = _parent;
        IF FOUND THEN
            INSERT INTO fcapp_orgunist(uid, name, parentid, path, hierarchylevel)
            VALUES(_uid, _name, parent_path || '/' || _uid, _level);
        END IF;
    END
$$ LANGUAGE 'plpgsql';
-- UPDATE fcapp_orgunits set parentid = get_ou_parent(id) WHERE id > 1;
-- UPDATE fcapp_orgunits set name = replace(name, ' District', '') where hierarchylevel = 3
-- UPDATE fcapp_orgunits set name = replace(name, ' Subcounty', '') where hierarchylevel=4;

CREATE TABLE fcapp_user_roles (
    id SERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    descr text DEFAULT '',
    cdate TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX fcapp_user_roles_idx1 ON fcapp_user_roles(name);

CREATE TABLE fcapp_user_role_permissions (
    id SERIAL NOT NULL PRIMARY KEY,
    user_role INTEGER NOT NULL REFERENCES fcapp_user_roles ON DELETE CASCADE ON UPDATE CASCADE,
    Sys_module TEXT NOT NULL, -- the name of the module - defined above this level
    sys_perms VARCHAR(16) NOT NULL,
    unique(sys_module,user_role)
);

CREATE TABLE fcapp_users (
    id bigserial NOT NULL PRIMARY KEY,
    user_role  INTEGER NOT NULL REFERENCES fcapp_user_roles ON DELETE RESTRICT ON UPDATE CASCADE, --(call agents, admin, service providers)
    firstname TEXT NOT NULL DEFAULT '',
    lastname TEXT NOT NULL DEFAULT '',
    username TEXT NOT NULL UNIQUE,
    telephone TEXT NOT NULL DEFAULT '', -- acts as the username
    password TEXT NOT NULL, -- blowfish hash of password
    email TEXT NOT NULL DEFAULT '',
    allowed_ips TEXT NOT NULL DEFAULT '127.0.0.1;::1', -- semi-colon separated list of allowed ip masks
    denied_ips TEXT NOT NULL DEFAULT '', -- semi-colon separated list of denied ip masks
    failed_attempts TEXT DEFAULT '0/'||to_char(now(),'yyyymmdd'),
    transaction_limit TEXT DEFAULT '0/'||to_char(now(),'yyyymmdd'),
    is_active BOOLEAN NOT NULL DEFAULT 't',
    is_system_user BOOLEAN NOT NULL DEFAULT 'f',
    last_login TIMESTAMPTZ,
    last_passwd_update TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX fcapp_users_idx1 ON fcapp_users(telephone);
CREATE INDEX fcapp_users_idx2 ON fcapp_users(username);

CREATE TABLE fcapp_audit_log (
        id BIGSERIAL NOT NULL PRIMARY KEY,
        logtype VARCHAR(32) NOT NULL DEFAULT '',
        actor TEXT NOT NULL,
        action text NOT NULL DEFAULT '',
        remote_ip INET,
        detail TEXT NOT NULL,
        created_by INTEGER REFERENCES fcapp_users(id), -- like actor id
        created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE fcapp_dhis2_indicator_mapping(
        id BIGSERIAL NOT NULL PRIMARY KEY,
        form TEXT NOT NULL DEFAULT '',
        cmd TEXT NOT NULL DEFAULT '',
        slug TEXT NOT NULL DEFAULT '',
        form_order INTEGER,
        shortname TEXT NOT NULL DEFAULT '',
        description TEXT NOT NULL DEFAULT '',
        dataset TEXT NOT NULL DEFAULT '',
        dataelement TEXT NOT NULL DEFAULT '',
        category_combo TEXT NOT NULL DEFAULT '',
        category_option_combo TEXT NOT NULL DEFAULT '',
        created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX fcapp_au_idx1 ON fcapp_audit_log(created);
CREATE INDEX fcapp_au_idx2 ON fcapp_audit_log(logtype);
CREATE INDEX fcapp_au_idx4 ON fcapp_audit_log(action);

INSERT INTO fcapp_user_roles(name, descr)
VALUES('Administrator','For the Administrators'), ('API User', 'For the API Users');

INSERT INTO fcapp_user_role_permissions(user_role, sys_module,sys_perms)
VALUES
        ((SELECT id FROM fcapp_user_roles WHERE name ='Administrator'),'Users','rw');

INSERT INTO fcapp_users(firstname,lastname,username,telephone,password,email,user_role,is_system_user)
VALUES
        ('Samuel','Sekiwere','admin', '+256753475676', crypt('admin',gen_salt('bf')),'sekiskylink@gmail.com',
        (SELECT id FROM fcapp_user_roles WHERE name ='Administrator'),'t'),
        ('Ivan','Muguya','ivan', '+256756253430', crypt('ivan',gen_salt('bf')),'ivanupsons@gmail.com',
        (SELECT id FROM fcapp_user_roles WHERE name ='API User'),'t');

CREATE OR REPLACE FUNCTION public.gen_code()
 RETURNS text
 LANGUAGE plpython3u
AS $function$
import string
import random
from uuid import uuid4


def id_generator(size=10, chars=string.ascii_lowercase + string.ascii_uppercase + string.digits):
    return random.choice(string.ascii_uppercase) + ''.join(random.choice(chars) for _ in range(size))

return id_generator()
$function$;

CREATE OR REPLACE FUNCTION fcapp_has_msisdn(contactid INT) RETURNS BOOLEAN AS
$delim$
    DECLARE
        c_id INTEGER;
    BEGIN
        SELECT id INTO c_id FROM contacts_contacturn WHERE contact_id = contactid;
        IF FOUND THEN
            RETURN TRUE;
        END IF;
        RETURN FALSE;
    END;
$delim$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION fcapp_has_field(contactid INT, field TEXT) RETURNS BOOLEAN AS
$delim$
    -- checks whether an contact has a passed field
    DECLARE
        res BOOLEAN := FALSE;
        field_uuid TEXT;
    BEGIN
        SELECT uuid INTO field_uuid FROM contacts_contactfield WHERE label = field;
        IF FOUND THEN
            SELECT 
                CASE WHEN fields ? field_uuid AND length(fields->field_uuid->>'text') > 0 THEN
                    TRUE ELSE FALSE END INTO res
            FROM 
                contacts_contact 
            WHERE
                id = contactid;
            RETURN res;
        END IF;
        RETURN FALSE;
    END;
$delim$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION fcapp_contactfield_value(contactid INT, field TEXT) RETURNS TEXT AS
$delim$
    DECLARE
    field_uuid TEXT;
    res TEXT;
    orgid INT;
    BEGIN
        SELECT org_id INTO orgid FROM contacts_contact WHERE id = contactid;
        IF FOUND THEN
            SELECT uuid INTO field_uuid FROM contacts_contactfield WHERE label = field AND org_id = orgid;
            IF FOUND THEN
                SELECT 
                    fields->field_uuid->>'text' INTO res
                FROM 
                    contacts_contact
                WHERE
                    id = contactid;
                RETURN res;
            END IF;
        END IF;
        RETURN res;
    END;
$delim$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION fcapp_delete_contactfield(contactid INT, field TEXT) RETURNS BOOLEAN AS
$delim$
DECLARE
    field_uuid TEXT;
    BEGIN
        SELECT uuid INTO field_uuid FROM contacts_contactfield WHERE label = field;
        IF FOUND THEN
            UPDATE contacts_contact SET fields = fields - field_uuid WHERE id = contactid;
            RETURN TRUE;
        END IF;
        RETURN FALSE;
    END;
$delim$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION fcapp_delete_contactfield_by_id(contactid INT, field_id INT) RETURNS BOOLEAN AS
$delim$
DECLARE
    field_uuid TEXT;
    BEGIN
        SELECT uuid INTO field_uuid FROM contacts_contactfield WHERE id = field_id;
        IF FOUND THEN
            UPDATE contacts_contact SET fields = fields - field_uuid WHERE id = contactid;
            RETURN TRUE;
        END IF;
        RETURN FALSE;
    END;
$delim$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION fcapp_get_secondary_receivers(contact text, 
    OUT contact_id int, OUT name text, OUT uuid text,
    OUT msisdn text, OUT contact_field int, OUT has_msisdn BOOLEAN, OUT has_hoh_msisdn BOOLEAN)
    RETURNS SETOF record
AS $$
    WITH t as (
    SELECT
        id, name, uuid,
        fcapp_contactfield_value(id, 'HoH MSISDN') as msisdn,
        (SELECT id FROM contacts_contactfield WHERE label= 'HoH MSISDN') as contact_field
    FROM contacts_contact
    UNION
        SELECT
        id, name, uuid,
            fcapp_contactfield_value(id, 'SecReceiver MSISDN') as msisdn,
            (SELECT id FROM contacts_contactfield WHERE label= 'SecReceiver MSISDN') as contact_field
        FROM contacts_contact
        
    ) 
    SELECT 
        id AS contact_id, name, uuid, msisdn, contact_field,
        fcapp_has_msisdn(id) AS has_msisdn, fcapp_has_field(id, 'HoH MSISDN') AS has_hoh_msisdn
    FROM 
        t 
    WHERE 
        substring(reverse(msisdn), 0, 9) = substring(reverse(contact), 0, 9);
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION fcapp_get_registered_contact_details(contact text, registred_by text, OUT contact_id int, OUT name text, OUT uuid text,
    OUT msisdn text, OUT contact_field int, OUT has_msisdn BOOLEAN)
    RETURNS SETOF record
AS $$
    WITH t AS
        (SELECT contact_id, string_value, contact_field_id FROM values_value
            WHERE
            contact_field_id IN (SELECT id FROM contacts_contactfield WHERE label IN('Registered By'))
            AND substring(reverse(string_value), 0, 9) = substring(reverse(registred_by), 0, 9)
        )
            SELECT a.id, a.name, a.uuid, t.string_value, t.contact_field_id,
                fcapp_has_msisdn(a.id) AS has_msisdn
            FROM contacts_contact a, t, contacts_contacturn b
            WHERE t.contact_id = a.id AND a.id = b.contact_id AND substring(reverse(b.path), 0, 9) = substring(reverse(contact), 0, 9);
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION fcapp_get_registered_contact_details(contact text, registred_by text, OUT contact_id int, OUT name text, OUT uuid text,
    OUT msisdn text, OUT contact_field int, OUT has_msisdn BOOLEAN)
    RETURNS SETOF record
AS $$
$$ LANGUAGE SQL;

CREATE TABLE fcapp_flow_data(
    id BIGSERIAL PRIMARY KEY NOT NULL,
    msisdn TEXT NOT NULL DEFAULT '',
    contact_uuid TEXT NOT NULL DEFAULT '',
    district INTEGER REFERENCES fcapp_locations(id),
    facility TEXT NOT NULL DEFAULT '',
    facilityuid TEXT NOT NULL DEFAULT '',
    subcounty TEXT,
    parish TEXT,
    village TEXT,
    report_type VARCHAR(16),
    week VARCHAR(8),
    month VARCHAR(8),
    quarter TEXT NOT NULL DEFAULT '',
    year INTEGER NOT NULL,
    "values" JSONB,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX fcapp_flow_data_contact_uuid ON fcapp_flow_data(contact_uuid);
CREATE INDEX fcapp_flow_data_district ON fcapp_flow_data(district);
CREATE INDEX fcapp_flow_data_msisdn ON fcapp_flow_data(msisdn);
CREATE INDEX fcapp_flow_data_created ON fcapp_flow_data(created);
CREATE INDEX fcapp_flow_data_updated ON fcapp_flow_data(updated);
CREATE INDEX fcapp_flow_data_report_type ON fcapp_flow_data(report_type);


CREATE OR REPLACE FUNCTION add_node(treeid INT, node_name TEXT, p_id INT) RETURNS BOOLEAN AS --p_id = id of node where to add
$delim$
    DECLARE
    new_lft INTEGER;
    lvl INTEGER;
    dummy TEXT;
    node_type INTEGER;
    child_type INTEGER;
    BEGIN
        IF node_name = '' THEN
            RAISE NOTICE 'Node name cannot be empty string';
            RETURN FALSE;
        END IF;
        SELECT level INTO lvl FROM fcapp_locationtype WHERE id = (SELECT type_id FROM fcapp_locations WHERE id = p_id);
        IF NOT FOUND THEN
            RAISE EXCEPTION 'Cannot add node: failed to find level';
        END IF;
        SELECT rght, type_id INTO new_lft, node_type FROM fcapp_locations WHERE id =  p_id AND tree_id = treeid;
        IF NOT FOUND THEN
            RAISE EXCEPTION 'No such node id= % ', p_id;
        END IF;

        SELECT id INTO child_type FROM fcapp_locationtype WHERE level = lvl + 1 AND tree_id = tree_id;
        IF NOT FOUND THEN
            RAISE EXCEPTION 'You cannot add to root node';
        END IF;

        SELECT name INTO dummy FROM fcapp_locations WHERE name = node_name
            AND tree_id = treeid AND type_id = child_type AND tree_parent_id = p_id;
        IF FOUND THEN
            RAISE NOTICE 'Node already exists : %', node_name;
            RETURN FALSE;
        END IF;

        UPDATE fcapp_locations SET lft = lft + 2 WHERE lft > new_lft AND tree_id = treeid;
        UPDATE fcapp_locations SET rght = rght + 2 WHERE rght >= new_lft AND tree_id = treeid;
        INSERT INTO fcapp_locations (name, lft, rght, tree_id,type_id, tree_parent_id)
        VALUES (node_name, new_lft, new_lft+1, treeid, child_type, p_id);
        RETURN TRUE;
    END;
$delim$ LANGUAGE plpgsql;
