package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/jlaffaye/ftp"
)

var backFilePath string = "back.json"

////指纹集合文件名称
var hashFileName string = "hash.zlbf"

//var  PATH = "";

//生成sql文件
var ORACLEBAKSQLPATHTL string = "oracle/backdirsql.tl"
var ORACLEBAKSQLPATH string = "bat/backdir.sql"

//生成执行sql文件
var ORACLEBAKPATHTL string = "oracle/backdir.tl"
var ORACLEBAKPATH string = "bat/backdir.bat"

//生成执行备份文件
var ORACLEBAKBATPATHTL string = "oracle/oracle.tl"
var ORACLEBAKBATPATH string = "bat/oracle.bat"

//定时任务
var SCHTASKSPATHTL string = "oracle/schtasks.tl"
var SCHTASKSPATH string = "bat/schtasks.bat"

//start.bat
var STARTPATHTL string = "oracle/start.tl"
var STARTPATH string = "bat/start.bat"

//设置文件
type Backinfo struct {
	TargetPath    string
	BackPath      []string
	FtpIp         string
	FtpUserName   string
	FtpPassWord   string
	OracleBakPath string
	OracleURL     string
}

func main() {

	//dir2 , err2 := getCurrentDirectory()
	//if err2 != nil {
	//	log.Println("获取当前目录失败" + err2.Error() )
	//
	//}

	//
	//log.Println("文件目录： " + dir2 )

	//

	info, err := readBackInfoContent()
	if err != nil {
		log.Println("读取配置文件错误文件内容错误： " + err.Error())

	}

	//initHashfile()

	if !checkFileIsExist(ORACLEBAKSQLPATH) {
		initbak(info)
	}

	cmd := os.Args[1]

	log.Println("cmd " + cmd)

	switch cmd {
	case "o":
		BakOracleBat(info.OracleBakPath)
	case "f":
		BakFiles(info)
	default:
		//log.Println("删除指纹文件错误文件内容错误： " + err.Error())
		BakOracleBat(info.OracleBakPath)
		BakFiles(info)
	}

	os.Exit(0)

}

//
func initHashfile() error {

	t := time.Now()

	if 1 == t.Day() { //第一天
		err := os.Remove(hashFileName) //删除文件test.txt
		if err != nil {
			log.Println("删除指纹文件错误文件内容错误： " + err.Error())
		}
	}

	//if err := createHashFile(); err != nil {
	//	log.Println("读取指纹文件错误文件内容错误： " + err.Error())
	//	return err
	//}

	return nil
}

