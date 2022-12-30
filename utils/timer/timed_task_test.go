package timer

import (
	"testing"
)

func TestNewTimerTask(t *testing.T) {
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
