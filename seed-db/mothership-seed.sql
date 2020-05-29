-- NOTE: This file is no longer a pure dump. Rather, it is manually updated.

--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.3
-- Dumped by pg_dump version 9.6.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: deployment; Type: SCHEMA; Schema: -; Owner: caas
--

CREATE USER caas;

CREATE SCHEMA deployment;

ALTER SCHEMA deployment OWNER TO caas;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- Name: cloud_id; Type: DOMAIN; Schema: public; Owner: caas
--

CREATE DOMAIN cloud_id AS character varying(16) NOT NULL;


ALTER DOMAIN cloud_id OWNER TO caas;

--
-- Name: physical_cluster_id; Type: DOMAIN; Schema: public; Owner: caas
--

CREATE DOMAIN physical_cluster_id AS character varying(32) NOT NULL;


ALTER DOMAIN physical_cluster_id OWNER TO caas;


--
-- Name: logical_cluster_id; Type: DOMAIN; Schema: public; Owner: caas
--

CREATE DOMAIN logical_cluster_id AS character varying(32) NOT NULL;


ALTER DOMAIN logical_cluster_id OWNER TO caas;


--
-- Name: physical_cluster_version; Type: DOMAIN; Schema: public; Owner: caas
--

CREATE DOMAIN physical_cluster_version AS character varying(32) NOT NULL;


ALTER DOMAIN physical_cluster_version OWNER TO caas;

--
-- Name: cp_component_id; Type: DOMAIN; Schema: public; Owner: caas
--

CREATE DOMAIN cp_component_id AS character varying(32) NOT NULL;


ALTER DOMAIN cp_component_id OWNER TO caas;

--
-- Name: environment_id; Type: DOMAIN; Schema: public; Owner: caas
--

CREATE DOMAIN environment_id AS character varying(32) NOT NULL;


ALTER DOMAIN environment_id OWNER TO caas;

--
-- Name: hash_function; Type: TYPE; Schema: public; Owner: caas
--

CREATE TYPE hash_function AS ENUM (
    'none',
    'bcrypt'
);


ALTER TYPE hash_function OWNER TO caas;

--
-- Name: k8s_cluster_id; Type: DOMAIN; Schema: public; Owner: caas
--

CREATE DOMAIN k8s_cluster_id AS character varying(32) NOT NULL;


ALTER DOMAIN k8s_cluster_id OWNER TO caas;

--
-- Name: network_id; Type: DOMAIN; Schema: public; Owner: caas
--

CREATE DOMAIN network_id AS character varying(32) NOT NULL;


ALTER DOMAIN network_id OWNER TO caas;

--
-- Name: region_id; Type: DOMAIN; Schema: public; Owner: caas
--

CREATE DOMAIN region_id AS character varying(32) NOT NULL;


ALTER DOMAIN region_id OWNER TO caas;

--
-- Name: sasl_mechanism; Type: TYPE; Schema: public; Owner: caas
--

CREATE TYPE sasl_mechanism AS ENUM (
    'PLAIN',
    'SCRAM-SHA-256',
    'SCRAM-SHA-512'
);


ALTER TYPE sasl_mechanism OWNER TO caas;

--
-- Name: account_id; Type: DOMAIN; Schema: public; Owner: caas
--

CREATE DOMAIN account_id AS character varying(32) NOT NULL;


ALTER DOMAIN account_id OWNER TO caas;

SET search_path = deployment, pg_catalog;

--
-- Name: account_num; Type: SEQUENCE; Schema: deployment; Owner: caas
--

CREATE SEQUENCE account_num
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE account_num OWNER TO caas;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: account; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE account (
    id public.account_id DEFAULT ('a-'::text || nextval('account_num'::regclass)) NOT NULL,
    name character varying(32) NOT NULL,
    config jsonb,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    deactivated boolean DEFAULT false NOT NULL,
    organization_id integer NOT NULL,
    internal boolean DEFAULT false NOT NULL
);


ALTER TABLE account OWNER TO caas;

CREATE UNIQUE INDEX account_name_is_unique ON deployment.account USING btree (name, organization_id) WHERE (deactivated = FALSE);

INSERT INTO account (id, name, organization_id) VALUES ('t0', 'Internal', 0);

--
-- Name: api_key; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE api_key (
    id integer NOT NULL,
    key character varying(128) NOT NULL,
    hashed_secret character varying(128) NOT NULL,
    hash_function public.hash_function DEFAULT 'bcrypt'::public.hash_function NOT NULL,
    sasl_mechanism public.sasl_mechanism DEFAULT 'PLAIN'::public.sasl_mechanism NOT NULL,
    cluster_id public.logical_cluster_id,
    deactivated boolean DEFAULT false NOT NULL,
    description text DEFAULT('') NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    user_id integer DEFAULT 0 NOT NULL
);


ALTER TABLE api_key OWNER TO caas;

--
-- Name: api_key_id_seq; Type: SEQUENCE; Schema: deployment; Owner: caas
--

CREATE SEQUENCE api_key_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE api_key_id_seq OWNER TO caas;

--
-- Name: api_key_id_seq; Type: SEQUENCE OWNED BY; Schema: deployment; Owner: caas
--

ALTER SEQUENCE api_key_id_seq OWNED BY api_key.id;


--
-- Name: cloud; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE cloud (
    id public.cloud_id NOT NULL,
    config jsonb,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    name TEXT DEFAULT '' NOT NULL
);


ALTER TABLE cloud OWNER TO caas;

--
-- Name: network_isolation_domain_num; Type: SEQUENCE; Schema: deployment; Owner: caas
--

CREATE SEQUENCE IF NOT EXISTS network_isolation_domain_num START WITH 1 INCREMENT 1 NO CYCLE NO MINVALUE NO MAXVALUE;

ALTER TABLE network_isolation_domain_num OWNER TO caas;

--
-- Name: network_isolation_domain; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE IF NOT EXISTS network_isolation_domain (
    id text PRIMARY KEY DEFAULT ('nid-' || nextval('network_isolation_domain_num')::text),
    description varchar(140) NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    deactivated timestamp without time zone DEFAULT NULL
);

ALTER TABLE network_isolation_domain OWNER TO caas;

--
-- Name: deployment_num; Type: SEQUENCE; Schema: deployment; Owner: caas
--

CREATE SEQUENCE IF NOT EXISTS deployment_num START WITH 1 INCREMENT 1 NO CYCLE NO MINVALUE NO MAXVALUE;

ALTER TABLE deployment_num OWNER TO caas;

--
-- Name: deployment; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE IF NOT EXISTS deployment (
    id                  text PRIMARY KEY DEFAULT ('deployment-' || nextval('deployment_num')::text),
    created             timestamp without time zone DEFAULT now() NOT NULL,
    modified            timestamp without time zone DEFAULT now() NOT NULL,
    deactivated         timestamp without time zone DEFAULT NULL,
    account_id          varchar(140) NOT NULL,
    network_access      jsonb,
    network_region_id   text DEFAULT NULL,
    sku                 varchar(140) NOT NULL,
    provider            jsonb DEFAULT '{}'::jsonb
);

ALTER TABLE deployment OWNER TO caas;

--
-- Name: physical_cluster_num; Type: SEQUENCE; Schema: deployment; Owner: caas
--

CREATE SEQUENCE physical_cluster_num
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;

  ALTER TABLE physical_cluster_num OWNER TO caas;

--
-- Name: physical_cluster; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE physical_cluster (
    id public.physical_cluster_id NOT NULL,
    k8s_cluster_id public.k8s_cluster_id,
    type character varying(32) NOT NULL,
    config jsonb,
    deactivated timestamp without time zone,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    status character varying(32) NOT NULL,
    is_schedulable boolean DEFAULT true NOT NULL,
    status_detail jsonb DEFAULT '{}'::jsonb NOT NULL,
    status_modified timestamp without time zone,
    status_received timestamp without time zone,
    last_initialized timestamp without time zone,
    last_deleted timestamp without time zone,
    network_isolation_domain_id text,
    sni_enabled bool DEFAULT false NOT NULL
);


ALTER TABLE physical_cluster OWNER TO caas;

ALTER TABLE deployment.physical_cluster ADD CONSTRAINT "fk-physical_cluster-network_isolation_domain" FOREIGN KEY ("network_isolation_domain_id") REFERENCES deployment.network_isolation_domain ("id") NOT VALID;

ALTER TABLE deployment.physical_cluster VALIDATE CONSTRAINT "fk-physical_cluster-network_isolation_domain";


CREATE SEQUENCE deployment.metadata_change_id_seq
    START WITH 2
    INCREMENT BY 1
    NO CYCLE
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE deployment.metadata_change_id_seq OWNER TO caas;

--
-- Name: logical_cluster_num; Type: SEQUENCE; Schema: deployment; Owner: caas
--

CREATE SEQUENCE logical_cluster_num
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;

ALTER TABLE logical_cluster_num OWNER TO caas;

--
-- Name: logical_cluster; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE logical_cluster (
    id public.logical_cluster_id NOT NULL,
    name character varying(64) NOT NULL,
    physical_cluster_id public.physical_cluster_id NOT NULL,
    type character varying(32) NOT NULL,
    account_id public.account_id,
    config jsonb,
    deactivated timestamp without time zone,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    status_detail jsonb NOT NULL default '{}'::jsonb,
    status_modified timestamp without time zone DEFAULT now() NOT NULL,
    deployment_id text,
    last_change_id bigint NOT NULL default 1
);

CREATE INDEX logical_cluster_last_change_id_idx on deployment.logical_cluster (last_change_id);

ALTER TABLE logical_cluster OWNER TO caas;

-- Function to capture inserts and updates to logical clusters (Create/Deactivate)
CREATE OR REPLACE FUNCTION deployment.update_last_change_id() RETURNS TRIGGER AS $body$
BEGIN
    -- Connect has status updates at the LC level which will needlessly create change records
    -- Filter out the stuff we don't need
    IF NEW.type in ('kafka', 'healthcheck', 'ksql') THEN
      NEW.last_change_id = nextval('deployment."metadata_change_id_seq"');
      NEW.modified = NOW();
    END IF;
    RETURN NEW;
END;
$body$
LANGUAGE plpgsql;

-- Trigger to capture inserts and updates to logical clusters (Create/Deactivate)
CREATE TRIGGER logical_cluster_last_change_id_trigger
 BEFORE INSERT OR UPDATE ON deployment.logical_cluster
 FOR EACH ROW
 WHEN (pg_trigger_depth() < 1) -- skip this trigger when row is modified by the other trigger
 EXECUTE PROCEDURE deployment.update_last_change_id();

