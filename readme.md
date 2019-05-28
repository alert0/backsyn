
{
  "TargetPath":"d:/filebak",
  "BackPath":["D:/newhztfiles","d:/oraclebak","D:/newhzt/tomcat/webapps/ROOT/src","D:/newhzt/tomcat/webapps/ROOT/WEB-INF/lib"],
  "FtpIp" : "xxxx:1414",
  "FtpUserName": "back",
  "FtpPassWord": "ssss",
  "OracleBakPath": "d:/oraclebak",
  "OracleURL": "newhzt@xxxx@orcl"
}

修改配置文件 back.json 

BackPath = 要备份的目录的，支持多文件夹。
OracleURL = name@password@dbname

oracle 每天备份一次，文件每个月1号完全备份一次，其他每天增量备份。


###mysql 备份


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


{{.Path}}  -u{{.UserName}} -p{{.PassWord}}   --databases  {{.DBName}} > %backupfile%   

exit


