# countbytes
Quietly pass stdin to stdout while outputting progress size to stderr - for use in bash pipelines

This package simply takes a streaming input from stdin and quietly 
passes it through to stdout while outputting the total bytes so far 
transferred as stderr. I wrote this to unobtrusively monitor 
transferal of large files via, e.g. rsync 
