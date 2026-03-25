package util

// 滚动窗口: 窗口大小固定, 每次滑动固定 step
func TumblingWindow[T any] (data []T, windowSize int) [][]T {
	var result [][]T
	n := len(data)
	for i := 0;i < n; i += windowSize {
		end := i + windowSize
		if end > n {
			end = n
		}
		result = append(result, data[i:end])
	}
	return result
}

// 滑动窗口
func SlidingWindow[T data] (data []T, windowSize int, step int) [][]T {
	var result [][]T
	n := len(data)
	for i := 0;i < n; i += step {
		end := i + windowSize
		if end > n {
			end = n
		}
		result = append(result, data[i:end])
	}
	return result
}

type Event struct {
	Timestamp int64
	Data string
}

// 会话窗口 Session Window
// 按【空闲间隔】切分
func SessionWindow(events []Event, gap int) [][]Event {
    var result [][]Event
	if len(events) == 0 {
		return result
	}
	currentWindow := []Event{events[0]}
	for i := 1;i < len(events); i++ {
		prev := events[i-1]
		curr := events[i]
		if curr.Timestamp -  prev.Timestamp > gap {
			result = append(result, currentWindow)
			currentWindow = []Event{curr}
		} else {
			currentWindow = append(currentWindow, curr)
		}
	}
	result = append(result, currentWindow)
	return result
}