func initbak(info Backinfo) error {

	OracleBakPath := strings.Replace(info.OracleBakPath, "/", "\\", -1)

	var err = TemplateSaveFile(ORACLEBAKSQLPATHTL, ORACLEBAKSQLPATH, OracleBakPath)
	if err != nil {
		log.Println("生成oracledir.sql 失败" + err.Error())
	}

	dir, err1 := getCurrentDirectory()

	if err1 != nil {
		log.Println("获取当前目录失败" + err1.Error())
	}

	//fmt.Println("ssssssssssss" + dir)

	oracledir := map[string]string{
		"Dir":           dir,
		"OracleBakPath": OracleBakPath,
	}

	err = TemplateSaveFile(ORACLEBAKPATHTL, ORACLEBAKPATH, oracledir)
	if err != nil {
		log.Println("生成oracledir.bat 失败" + err.Error())
	}

	var oracledatatmp []string = strings.Split(info.OracleURL, "@")

	if oracledatatmp == nil || len(oracledatatmp) < 2 {

		log.Println("读取oracle配置信息失败")

	}

	oracleddata := map[string]string{
		"OracleBakPath": OracleBakPath,
		"UserName":      oracledatatmp[0],
		"PassWord":      oracledatatmp[1],
		"DBName":        oracledatatmp[2],
	}

	err = TemplateSaveFile(ORACLEBAKBATPATHTL, ORACLEBAKBATPATH, oracleddata)
	if err != nil {
		log.Println("生成oracle.bat 失败" + err.Error())
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	baktime := fmt.Sprintf("0%d:%d%d", r.Intn(6), r.Intn(5), r.Intn(9))
	log.Println(baktime)

	schtasks := map[string]string{
		"dir":  dir,
		"time": baktime,
	}

	err = TemplateSaveFile(SCHTASKSPATHTL, SCHTASKSPATH, schtasks)
	if err != nil {
		log.Println("生成schtasks.bat 失败" + err.Error())
	}

	err = TemplateSaveFile(STARTPATHTL, STARTPATH, dir)
	if err != nil {
		log.Println("生成start.bat 失败" + err.Error())
	}

	err = execu(SCHTASKSPATH)

	if err != nil {
		log.Println("运行schtasks.bat 失败" + err.Error())
		return err
	}

	return nil

}

//oracleback
func BakOracleBat(oraclepath string) error {

	dir, err := getCurrentDirectory()
	if err != nil {
		log.Println("获取当前目录失败" + err.Error())
		return err
	}

	if !checkFileIsExist(filepath.Join(dir, oraclepath)) {
		err := execu(ORACLEBAKPATH)
		if err != nil {
			log.Println("运行文件失败" + ORACLEBAKPATH + err.Error())
			return err
		}
	}

	err = execu(ORACLEBAKBATPATH)
	if err != nil {
		log.Println("运行文件失败" + ORACLEBAKBATPATH + err.Error())
		return err
	}

	return nil
}

func BakFiles(info Backinfo) error {

	var xcplasttime = time.Now().AddDate(0, 0, -1).Format("01-02-2006")

	var lasttime = time.Now().Format("2006-01-02")

	var lastmoth = time.Now().Format("2006-01")

	if !checkFileIsExist(hashFileName) {
		if err := createHashFile(); err != nil {
			log.Println("读取指纹文件错误文件内容错误： " + err.Error())
			//return err
		}
		xcplasttime = "01-02-2006"
	}

	if err := tarpath(info, lasttime, xcplasttime); err != nil {
		log.Println("复制文件失败" + err.Error())
	}

	if err := zipfiles(info.TargetPath, lasttime); err != nil {
		log.Println("压缩文件失败" + err.Error())
	}

	var remoteSavePath = lastmoth + "^" + strings.Replace(get_external(), ".", "-", -1)
	var oracledatatmp []string = strings.Split(info.OracleURL, "@")

	files, _ := ioutil.ReadDir(info.TargetPath)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			fmt.Println(file.Name())
			ftpUploadFile(info.FtpIp, info.FtpUserName, info.FtpPassWord, filepath.Join(info.TargetPath, file.Name()), remoteSavePath, oracledatatmp[0]+file.Name())
		}
	}

	//var localFile = filepath.Join(info.TargetPath, lasttime+".7z")
	//
	//var oracledatatmp []string = strings.Split(info.OracleURL, "@")
	//

	//
	//log.Println("压缩文件", remoteSavePath, lasttime+".7z")
	//
	//var err = ftpUploadFile(info.FtpIp, info.FtpUserName, info.FtpPassWord, localFile, remoteSavePath, lasttime+oracledatatmp[0]+".7z")
	//
	//if err != nil {
	//	log.Println("上传ftp文件失败" + err.Error())
	//	//ftpUploadFile(info.FtpIp, info.FtpUserName, info.FtpPassWord, localFile, remoteSavePath, lasttime+oracledatatmp[0]+".7z")
	//}

	os.RemoveAll(info.TargetPath)

	//return err
	return nil

}

//读取back.json的配置文件
func readBackInfoContent() (Backinfo, error) {

	file, err := os.Open(backFilePath)

	if err != nil {
		log.Println("读取指纹文件内容错误： " + err.Error())
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Println("close指纹文件 " + backFilePath + " 失败： " + err.Error())
		}
	}()

	jsonContent, err := ioutil.ReadAll(file)

	if err != nil {
		log.Println("读取指纹文件内容错误： " + err.Error())
	}
	//content := make(&backinfo)
	var backinfo Backinfo
	err = json.Unmarshal(jsonContent, &backinfo)
	if err != nil {
		log.Println("指纹文件 json Unmarshal 失败： " + err.Error())
	}

	fmt.Println(backinfo.BackPath[0])

	return backinfo, err

}