--
-- Name: cp_component; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE cp_component (
    id public.cp_component_id NOT NULL,
    default_version public.physical_cluster_version DEFAULT '0.0.7'::character varying,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE cp_component OWNER TO caas;

--
-- Name: environment; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE environment (
    id public.environment_id NOT NULL,
    config jsonb,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE environment OWNER TO caas;

--
-- Name: k8s_cluster_num; Type: SEQUENCE; Schema: deployment; Owner: caas
--

CREATE SEQUENCE k8s_cluster_num
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE k8s_cluster_num OWNER TO caas;

--
-- Name: k8s_cluster; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE k8s_cluster (
    id public.k8s_cluster_id DEFAULT ('k8s-'::text || nextval('k8s_cluster_num'::regclass)) NOT NULL,
    config jsonb,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    deactivated timestamp without time zone,
    network_region_id text
);


ALTER TABLE k8s_cluster OWNER TO caas;

--
-- Name: organization; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE organization (
    id integer NOT NULL,
    name character varying(32) NOT NULL,
    deactivated boolean DEFAULT false NOT NULL,
    plan jsonb NOT NULL DEFAULT('{}'),
    saml jsonb NOT NULL DEFAULT('{}'),
    sso jsonb NOT NULL DEFAULT('{}'),
    marketplace jsonb NOT NULL DEFAULT('{}'),
    resource_id TEXT NOT NULL DEFAULT '',
    audit_log jsonb NOT NULL DEFAULT('{}'),
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE organization OWNER TO caas;

COPY organization (id, name, plan) FROM stdin;
0	Internal	{"billing": {"email": "caas-team@confluent.io", "method": "MANUAL", "interval": "MONTHLY", "accrued_this_cycle": "0", "stripe_customer_id": ""}, "tax_address": {"zip": "", "city": "", "state": "", "country": "", "street1": "", "street2": ""}, "product_level": "TEAM", "referral_code": ""}
\.


--
-- Name: organization_id_seq; Type: SEQUENCE; Schema: deployment; Owner: caas
--

CREATE SEQUENCE organization_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE organization_id_seq OWNER TO caas;

--
-- Name: organization_id_seq; Type: SEQUENCE OWNED BY; Schema: deployment; Owner: caas
--

ALTER SEQUENCE organization_id_seq OWNED BY organization.id;

--
-- Name: entitlement; Type: TABLE; Schema: deployment;  Owner: caas
--

CREATE TABLE entitlement (
    id integer NOT NULL,
    external_id character varying(100) NOT NULL,
    name character varying(100) NOT NULL,
    customer_id character varying(100) NOT NULL DEFAULT '',
    product_id character varying(100) NOT NULL,
    plan_id character varying(100) NOT NULL,
    state character varying(32) NOT NULL,
    external_state character varying(50) NOT NULL,
    usage_reporting_id character varying(100) NOT NULL,
    organization_id integer NOT NULL,
    deactivated boolean DEFAULT false NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    started timestamp without time zone,
    ended timestamp without time zone,
    config jsonb NOT NULL DEFAULT('{}'),
    state_description jsonb DEFAULT ('{}'::jsonb),
    offer_type character varying(32),
    parent_id integer check (parent_id != id)
);

ALTER TABLE entitlement OWNER TO caas;

--
-- Name: entitlement_id_seq; Type: SEQUENCE; Schema: deployment; Owner: caas
--

CREATE SEQUENCE entitlement_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE entitlement_id_seq OWNER TO caas;

--
-- Name: entitlement_id_seq; Type: SEQUENCE OWNED BY; Schema: deployment; Owner: caas
--

ALTER SEQUENCE entitlement_id_seq OWNED BY entitlement.id;

--
-- Name: entitlement_external_id_is_unique; Type: INDEX; Schema: deployment; Owner: caas
--

CREATE UNIQUE INDEX entitlement_external_id_is_unique ON entitlement USING btree (external_id) WHERE (deactivated = false);

--
-- Name: entitlement_customer_id; Type: INDEX; Schema: deployment; Owner: -
--

CREATE INDEX entitlement_customer_id ON deployment.entitlement USING btree (customer_id) WHERE (deactivated = false);

--
-- Name: entitlement_organization_id; Type: INDEX; Schema: deployment; Owner: -
--

CREATE INDEX entitlement_organization_id ON deployment.entitlement USING btree (organization_id) WHERE (deactivated = false);

--
-- Name: marketplace_listener_errors; Type: TABLE; Schema: deployment; Owner: -
--

CREATE TABLE deployment.marketplace_listener_errors (
    id integer NOT NULL,
    event_id character varying(100) NOT NULL,
    marketplace_partner character varying(36) NOT NULL,
    event_created timestamp without time zone DEFAULT now() NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    error jsonb DEFAULT '{}'::jsonb NOT NULL,
    data_payload jsonb DEFAULT '{}'::jsonb NOT NULL
);


--
-- Name: marketplace_listener_errors_id_seq; Type: SEQUENCE; Schema: deployment; Owner: -
--

CREATE SEQUENCE deployment.marketplace_listener_errors_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: marketplace_listener_errors_id_seq; Type: SEQUENCE OWNED BY; Schema: deployment; Owner: -
--

ALTER SEQUENCE deployment.marketplace_listener_errors_id_seq OWNED BY deployment.marketplace_listener_errors.id;

--
-- Name: marketplace_listener_errors id; Type: DEFAULT; Schema: deployment; Owner: -
--

ALTER TABLE ONLY deployment.marketplace_listener_errors ALTER COLUMN id SET DEFAULT nextval('deployment.marketplace_listener_errors_id_seq'::regclass);

--
-- Name: marketplace_listener_errors marketplace_listener_errors_pkey; Type: CONSTRAINT; Schema: deployment; Owner: -
--

ALTER TABLE ONLY deployment.marketplace_listener_errors
    ADD CONSTRAINT marketplace_listener_errors_pkey PRIMARY KEY (id);

--
-- Name: coupon; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE coupon (
    id TEXT NOT NULL,
    coupon_type INTEGER DEFAULT 0 NOT NULL,
    amount_off INTEGER DEFAULT 0 NOT NULL,
    percent_off INTEGER DEFAULT 0 NOT NULL,
    redeem_by TIMESTAMP WITHOUT TIME zone,
    times_redeemed INTEGER DEFAULT 0 NOT NULL,
    max_redemptions INTEGER DEFAULT 0 NOT NULL,
    duration_in_months INTEGER DEFAULT 0 NOT NULL,
    deactivated BOOL DEFAULT FALSE NOT NULL,
    created TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL,
    modified TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL
);

ALTER TABLE coupon OWNER TO caas;

--
-- Name: coupon_id_seq; Type: SEQUENCE; Schema: deployment; Owner: caas
--

CREATE SEQUENCE coupon_id_seq;

ALTER TABLE coupon_id_seq OWNER TO caas;

--
-- Name: coupon_id_seq; Type: SEQUENCE OWNED BY; Schema: deployment; Owner: caas
--

ALTER SEQUENCE coupon_id_seq OWNED BY coupon.id;

--
-- Name: event; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE event (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    resource_type integer,
    resource_id TEXT NOT null,
    action INTEGER NOT NULL,
    data jsonb NOT NULL DEFAULT('{}'),
    created TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL
);

CREATE INDEX IF NOT EXISTS event_organization_id_created_idx ON deployment.event (organization_id, created);

ALTER TABLE event OWNER TO caas;

--
-- Name: organization_eventer_state; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE organization_eventer_state (
    last_read_org_id INTEGER
);

ALTER TABLE organization_eventer_state OWNER TO caas;

--
-- Name: region; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE region (
    id public.region_id NOT NULL,
    cloud public.cloud_id,
    config jsonb,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    name TEXT DEFAULT '' NOT NULL
);


ALTER TABLE region OWNER TO caas;

--
-- Name: users; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE users (
    id integer NOT NULL,
    resource_id TEXT NOT NULL,
    email character varying(128) NOT NULL,
    service_name character varying(64) DEFAULT '' NOT NULL,
    service_description character varying(128) DEFAULT '' NOT NULL,
    service_account boolean DEFAULT false NOT NULL,
    first_name character varying(32) NOT NULL,
    last_name character varying(32) NOT NULL,
    password character varying(64) NOT NULL,
    organization_id integer NOT NULL,
    deactivated boolean DEFAULT false NOT NULL,
    sso jsonb DEFAULT '{}'::jsonb NOT NULL,
    internal boolean DEFAULT false NOT NULL,
    verified timestamp without time zone DEFAULT timestamp '1970-01-01 00:00:00.00000'  NOT NULL,
    password_changed timestamp without time zone DEFAULT now() NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    preferences jsonb DEFAULT '{}'::jsonb NOT NULL
);

CREATE UNIQUE INDEX users_email_one_active ON deployment.users USING btree (email) WHERE (deactivated = FALSE);
CREATE UNIQUE INDEX users_service_name_is_unique ON deployment.users USING btree (service_name, organization_id) WHERE (service_account = TRUE AND deactivated = FALSE);

ALTER TABLE users
  ADD CONSTRAINT users_resource_id_uniq UNIQUE (resource_id);

ALTER TABLE users OWNER TO caas;

INSERT INTO users (id, resource_id, email, first_name, last_name, password, organization_id) VALUES (0, 'u-000000', 'caas-team+internal@confluent.io', '', '', '$2a$10$P05oytzmNEfvpSWa0l.RpOa7vEund/Hdt0w88kaRLMIWvpSNGB33S', 0);

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: deployment; Owner: caas
--

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE users_id_seq OWNER TO caas;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: deployment; Owner: caas
--

ALTER SEQUENCE users_id_seq OWNED BY users.id;


--
-- Name: users_resource_id_seq; Type: SEQUENCE; Schema: deployment; Owner: caas
--

CREATE SEQUENCE users_resource_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE users_resource_id_seq OWNER TO caas;

--
-- Name: api_key id; Type: DEFAULT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY api_key ALTER COLUMN id SET DEFAULT nextval('api_key_id_seq'::regclass);

--
-- Name: billing_job; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE billing_job (
    id SERIAL PRIMARY KEY,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    month timestamp without time zone DEFAULT now() NOT NULL,
    status jsonb NOT NULL DEFAULT('{}'),
    charges jsonb NOT NULL DEFAULT('{}')
);


ALTER TABLE billing_job OWNER TO caas;

--
-- Name: roll; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE roll (
    id SERIAL PRIMARY KEY,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    deactivated timestamp without time zone DEFAULT NULL,
    status jsonb NOT NULL DEFAULT('{}'),
    request jsonb NOT NULL DEFAULT('{}'),
    clusters jsonb NOT NULL DEFAULT('{}'),
    operation integer NOT NULL DEFAULT 0
);

ALTER TABLE roll OWNER TO caas;

--
-- Name: secret; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE secret (
    id SERIAL PRIMARY KEY,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    deactivated timestamp without time zone DEFAULT NULL,
    type TEXT DEFAULT '' NOT NULL,
    config jsonb DEFAULT '{}' NOT NULL
);

ALTER TABLE secret OWNER TO caas;

--
-- Name: task; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE task (
    id SERIAL PRIMARY KEY,
    run_date timestamp without time zone DEFAULT now() NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    start_time timestamp without time zone DEFAULT now() NOT NULL,
    end_time timestamp without time zone DEFAULT now() NOT NULL,
    type integer NOT NULL,
    status integer NOT NULL,
    message text DEFAULT('') NOT NULL,
    sub_tasks jsonb NOT NULL DEFAULT('{}')
);

ALTER TABLE task OWNER TO caas;

--
-- Name: usage; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE usage (
    id SERIAL PRIMARY KEY,
    logical_cluster_id public.logical_cluster_id,
    month TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL,
    metrics jsonb NOT NULL DEFAULT('{}'),
    modified TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL,
    created TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL
);

ALTER TABLE usage OWNER TO caas;

--
-- Name: promo_code; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE promo_code (
                       id SERIAL PRIMARY KEY,
                       code VARCHAR(50) UNIQUE NOT NULL,
                       amount BIGINT NOT NULL,
                       organization_id INT,
                       code_validity_start_date TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL,
                       code_validity_end_date TIMESTAMP WITHOUT TIME zone NOT NULL,
                       credit_validity_days INT NOT NULL,
                       max_uses INT NOT NULL DEFAULT (1) CHECK (max_uses > 0),
                       is_enabled BOOLEAN NOT NULL,
                       created TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL,
                       created_by character varying(128) NOT NULL,
                       modified TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL,
                       modified_by character varying(128) NOT NULL
);

CREATE UNIQUE INDEX promo_code_code_index on deployment.promo_code (code);
CREATE INDEX promo_code_created_by_index on deployment.promo_code (created_by);

ALTER TABLE promo_code OWNER TO caas;

--
-- Name: promo_code_claim; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE promo_code_claim (
                            id SERIAL PRIMARY KEY,
                            promo_code_id INT NOT NULL,
                            credit_expiration TIMESTAMP WITHOUT TIME zone NOT NULL,
                            organization_id INT NOT NULL,
                            amount_remaining BIGINT NOT NULL,
                            created TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL,
                            created_by INT NOT NULL,
                            modified TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL,
                            FOREIGN KEY (promo_code_id) REFERENCES promo_code(id)
);

CREATE INDEX promo_code_claim_organization_id_index on deployment.promo_code_claim (organization_id);
CREATE INDEX promo_code_claim_promo_code_id_index on deployment.promo_code_claim (promo_code_id);
CREATE UNIQUE INDEX promo_code_claim_organization_and_code_index on deployment.promo_code_claim(organization_id, promo_code_id);

ALTER TABLE promo_code_claim OWNER TO caas;

--
-- Name: billing_order; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE billing_order (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER NOT NULL,
    universal_id VARCHAR,
    commit_total BIGINT DEFAULT 0 NOT NULL,
    prepaid_amount BIGINT DEFAULT 0 NOT NULL,
    created_date TIMESTAMP WITHOUT TIME ZONE,
    discount DECIMAL DEFAULT 0 NOT NULL,
    start_date TIMESTAMP WITHOUT TIME ZONE,
    end_date TIMESTAMP WITHOUT TIME ZONE,
    created TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL,
    modified TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL,
    currency VARCHAR(16) NOT NULL,
    status INTEGER,
    billing_cycle INTEGER DEFAULT 0 NOT NULL
);

ALTER TABLE billing_order OWNER TO caas;

--
-- Name: billing_order_organization_id_universal_id_uniq; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY deployment.billing_order ADD CONSTRAINT billing_order_organization_id_universal_id_uniq UNIQUE (organization_id, universal_id);

--
-- Name: price; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE price (
    id SERIAL PRIMARY KEY,
    price_table jsonb NOT NULL DEFAULT('{}'),
    kafka_prices jsonb NOT NULL DEFAULT('{}'),
    connect_prices jsonb NOT NULL DEFAULT('{}'),
    support_prices jsonb NOT NULL DEFAULT('{}'),
    multipliers jsonb NOT NULL DEFAULT('{}'),
    effective_date TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL,
    modified TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL,
    created TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL,
    type INTEGER DEFAULT 1 NOT NULL,
    order_universal_id VARCHAR
);

ALTER TABLE price OWNER TO caas;

INSERT INTO deployment.price (price_table, multipliers, effective_date)
VALUES ('{}', '{
  "aws": {
    "eu-west-1": 1,
    "eu-west-2": 1,
    "eu-west-3": 1,
    "sa-east-1": 1,
    "us-east-1": 1,
    "us-east-2": 1,
    "us-west-1": 1,
    "us-west-2": 1,
    "ap-south-1": 1,
    "ca-central-1": 1,
    "eu-central-1": 1,
    "ap-northeast-1": 1,
    "ap-northeast-2": 1,
    "ap-southeast-1": 1,
    "ap-southeast-2": 1
  },
  "gcp": {
    "us-east1": 1,
    "us-east4": 1,
    "us-west1": 1,
    "us-west2": 1,
    "asia-east1": 1,
    "asia-east2": 1,
    "asia-south1": 1,
    "us-central1": 1,
    "europe-west1": 1,
    "europe-west2": 1,
    "europe-west3": 1,
    "europe-west4": 1,
    "europe-north1": 1,
    "asia-northeast1": 1,
    "asia-southeast1": 1,
    "southamerica-east1": 1,
    "australia-southeast1": 1,
    "northamerica-northeast1": 1
  },
  "azure": {
    "eastus": 1,
    "eastus2": 1,
    "uksouth": 1,
    "westus2": 1,
    "centralus": 1,
    "japaneast": 1,
    "westeurope": 1,
    "northeurope": 1,
    "francecentral": 1,
    "southeastasia": 1
  }
}', '2015-02-01 00:00:00');

UPDATE deployment.price SET kafka_prices = '{
  "KafkaPartition": {
    "prices": {
      "aws:high:custom:internet:internet": 0,
      "aws:high:custom:peered-vpc:internet": 0,
      "aws:high:dedicated:internet:internet": 0,
      "aws:high:dedicated:peered-vpc:internet": 0,
      "aws:high:dedicated:private-link:internet": 0,
      "aws:high:dedicated:transit-gateway:internet": 0,
      "aws:high:standard:internet:internet": 0.008,
      "aws:high:standard_v2:internet:internet": 0.008,
      "aws:low:basic:internet:internet": 0.008,
      "aws:low:custom:internet:internet": 0,
      "aws:low:custom:peered-vpc:internet": 0,
      "aws:low:dedicated:internet:internet": 0,
      "aws:low:dedicated:peered-vpc:internet": 0,
      "aws:low:dedicated:private-link:internet": 0,
      "aws:low:dedicated:transit-gateway:internet": 0,
      "aws:low:standard:internet:internet": 0.008,
      "aws:low:standard_v2:internet:internet": 0.008,
      "azure:high:custom:internet:internet": 0,
      "azure:high:custom:peered-vpc:internet": 0,
      "azure:high:dedicated:internet:internet": 0,
      "azure:high:dedicated:peered-vpc:internet": 0,
      "azure:high:standard:internet:internet": 0.008,
      "azure:high:standard_v2:internet:internet": 0.008,
      "azure:low:basic:internet:internet": 0.008,
      "azure:low:custom:internet:internet": 0,
      "azure:low:custom:peered-vpc:internet": 0,
      "azure:low:dedicated:internet:internet": 0,
      "azure:low:dedicated:peered-vpc:internet": 0,
      "azure:low:standard:internet:internet": 0.008,
      "azure:low:standard_v2:internet:internet": 0.008,
      "gcp:high:custom:internet:internet": 0,
      "gcp:high:custom:peered-vpc:internet": 0,
      "gcp:high:dedicated:internet:internet": 0,
      "gcp:high:dedicated:peered-vpc:internet": 0,
      "gcp:high:standard:internet:internet": 0.008,
      "gcp:high:standard_v2:internet:internet": 0.008,
      "gcp:low:basic:internet:internet": 0.008,
      "gcp:low:custom:internet:internet": 0,
      "gcp:low:custom:peered-vpc:internet": 0,
      "gcp:low:dedicated:internet:internet": 0,
      "gcp:low:dedicated:peered-vpc:internet": 0,
      "gcp:low:standard:internet:internet": 0.008,
      "gcp:low:standard_v2:internet:internet": 0.008
    },
    "unit": "Partition-hour"
  },
  "KafkaBase": {
    "prices": {
      "gcp:high:standard:internet:internet": 0,
      "gcp:low:standard:internet:internet": 0,
      "azure:high:dedicated:internet:internet": 0,
      "gcp:low:standard_v2:internet:internet": 1.5,
      "aws:low:basic:internet:internet": 0,
      "aws:low:custom:peered-vpc:internet": 4.0063,
      "aws:low:standard:internet:internet": 0,
      "aws:high:standard_v2:internet:internet": 1.5,
      "azure:low:standard:internet:internet": 0,
      "azure:high:custom:internet:internet": 3.9362,
      "aws:low:dedicated:peered-vpc:internet": 0,
      "azure:high:dedicated:peered-vpc:internet": 0,
      "aws:high:custom:internet:internet": 4.6869,
      "azure:high:standard_v2:internet:internet": 1.5,
      "aws:low:dedicated:private-link:internet": 0,
      "azure:high:standard:internet:internet": 0,
      "gcp:high:dedicated:internet:internet": 0,
      "aws:high:dedicated:private-link:internet": 0,
      "azure:low:custom:internet:internet": 2.7265,
      "aws:high:standard:internet:internet": 0,
      "gcp:low:dedicated:peered-vpc:internet": 0,
      "azure:low:standard_v2:internet:internet": 1.5,
      "gcp:low:custom:internet:internet": 3.2994,
      "gcp:high:custom:internet:internet": 4.7631,
      "aws:high:dedicated:internet:internet": 0,
      "azure:low:dedicated:internet:internet": 0,
      "aws:low:standard_v2:internet:internet": 1.5,
      "azure:low:basic:internet:internet": 0,
      "gcp:low:basic:internet:internet": 0,
      "azure:high:custom:peered-vpc:internet": 4.5502,
      "gcp:high:dedicated:peered-vpc:internet": 0,
      "aws:high:dedicated:peered-vpc:internet": 0,
      "gcp:high:standard_v2:internet:internet": 1.5,
      "azure:low:dedicated:peered-vpc:internet": 0,
      "gcp:low:dedicated:internet:internet": 0,
      "gcp:low:custom:peered-vpc:internet": 4.0438,
      "aws:low:custom:internet:internet": 3.2506,
      "aws:low:dedicated:internet:internet": 0,
      "gcp:high:custom:peered-vpc:internet": 5.5074,
      "aws:high:custom:peered-vpc:internet": 5.4426,
      "aws:high:dedicated:transit-gateway:internet": 0,
      "azure:low:custom:peered-vpc:internet": 3.3405,
      "aws:low:dedicated:transit-gateway:internet": 0
    },
    "unit": "Hour"
  },
  "KSQLNumCSUs": {
    "prices": {
      "gcp:high:standard:internet:internet": 0.2222,
      "gcp:low:standard:internet:internet": 0.2222,
      "azure:high:dedicated:internet:internet": 0.2222,
      "gcp:low:standard_v2:internet:internet": 0.2222,
      "aws:low:basic:internet:internet": 0.2222,
      "aws:low:custom:peered-vpc:internet": 0.2222,
      "aws:low:standard:internet:internet": 0.2222,
      "aws:high:standard_v2:internet:internet": 0.2222,
      "azure:low:standard:internet:internet": 0.2222,
      "azure:high:custom:internet:internet": 0.2222,
      "aws:low:dedicated:peered-vpc:internet": 0.2222,
      "azure:high:dedicated:peered-vpc:internet": 0.2222,
      "aws:high:custom:internet:internet": 0.2222,
      "azure:high:standard_v2:internet:internet": 0.2222,
      "aws:low:dedicated:private-link:internet": 0.2222,
      "azure:high:standard:internet:internet": 0.2222,
      "gcp:high:dedicated:internet:internet": 0.2222,
      "aws:high:dedicated:private-link:internet": 0.2222,
      "azure:low:custom:internet:internet": 0.2222,
      "aws:high:standard:internet:internet": 0.2222,
      "gcp:low:dedicated:peered-vpc:internet": 0.2222,
      "azure:low:standard_v2:internet:internet": 0.2222,
      "gcp:low:custom:internet:internet": 0.2222,
      "gcp:high:custom:internet:internet": 0.2222,
      "aws:high:dedicated:internet:internet": 0.2222,
      "azure:low:dedicated:internet:internet": 0.2222,
      "aws:low:standard_v2:internet:internet": 0.2222,
      "azure:low:basic:internet:internet": 0.2222,
      "gcp:low:basic:internet:internet": 0.2222,
      "azure:high:custom:peered-vpc:internet": 0.2222,
      "gcp:high:dedicated:peered-vpc:internet": 0.2222,
      "aws:high:dedicated:peered-vpc:internet": 0.2222,
      "gcp:high:standard_v2:internet:internet": 0.2222,
      "azure:low:dedicated:peered-vpc:internet": 0.2222,
      "gcp:low:dedicated:internet:internet": 0.2222,
      "gcp:low:custom:peered-vpc:internet": 0.2222,
      "aws:low:custom:internet:internet": 0.2222,
      "aws:low:dedicated:internet:internet": 0.2222,
      "gcp:high:custom:peered-vpc:internet": 0.2222,
      "aws:high:custom:peered-vpc:internet": 0.2222,
      "aws:high:dedicated:transit-gateway:internet": 0.2222,
      "azure:low:custom:peered-vpc:internet": 0.2222,
      "aws:low:dedicated:transit-gateway:internet": 0.2222
    },
    "unit": "CSU-hour"
  },
  "KafkaNetworkRead": {
    "prices": {
      "gcp:high:standard:internet:internet": 0.11,
      "gcp:low:standard:internet:internet": 0.11,
      "azure:high:dedicated:internet:internet": 0.014,
      "gcp:low:standard_v2:internet:internet": 0.04,
      "aws:low:basic:internet:internet": 0.13,
      "aws:low:custom:peered-vpc:internet": 0.0364,
      "aws:low:standard:internet:internet": 0.13,
      "aws:high:standard_v2:internet:internet": 0.06,
      "azure:low:standard:internet:internet": 0.24,
      "azure:high:custom:internet:internet": 0.0227,
      "aws:low:dedicated:peered-vpc:internet": 0.032,
      "azure:high:dedicated:peered-vpc:internet": 0.014,
      "aws:high:custom:internet:internet": 0.0523,
      "azure:high:standard_v2:internet:internet": 0.05,
      "aws:low:dedicated:private-link:internet": 0.032,
      "azure:high:standard:internet:internet": 0.24,
      "gcp:high:dedicated:internet:internet": 0.008,
      "aws:high:dedicated:private-link:internet": 0.032,
      "azure:low:custom:internet:internet": 0.0227,
      "aws:high:standard:internet:internet": 0.13,
      "gcp:low:dedicated:peered-vpc:internet": 0.008,
      "azure:low:standard_v2:internet:internet": 0.05,
      "gcp:low:custom:internet:internet": 0.0091,
      "gcp:high:custom:internet:internet": 0.0091,
      "aws:high:dedicated:internet:internet": 0.046,
      "azure:low:dedicated:internet:internet": 0.014,
      "aws:low:standard_v2:internet:internet": 0.06,
      "azure:low:basic:internet:internet": 0.12,
      "gcp:low:basic:internet:internet": 0.11,
      "azure:high:custom:peered-vpc:internet": 0.0227,
      "gcp:high:dedicated:peered-vpc:internet": 0.008,
      "aws:high:dedicated:peered-vpc:internet": 0.032,
      "gcp:high:standard_v2:internet:internet": 0.04,
      "azure:low:dedicated:peered-vpc:internet": 0.014,
      "gcp:low:dedicated:internet:internet": 0.008,
      "gcp:low:custom:peered-vpc:internet": 0.0091,
      "aws:low:custom:internet:internet": 0.0523,
      "aws:low:dedicated:internet:internet": 0.046,
      "gcp:high:custom:peered-vpc:internet": 0.0091,
      "aws:high:custom:peered-vpc:internet": 0.0364,
      "aws:high:dedicated:transit-gateway:internet": 0.112,
      "azure:low:custom:peered-vpc:internet": 0.0227,
      "aws:low:dedicated:transit-gateway:internet": 0.112
    },
    "unit": "GB"
  },
  "KafkaNumCKUs": {
    "prices": {
      "gcp:high:standard:internet:internet": 0,
      "gcp:low:standard:internet:internet": 0,
      "azure:high:dedicated:internet:internet": 2.941,
      "gcp:low:standard_v2:internet:internet": 0,
      "aws:low:basic:internet:internet": 0,
      "aws:low:custom:peered-vpc:internet": 0.2394,
      "aws:low:standard:internet:internet": 0,
      "aws:high:standard_v2:internet:internet": 0,
      "azure:low:standard:internet:internet": 0,
      "azure:high:custom:internet:internet": 0.3024,
      "aws:low:dedicated:peered-vpc:internet": 3.46,
      "azure:high:dedicated:peered-vpc:internet": 2.941,
      "aws:high:custom:internet:internet": 0.2394,
      "azure:high:standard_v2:internet:internet": 0,
      "aws:low:dedicated:private-link:internet": 3.46,
      "azure:high:standard:internet:internet": 0,
      "gcp:high:dedicated:internet:internet": 2.422,
      "aws:high:dedicated:private-link:internet": 3.46,
      "azure:low:custom:internet:internet": 0.3024,
      "aws:high:standard:internet:internet": 0,
      "gcp:low:dedicated:peered-vpc:internet": 2.422,
      "azure:low:standard_v2:internet:internet": 0,
      "gcp:low:custom:internet:internet": 0.3659,
      "gcp:high:custom:internet:internet": 0.3659,
      "aws:high:dedicated:internet:internet": 3.46,
      "azure:low:dedicated:internet:internet": 2.941,
      "aws:low:standard_v2:internet:internet": 0,
      "azure:low:basic:internet:internet": 0,
      "gcp:low:basic:internet:internet": 0,
      "azure:high:custom:peered-vpc:internet": 0.3024,
      "gcp:high:dedicated:peered-vpc:internet": 2.422,
      "aws:high:dedicated:peered-vpc:internet": 3.46,
      "gcp:high:standard_v2:internet:internet": 0,
      "azure:low:dedicated:peered-vpc:internet": 2.941,
      "gcp:low:dedicated:internet:internet": 2.422,
      "gcp:low:custom:peered-vpc:internet": 0.3659,
      "aws:low:custom:internet:internet": 0.2394,
      "aws:low:dedicated:internet:internet": 3.46,
      "gcp:high:custom:peered-vpc:internet": 0.3659,
      "aws:high:custom:peered-vpc:internet": 0.2394,
      "aws:high:dedicated:transit-gateway:internet": 3.46,
      "azure:low:custom:peered-vpc:internet": 0.3024,
      "aws:low:dedicated:transit-gateway:internet": 3.46
    },
    "unit": "CKU-hour"
  },
  "KafkaStorage": {
    "prices": {
      "gcp:high:standard:internet:internet": 0.00013889,
      "gcp:low:standard:internet:internet": 0.00013889,
      "azure:high:dedicated:internet:internet": 0.00015556,
      "gcp:low:standard_v2:internet:internet": 0.00013889,
      "aws:low:basic:internet:internet": 0.00013889,
      "aws:low:custom:peered-vpc:internet": 0.00015778,
      "aws:low:standard:internet:internet": 0.00013889,
      "aws:high:standard_v2:internet:internet": 0.00013889,
      "azure:low:standard:internet:internet": 0.00013889,
      "azure:high:custom:internet:internet": 0.0002525,
      "aws:low:dedicated:peered-vpc:internet": 0.00013889,
      "azure:high:dedicated:peered-vpc:internet": 0.00015556,
      "aws:high:custom:internet:internet": 0.00015778,
      "azure:high:standard_v2:internet:internet": 0.00013889,
      "aws:low:dedicated:private-link:internet": 0.00013889,
      "azure:high:standard:internet:internet": 0.00013889,
      "gcp:high:dedicated:internet:internet": 0.00012444,
      "aws:high:dedicated:private-link:internet": 0.00013889,
      "azure:low:custom:internet:internet": 0.0002525,
      "aws:high:standard:internet:internet": 0.00013889,
      "gcp:low:dedicated:peered-vpc:internet": 0.00012444,
      "azure:low:standard_v2:internet:internet": 0.00013889,
      "gcp:low:custom:internet:internet": 0.00014139,
      "gcp:high:custom:internet:internet": 0.00014139,
      "aws:high:dedicated:internet:internet": 0.00013889,
      "azure:low:dedicated:internet:internet": 0.00015556,
      "aws:low:standard_v2:internet:internet": 0.00013889,
      "azure:low:basic:internet:internet": 0.00013889,
      "gcp:low:basic:internet:internet": 0.00013889,
      "azure:high:custom:peered-vpc:internet": 0.0002525,
      "gcp:high:dedicated:peered-vpc:internet": 0.00012444,
      "aws:high:dedicated:peered-vpc:internet": 0.00013889,
      "gcp:high:standard_v2:internet:internet": 0.00013889,
      "azure:low:dedicated:peered-vpc:internet": 0.00015556,
      "gcp:low:dedicated:internet:internet": 0.00012444,
      "gcp:low:custom:peered-vpc:internet": 0.00014139,
      "aws:low:custom:internet:internet": 0.00015778,
      "aws:low:dedicated:internet:internet": 0.00013889,
      "gcp:high:custom:peered-vpc:internet": 0.00014139,
      "aws:high:custom:peered-vpc:internet": 0.00015778,
      "aws:high:dedicated:transit-gateway:internet": 0.00013889,
      "azure:low:custom:peered-vpc:internet": 0.0002525,
      "aws:low:dedicated:transit-gateway:internet": 0.00013889
    },
    "unit": "GB-hour"
  },
  "KafkaNetworkWrite": {
    "prices": {
      "gcp:high:standard:internet:internet": 0.22,
      "gcp:low:standard:internet:internet": 0.11,
      "azure:high:dedicated:internet:internet": 0.062,
      "gcp:low:standard_v2:internet:internet": 0.04,
      "aws:low:basic:internet:internet": 0.13,
      "aws:low:custom:peered-vpc:internet": 0.0364,
      "aws:low:standard:internet:internet": 0.13,
      "aws:high:standard_v2:internet:internet": 0.13,
      "azure:low:standard:internet:internet": 0.22,
      "azure:high:custom:internet:internet": 0.2045,
      "aws:low:dedicated:peered-vpc:internet": 0.032,
      "azure:high:dedicated:peered-vpc:internet": 0.062,
      "aws:high:custom:internet:internet": 0.1159,
      "azure:high:standard_v2:internet:internet": 0.12,
      "aws:low:dedicated:private-link:internet": 0.032,
      "azure:high:standard:internet:internet": 0.48,
      "gcp:high:dedicated:internet:internet": 0.034,
      "aws:high:dedicated:private-link:internet": 0.088,
      "azure:low:custom:internet:internet": 0.0227,
      "aws:high:standard:internet:internet": 0.28,
      "gcp:low:dedicated:peered-vpc:internet": 0.01,
      "azure:low:standard_v2:internet:internet": 0.05,
      "gcp:low:custom:internet:internet": 0.0298,
      "gcp:high:custom:internet:internet": 0.0571,
      "aws:high:dedicated:internet:internet": 0.102,
      "azure:low:dedicated:internet:internet": 0.014,
      "aws:low:standard_v2:internet:internet": 0.06,
      "azure:low:basic:internet:internet": 0.12,
      "gcp:low:basic:internet:internet": 0.11,
      "azure:high:custom:peered-vpc:internet": 0.2045,
      "gcp:high:dedicated:peered-vpc:internet": 0.034,
      "aws:high:dedicated:peered-vpc:internet": 0.088,
      "gcp:high:standard_v2:internet:internet": 0.11,
      "azure:low:dedicated:peered-vpc:internet": 0.014,
      "gcp:low:dedicated:internet:internet": 0.01,
      "gcp:low:custom:peered-vpc:internet": 0.0298,
      "aws:low:custom:internet:internet": 0.0523,
      "aws:low:dedicated:internet:internet": 0.046,
      "gcp:high:custom:peered-vpc:internet": 0.0571,
      "aws:high:custom:peered-vpc:internet": 0.1,
      "aws:high:dedicated:transit-gateway:internet": 0.088,
      "azure:low:custom:peered-vpc:internet": 0.0227,
      "aws:low:dedicated:transit-gateway:internet": 0.032
    },
    "unit": "GB"
  }
}';

UPDATE deployment.price SET price_table = jsonb_set(price_table, '{cluster}', kafka_prices, TRUE);

UPDATE deployment.price SET support_prices = '{
  "support-cloud-basic": {
    "min_price": 0
  },
  "support-cloud-premier": {
    "min_price": 10000,
    "usage_price": [
      {
        "rate": 0.1,
        "max_range": 100000,
        "min_range": 0
      },
      {
        "rate": 0.08,
        "max_range": 500000,
        "min_range": 100000
      },
      {
        "rate": 0.06,
        "max_range": 1000000,
        "min_range": 500000
      },
      {
        "rate": 0.03,
        "max_range": -1,
        "min_range": 1000000
      }
    ]
  },
  "support-cloud-business": {
    "min_price": 1000,
    "usage_price": [
      {
        "rate": 0.1,
        "max_range": 50000,
        "min_range": 0
      },
      {
        "rate": 0.08,
        "max_range": 100000,
        "min_range": 50000
      },
      {
        "rate": 0.06,
        "max_range": 1000000,
        "min_range": 100000
      },
      {
        "rate": 0.03,
        "max_range": -1,
        "min_range": 1000000
      }
    ]
  },
  "support-cloud-developer": {
    "min_price": 29,
    "usage_price": [
      {
        "rate": 0.05,
        "max_range": -1,
        "min_range": 0
      }
    ]
  }
}';

UPDATE deployment.price SET price_table = jsonb_set(price_table, '{support}', support_prices, TRUE);

UPDATE deployment.price SET connect_prices = '{
  "ConnectNumRecords": {
    "prices": {
      "aws:dedicated:internet:AzureBlobSink": 0,
      "aws:dedicated:internet:GcsSink": 0,
      "aws:dedicated:internet:S3_SINK": 0,
      "aws:dedicated:peered-vpc:AzureBlobSink": 0,
      "aws:dedicated:peered-vpc:GcsSink": 0,
      "aws:dedicated:peered-vpc:S3_SINK": 0,
      "aws:standard_v2:internet:AzureBlobSink": 0,
      "aws:standard_v2:internet:GcsSink": 0,
      "aws:standard_v2:internet:S3_SINK": 0,
      "aws:standard_v2:peered-vpc:AzureBlobSink": 0,
      "aws:standard_v2:peered-vpc:GcsSink": 0,
      "aws:standard_v2:peered-vpc:S3_SINK": 0,
      "azure:dedicated:internet:AzureBlobSink": 0,
      "azure:dedicated:internet:GcsSink": 0,
      "azure:dedicated:internet:S3_SINK": 0,
      "azure:dedicated:peered-vpc:AzureBlobSink": 0,
      "azure:dedicated:peered-vpc:GcsSink": 0,
      "azure:dedicated:peered-vpc:S3_SINK": 0,
      "azure:standard_v2:internet:AzureBlobSink": 0,
      "azure:standard_v2:internet:GcsSink": 0,
      "azure:standard_v2:internet:S3_SINK": 0,
      "azure:standard_v2:peered-vpc:AzureBlobSink": 0,
      "azure:standard_v2:peered-vpc:GcsSink": 0,
      "azure:standard_v2:peered-vpc:S3_SINK": 0,
      "gcp:dedicated:internet:AzureBlobSink": 0,
      "gcp:dedicated:internet:GcsSink": 0,
      "gcp:dedicated:internet:S3_SINK": 0,
      "gcp:dedicated:peered-vpc:AzureBlobSink": 0,
      "gcp:dedicated:peered-vpc:GcsSink": 0,
      "gcp:dedicated:peered-vpc:S3_SINK": 0,
      "gcp:standard_v2:internet:AzureBlobSink": 0,
      "gcp:standard_v2:internet:GcsSink": 0,
      "gcp:standard_v2:internet:S3_SINK": 0,
      "gcp:standard_v2:peered-vpc:AzureBlobSink": 0,
      "gcp:standard_v2:peered-vpc:GcsSink": 0,
      "gcp:standard_v2:peered-vpc:S3_SINK": 0
    },
    "unit": "Record"
  },
  "ConnectNumTasks": {
    "prices": {
      "aws:dedicated:internet:AzureBlobSink": 0.0347,
      "aws:dedicated:internet:GcsSink": 0.0347,
      "aws:dedicated:internet:S3_SINK": 0.0347,
      "aws:dedicated:peered-vpc:AzureBlobSink": 0.0347,
      "aws:dedicated:peered-vpc:GcsSink": 0.0347,
      "aws:dedicated:peered-vpc:S3_SINK": 0.0347,
      "aws:standard_v2:internet:AzureBlobSink": 0.0347,
      "aws:standard_v2:internet:GcsSink": 0.0347,
      "aws:standard_v2:internet:S3_SINK": 0.0347,
      "aws:standard_v2:peered-vpc:AzureBlobSink": 0.0347,
      "aws:standard_v2:peered-vpc:GcsSink": 0.0347,
      "aws:standard_v2:peered-vpc:S3_SINK": 0.0347,
      "azure:dedicated:internet:AzureBlobSink": 0.0347,
      "azure:dedicated:internet:GcsSink": 0.0347,
      "azure:dedicated:internet:S3_SINK": 0.0347,
      "azure:dedicated:peered-vpc:AzureBlobSink": 0.0347,
      "azure:dedicated:peered-vpc:GcsSink": 0.0347,
      "azure:dedicated:peered-vpc:S3_SINK": 0.0347,
      "azure:standard_v2:internet:AzureBlobSink": 0.0347,
      "azure:standard_v2:internet:GcsSink": 0.0347,
      "azure:standard_v2:internet:S3_SINK": 0.0347,
      "azure:standard_v2:peered-vpc:AzureBlobSink": 0.0347,
      "azure:standard_v2:peered-vpc:GcsSink": 0.0347,
      "azure:standard_v2:peered-vpc:S3_SINK": 0.0347,
      "gcp:dedicated:internet:AzureBlobSink": 0.0347,
      "gcp:dedicated:internet:GcsSink": 0.0347,
      "gcp:dedicated:internet:S3_SINK": 0.0347,
      "gcp:dedicated:peered-vpc:AzureBlobSink": 0.0347,
      "gcp:dedicated:peered-vpc:GcsSink": 0.0347,
      "gcp:dedicated:peered-vpc:S3_SINK": 0.0347,
      "gcp:standard_v2:internet:AzureBlobSink": 0.0347,
      "gcp:standard_v2:internet:GcsSink": 0.0347,
      "gcp:standard_v2:internet:S3_SINK": 0.0347,
      "gcp:standard_v2:peered-vpc:AzureBlobSink": 0.0347,
      "gcp:standard_v2:peered-vpc:GcsSink": 0.0347,
      "gcp:standard_v2:peered-vpc:S3_SINK": 0.0347
    },
    "unit": "Task-hour"
  },
  "ConnectThroughput": {
    "prices": {
      "aws:dedicated:internet:AzureBlobSink": 0.03,
      "aws:dedicated:internet:GcsSink": 0.03,
      "aws:dedicated:internet:S3_SINK": 0.03,
      "aws:dedicated:peered-vpc:AzureBlobSink": 0.03,
      "aws:dedicated:peered-vpc:GcsSink": 0.03,
      "aws:dedicated:peered-vpc:S3_SINK": 0.03,
      "aws:standard_v2:internet:AzureBlobSink": 0.03,
      "aws:standard_v2:internet:GcsSink": 0.03,
      "aws:standard_v2:internet:S3_SINK": 0.03,
      "aws:standard_v2:peered-vpc:AzureBlobSink": 0.03,
      "aws:standard_v2:peered-vpc:GcsSink": 0.03,
      "aws:standard_v2:peered-vpc:S3_SINK": 0.03,
      "azure:dedicated:internet:AzureBlobSink": 0.03,
      "azure:dedicated:internet:GcsSink": 0.03,
      "azure:dedicated:internet:S3_SINK": 0.03,
      "azure:dedicated:peered-vpc:AzureBlobSink": 0.03,
      "azure:dedicated:peered-vpc:GcsSink": 0.03,
      "azure:dedicated:peered-vpc:S3_SINK": 0.03,
      "azure:standard_v2:internet:AzureBlobSink": 0.03,
      "azure:standard_v2:internet:GcsSink": 0.03,
      "azure:standard_v2:internet:S3_SINK": 0.03,
      "azure:standard_v2:peered-vpc:AzureBlobSink": 0.03,
      "azure:standard_v2:peered-vpc:GcsSink": 0.03,
      "azure:standard_v2:peered-vpc:S3_SINK": 0.03,
      "gcp:dedicated:internet:AzureBlobSink": 0.03,
      "gcp:dedicated:internet:GcsSink": 0.03,
      "gcp:dedicated:internet:S3_SINK": 0.03,
      "gcp:dedicated:peered-vpc:AzureBlobSink": 0.03,
      "gcp:dedicated:peered-vpc:GcsSink": 0.03,
      "gcp:dedicated:peered-vpc:S3_SINK": 0.03,
      "gcp:standard_v2:internet:AzureBlobSink": 0.03,
      "gcp:standard_v2:internet:GcsSink": 0.03,
      "gcp:standard_v2:internet:S3_SINK": 0.03,
      "gcp:standard_v2:peered-vpc:AzureBlobSink": 0.03,
      "gcp:standard_v2:peered-vpc:GcsSink": 0.03,
      "gcp:standard_v2:peered-vpc:S3_SINK": 0.03
    },
    "unit": "GB"
  }
}';

