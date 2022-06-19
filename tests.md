```
recheckBaggageAfter = false
```

# Case 1
```
Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: false}]
Перелеты в Дельте: [{RecheckBaggage: false}, {RecheckBaggage: false}]
```

# Case 2
```
Перелеты Партнера: [{RecheckBaggage: true}, {RecheckBaggage: false}]
Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}]
```

# Case 3
```
Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
Перелеты в Дельте: [{RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
```

# Case 4
```
Перелеты Партнера: [{RecheckBaggage: true}, {RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
```

----------------------------------------------
```
recheckBaggageAfter = true
```

# Case 1
```
Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: false}]
Перелеты в Дельте: [{RecheckBaggage: false}, {RecheckBaggage: false}]
```

# Case 2
```
Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: true}]
Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}]
```

# Case 3
```
Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}, {RecheckBaggage: false}]
```

# Case 4
```
Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}, {RecheckBaggage: true}]
Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
```
