# golang-moving-average
Moving average implementation for Go

## Usage 
```
import (
  "time"
  movingaverage "github.com/naeemkhan12/golang-moving-average"
)

ma := movingaverage.New(5) // 5 is the window size
ma.Add(movingaverage.Values{10,time.Now()})
ma.Add(movingaverage.Values{15,time.Now()})
ma.Add(movingaverage.Values{20,time.Now()})
ma.Add(movingaverage.Values{1,time.Now()})
ma.Add(movingaverage.Values{2,time.Now()})
ma.Add(movingaverage.Values{5,time.Now()}) // This one will effectively overwrite the first value (10 in this example)
avg := ma.Avg() 
```

## Partially used windows
In case you define a window of let's say 5 and only put in 2 values, the average will be based on those 2 values.

Window 5 - Values: 2, 2  - Average: 2 (not 0.8)