UPDATE deployment.price SET price_table = jsonb_set(price_table, '{connect}', connect_prices, TRUE);

--
-- Name: price_audit_log; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE price_audit_log (
    id SERIAL PRIMARY KEY,
    rate_card_id INTEGER NOT NULL REFERENCES price(id),
    field_changed VARCHAR(50) NOT NULL,
    previous_value jsonb,
    current_value jsonb,
    operation INTEGER DEFAULT 0 NOT NULL,
    operation_time TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL,
    email character varying(128)
);

ALTER TABLE price_audit_log OWNER TO caas;

CREATE INDEX index_price_audit_log_rate_card_id ON deployment.price_audit_log (rate_card_id);

--
-- Name: billing_invoice; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE billing_invoice (
    id SERIAL PRIMARY KEY,
    organization_id integer NOT NULL,
    total BIGINT DEFAULT 0 NOT NULL,
    lines jsonb NOT NULL DEFAULT('{}'),
    billing_method INTEGER DEFAULT 0 NOT NULL,
    currency TEXT DEFAULT '' NOT NULL,
    from_date timestamp without time zone,
    to_date timestamp without time zone,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    status INTEGER
);

ALTER TABLE billing_invoice OWNER TO caas;

CREATE INDEX index_billing_invoice_organization_id ON deployment.billing_invoice (organization_id);
CREATE INDEX index_billing_invoice_not_sent_billing_method_from_date ON deployment.billing_invoice (billing_method,from_date) WHERE status <> 3;

