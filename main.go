package main

import (
	//"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
	//"path/filepath"
	//"bufio"
	//"strconv"
	//"container/list"
	//"sort"
	//"path/filepath"
	"sync"
)

//学生成绩结构体
type StuScore struct {
	//姓名
	name string
	//成绩
	score int
}

type StuScores []StuScore

//Len()
func (s StuScores) Len() int {
	return len(s)
}

//Less():成绩将有低到高排序
func (s StuScores) Less(i, j int) bool {
	return s[i].score < s[j].score
}

//Swap()
func (s StuScores) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var locker = new(sync.Mutex)
var condn = sync.NewCond(locker)

func test(x int) {
	condn.L.Lock() // 获取锁
	condn.Wait()   // 等待通知  暂时阻塞
	fmt.Println(x)
	//time.Sleep(time.Second * 1)
	condn.L.Unlock() // 释放锁，不释放的话将只会有一次输出
}
func main() {
	for i := 0; i < 2; i++ {
		go test(i)
	}
	fmt.Println("start all")
	condn.Broadcast() //  下发广播给所有等待的goroutine
	time.Sleep(time.Second * 5)
	/*stus := StuScores{
		{"alan", 95},
		{"hikerell", 91},
		{"acmfly", 96},
		{"leao", 90}}

	fmt.Println("Default:")
	//原始顺序
	for _, v := range stus {
		fmt.Println(v.name, ":", v.score)
	}
	fmt.Println()
	//StuScores已经实现了sort.Interface接口
	sort.Sort(stus)

	fmt.Println("Sorted:")
	//排好序后的结构
	for _, v := range stus {
		fmt.Println(v.name, ":", v.score)
	}

	//判断是否已经排好顺序，将会打印true
	fmt.Println("IS Sorted?", sort.IsSorted(stus))

	list := list.New()
	list.PushBack(1)
	list.PushBack(2)

	//fmt.Printf("len: %v\n", list.Len())
	//fmt.Printf("first: %#v\n", list.Front().Value)
	//fmt.Printf("second: %#v\n", list.Front().Next())*/
}

func Pipe() {
	io.Copy(os.Stdin, os.Stdout)
	fmt.Println("EOF-----")
}

func PipeWrite(pipeWriter *io.PipeWriter) {
	var (
		i   = 0
		err error
		n   int
	)
	data := []byte("Go语言学习园地")
	for _, err = pipeWriter.Write(data); err == nil; n, err = pipeWriter.Write(data) {
		i++
		if i == 3 {
			pipeWriter.CloseWithError(errors.New("输出3次后结束"))
		}
	}
	fmt.Println("close 后输出的字节数：", n, " error：", err)
}

func PipeRead(pipeReader *io.PipeReader) {
	var (
		err error
		n   int
	)
	data := make([]byte, 1024)
	for n, err = pipeReader.Read(data); err == nil; n, err = pipeReader.Read(data) {
		fmt.Printf("%s\n", data[:n])
	}
	fmt.Println("writer 端 closewitherror 后：", err)
}

func read_dir(dir string, dep int) error {
	dep += 1
	sp := "   "
	f_header := "|--"
	//d_header := "|"
	dir_sp := string(os.PathSeparator)
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, i := range rd {
		if !i.IsDir() {
			fmt.Println(strings.Repeat(sp, dep) + f_header + i.Name())
			continue
		}
		fmt.Println(strings.Repeat(sp, dep) + f_header + i.Name())
		read_dir(dir+dir_sp+i.Name(), dep)
	}
	return nil
}

/*import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var sema = make(chan struct{}, 2000)

func main() {
	t := time.Tick(time.Second * 5)
	for {
		select {
		case <-t:
			print_root()
		}
	}
}
func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

func print_root() {
	var n sync.WaitGroup
	fmt.Println(time.Now())
	// Determine the initial directories.
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse the file tree.
	fileSizes := make(chan int64)
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes)
	}

	go func() {
		n.Wait()
		close(fileSizes)
	}()

	// Print the results.
	var nfiles, nbytes int64
	for size := range fileSizes {
		nfiles++
		nbytes += size
	}
	printDiskUsage(nfiles, nbytes)
	fmt.Println(time.Now())
}

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	sema <- struct{}{} // acquire token
	defer func() { <-sema }()
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

/*type Address struct {
	Type    string
	City    string
	Country string
}

type Card struct {
	Name      string
	Age       int
	Addresses []*Address
}

func main() {
	pa := &Address{"private", "Shanghai", "China"}
	pu := &Address{"work", "Beijing", "China"}
	c := Card{"Xin", 32, []*Address{pa, pu}}

	js, _ := json.Marshal(c)
	fmt.Printf("Json: %s", js)
}*/

