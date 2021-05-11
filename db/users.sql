--
-- PostgreSQL database dump
--

-- Dumped from database version 10.11
-- Dumped by pg_dump version 10.11

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: fcapp_user_role_permissions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.fcapp_user_role_permissions (
    id integer NOT NULL,
    user_role integer NOT NULL,
    sys_module text NOT NULL,
    sys_perms character varying(16) NOT NULL
);


ALTER TABLE public.fcapp_user_role_permissions OWNER TO postgres;

--
-- Name: fcapp_user_role_permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.fcapp_user_role_permissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.fcapp_user_role_permissions_id_seq OWNER TO postgres;

--
-- Name: fcapp_user_role_permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.fcapp_user_role_permissions_id_seq OWNED BY public.fcapp_user_role_permissions.id;


--
-- Name: fcapp_user_roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.fcapp_user_roles (
    id integer NOT NULL,
    name text NOT NULL,
    descr text DEFAULT ''::text,
    cdate timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.fcapp_user_roles OWNER TO postgres;

--
-- Name: fcapp_user_roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.fcapp_user_roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.fcapp_user_roles_id_seq OWNER TO postgres;

--
-- Name: fcapp_user_roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.fcapp_user_roles_id_seq OWNED BY public.fcapp_user_roles.id;


--
-- Name: fcapp_users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.fcapp_users (
    id bigint NOT NULL,
    user_role integer NOT NULL,
    firstname text DEFAULT ''::text NOT NULL,
    lastname text DEFAULT ''::text NOT NULL,
    username text NOT NULL,
    telephone text DEFAULT ''::text NOT NULL,
    password text NOT NULL,
    email text DEFAULT ''::text NOT NULL,
    allowed_ips text DEFAULT '127.0.0.1;::1'::text NOT NULL,
    denied_ips text DEFAULT ''::text NOT NULL,
    failed_attempts text DEFAULT ('0/'::text || to_char(now(), 'yyyymmdd'::text)),
    transaction_limit text DEFAULT ('0/'::text || to_char(now(), 'yyyymmdd'::text)),
    is_active boolean DEFAULT true NOT NULL,
    is_system_user boolean DEFAULT false NOT NULL,
    last_login timestamp with time zone,
    last_passwd_update timestamp with time zone DEFAULT now() NOT NULL,
    created timestamp without time zone DEFAULT now() NOT NULL,
    updated timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.fcapp_users OWNER TO postgres;

--
-- Name: fcapp_users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.fcapp_users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.fcapp_users_id_seq OWNER TO postgres;

--
-- Name: fcapp_users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.fcapp_users_id_seq OWNED BY public.fcapp_users.id;


--
-- Name: fcapp_user_role_permissions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fcapp_user_role_permissions ALTER COLUMN id SET DEFAULT nextval('public.fcapp_user_role_permissions_id_seq'::regclass);


--
-- Name: fcapp_user_roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fcapp_user_roles ALTER COLUMN id SET DEFAULT nextval('public.fcapp_user_roles_id_seq'::regclass);


--
-- Name: fcapp_users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fcapp_users ALTER COLUMN id SET DEFAULT nextval('public.fcapp_users_id_seq'::regclass);


--
-- Data for Name: fcapp_user_role_permissions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.fcapp_user_role_permissions (id, user_role, sys_module, sys_perms) FROM stdin;
\.


--
-- Data for Name: fcapp_user_roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.fcapp_user_roles (id, name, descr, cdate) FROM stdin;
1	Administrator	For the Administrators	2018-12-10 18:52:18.353825+03
2	API User	For the API Users	2018-12-10 18:52:18.353825+03
\.


--
-- Data for Name: fcapp_users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.fcapp_users (id, user_role, firstname, lastname, username, telephone, password, email, allowed_ips, denied_ips, failed_attempts, transaction_limit, is_active, is_system_user, last_login, last_passwd_update, created, updated) FROM stdin;
1	1	Samuel	Sekiwere	admin	+256753475676	$2a$06$nmxCe6rjlFTQXE5QtEviDeDoy29NGocTVNzGHHPub3pSVngcsKude	sekiskylink@gmail.com	127.0.0.1;::1		0/20181210	0/20181210	t	t	\N	2018-12-10 18:53:45.796145+03	2018-12-10 18:53:45.796145	2018-12-10 18:53:45.796145+03
2	2	Ivan	Muguya	ivan	+256756253430	$2a$06$nSDJTIsy8HO4G9Jc8gkSlOApIrn79KO8IfSzXmxVuyFXB.rJZFu9C	ivanupsons@gmail.com	127.0.0.1;::1		0/20181210	0/20181210	t	t	\N	2018-12-10 18:53:45.796145+03	2018-12-10 18:53:45.796145	2018-12-10 18:53:45.796145+03
\.


--
-- Name: fcapp_user_role_permissions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.fcapp_user_role_permissions_id_seq', 1, false);


--
-- Name: fcapp_user_roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.fcapp_user_roles_id_seq', 2, true);


--
-- Name: fcapp_users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.fcapp_users_id_seq', 2, true);


--
-- Name: fcapp_user_role_permissions fcapp_user_role_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fcapp_user_role_permissions
    ADD CONSTRAINT fcapp_user_role_permissions_pkey PRIMARY KEY (id);


--
-- Name: fcapp_user_role_permissions fcapp_user_role_permissions_sys_module_user_role_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fcapp_user_role_permissions
    ADD CONSTRAINT fcapp_user_role_permissions_sys_module_user_role_key UNIQUE (sys_module, user_role);


--
-- Name: fcapp_user_roles fcapp_user_roles_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fcapp_user_roles
    ADD CONSTRAINT fcapp_user_roles_name_key UNIQUE (name);


--
-- Name: fcapp_user_roles fcapp_user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fcapp_user_roles
    ADD CONSTRAINT fcapp_user_roles_pkey PRIMARY KEY (id);


--
-- Name: fcapp_users fcapp_users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fcapp_users
    ADD CONSTRAINT fcapp_users_pkey PRIMARY KEY (id);


--
-- Name: fcapp_users fcapp_users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fcapp_users
    ADD CONSTRAINT fcapp_users_username_key UNIQUE (username);


--
-- Name: fcapp_user_roles_idx1; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fcapp_user_roles_idx1 ON public.fcapp_user_roles USING btree (name);


--
-- Name: fcapp_users_idx1; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fcapp_users_idx1 ON public.fcapp_users USING btree (telephone);


--
-- Name: fcapp_users_idx2; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fcapp_users_idx2 ON public.fcapp_users USING btree (username);


--
-- Name: fcapp_user_role_permissions fcapp_user_role_permissions_user_role_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fcapp_user_role_permissions
    ADD CONSTRAINT fcapp_user_role_permissions_user_role_fkey FOREIGN KEY (user_role) REFERENCES public.fcapp_user_roles(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: fcapp_users fcapp_users_user_role_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fcapp_users
    ADD CONSTRAINT fcapp_users_user_role_fkey FOREIGN KEY (user_role) REFERENCES public.fcapp_user_roles(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--