--
-- Name: invoice; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE invoice (
    id SERIAL PRIMARY KEY,
    organization_id integer NOT NULL,
    total BIGINT DEFAULT 0 NOT NULL,
    lines jsonb NOT NULL DEFAULT('{}'),
    billing_method INTEGER DEFAULT 0 NOT NULL,
    currency TEXT DEFAULT '' NOT NULL,
    from_date timestamp without time zone,
    to_date timestamp without time zone,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    status INTEGER
);

ALTER TABLE invoice OWNER TO caas;

--
-- Name: credit; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE credit (
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL DEFAULT '',
    type INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    active_date TIMESTAMP WITHOUT TIME zone NOT NULL,
    expire_date TIMESTAMP WITHOUT TIME zone NOT NULL,
    created TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL,
    modified TIMESTAMP WITHOUT TIME zone DEFAULT now() NOT NULL,
    deactivated boolean DEFAULT false NOT NULL
);

ALTER TABLE credit OWNER TO caas;

--
-- Name: connect_task_usage_seq; Type: SEQUENCE; Schema: deployment; Owner: caas
--

CREATE SEQUENCE connect_task_usage_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE connect_task_usage_seq OWNER TO caas;

--
-- Name: connect_task_usage; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE connect_task_usage (
    id integer NOT NULL PRIMARY KEY,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    tasks_used integer,
    organization_id integer NOT NULL,
    task_limit_config jsonb DEFAULT '{}'::jsonb NOT NULL
);