/*package main

import (
	//"fmt"
	//"gopkg.in/couchbase/gocb.v1"
	//"github.com/xlsx"
	"encoding/csv"
	"os"
)

func main() {
	f, err := os.Create("haha2.xls")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f)
	w.Write([]string{"编号", "姓名", "年龄"})
	w.Write([]string{"1", "张三", "23"})
	w.Write([]string{"2", "李四", "24"})
	w.Write([]string{"3", "王五", "25"})
	w.Write([]string{"4", "赵六", "26"})
	w.Flush()

	/*var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("aa")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	fmt.Println(row)
	cell = row.AddCell()
	cell.Value = "adsasdasd"
	err = file.Save("aa.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}*/

/*excelFileName := "0827.xlsx"
xlFile, err := xlsx.OpenFile(excelFileName)
fmt.Println(xlFile)
if err != nil {
	fmt.Println("error!!")
}
for _, sheet := range xlFile.Sheets {
	for _, row := range sheet.Rows {
		for _, cell := range row.Cells {
			fmt.Println(cell)
			//fmt.Printf("%s\n", cell.String())
		}
	}
}*/
//}

/*func main() {
myCluster, _ := gocb.Connect("couchbase://120.92.5.17")
//fmt.Println(myCluster.ConnectTimeout())
myBucket, err := myCluster.OpenBucket("cp_cb", "caimikj")
fmt.Println(myBucket.DurabilityTimeout())
fmt.Println(myBucket.DurabilityPollTimeout())
myBucket.SetDurabilityPollTimeout(50000000000)
fmt.Println(myBucket.DurabilityPollTimeout())
fmt.Println(err)
fmt.Println("*************************************************")
var value interface{}
cas, err := myBucket.Get("nihao", &value)
fmt.Println(err)
fmt.Printf("Got value `%+v` with CAS `%08x`\n", value, cas)
/*value := "test value"
cas, err := myBucket.Insert("document_name", &value, 0)
fmt.Println(err)
fmt.Printf("Inserted document CAS is `%08x`\n", cas)*/
//var beer map[string]interface{}
//myBucket.Append("test", "niaho")
//cas, _ := myBucket.Get("aass_brewery-juleol", &beer)

//beer["comment"] = "Random beer from Norway"

//myBucket.Replace("aass_brewery-juleol", &beer, cas, 0)
//}

/*import (
	"fmt"
	"time"
	//"time"
	//"strings"
	//"time"
	//"strings"
	//"sync/atomic"
)

func main() {
	fmt.Println(time.Now().Unix())
}

/*func main() {
	var t_init int = 0
	var t_done int = 3
	var ops int64 = 0
	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(5 * time.Second)
			atomic.AddInt64(&ops, int64(1))
		}()
	}
	//time out
	go func() {
		for {
			t_init += 1
			time.Sleep(1 * time.Second)
			fmt.Println(t_init)
		}
	}()
	for {
		if t_done == t_init {
			fmt.Println("time out !!")
			break
		}
		time.Sleep(1 * time.Second)
		fmt.Println("sleep")
		if ops == int64(10) {
			fmt.Println("done")
			break
		}
	}
}*/

/*package main

import (
	"log"
	//"fmt"
	"github.com/coreos/go-etcd/etcd"
	//"time"
	"strings"
)

func main() {
	machines := []string{"http://123.59.164.180:2379"}
	client := etcd.NewClient(machines)

	//if _, err := client.Set("/foo", "bar", 0); err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(client.Delete("/foo", true))
	/*a, _ := client.Get("foo", false, false)
	log.Println(a.Node.Value)
	log.Println(client.Get("foo", false, false))*/
//client.Close()
//a:= make(chan *etcd.Response)
/*a := make(chan bool)
b := make(chan *etcd.Response)
//d := 0
//var c *etcd.Response
var p_value *etcd.Response
go func() {
	p_value, _ = client.Watch("/part_admin_info/10.136.9.73/process", 0, false, b, a)
}()
for {
	//if d == 1 {
	//	break
	//}
	time.Sleep(1 * time.Second)
	log.Println("sleep!!")
	select {
	case c := <-b:

		log.Println("recive!!")
		log.Println(c.Node.Value)
		if c.Node.Value == "ok" {
			a <- false
			log.Println(p_value)
		}
	}
	//log.Println(p_value.Node.Value)
}*/

/*	p_value, err := client.Watch("/part_admin_info/10.136.9.73/process", 0, false, nil, nil)
	if err == nil {
		at1 := strings.Split(p_value.Node.Value, "_")
		if len(at1) > 0 {
			if at1[0] == "ok" || at1[0] == "cok" {
				log.Println(p_value.Node.Value)
				//rt1, _ := client.Get(log_data+"tmp_"+nt1+".log", false, false)
				//c.Ctx.WriteString(rt1.Node.Value)
				return
			}

		}
	}

	//fmt.Println(time.Now().Format("2006-01-02T15:04:05"))
}*/