//
//创建指纹文件
func createHashFile() error {
	defer func() {
		if err_p := recover(); err_p != nil {
			fmt.Println("createHashFile模块出错")
		}
	}()

	var hashFilePath = filepath.Join(hashFileName)

	if !checkFileIsExist(hashFilePath) {

		err := ioutil.WriteFile(hashFilePath, []byte("{}"), 0777)
		if err != nil {
			log.Println("创建指纹文件失败： " + err.Error())
			return err
		}
	}
	return nil
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//func backpath(list []string ,target string ) error {
//
//	var hashMapContent ,err  =  readFileContent()  // 读取指纹
//
//	if(err != nil) {
//		return err
//	}
//
//
//	for _, value := range list {
//		//var targetPath = ""
//		filepath.Walk(value, func(path string, f os.FileInfo, err error) error {
//
//			var partPath, _ = filepath.Rel(value, path)
//
//			var targetPath = filepath.Join(target, partPath)
//
//			//path:原始文件地址，targetPath:备份文件地址
//			//每个path都需要去比对md5文件中做比对，判断文件是否被修改过
//			//如果文件是个目录则不写入指纹文件
//			if f.IsDir() {
//				copyFile(path, targetPath)
//			} else {
//				md5 := makeFileMd5(path) //获取文件md5
//				isUpdate := comparedFileMd5(hashMapContent, md5, path)
//				//如果修改过则复制文件，并更新md5文件
//				if isUpdate {
//					copyFile(path, targetPath)
//				}
//				//如果没有修改过则不执行任何操作
//			}
//			return nil
//		})
//	}
//
//	var  hashFilePath = filepath.Join(hashFileName)
//	writeFileContent(hashMapContent, hashFilePath)
//	//释放读取的指纹文件内存
//	hashMapContent = nil
//
//	return  nil
//
//}
//
func tarpath(backinfo Backinfo, lasttime, time string) error {

	//var hashMapContent, err = readFileContent() // 读取指纹
	//
	//if err != nil {
	//	return err
	//}

	for _, value := range backinfo.BackPath {

		////var targetPath = ""
		//filepath.Walk(value, func(path string, f os.FileInfo, err error) error {
		//

		var err = xcopy(value, backinfo.TargetPath+"/"+lasttime+"/", time)
		if err != nil {
			return err
		}
		log.Println("执行xcopy： " + value + backinfo.TargetPath + "/" + lasttime + "/")
		//var partPath, _ = filepath.Rel(value, path)
		//
		//var targetPath = filepath.Join(backinfo.TargetPath, time, partPath)
		//
		////path:原始文件地址，targetPath:备份文件地址
		////每个path都需要去比对md5文件中做比对，判断文件是否被修改过
		////如果文件是个目录则不写入指纹文件
		//if f.IsDir() {
		//	copyFile(path, targetPath)
		//} else {
		//	md5 := makeFileMd5(path) //获取文件md5
		//	isUpdate := comparedFileMd5(hashMapContent, md5, path)
		//	//如果修改过则复制文件，并更新md5文件
		//	if isUpdate {
		//		copyFile(path, targetPath)
		//	}
		//	//如果没有修改过则不执行任何操作
		//}
		//return nil
		//})
	}

	//var hashFilePath = filepath.Join(hashFileName)
	//err = writeFileContent(hashMapContent, hashFilePath)
	//if err != nil {
	//	log.Println("写入指纹文件失败" + err.Error())
	//return err
	//
	//}
	////释放读取的指纹文件内存
	//hashMapContent = nil

	return nil

}

func zipfiles(targetPath string, time string) error {

	//压缩文件
	var zip = filepath.Join(targetPath, time+".7z")
	var ziptargetPath = filepath.Join(targetPath, time)

	var err = compress7zip(ziptargetPath, zip)

	if err != nil {
		return err
	}

	////压缩文件
	//var hashFilePath = filepath.Join(hashFileName)
	//err = compress7zip(hashFilePath, ziptargetPath)
	//
	//if err != nil {
	//	return err
	//}

	return nil
}

//
//func copyFile(basePath, targetPath string) {
//	defer func() {
//		if err_p := recover(); err_p != nil {
//			fmt.Println("copyFile模块出错")
//		}
//	}()
//
//	baseStat, err := os.Stat(basePath)
//	if err != nil {
//		log.Panicln("需要备份的文件检测失败，文件出现问题，无法复制")
//		return
//	}
//	//targetStat, err := os.Stat(targetPath)
//	_, err = os.Stat(targetPath)
//	if err != nil {
//		if os.IsNotExist(err) {
//			//如果目标文件不存在
//			if baseStat.IsDir() {
//				//如果缺失的是一个空目录
//				errMkDir := os.MkdirAll(targetPath, 0777)
//				if errMkDir != nil {
//					log.Println("创建目录 " + targetPath + " 失败")
//				}
//			} else {
//				//如果缺失的是一个文件,则复制文件
//				copyFileContent(basePath, targetPath)
//			}
//		} else {
//			return
//		}
//	} else {
//		//如果目标文件存在
//		if baseStat.IsDir() {
//			//如果是一个空目录
//		} else {
//			//如果是一个文件，则复制文件
//			copyFileContent(basePath, targetPath)
//		}
//	}
//
//}
//
////复制文件内容
//func copyFileContent(basePath, targetPath string) {
//	defer func() {
//		if err := recover(); err != nil {
//			fmt.Println("copyFileContent模块出错")
//		}
//	}()
//
//	baseFile, err := os.Open(basePath)
//	if err != nil {
//		log.Println("读取文件 " + basePath + " 失败")
//		return
//	}
//	defer func() {
//		err := baseFile.Close()
//		if err != nil {
//			log.Println("close文件 " + basePath + " 失败： " + err.Error())
//		}
//	}()
//	targetFile, err := os.Create(targetPath)
//	if err != nil {
//		log.Println("创建文件 " + targetPath + " 失败： " + err.Error())
//		return
//	}
//	defer func() {
//		err := targetFile.Close()
//		if err != nil {
//			log.Println("close文件 " + targetPath + " 失败： " + err.Error())
//		}
//	}()
//	copyData, err := io.Copy(targetFile, baseFile)
//	if err != nil {
//		log.Println("复制文件文件 " + basePath + " 失败： " + err.Error())
//	}
//	fmt.Println("正在复制文件： " + basePath + " 大小为： " + strconv.FormatInt(copyData, 10))
//}
//
////读取整个指纹文件到内存
//func readFileContent() (*map[string]uint32, error) {
//	defer func() {
//		if err := recover(); err != nil {
//			log.Println("readFileContent模块出错")
//		}
//	}()
//	var hashFilePath = filepath.Join(hashFileName)
//
//	file, err := os.Open(hashFilePath)
//
//	if err != nil {
//		log.Println("读取指纹文件内容错误： " + err.Error())
//	}
//
//	defer func() {
//		err := file.Close()
//		if err != nil {
//			log.Println("close指纹文件 " + hashFilePath + " 失败： " + err.Error())
//		}
//	}()
//
//	jsonContent, err := ioutil.ReadAll(file)
//
//	if err != nil {
//		log.Println("读取指纹文件内容错误： " + err.Error())
//	}
//	content := make(map[string]uint32)
//	err = json.Unmarshal(jsonContent, &content)
//	if err != nil {
//		log.Println("指纹文件 json Unmarshal 失败： " + err.Error())
//	}
//	return &content, err
//}
//
////生成文件的md5
//func makeFileMd5(filePath string) uint32 {
//
//	defer func() {
//		if err_p := recover(); err_p != nil {
//			fmt.Println("makeFileMd5模块出错")
//		}
//	}()
//
//	b, err := ioutil.ReadFile(filePath)
//
//	if err != nil {
//
//		log.Println("makefilemd5读取文件失败： " + err.Error())
//
//		//return 0
//		//_, ok := err.(*os.PathError)
//		//if ok {
//		//	log.Println("指纹文件 json Unmarshal 失败： " + err.Error())
//		//} else {
//		//	log.Println("指纹文件 json Unmarshal 失败： " + err.Error())
//		//}
//	}
//
//	return adler32.Checksum(b)
//}
//
////比对指纹文件中的md5和新读取文件的md5
//func comparedFileMd5(mapContent *map[string]uint32, md5 uint32, path string) bool {
//	defer func() {
//		if err := recover(); err != nil {
//			log.Println("comparedFileMd5模块出错")
//		}
//	}()
//
//	if contentMd5, ok := (*mapContent)[path]; ok {
//		//如果md5存在，且不相同，则代表文件更新过，更新md5值，并且复制文件
//		if md5 != contentMd5 {
//			(*mapContent)[path] = md5
//			return true
//		} else {
//			return false
//		}
//	} else {
//		//如果md5不存在，则写入新的path
//		(*mapContent)[path] = md5
//		return true
//	}
//}
//
////写入指纹文件
//func writeFileContent(mapContent *map[string]uint32, path string) error {
//	defer func() {
//		if err := recover(); err != nil {
//			log.Println("writeFileContent模块出错")
//		}
//	}()
//
//	jsonContent, err := json.Marshal(*mapContent)
//
//	if err != nil {
//		log.Println("指纹文件 json Marshal 失败： " + err.Error())
//		return err
//	}
//	err = ioutil.WriteFile(path, jsonContent, 0777)
//	if err != nil {
//		log.Println("写入指纹文件失败： " + err.Error())
//		return err
//	}
//	return nil
//}

//调用7zip压缩
func compress7zip(frm, dst string) error {
	cmd := exec.Command("7z/7z.exe", "a", "-mx=1", "-v5g", dst, frm)
	//cmd.Args = []string{"a",dst,frm};
	//cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("执行7zip压缩命令错误： " + err.Error())
		//log.Fatal(err)
		return err
	}
	fmt.Printf("in all caps: %q\n", out.String())
	return nil
}

