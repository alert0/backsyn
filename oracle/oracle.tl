@echo off


set BACKUP_DIR={{.OracleBakPath}}

rem delete 30days files

rd /s /Q %BACKUP_DIR%
md %BACKUP_DIR%

cd %BACKUP_DIR%

rem backup schemas

set ORACLE_USERNAME={{.UserName}}
set ORACLE_PASSWORD={{.PassWord}}
set ORACLE_DB={{.DBName}}
set BACK_OPTION="SCHEMAS={{.UserName}}"
set backupfile=%ORACLE_USERNAME%_%date:~0,4%-%date:~5,2%-%date:~8,2%.dmp
set logfile=%ORACLE_USERNAME%_%date:~0,4%-%date:~5,2%-%date:~8,2%.log
expdp %ORACLE_USERNAME%/%ORACLE_PASSWORD%@%ORACLE_DB%  directory=bak  dumpfile=%backupfile%   logfile=%logfile%  schemas=%ORACLE_USERNAME%


exit