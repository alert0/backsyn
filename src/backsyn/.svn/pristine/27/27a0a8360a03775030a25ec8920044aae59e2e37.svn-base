{{.Dir}}\FastCopy\FastCopy.exe /cmd=diff  /auto_close  /force_close  /from_date="-1D"  {{.BackPath}} {{.OracleBakPath}} /to="{{.ToPath}}" 

set backupfilezip=%date:~0,4%-%date:~5,2%-%date:~8,2%.7z

{{.Dir}}\7z\7z.exe a -mx=1 -v5g   {{.ToPath}}\%backupfilezip%  {{.ToPath}} 

