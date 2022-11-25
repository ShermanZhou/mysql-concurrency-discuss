# Summary

Run two transactions of Mysql updates, to catch if there is a potential deadlock.

The original question is here
[StackExchange](https://dba.stackexchange.com/questions/320086/why-is-mysql-delete-on-foreign-key-index-locking-a-record-not-part-of-the-index/320128#320128)

# Build
- Make sure relatively new go compiler is installed.
- In the project path, run `go install` to install and build the binary.
- Change connection string to your develop database according to the Stack Exchange's schema.

# Run
run the binary with -p your-db-password
```
# under path of built executable e.g.
./msql-concurrent -p some-password
```
