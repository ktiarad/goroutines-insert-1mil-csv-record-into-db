CREATE DATABASE goinsertmil;

\c goinsertmil

CREATE TABLE domains (
    id serial primary key,
    global_rank int,
    tld_rank int,
    domain varchar(255),
    tld varchar(255),
    ref_sub_nets int,
    ref_ips int,
    idn_domain varchar(255),
    idn_tld varchar(255),
    prev_global_rank int,
    prev_tld_rank int,
    prev_ref_sub_nets int,
    prev_ref_ips int
);