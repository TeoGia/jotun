# jotun
A simple JVM heap size monitoring tool

CirceCI latest build status: [![CircleCI](https://circleci.com/gh/TeoGia/jotun.svg?style=svg&circle-token=5ad98b6e111e7e48d491de0b56d5b3348f1e86b3)](https://circleci.com/gh/TeoGia/jotun)


## Introduction
During my short career as a DevOps engineer, I came across a JVM heap monitoring task and soon found out that there was no easy way of getting the heap usage of a running JAVA process without wrapping bash calls to jstat and using bc in order to just get a plain string output. No easy way to integrate into other DevOps scripts using a starndarized Data output format. Well, Jotun is made up to do just that. Get real time info on the heap usage of any JAVA process plus some overall RAM usage stats and then output it neatly in JSON format for easy integration to other custom tools.

## Installation
//TODO 
For the time being juton is only available for UNIX operating systems (Linux & MacOS)
//todo add installation guide

## Usage
>Monitor a single PID

To monitor the heap usage of a single JAVA use:

`./jotun -p PID_of_the_process_you_wish_to_monitor`

The output will be something like that:

`{"Pid":"4207","PidName":"javaApp.jar","Heap":"27846.90","Format":"kB","AvailableRAM":"16316300.00","FreeRAM":"12117944.00"}`

### Important!!
Keep in mind that the output's default unit if measurement is kB. (for fields: Heap, AvailableRAM and FreeRAM).
If you wish to display the output in MB, Bytes or GB, you can use the -h flag like this:

`./jotun -p PID_of_the_process_you_wish_to_monitor -h kB`

Currently only B, kB, MB and GB representations are supported. (kB is the default one).

The output is rounded up to 2 decimal points for any chosen representation.


>Monitor a PID list

To monitor a PID list of Java processes, you can use the --pid-list flag like this:

`./jotun --pidlist 123,456,789` Note that the separator is ,

The JSON output would be like this:

`{"PidLlist":[{"Pid":"456","PidName":"java3.jar","Heap":"27846.90","Format":"kB","AvailableRAM":"16316300.00","FreeRAM":"12106400.00"},{"Pid":"123","PidName":"java2.jar","Heap":"27846.90","Format":"kB","AvailableRAM":"16316300.00","FreeRAM":"12106400.00"},{"Pid":"789","PidName":"java3.jar","Heap":"27846.90","Format":"kB","AvailableRAM":"16316300.00","FreeRAM":"12106400.00"}]}`



>Monitor all JAVA running processes

To monitor all JAVA processes running on the machine, you will need the --all flag. No need to specify any pid, jotun will find everything itself.
You can use it like this:
`./jotun --all`

The output will have the same format and schema as the one above for the pid list.
Example:

`{"PidLlist":[{"Pid":"456","PidName":"java3.jar","Heap":"27846.90","Format":"kB","AvailableRAM":"16316300.00","FreeRAM":"12106400.00"},{"Pid":"123","PidName":"java2.jar","Heap":"27846.90","Format":"kB","AvailableRAM":"16316300.00","FreeRAM":"12106400.00"},{"Pid":"789","PidName":"java3.jar","Heap":"27846.90","Format":"kB","AvailableRAM":"16316300.00","FreeRAM":"12106400.00"}]}`


## Conclusion

Jotun is a free program under the MIT license and will remain that way forever. Please feel free to report any issues or development requests through github's ticketing system. Feel free to Fork and enjoy yourself as well.

Hope you'll find it useful
