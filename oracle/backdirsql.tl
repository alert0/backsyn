create or replace directory bak as '{{.}}';
grant read,write on directory bak to public;
exit