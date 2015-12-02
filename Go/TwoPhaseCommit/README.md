# Two phase commit simulator   
  
### Usage:  
```
go run commit.go [num of cohorts]  
```
Press Ctrl + C to stop the program  
  
### Description:  
Input the number of cohorts, coordinator will send commit request to all cohorts.  
The n-th cohort starts with state = n, the commit request is to add a number to the current state of a cohort.  
(1) The coordinator sends "query to commit" to every cohort.  
(2) Upon receive, each cohort add the number in "query to commit" to its current state. If the cohort finish successfully, it will send a "agreement" to the coordinator. However, the likelihood of failure of each cohort is 20%, if a cohort failed, it will send "abort" to coordinator.  
  
#### Success:  
If the coordinator received an agreement message from all cohorts:  
(1) Coordinator sends "commit" to all cohorts.  
(2) Each cohort sends an "acknowledgment" to the coordinator upon receive "commit".  
(3) The coordinator sleep for 2 seconds when all "acknowledgment" have been received.  
  
#### Failure:  
If any cohort send "abort" to the coordinator:  
(1) The coordinator sends "rollback" to all cohorts.  
(2) Each cohort undoes the addition using the undo log.  
(3) Each cohort sends "acknowledgement" to the coordinator.  
(4) The coordinator sleep for 2 seconds when all "acknowledgment" have been received.  
(5) Coordinator starts new transaction.  
