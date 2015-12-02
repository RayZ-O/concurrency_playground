# Lamport clock  
  
### Usage:  
go run lamport.go [num of clocks]  
Press Ctrl + C to stop the program  
  
### Description:  
1. Input the number of clocks, the n-th clock starts with timestamps = n  
2. Each clock will periodically send message to a random clock, when a clock receives a message, compare the timestamp in the message with its timestamp. if the timestamp in the message is later than the receiver's timestamp, the receiver update its timestamp.  
  
# Vector clock  
  
### Usage:  
go run vector.go [num of clocks]  
Press Ctrl + C to stop the program  
  
### Description:  
Input the number of clocks, all elements in the time vector of each clock are initialize to 0  
(1) Each clock will periodically increment its own logical clock in the vector by 1   
(2) Each clock will periodically send its entire vector to a random clock  
(3) Upon receive, the clock increments its own logical clock in the vector by 1 and updates each element in its vector by taking   
    max(value in its own vector, value in the vector in the message)  
