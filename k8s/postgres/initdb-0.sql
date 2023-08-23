--
CREATE ROLE "user" WITH
  LOGIN
  PASSWORD 'SCRAM-SHA-256$4096:Pg6vHBGELv4BGswAdacNvA==$kyU1Z+FbiuNC3ZQhFyWdd6CB/cgVMWiYCfsPhh4abuQ=:zC+ex3AVvAjeYl3a6xsz8MPCMLJIdcvj5xOAi2+dLD0=';

CREATE DATABASE "database" WITH
    OWNER "user";

COMMENT ON
    DATABASE "database"
    IS 'sample database for name service';
