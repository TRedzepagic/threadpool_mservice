# threadpool_mservice
Thread pooling and microservices practice project
## Cloning
You can "git clone" my repo with (Entire repository):

```
"git clone https://github.com/TRedzepagic/threadpool_mservice.git"
```
## To-Do
-   Configuration file?
## Info
This is a thread pool - like implementation in Go. We have a queue of tasks and workers/threads/goroutines vying for them (the threads wait for tasks).
Currently the program supports multiple workers/threads/goroutines working on a single queue of tasks. The program pings specified hosts, and if the hosts do not respond, it will send an e-mail. These tasks are separate and they all go into the aforementioned queue of tasks (FIFO).
