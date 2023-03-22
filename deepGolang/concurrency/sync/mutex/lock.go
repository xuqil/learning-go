package mutex

type Lock struct {
	state int // 锁状态
}

func (l *Lock) Lock() {
	i := 0
	// 这个过程称为自旋，自旋 10 次就退出自旋
	for locked := CAS(UN_LOCK, LOCKED); !locked && i < 10 {
		i++
	}

	if locked {
		return
	}

    // 将自己的线程或协程加入阻塞队列，等待唤醒
	enqueue()
}
