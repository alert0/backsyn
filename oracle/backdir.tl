mkdir {{.OracleBakPath}}
sqlplus / as sysdba @{{.Dir}}\bat\backdir.sql
exit
