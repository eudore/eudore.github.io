package main
  
import (
    "os"
    "fmt"
    "net"
    "flag"
    "time"
    "crypto/md5"
    "io"
    "strings"
    "strconv"
)
  
type sysFileInfo struct{
    fName       string
    fSize       int64
    fMtime      time.Time
    fPerm       os.FileMode
    fMd5        string
    fType       bool
}
  
var (
    listenPort = flag.String( "port","9081","server listen port" )
    syncFile = flag.String( "file","","transfer file" )
    syncHost = flag.String( "host","127.0.0.1","server host" )
    syncSer = flag.Bool( "d",false,"server mode")
    syncFold = flag.String( "dir","/tmp/gosync/","recive sync fold ")
)
  
func main(){
    flag.Parse()
    if *syncSer { 
        servPort:=fmt.Sprintf( ":%s",*listenPort )
        l,err := net.Listen( "tcp",servPort )
        if err != nil{
           fmt.Println( "net failed",err )
        }
        err = os.MkdirAll( *syncFold , 0755)
        if err != nil{
           fmt.Println( err )
        }
        fmt.Println( "Start Service" )
        Serve( l )
     }else{ 
        destination:=fmt.Sprintf( "%s:%s",*syncHost,*listenPort )
        clientSend( *syncFile,destination)
     }
}
  
func clientSend(files string,destination string){
    fInfo:=getFileInfo( files)
    newName :=fmt.Sprintf( "%s",fInfo.fName)
    cmdLine:=  fmt.Sprintf( "upload %s %d %d %d %s " ,newName,fInfo.fMtime.Unix(),fInfo.fPerm,fInfo.fSize,fInfo.fMd5)
    cn,err:=net.Dial( "tcp", destination)
    if err !=nil {
        fmt.Println( "connect error",err )
        return
    }
    defer cn.Close()
    cn.Write( []byte( cmdLine ) )
    cn.Write( []byte( "\r\n" ) )
    fileHandle,err := os.Open( files )
    if err != nil {
        fmt.Println("open ERROR",err)
        return
    }
    io.Copy( cn,fileHandle)
    for{
        buffer :=make( []byte,1024)
        num,err := cn.Read(buffer)
        if err == nil && num > 0{
            fmt.Println(  string(buffer[ :num ]) )
            break
        }
    }
}
  
func getFileInfo( filename string) *sysFileInfo{
    fi,err:= os.Lstat( filename )
    if err != nil {
        fmt.Println("info ERROR",err)
        return nil
    }
    fileHandle,err := os.Open( filename )
    if err != nil {
        fmt.Println("open ERROR",err)
        return nil
    }
  
    h := md5.New()
    _,err = io.Copy( h,fileHandle )
    fileInfo := & sysFileInfo {
        fName : fi.Name(),
        fSize : fi.Size(),
        fPerm : fi.Mode().Perm(),
        fMtime: fi.ModTime(),
        fType : fi.IsDir(),
        fMd5  : fmt.Sprintf( "%x", h.Sum( nil )),
    }
        return fileInfo
}
  
func Serve( l net.Listener) {
    for{
        conn,err := l.Accept()
        if err != nil{
            if ne,ok := err.( net.Error );ok && ne.Temporary(){
                continue
            }
            fmt.Println( "network error",err )
        }
        go Handler(conn)
    }
}
  
func Handler( conn net.Conn) {
    defer conn.Close()
    state := 0
    var cmd *sysFileInfo
    var fSize int64
    var tempFileName string
    var n int64
    for {
        buffer :=make( []byte,2048)
        num,err := conn.Read(buffer)
        numLen:=int64( num )
        if err != nil && err != io.EOF {
            fmt.Println( "cannot read",err )
        }
        n=0
        if state  == 0 {
            n,cmd = cmdParse( buffer[:num] )
            tempFileName = fmt.Sprintf( "%s.newsync",cmd.fName)
            fSize = cmd.fSize
            state = 1
        }
        if state == 1 {
            last := numLen
            if fSize <= numLen-n {
                last = fSize + n
                state = 2
            }
            err = writeToFile( buffer[int( n ):int( last )],tempFileName,cmd.fPerm )
            if err != nil{
                fmt.Println( "read num error : ",err )
            }
            fSize -=last-n
            if state == 2{
                os.Remove( cmd.fName)
                err = os.Rename( tempFileName,cmd.fName)
                if err != nil{
                    fmt.Println( "rename ",tempFileName," to ",cmd.fName," failed" )
                }
                err = os.Chtimes( cmd.fName,time.Now(),cmd.fMtime )
                if err != nil{
                    fmt.Println( "change the mtime error ",err )
                }
                fileHandle,err := os.Open( cmd.fName)
                if err != nil {
                    fmt.Println("open ERROR",err)
                }
                h := md5.New()
                io.Copy( h,fileHandle )
                newfMd5 := fmt.Sprintf( "%x", h.Sum( nil ))
                if newfMd5 == cmd.fMd5{
                    sendInfo:=fmt.Sprintf("%s sync success",cmd.fName)
                    conn.Write([]byte(sendInfo))
                }else{
                    sendInfo:=fmt.Sprintf("%s sync failed",cmd.fName)
                    conn.Write([]byte(sendInfo))
                }
            }
        }
    }
}
  
func cmdParse( infor []byte) ( int64 , *sysFileInfo) {
    var i int64
    for i=0;i<int64(len(infor));i++ {
       if infor[i] == '\n' && infor[i-1]  == '\r' {
           cmdLine:=strings.Split( string( infor[:i-1] ) ," ")
           fileName := fmt.Sprintf( "%s/%s",*syncFold,cmdLine[ 1 ] )
           filePerm, _ := strconv.Atoi( cmdLine[ 3 ])
           fileMtime,_:= strconv.ParseInt( cmdLine[ 2 ],10,64 )
           fileSize,_:= strconv.ParseInt( cmdLine[ 4 ],10,64)
           fileInfo := & sysFileInfo {
                fName : fileName,
                fMtime: time.Unix( fileMtime,0 ),
                fPerm : os.FileMode(filePerm),
                fSize : fileSize,
                fMd5  : string(cmdLine[ 5 ]),
           }
           return i+1,fileInfo
       }
    }
       return 0,nil
}
  
func writeToFile( data []byte ,fileName string,perm os.FileMode)  error{
    writeFile,err := os.OpenFile( fileName,os.O_RDWR | os.O_APPEND | os.O_CREATE ,perm)
    if err != nil{
        fmt.Println( "write file error:",err )
        return err
    }
    defer writeFile.Close()
    _,err = writeFile.Write( data )
    if err != nil{
       fmt.Println( "write file error",err )
       return err
    }
    return nil
}