//调用7zip压缩
func xcopy(frm, dst, time string) error {

	frm = strings.Replace(frm, "/", "\\", -1)
	dst = strings.Replace(dst, "/", "\\", -1)
	cmd := exec.Command("xcopy", frm, dst, "/s", "/e", "/y", "/d:"+time)
	//cmd.Args = []string{"a",dst,frm};
	//cmd.Stdin = strings.NewReader("some input")
	//var out bytes.Buffer
	//cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("执行xcopy压缩命令错误： " + err.Error())
		//log.Fatal(err)
		return err
	}
	//fmt.Printf("in all caps: %q\n", out.String())
	return nil
}

func execu(name string) error {

	dir, err := getCurrentDirectory()
	if err != nil {
		log.Println("获取当前目录失败" + err.Error())
		return err
	}

	cmd := exec.Command("cmd", "/C", dir+"\\"+name)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}
	//fmt.Println("Result: " + out.String())
	return nil
}

//
//// 参数frm可以是文件或目录，不会给dst添加.zip扩展名
//func compress(frm, dst string) error {
//	buf := bytes.NewBuffer(make([]byte, 0, 10*1024*1024)) // 创建一个读写缓冲
//	myzip := zip.NewWriter(buf)                           // 用压缩器包装该缓冲
//	// 用Walk方法来将所有目录下的文件写入zip
//	err := filepath.Walk(frm, func(path string, info os.FileInfo, err error) error {
//		var file []byte
//		if err != nil {
//			return filepath.SkipDir
//		}
//		header, err := zip.FileInfoHeader(info) // 转换为zip格式的文件信息
//		if err != nil {
//			return filepath.SkipDir
//		}
//		header.Name, _ = filepath.Rel(filepath.Dir(frm), path)
//		if !info.IsDir() {
//			// 确定采用的压缩算法（这个是内建注册的deflate）
//			header.Method = 8
//			file, err = ioutil.ReadFile(path) // 获取文件内容
//			if err != nil {
//				return filepath.SkipDir
//			}
//		} else {
//			file = nil
//		}
//		// 上面的部分如果出错都返回filepath.SkipDir
//		// 下面的部分如果出错都直接返回该错误
//		// 目的是尽可能的压缩目录下的文件，同时保证zip文件格式正确
//		w, err := myzip.CreateHeader(header) // 创建一条记录并写入文件信息
//		if err != nil {
//			return err
//		}
//		_, err = w.Write(file) // 非目录文件会写入数据，目录不会写入数据
//		if err != nil {        // 因为目录的内容可能会修改
//			return err // 最关键的是我不知道咋获得目录文件的内容
//		}
//		return nil
//	})
//	if err != nil {
//		return err
//	}
//	myzip.Close()               // 关闭压缩器，让压缩器缓冲中的数据写入buf
//	file, err := os.Create(dst) // 建立zip文件
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//	_, err = buf.WriteTo(file) // 将buf中的数据写入文件
//	if err != nil {
//		return err
//	}
//	return nil
//}

