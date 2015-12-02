# MapReduce simulator  
  
### Usage:  
```
go run mapper.go [input file name] | sort | go run reducer.go  
```
e.g.  
```
go run mapper.go word.txt | sort | go run reducer.go  
```  
  
### Description:  
Simple word count simulator implemented in Go  
(1) All heading and trailing non-alphanumeric characters are removed   
(2) All words are converted to lower case  
(3) Using linux command line tool sort to sort the mapper output  