ALTER TABLE connect_task_usage OWNER TO caas;

--
-- Name: connect_task_usage id; Type: DEFAULT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY connect_task_usage ALTER COLUMN id SET DEFAULT nextval('connect_task_usage_seq'::regclass);

--
-- Name: connect_task_usage organization_id_uniq; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY connect_task_usage
    ADD CONSTRAINT organization_id_uniq UNIQUE (organization_id);

--
-- Name: connect_plugin; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE connect_plugin (
    id serial PRIMARY KEY,
    name character varying(64) UNIQUE NOT NULL,
    clouds text[] NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    pa_date timestamp without time zone,
    plugin jsonb DEFAULT '{}'::jsonb NOT NULL,
    display jsonb DEFAULT '{}'::jsonb NOT NULL,
    validation_parameters jsonb DEFAULT '{}'::jsonb NOT NULL
);

ALTER TABLE connect_plugin OWNER TO caas;

INSERT INTO deployment.connect_plugin (name, clouds, plugin, display)
VALUES ('S3_SINK', '{"aws", "gcp"}', '{
  "class": "S3_SINK",
  "type": "sink",
  "version": "0.1.0"
}', '{
  "image_url": "https://cdn.worldvectorlogo.com/logos/aws-s3.svg",
  "product_maturity_phase": 3,
  "display_name": "Amazon S3 Sink"
}');
INSERT INTO deployment.connect_plugin (name, clouds, plugin, display)
VALUES ('GcsSink', '{"gcp"}', '{
  "class": "GcsSink",
  "type": "sink",
  "version": "0.1.0"
}', '{
  "image_url": "https://api.hub.confluent.io/api/plugins/confluentinc/kafka-connect-gcs/versions/5.0.3/assets/googlecloud.png",
  "product_maturity_phase": 3,
  "display_name": "Google Cloud Storage Sink"
}');
INSERT INTO deployment.connect_plugin (name, clouds, plugin, display)
VALUES ('SpannerSink', '{"gcp"}', '{
  "class": "SpannerSink",
  "type": "sink",
  "version": "0.1.0"
}', '{
  "image_url": "https://ccloud-connector-images.s3-us-west-2.amazonaws.com/GCP-Spanner-Logo.svg.png",
  "product_maturity_phase": 2,
  "display_name": "Google Cloud Spanner Sink"
}');
INSERT INTO deployment.connect_plugin (name, clouds, plugin, display)
VALUES ('AzureBlobSink', '{"azure"}', '{
  "class": "AzureBlobSink",
  "type": "sink",
  "version": "0.1.0"
}', '{
  "image_url": "https://www.drupal.org/files/styles/grid-3-2x/public/project-images/azure-storage-blob.png",
  "product_maturity_phase": 2,
  "display_name": "Azure Blob Storage Sink"
}');
INSERT INTO deployment.connect_plugin (name, clouds, plugin, display)
VALUES ('AzureDataLakeGen2Sink', '{"azure"}', '{
  "class": "AzureDataLakeGen2Sink",
  "type": "sink",
  "version": "0.1.0"
}', '{
  "image_url": "https://ccloud-connector-images.s3-us-west-2.amazonaws.com/Azure-DataLake-icon.png",
  "product_maturity_phase": 2,
  "display_name": "Azure Data Lake Storage Gen2 Sink"
}');
INSERT INTO deployment.connect_plugin (name, clouds, plugin, display)
VALUES ('AzureEventHubsSource', '{"aws","azure","gcp"}', '{
  "class": "AzureEventHubsSource",
  "type": "source",
  "version": "0.1.0"
}', '{
  "image_url": "https://ccloud-connector-images.s3-us-west-2.amazonaws.com/EventsHublogo.png",
  "product_maturity_phase": 2,
  "display_name": "Azure Event Hubs Source"
}');
INSERT INTO deployment.connect_plugin (name, clouds, plugin, display)
VALUES ('BigQuerySink', '{"gcp"}', '{
  "class": "BigQuerySink",
  "type": "sink",
  "version": "0.1.0"
}', '{
  "image_url": "https://d1i4a15mxbxib1.cloudfront.net/api/plugins/wepay/kafka-connect-bigquery/versions/1.1.2/assets/BigQuery.png",
  "product_maturity_phase": 2,
  "display_name": "Google BigQuery Sink"
}');
INSERT INTO deployment.connect_plugin (name, clouds, plugin, display)
VALUES ('PostgresSource', '{"aws","azure","gcp"}', '{
  "class": "PostgresSource",
  "type": "source",
  "version": "0.1.0"
}', '{
  "image_url": "https://ccloud-connector-images.s3-us-west-2.amazonaws.com/elephant.png",
  "product_maturity_phase": 2,
  "display_name": "Postgres Source"
}');
INSERT INTO deployment.connect_plugin (name, clouds, plugin, display)
VALUES ('MySqlSource', '{"aws","azure","gcp"}', '{
  "class": "MySqlSource",
  "type": "source",
  "version": "0.1.0"
}', '{
  "image_url": "https://cdn.worldvectorlogo.com/logos/mysql.svg",
  "product_maturity_phase": 2,
  "display_name": "MySQL Source"
}');
INSERT INTO deployment.connect_plugin (name, clouds, plugin, display)
VALUES ('KinesisSource', '{"aws","azure","gcp"}', '{
  "class": "KinesisSource",
  "type": "source",
  "version": "0.1.0"
}', '{
  "image_url": "https://ccloud-connector-images.s3-us-west-2.amazonaws.com/Amazon-Kinesis%404x.png",
  "product_maturity_phase": 2,
  "display_name": "Kinesis Source"
}');
INSERT INTO deployment.connect_plugin (name, clouds, plugin, display)
VALUES ('OracleDatabaseSource', '{"aws","azure","gcp"}', '{
  "class": "OracleDatabaseSource",
  "type": "source",
  "version": "0.1.0"
}', '{
  "image_url": "https://images.youracclaim.com/images/de9f4975-5ac3-4d39-914e-f733121683e1/Oracle_Org.png",
  "product_maturity_phase": 2,
  "display_name": "Oracle Database Source"
}');
INSERT INTO deployment.connect_plugin (name, clouds, plugin, display)
VALUES ('MicrosoftSqlServerSource', '{"aws","azure","gcp"}', '{
  "class": "MicrosoftSqlServerSource",
  "type": "source",
  "version": "0.1.0"
}', '{
  "image_url": "https://ccloud-connector-images.s3-us-west-2.amazonaws.com/Microsoft+SQL+Server.svg",
  "product_maturity_phase": 2,
  "display_name": "Microsoft SQL Server Source"
}');
INSERT INTO deployment.connect_plugin (name, clouds, plugin, display)
VALUES ('PubSubSource', '{"aws","azure","gcp"}', '{
  "class": "PubSubSource",
  "type": "source",
  "version": "0.1.0"
}', '{
  "image_url": "https://ccloud-connector-images.s3-us-west-2.amazonaws.com/Cloud+PubSub.png",
  "product_maturity_phase": 2,
  "display_name": "Google Cloud Pub/Sub Source"
}');
INSERT INTO deployment.connect_plugin (name, clouds, plugin, display)
VALUES ('RedshiftSink', '{"aws","azure","gcp"}', '{
  "class": "RedshiftSink",
  "type": "sink",
  "version": "0.1.0"
}', '{
  "image_url": "https://ccloud-connector-images.s3-us-west-2.amazonaws.com/Amazon-Redshift%404x.png",
  "product_maturity_phase": 2,
  "display_name": "Amazon Redshift Sink"
}');

