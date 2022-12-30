package timer

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	Logger = log.New(os.Stderr, "INFO -", 13)

	stime int = 0
)

type Tag struct {
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// 清空所有文章
func CleanAllArticle() error {
	if time.Now().UnixMilli()%7 == 0 {
		return fmt.Errorf("err times 7")
	}

	return nil
}

// 清空所有标签
func CleanAllTag() (bool, error) {
	if time.Now().UnixMilli()%3 == 0 {
		return false, fmt.Errorf("err times 3")
	}

	return true, nil
}
func DoCron() {
	log.Println("Starting...")

	//会根据本地时间创建一个新（空白）的 Cron job runner
	c := cron.New()

	//AddFunc 会向 Cron job runner 添加一个 func ，以按给定的时间表运行
	Id, err := c.AddFunc("* * * * * *", func() {
		Logger.Println("Run models.CleanAllTag...")
		_, _ = CleanAllTag()
	})
	Logger.Println("Add cleanAllTag", Id, err)
	Id, err = c.AddFunc("* * * * * *", func() {
		Logger.Println("Run models.CleanAllArticle...")
		_ = CleanAllArticle()
	})
	Logger.Println("Add CleanAllArticle", Id, err)

	//在当前执行的程序中启动 Cron 调度程序。其实这里的主体是 goroutine + for + select + timer 的调度控制哦
	c.Start()

	//会创建一个新的定时器，持续你设定的时间 d 后发送一个 channel 消息
	t1 := time.NewTimer(time.Second * 2)

	for true {
		//阻塞 select 等待 channel
		select {

		case <-t1.C:
			//会重置定时器，让它重新开始计时
			t1.Reset(time.Second * 10)
		}
	}
}

var Job = mockJob{}

type mockJob struct{}

func (job mockJob) Run() {
	runFunc()
}

func MockFunc() {
	time.Sleep(time.Second * time.Duration(5))
	fmt.Printf("mock %vs...\n", 5)
}
func PingFunc() {
	time.Sleep(time.Second * time.Duration(3))
	fmt.Printf("ping %vs...\n", 3)
}
func runFunc() {
	time.Sleep(time.Second * time.Duration(1))
	fmt.Printf("runFunc %vs...\n", 1)
}
func CronsTask() {

	tm := NewTimerTask()
	_tm := tm.(*Task)
	{
		Ids, err := tm.AddTaskByFunc("func", "@every 1s", MockFunc)
		Logger.Println(Ids, err)
		c, ok := _tm.TaskList["func"]
		if !ok {
			Logger.Println("Error no find func")
		} else {
			Logger.Println("cron find func", c, ok)

		}
	}

	{
		Ids, err := tm.AddTaskByFunc("ping", "@every 1s", PingFunc)
		Logger.Println(Ids, err)
		c, ok := _tm.TaskList["ping"]
		if !ok {
			Logger.Println("Error no find func")
		} else {
			Logger.Println("cron find func", c, ok)

		}
	}
	{
		Ids, err := tm.AddTaskByJob("job", "@every 1s", Job)
		Logger.Println(Ids, err)
		c, ok := _tm.TaskList["job"]
		if !ok {
			Logger.Println("Error no find job")
		} else {
			Logger.Println("cron find func", c, ok)

		}
	}

	{
		f, ok := tm.FindCron("func")
		if !ok {
			Logger.Println("Error no find func")
		} else {
			Logger.Printf("find f:%v\n", f)
		}
		f, ok = tm.FindCron("job")
		if !ok {
			Logger.Println("Error no find job")
		} else {
			Logger.Printf("find f:%v\n", f)
		}
		f, ok = tm.FindCron("none")
		if ok {
			Logger.Println("Error find none")
		} else {
			Logger.Printf("find f:%v\n", f)
		}
	}
	{
		tm.Clear("func")
		f, ok := tm.FindCron("func")
		if ok {
			Logger.Println("Error find func", f)
		}
	}

}

func CheckSelf() {
	CronsTask()
	//   select{}
	t1 := time.NewTimer(time.Second * 10)

	for {
		//阻塞 select 等待 channel
		select {

		case newTime := <-t1.C:
			//会重置定时器，让它重新开始计时
			t1.Reset(time.Second * 10)
			Logger.Printf("new time:%v\n", newTime)
		}
	}
}
