# Publish Subscriber System Simulator in Go

### Usage:
```
go run main.go  
```

### Description:
There are 3 publishers in the simulator:  
"sport" Publishes every 2 seconds  
"hacker" Publishes every 5 seconds  
"travell" Publishes every 7 seconds  

And there are 3 subscribers:  
Subscribers 1 subscribes "sport"  
Subscribers 2 subscribes "sport", "hacker"  
Subscribers 3 subscribes "sport", "hacker", "travell"  