--
-- Name: connect_error_message_mappings; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE connect_error_message_mappings (
    id serial PRIMARY KEY,
    error_message varchar NOT NULL,
    user_message varchar NOT NULL,
    connector_type varchar NOT NULL
);

ALTER TABLE connect_error_message_mappings OWNER TO caas;

--
-- Name: feature_opt_ins; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE feature_opt_ins (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    feature INTEGER NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL
);

ALTER TABLE feature_opt_ins OWNER TO caas;

--
-- Name: billing_record; Type: TABLE; Schema: deployment; Owner: caas
--
CREATE TABLE deployment.billing_record (
    id text PRIMARY KEY,
    organization_id integer NOT NULL,
    transaction_id text,
    logical_cluster_id text,
    invoice_lines jsonb NOT NULL DEFAULT('{}'),
    metrics jsonb NOT NULL DEFAULT('{}'),
    amount BIGINT NOT NULL,
    type integer NOT NULL,
    window_size integer NOT NULL,
    timestamp integer NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    deactivated boolean DEFAULT false NOT NULL
);

ALTER TABLE billing_record OWNER TO caas;

CREATE INDEX index_billing_records_logical_cluster_id ON deployment.billing_record (logical_cluster_id);
CREATE INDEX index_billing_records_organization_id ON deployment.billing_record (organization_id);
CREATE INDEX index_billing_records_transaction_id ON deployment.billing_record (transaction_id);

--
-- Name: organization id; Type: DEFAULT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY organization ALTER COLUMN id SET DEFAULT nextval('organization_id_seq'::regclass);

--
-- Name: entitlement id; Type: DEFAULT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY entitlement ALTER COLUMN id SET DEFAULT nextval('entitlement_id_seq'::regclass);

--
-- Name: users id; Type: DEFAULT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);

--
-- Name: usage_metrics_errors; Type: TABLE; Schema: deployment; Owner: -
--

CREATE TABLE deployment.usage_metrics_errors (
    id integer NOT NULL,
    operation_id character varying(40) NOT NULL,
    organization_id integer NOT NULL,
    product_level character varying(20) NOT NULL,
    marketplace_partner character varying(10) NOT NULL,
    lines jsonb DEFAULT '{}'::jsonb NOT NULL,
    start_time timestamp without time zone NOT NULL,
    end_time timestamp without time zone NOT NULL,
    created timestamp without time zone NOT NULL,
    modified timestamp without time zone NOT NULL,
    error jsonb DEFAULT '{}'::jsonb NOT NULL,
    sent_to_marketplace boolean NOT NULL,
    version character varying(20) DEFAULT '' NOT NULL
);

--
-- Name: usage_metrics_errors; Type: TABLE; Schema: deployment; Owner: caas
--

ALTER TABLE deployment.usage_metrics_errors OWNER TO caas;

--
-- Name: usage_metrics_errors_id_seq; Type: SEQUENCE; Schema: deployment; Owner: -
--

CREATE SEQUENCE deployment.usage_metrics_errors_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

--
-- Name: usage_metrics_errors_id_seq; Type: SEQUENCE; Schema: deployment; Owner: caas
--

ALTER TABLE deployment.usage_metrics_errors_id_seq OWNER TO caas;

--
-- Name: usage_metrics_errors_id_seq; Type: SEQUENCE OWNED BY; Schema: deployment; Owner: caas
--

ALTER SEQUENCE deployment.usage_metrics_errors_id_seq OWNED BY deployment.usage_metrics_errors.id;

--
-- Name: usage_metrics_errors id; Type: DEFAULT; Schema: deployment; Owner: -
--

ALTER TABLE ONLY deployment.usage_metrics_errors ALTER COLUMN id SET DEFAULT nextval('deployment.usage_metrics_errors_id_seq'::regclass);

--
-- Name: usage_metrics_errors usage_metrics_errors_pkey; Type: CONSTRAINT; Schema: deployment; Owner: -
--

ALTER TABLE ONLY deployment.usage_metrics_errors
    ADD CONSTRAINT usage_metrics_errors_pkey PRIMARY KEY (id);

--
-- Data for Name: account; Type: TABLE DATA; Schema: deployment; Owner: caas
--

COPY account (id, name, config, created, modified, deactivated, organization_id) FROM stdin;
\.

--
-- Name: api_key_id_seq; Type: SEQUENCE SET; Schema: deployment; Owner: caas
--

SELECT pg_catalog.setval('api_key_id_seq', 2, true);

--
-- Data for Name: cp_component; Type: TABLE DATA; Schema: deployment; Owner: caas
--

COPY cp_component (id, default_version, created, modified) FROM stdin;
kafka	0.3.0	2017-06-22 13:50:24.580803	2017-06-22 13:50:24.580803
zookeeper	0.3.0	2017-06-22 13:50:24.580803	2017-06-22 13:50:24.580803
\.


--
-- Data for Name: environment; Type: TABLE DATA; Schema: deployment; Owner: caas
--

COPY environment (id, config, created, modified) FROM stdin;
devel	{}	2017-06-08 23:18:32.009539	2017-08-19 01:23:42.349148
\.

INSERT INTO environment (id, config, created, modified) VALUES ('private', '{}', now(), now());

--
-- Name: k8s_cluster_num; Type: SEQUENCE SET; Schema: deployment; Owner: caas
--

SELECT pg_catalog.setval('k8s_cluster_num', 1, true);

--
-- Name: organization_id_seq; Type: SEQUENCE SET; Schema: deployment; Owner: caas
--

SELECT pg_catalog.setval('organization_id_seq', 592, true);

INSERT INTO deployment.cloud (id, config, created, modified, name)
VALUES
    ('aws', '{"glb_dns_domain": "aws.glb.devel.cpdev.cloud", "dns_domain": "aws.devel.cpdev.cloud", "internal_dns_domain": "aws.internal.devel.cpdev.cloud"}', now(), now(), 'Amazon Web Services'),
    ('gcp', '{"glb_dns_domain": "gcp.glb.devel.cpdev.cloud", "dns_domain": "gcp.devel.cpdev.cloud", "internal_dns_domain": "gcp.internal.devel.cpdev.cloud"}', now(), now(), 'Google Cloud Platform'),
    ('azure', '{"glb_dns_domain": "azure.glb.devel.cpdev.cloud", "dns_domain": "azure.devel.cpdev.cloud", "internal_dns_domain": "azure.internal.devel.cpdev.cloud"}', now(), now(), 'Azure');

--
-- Data for Name: region; Type: TABLE DATA; Schema: deployment; Owner: caas
--

COPY region (id, cloud, config, created, modified, name) FROM stdin;
us-west-2	aws	{"docker": {"repo": "037803949979.dkr.ecr.us-west-2.amazonaws.com", "image_prefix": "confluentinc"}}	2017-06-22 13:50:24.567898	2017-06-22 13:50:24.567898	US West (Oregon)
us-central1	gcp	{"docker": {"repo": "us.gcr.io", "image_prefix": "cc-devel"}}	2017-06-22 13:50:24.567898	2017-06-22 13:50:24.567898	US Central
centralus	azure	{"docker": {"repo": "cclouddevel.azurecr.io", "image_prefix": "confluentinc"}}	2017-06-22 13:50:24.567898	2017-06-22 13:50:24.567898	Central US
\.


--
-- Name: account_num; Type: SEQUENCE SET; Schema: deployment; Owner: caas
--

SELECT pg_catalog.setval('account_num', 589, true);


--
-- Data for Name: users; Type: TABLE DATA; Schema: deployment; Owner: caas
--

COPY users (id, email, first_name, last_name, organization_id, deactivated, created, modified) FROM stdin;
\.


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: deployment; Owner: caas
--

SELECT pg_catalog.setval('users_id_seq', 2, true);


--
-- Name: api_key api_key_pkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY api_key
    ADD CONSTRAINT api_key_pkey PRIMARY KEY (id);


--
-- Name: cloud cloud_pkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY cloud
    ADD CONSTRAINT cloud_pkey PRIMARY KEY (id);


--
-- Name: physical_cluster physical_cluster_pkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY physical_cluster
    ADD CONSTRAINT physical_cluster_pkey PRIMARY KEY (id);

--
-- Name: logical_cluster logical_cluster_pkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY logical_cluster
    ADD CONSTRAINT logical_cluster_pkey PRIMARY KEY (id);


