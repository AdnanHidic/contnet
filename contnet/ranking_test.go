package contnet

import (
    "testing"
    "time"
    "log"
)

func TestRanking(t *testing.T) {
    // base time
    baseTime := time.Now()
    // test times will consists of times advanced +i up to 72 hours after now
    testTimes := []time.Time {}
    for i:=0;i<4;i++ {
        testTimes = append(testTimes, baseTime.Add(time.Duration(i)*time.Hour))
    }

    quality := 0.5
    popularity := 100.0

    content := Object.Content.New(1,nil,baseTime, quality, popularity)

    for i:=0;i<len(testTimes);i++ {
        age := __age(testTimes[i],*content)
        log.Print(testTimes[i].Sub(age), baseTime.Sub(age))
    }
}
