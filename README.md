# components-tests
Repo with tests for components

## Presequencies

* Golang > 1.17

1. Run `make install` to install required libraries
2. `cott` requires folder for logs files. Config log folder:
   1. create log directory `sudo mkdir /var/log/cott`
   2. set folder permissions `sudo chown $USER /var/log/cott`