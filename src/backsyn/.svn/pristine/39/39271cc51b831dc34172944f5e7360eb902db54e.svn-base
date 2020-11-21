@echo off


forfiles /p {{.OracleBakPath}} /m *.dmp -d -3  /C "cmd /c del /f @path\"

forfiles /p {{.OracleBakPath}} /m *.7z -d -3  /C "cmd /c del /f @path\"

set BACKUP_DIR={{.OracleBakPath}}

rem delete 30days files

rem rd /s /Q %BACKUP_DIR%

md %BACKUP_DIR%

cd %BACKUP_DIR%

rem backup schemas


set backupfile={{.OracleBakPath}}\%date:~0,4%-%date:~5,2%-%date:~8,2%.dmp
set backupfilezip=%date:~0,4%-%date:~5,2%-%date:~8,2%.7z


{{.UserName}}

{{.Dir}}\7z\7z.exe a  {{.OracleBakPath}}\%backupfilezip%  %backupfile% 


del /f   %backupfile% 


exit