--
-- Name: cp_component cp_component_pkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY cp_component
    ADD CONSTRAINT cp_component_pkey PRIMARY KEY (id);


--
-- Name: environment environment_pkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY environment
    ADD CONSTRAINT environment_pkey PRIMARY KEY (id);


--
-- Name: k8s_cluster k8s_cluster_pkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY k8s_cluster
    ADD CONSTRAINT k8s_cluster_pkey PRIMARY KEY (id);

--
-- Name: organization organization_pkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY organization
    ADD CONSTRAINT organization_pkey PRIMARY KEY (id);


--
-- Name: entitlement entitlement_pkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY entitlement
    ADD CONSTRAINT entitlement_pkey PRIMARY KEY (id);


--
-- Name: region region_pkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY region
    ADD CONSTRAINT region_pkey PRIMARY KEY (id);


--
-- Name: account account_pkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY account
    ADD CONSTRAINT account_pkey PRIMARY KEY (id);

--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: event event_organization_id_fkey; Type: FK CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY event
    ADD CONSTRAINT event_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES organization(id);

--
-- Name: event event_user_id_fkey; Type: FK CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY event
    ADD CONSTRAINT event_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);


--
-- Name: api_key api_key_user_id_fkey; Type: FK CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY api_key
    ADD CONSTRAINT api_key_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);


--
-- Name: api_key api_key_account_id_fkey; Type: FK CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY api_key
    ADD CONSTRAINT api_key_cluster_id_fkey FOREIGN KEY (cluster_id) REFERENCES logical_cluster(id);


--
-- Name: logical_cluster logical_cluster_account_id_fkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY logical_cluster
    ADD CONSTRAINT logical_cluster_account_id_fkey FOREIGN KEY (account_id) REFERENCES account(id);

-- Name: logical_cluster logical_cluster_deployment_id_fkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY logical_cluster
    ADD CONSTRAINT logical_cluster_deployment_id_fkey FOREIGN KEY (deployment_id) REFERENCES deployment(id);


--
-- Name: logical_cluster logical_cluster_physical_cluster_id_fkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY logical_cluster
    ADD CONSTRAINT logical_cluster_physical_cluster_id_fkey FOREIGN KEY (physical_cluster_id) REFERENCES physical_cluster(id);


--
-- Name: physical_cluster physical_cluster_k8s_cluster_id_fkey; Type: FK CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY physical_cluster
    ADD CONSTRAINT physical_cluster_k8s_cluster_id_fkey FOREIGN KEY (k8s_cluster_id) REFERENCES k8s_cluster(id);

--
-- Name: region region_cloud_fkey; Type: FK CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY region
    ADD CONSTRAINT region_cloud_fkey FOREIGN KEY (cloud) REFERENCES cloud(id);


--
-- Name: account account_organization_id_fkey; Type: FK CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY account
    ADD CONSTRAINT account_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES organization(id);


--
-- Name: users users_organization_id_fkey; Type: FK CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES organization(id);


--
-- Name: connect_task_usage connect_task_usage_organization_id_fkey; Type: CONSTRAINT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY connect_task_usage
    ADD CONSTRAINT connect_task_usage_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES organization(id) NOT VALID;

ALTER TABLE connect_task_usage VALIDATE CONSTRAINT "connect_task_usage_organization_id_fkey";

--
-- Name: support_plan_history; Type: TABLE; Schema: deployment; Owner: caas
--