func ftpUploadFile(ftpserver, ftpuser, pw, localFile, remoteSavePath, saveName string) error {
	ftp, err := ftp.Connect(ftpserver)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	err = ftp.Login(ftpuser, pw)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	//
	//err = ftp.Delete(remoteSavePath +"/"+ saveName)
	//
	//if err != nil {
	//	fmt.Println("删除文件失败")
	//}

	err = ftp.ChangeDir(remoteSavePath)
	if err != nil {
		fmt.Println(err)

		//return  nil
		err = ftp.MakeDir(remoteSavePath)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		ftp.ChangeDir(remoteSavePath)
	}

	dir, err := ftp.CurrentDir()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println(dir)

	file, err := os.Open(localFile)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()
	err = ftp.Stor(saveName, file)
	if err != nil {
		//ftp.Delete(saveName)
		fmt.Println(err)
		return err
	}

	ftp.Logout()

	ftp.Quit()

	log.Println("success upload file:", localFile)

	return err
}

func get_external() string {
	data := ""
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		data = GetIntranetIp()
	}
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(resp.Body)
	//s := buf.String()
	data = strings.Replace(string(content), ".", "-", -1)
	data = strings.Replace(data, ":", "-", -1)
	return data
}

func GetIntranetIp() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {

				return ipnet.IP.String()
				//fmt.Println("ip:", ipnet.IP.String())
			}

		}
	}
	return ""
}

func getCurrentDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return dir, nil
}

func TemplateSaveFile(tlpath, savepath string, datas interface{}) error {

	tmpl, err := TemplateInit(tlpath)

	if err != nil {
		log.Println("加载模板：" + tlpath + err.Error())
		return err
	}

	data, err := TemplateExecute(tmpl, datas)

	if err != nil {
		log.Println("生成模板：" + tlpath + err.Error())
		return err
	}

	f, err := os.Create(savepath) //创建文件

	if err != nil {
		log.Println("打开：" + savepath + err.Error())
		return err
	}

	defer f.Close()

	n, err := io.WriteString(f, data) //写入文件(字符串)

	fmt.Printf("写入 %d 个字节n", n)

	if err != nil {
		log.Println("生成：" + savepath + err.Error())
		return err
	}

	return nil
}

func TemplateInit(templatePath string) (*template.Template, error) {

	content, err := ioutil.ReadFile(templatePath)

	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("test").Parse(string(content))
	if err != nil {
		fmt.Printf("file error: %v", err)
		return nil, err
	}
	return tmpl, nil
}

func TemplateExecute(tmpl *template.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, data)
	if err != nil {
		fmt.Printf("tmplate error: %v", err)
		return "", err
	}
	return buf.String(), nil
}
