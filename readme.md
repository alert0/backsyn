
{
  "TargetPath":"d:/filebak",
  "BackPath":["D:/newhztfiles","d:/oraclebak","D:/newhzt/tomcat/webapps/ROOT/src","D:/newhzt/tomcat/webapps/ROOT/WEB-INF/lib"],
  "FtpIp" : "xxxx:1414",
  "FtpUserName": "back",
  "FtpPassWord": "qwerty123",
  "OracleBakPath": "d:/oraclebak",
  "OracleURL": "newhzt@xxxx@orcl"
}

修改配置文件 back.json 

BackPath = 要备份的目录的，支持多文件夹。
OracleURL = name@password@dbname

oracle 每天备份一次，文件每个月1号完全备份一次，其他每天增量备份。