CREATE TABLE support_plan_history (
    id integer NOT NULL,
    organization_id integer NOT NULL,
    plan_sku character varying(128) NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    modified timestamp without time zone DEFAULT now() NOT NULL,
    effective_date timestamp without time zone DEFAULT now() NOT NULL,
    deleted boolean DEFAULT false NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE support_plan_history OWNER TO caas;

--
-- Name: support_plan_history_id_seq; Type: SEQUENCE; Schema: deployment; Owner: caas
--

CREATE SEQUENCE support_plan_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE support_plan_history_id_seq OWNER TO caas;

ALTER SEQUENCE support_plan_history_id_seq OWNED BY support_plan_history.id;

--
-- Name: support_plan_history id; Type: DEFAULT; Schema: deployment; Owner: caas
--

ALTER TABLE ONLY support_plan_history ALTER COLUMN id SET DEFAULT nextval('support_plan_history_id_seq'::regclass);

--
-- Name: secret_index; Type: TABLE; Schema: deployment; Owner: cass
--

CREATE TABLE secret_physical_cluster_map (
	secret_id integer REFERENCES secret(id) ON DELETE CASCADE,
	physical_cluster_id public.physical_cluster_id REFERENCES physical_cluster(id) ON DELETE CASCADE,
	PRIMARY KEY (secret_id, physical_cluster_id)
);

ALTER TABLE secret_physical_cluster_map OWNER TO caas;

--
-- Name: secret_physical_cluster_map_idx; Type: INDEX; Schema: deployment; Owner: cass
--

CREATE INDEX secret_physical_cluster_map_idx ON secret_physical_cluster_map (physical_cluster_id);


--
-- Name: feature_requests id; Type: DEFAULT; Schema: deployment; Owner: caas
--

CREATE TABLE feature_requests (
    id integer PRIMARY KEY,
    type TEXT DEFAULT '' NOT NULL,
    request TEXT DEFAULT '' NOT NULL,
    organization_id integer NOT NULL,
    user_id integer NOT NUlL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    cloud public.cloud_id,
    region public.region_id
);

ALTER TABLE ONLY feature_requests
    ADD CONSTRAINT feature_requests_organization_id_fkey
    FOREIGN KEY (organization_id)
    REFERENCES organization(id);


ALTER TABLE ONLY feature_requests
    ADD CONSTRAINT feature_requests_user_id_fkey
    FOREIGN KEY (user_id)
    REFERENCES users(id);

ALTER TABLE feature_requests OWNER TO caas;

--
-- Name: feature_requests_id_seq; Type: SEQUENCE; Schema: deployment; Owner: caas
--

CREATE SEQUENCE feature_requests_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE feature_requests_id_seq OWNER TO caas;

--
-- Name: feature_requests_id_seq; Type: SEQUENCE OWNED BY; Schema: deployment; Owner: caas
--

ALTER SEQUENCE feature_requests_id_seq OWNED BY feature_requests.id;

ALTER TABLE ONLY feature_requests ALTER COLUMN id SET DEFAULT nextval('feature_requests_id_seq'::regclass);


--
-- Name: control_plane; Type: SCHEMA; Schema: -; Owner: caas
--

CREATE SCHEMA control_plane;
ALTER SCHEMA control_plane OWNER TO caas;

--
-- Name: upgrade_request; Type: TABLE; Schema: control_plane; Owner: caas
--
CREATE TABLE control_plane.upgrade_request (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL,
    description text NOT NULL, -- preferably a Jira for audit
    request_created timestamp without time zone DEFAULT now() NOT NULL,
    last_modified timestamp without time zone DEFAULT now() NOT NULL,
    cluster_type character varying(32) NOT NULL, -- cluster type
    clusters jsonb DEFAULT '{}'::jsonb NOT NULL, -- ids of the clusters to be upgraded
    options jsonb DEFAULT '{}'::jsonb, -- future support for options/flags like dry run
    status text NOT NULL,
    status_detail text
);
ALTER TABLE control_plane.upgrade_request OWNER TO caas;

--
-- Name: upgrade_task; Type: TABLE; Schema: control_plane; Owner: caas
--
CREATE TABLE control_plane.upgrade_task (
    cluster_id public.physical_cluster_id NOT NULL,
    update_id text, -- A change id or roll id from scheduler
    dedicated_cluster boolean,
    upgrade_id INTEGER NOT NULL,
    cluster_type character varying(32) NOT NULL,
    weight INTEGER, -- used to determine a chunk of similar clusters to upgrade
    upgrade_triggered timestamp without time zone DEFAULT timestamp '1970-01-01 00:00:00.00000'  NOT NULL,
    status text NOT NULL,
    last_modified timestamp without time zone DEFAULT now() NOT NULL,
    target_psc_version text,
    status_detail text,
    PRIMARY KEY (upgrade_id, cluster_id),
    FOREIGN KEY (upgrade_id) REFERENCES control_plane.upgrade_request (id),
    FOREIGN KEY (cluster_id) REFERENCES deployment.physical_cluster (id)
);
ALTER TABLE control_plane.upgrade_task OWNER TO caas;

--
-- Name: skip_upgrade_rules; Type: TABLE; Schema: control_plane; Owner: caas
--
-- List of cluster ids that need to be skipped because they are
-- priority/sensitive customers or we do not want upgrades for them
CREATE TABLE control_plane.skip_upgrade_rules (
    id SERIAL PRIMARY KEY,
    cluster_id public.physical_cluster_id NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    activated timestamp without time zone DEFAULT now()  NOT NULL,
    deactivated timestamp without time zone DEFAULT '2035-12-31 00:00:00',
    description text,
    FOREIGN KEY (cluster_id) REFERENCES deployment.physical_cluster (id)
);
ALTER TABLE control_plane.skip_upgrade_rules OWNER TO caas;

--
-- Storage Class
--

CREATE SEQUENCE IF NOT EXISTS deployment.storage_class_num START WITH 1 INCREMENT 1 NO CYCLE NO MINVALUE NO MAXVALUE;
ALTER TABLE deployment.storage_class_num OWNER TO caas;

CREATE TABLE IF NOT EXISTS deployment.storage_class (
      id              text PRIMARY KEY DEFAULT ('sc-' || nextval('deployment.storage_class_num')::text),
      encryption_key_id TEXT NOT NULL,
      physical_cluster_id TEXT REFERENCES deployment.physical_cluster(id) ON DELETE CASCADE NOT NULL,
      account_id TEXT references deployment.account(id) NOT NULL,
      created timestamp without time zone DEFAULT now() NOT NULL,
      deactivated timestamp without time zone DEFAULT NULL
);

ALTER TABLE deployment.storage_class OWNER TO caas;
CREATE INDEX CONCURRENTLY IF NOT EXISTS storage_class_account_id_idx ON deployment.storage_class (account_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS storage_class_physical_cluster_id_idx ON deployment.storage_class (physical_cluster_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS storage_class_encryption_key_id_idx ON deployment.storage_class (encryption_key_id);

--
-- Cloud Service Accounts.
--

CREATE TABLE IF NOT EXISTS deployment.cloud_service_account (
      id TEXT PRIMARY KEY,
      cloud_id TEXT references deployment.cloud(id) NOT NULL
);

ALTER TABLE deployment.cloud_service_account OWNER TO caas;
CREATE INDEX CONCURRENTLY IF NOT EXISTS cloud_service_account_cloud_id_idx ON deployment.cloud_service_account (cloud_id);

--
-- NetworkRegion
--

CREATE SEQUENCE IF NOT EXISTS deployment.network_region_num START WITH 1 INCREMENT 1 NO CYCLE NO MINVALUE NO MAXVALUE;
ALTER TABLE deployment.network_region_num OWNER TO caas;

CREATE TABLE IF NOT EXISTS deployment.network_region (
    id              text PRIMARY KEY DEFAULT ('nr-' || nextval('deployment.network_region_num')::text),
    name            text DEFAULT NULL,
    description     text DEFAULT NULL,
    requested_cidr  cidr NOT NULL,
    provider        jsonb NOT NULL DEFAULT '{}'::jsonb,
    account_id      text NOT NULL,
    service_network jsonb DEFAULT '{}'::jsonb,
    site_name       text DEFAULT NULL,
    status          jsonb DEFAULT '{}'::jsonb,
    sni_enabled     boolean NOT NULL,
    region_id       text NOT NULL,
    created         timestamp without time zone DEFAULT now() NOT NULL,
    modified        timestamp without time zone DEFAULT now() NOT NULL,
    deactivated     timestamp without time zone DEFAULT NULL,
    dedicated       boolean DEFAULT True
);
ALTER TABLE deployment.network_region OWNER TO caas;

ALTER TABLE ONLY deployment.network_region ADD CONSTRAINT network_region_account_id_fkey FOREIGN KEY (account_id) REFERENCES deployment.account(id) NOT VALID;
ALTER TABLE ONLY deployment.network_region VALIDATE CONSTRAINT network_region_account_id_fkey;
ALTER TABLE ONLY deployment.deployment ADD CONSTRAINT deployment_network_region_id_fkey FOREIGN KEY (network_region_id) REFERENCES deployment.network_region(id) NOT VALID;
ALTER TABLE ONLY deployment.deployment VALIDATE CONSTRAINT deployment_network_region_id_fkey;

ALTER TABLE ONLY deployment.network_region ADD CONSTRAINT network_region_region_id_fkey FOREIGN KEY (region_id) REFERENCES deployment.region(id) NOT VALID;
ALTER TABLE ONLY deployment.network_region VALIDATE CONSTRAINT network_region_region_id_fkey;

CREATE INDEX CONCURRENTLY IF NOT EXISTS network_region_account_id_idx ON deployment.network_region (account_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS network_region_region_id_idx ON deployment.network_region (region_id);
CREATE UNIQUE INDEX IF NOT EXISTS network_region_shared_idx ON deployment.network_region ((provider->>'cloud'), (provider->>'region'), dedicated) WHERE NOT dedicated;

--
-- NetworkConfig
--

CREATE SEQUENCE IF NOT EXISTS deployment.network_config_num START WITH 1 INCREMENT 1 NO CYCLE NO MINVALUE NO MAXVALUE;
ALTER TABLE deployment.network_config_num OWNER TO caas;

CREATE TABLE IF NOT EXISTS deployment.network_config (
    id                  text PRIMARY KEY DEFAULT ('nc-' || nextval('deployment.network_config_num')::text),
    network_region_id   text NOT NULL,
    account_id          text NOT NULL,
    status              jsonb NOT NULL DEFAULT '{}'::jsonb,
    type                text NOT NULL,
    config              jsonb DEFAULT '{}'::jsonb,
    created             timestamp without time zone DEFAULT now() NOT NULL,
    modified            timestamp without time zone DEFAULT now() NOT NULL,
    deactivated         timestamp without time zone DEFAULT NULL,
    shared              boolean DEFAULT False
);
ALTER TABLE deployment.network_config OWNER TO caas;

ALTER TABLE ONLY deployment.network_config ADD CONSTRAINT network_config_network_region_id_fkey FOREIGN KEY (network_region_id) REFERENCES deployment.network_region(id) NOT VALID;
ALTER TABLE ONLY deployment.network_config VALIDATE CONSTRAINT network_config_network_region_id_fkey;

ALTER TABLE ONLY deployment.network_config ADD CONSTRAINT network_config_account_id_fkey FOREIGN KEY (account_id) REFERENCES deployment.account(id) NOT VALID;
ALTER TABLE ONLY deployment.network_config VALIDATE CONSTRAINT network_config_account_id_fkey;

CREATE INDEX CONCURRENTLY IF NOT EXISTS network_config_account_id_idx ON deployment.network_config (account_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS network_config_network_region_id_idx ON deployment.network_config (network_region_id);

--
-- NetworkAccess
--

CREATE SEQUENCE IF NOT EXISTS deployment.network_access_num START WITH 1 INCREMENT 1 NO CYCLE NO MINVALUE NO MAXVALUE;
ALTER TABLE deployment.network_access_num OWNER TO caas;

CREATE TABLE IF NOT EXISTS deployment.network_access (
    id                  text PRIMARY KEY DEFAULT ('na-' || nextval('deployment.network_access_num')::text),
    deployment_id       text NOT NULL,
    type                text NOT NULL,
    config              jsonb DEFAULT '{}'::jsonb,
    network_config_id   text DEFAULT NULL,
    account_id          text NOT NULL,
    created             timestamp without time zone DEFAULT now() NOT NULL,
    modified            timestamp without time zone DEFAULT now() NOT NULL,
    deactivated         timestamp without time zone DEFAULT NULL,
    last_change_id      bigint NOT NULL default 1,
    UNIQUE(deployment_id, type, network_config_id)
);
ALTER TABLE deployment.network_access OWNER TO caas;

ALTER TABLE ONLY deployment.network_access ADD CONSTRAINT network_access_deployment_id_fkey FOREIGN KEY (deployment_id) REFERENCES deployment.deployment(id) NOT VALID;
ALTER TABLE ONLY deployment.network_access VALIDATE CONSTRAINT network_access_deployment_id_fkey;

ALTER TABLE ONLY deployment.network_access ADD CONSTRAINT network_access_network_config_id_fkey FOREIGN KEY (network_config_id) REFERENCES deployment.network_config(id) NOT VALID;
ALTER TABLE ONLY deployment.network_access VALIDATE CONSTRAINT network_access_network_config_id_fkey;

ALTER TABLE ONLY deployment.network_access ADD CONSTRAINT network_access_account_id_fkey FOREIGN KEY (account_id) REFERENCES deployment.account(id) NOT VALID;
ALTER TABLE ONLY deployment.network_access VALIDATE CONSTRAINT network_access_account_id_fkey;

CREATE INDEX IF NOT EXISTS network_access_deployment_id_idx ON deployment.network_access (deployment_id);
CREATE INDEX IF NOT EXISTS network_access_account_id_idx ON deployment.network_access (account_id);
CREATE INDEX IF NOT EXISTS network_access_network_config_id_idx ON deployment.network_access (network_config_id);

-- Function to capture inserts and updates to network_access
CREATE OR REPLACE FUNCTION deployment.update_lc_last_change_id_from_network_access() RETURNS TRIGGER AS $body$
BEGIN
    NEW.last_change_id = nextval('deployment."metadata_change_id_seq"');
    UPDATE deployment.logical_cluster
      SET last_change_id = NEW.last_change_id, modified = NOW()
      WHERE deployment_id = NEW.deployment_id AND type in ('kafka', 'healthcheck', 'ksql');
    RETURN NEW;
END;
$body$
LANGUAGE plpgsql;

-- Trigger to capture inserts and updates to network_access
CREATE TRIGGER network_access_last_change_id_trigger
 BEFORE INSERT OR UPDATE ON deployment.network_access
 FOR EACH ROW EXECUTE PROCEDURE deployment.update_lc_last_change_id_from_network_access();

-- View to simplify JDBC connector config - https://confluentinc.atlassian.net/wiki/spaces/CAAS/pages/1185121071/Logical+Network+Access+Sync+Pipeline
CREATE VIEW deployment.logical_network_access AS SELECT lc.id AS id, lc.type AS type, lc.deactivated AS deactivated, na.network_access, d.network_region_id AS network_region_id, pc.id AS physical_cluster_id, pc.config->'kafka'->>'kafka_api_id' AS pkac_id, pc.k8s_cluster_id AS k8s_cluster_id, lc.modified AS modified, lc.last_change_id AS change_id
 FROM deployment.logical_cluster AS lc JOIN deployment.deployment AS d ON lc.deployment_id = d.id
 JOIN deployment.physical_cluster AS pc ON lc.physical_cluster_id = pc.id
 LEFT JOIN LATERAL(
  SELECT jsonb_object_agg(type, config) AS network_access
  FROM deployment.network_access
  WHERE deployment_id = d.id AND deactivated IS NULL
 ) AS na ON TRUE
 WHERE lc.type in ('kafka', 'healthcheck', 'ksql') AND d.network_region_id IS NOT NULL;

--
-- K8sCluster
--

ALTER TABLE deployment.k8s_cluster ADD CONSTRAINT k8s_cluster_network_region_id_fkey FOREIGN KEY (network_region_id) REFERENCES deployment.network_region (id) NOT VALID;
ALTER TABLE deployment.k8s_cluster VALIDATE CONSTRAINT k8s_cluster_network_region_id_fkey;

--
-- Zone
--

CREATE SEQUENCE IF NOT EXISTS deployment.zone_num START WITH 1 INCREMENT 1 NO CYCLE NO MINVALUE NO MAXVALUE;
ALTER TABLE deployment.zone_num OWNER TO caas;

CREATE TABLE IF NOT EXISTS deployment.zone (
    id              text PRIMARY KEY DEFAULT ('zone-' || nextval('deployment.zone_num')::text),
    zone_id         text NOT NULL,
    name            text NOT NULL,
    region_id       text NOT NULL,
    scheduleable    boolean DEFAULT true NOT NULL,
    created         timestamp without time zone DEFAULT now() NOT NULL,
    modified        timestamp without time zone DEFAULT now() NOT NULL,
    deactivated     timestamp without time zone DEFAULT NULL
);
ALTER TABLE deployment.zone OWNER TO caas;

ALTER TABLE ONLY deployment.zone ADD CONSTRAINT zone_region_id_fkey FOREIGN KEY (region_id) REFERENCES deployment.region(id) NOT VALID;
ALTER TABLE ONLY deployment.zone VALIDATE CONSTRAINT zone_region_id_fkey;

CREATE INDEX CONCURRENTLY IF NOT EXISTS zone_region_id_idx ON deployment.zone (region_id);

INSERT INTO deployment.zone (zone_id, name, region_id, scheduleable, created, modified)
VALUES
    ('usw2-az1', 'us-west-2a', 'us-west-2', true, now(), now()),
    ('usw2-az2', 'us-west-2b', 'us-west-2', true, now(), now()),
    ('usw2-az3', 'us-west-2c', 'us-west-2', true, now(), now()),
    ('us-central1-b', 'us-central1-b', 'us-central1', true, now(), now()),
    ('centralus', 'centralus', 'centralus', true, now(), now());

INSERT INTO deployment.network_region (id, requested_cidr, region_id, provider, account_id, service_network, status, sni_enabled, created)
VALUES
    ('nr-1', '10.0.0.0/16', 'us-west-2', '{"cloud": "aws", "region": "us-west-2", "zones": [{"name": "us-west-2a", "zone_id": "usw2-az1"}, {"name": "us-west-2b", "zone_id": "usw2-az2"}, {"name": "us-west-2c", "zone_id": "usw2-az3"}]}', 't0', '{"aws": {"account_id": "037803949979", "vpc_id": "vpc-958feff3"}}', '{"type": "READY"}', False, now()),
    ('nr-2', '10.1.0.0/16', 'us-west-2', '{"cloud": "aws", "region": "us-west-2", "zones": [{"name": "us-west-2a", "zone_id": "usw2-az1"}, {"name": "us-west-2b", "zone_id": "usw2-az2"}, {"name": "us-west-2c", "zone_id": "usw2-az3"}]}', 't0', '{"aws": {"account_id": "037803949979", "vpc_id": "vpc-abcdef12"}}', '{"type": "READY"}', True, now()),
    ('nr-3', '10.2.0.0/16', 'us-west-2', '{"cloud": "aws", "region": "us-west-2", "zones": [{"name": "us-west-2a", "zone_id": "usw2-az1"}]}', 't0', '{"aws": {"account_id": "037803949979", "vpc_id": "vpc-eff08497"}}', '{"type": "READY"}', False, now()),
    ('nr-4', '10.3.0.0/16', 'us-central1', '{"cloud": "gcp", "region": "us-central1", "zones": [{"name": "us-central1-b", "zone_id": "us-central1-b"}]}', 't0', '{"gcp": {"project_id": "cc-devel", "vpc_network_name": "k8s-test"}}', '{"type": "READY"}', False, now()),
    ('nr-5', '10.4.0.0/16', 'centralus', '{"cloud": "azure", "region": "centralus", "zones": [{"name": "centralus", "zone_id": "centralus"}]}', 't0', '{"azure": {"subscription_id": "a1-b2-c3-d4-e5", "vnet_id": "v-1"}}', '{"type": "READY"}', False, now());

INSERT INTO deployment.k8s_cluster (id, network_region_id, config, created, modified)
VALUES
    ('k8s2', 'nr-1', '{"name": "k8s-mothership.us-west-2", "caas_version": "0.6.10", "is_schedulable": true, "img_pull_policy": "IfNotPresent"}', now(), now()),
    ('k8s3', 'nr-2', '{"name": "k8s3.us-west-2", "caas_version": "0.6.10", "is_schedulable": true, "img_pull_policy": "IfNotPresent"}', now(), now()),
    ('k8s4', 'nr-4', '{"name": "k8s4.us-central1", "caas_version": "0.6.10", "is_schedulable": true, "img_pull_policy": "IfNotPresent"}', now(), now()),
    ('k8s5', 'nr-3', '{"name": "k8s5.us-west-2", "caas_version": "0.6.10", "is_schedulable": true, "img_pull_policy": "IfNotPresent"}', now(), now()),
    ('k8s6', 'nr-5', '{"name": "k8s6.centralus", "caas_version": "0.6.10", "is_schedulable": true, "img_pull_policy": "IfNotPresent"}', now(), now());

--
-- PostgreSQL database dump complete
--
