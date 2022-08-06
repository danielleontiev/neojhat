# neojhat

Command-line utility for analyzing JVM .hprof heap dumps. It's optimized for working
with huge dump files that do not fit into computer RAM.

## Install


### Compile and install with `go install`

1. Make sure you have `Go` installed. Refer to https://go.dev/dl/ for installation instructions.

2. Run the command if you would like to get the latest development version

   ```sh
   go install github.com/danielleontiev/neojhat@latest
   ```

   or run this command to obtain the latest stable release
   ```sh
   go install github.com/danielleontiev/neojhat@v0.2.0
   ```

3. Add `$USER/go/bin` to your `$PATH` or run directly from the directory

### Download pre-compiled release binary from Releases page

Go to [Releases](https://github.com/danielleontiev/neojhat/releases) and pick the binary from there.

### Obtaining pre-compiled development version binary

Go to [Actions](https://github.com/danielleontiev/neojhat/actions), select the
latest successful workflow and download `neojhat` from **Artifacts** section.

## Help

```sh
neojhat --help
```

```
neojhat v0.2.0
neojhat (threads|summary|objects)

Usage of threads:
  -hprof string
        path to .hprof file (required)
  -local-vars
        show local variables
  -no-color
        disable color output
  -non-interactive
        disable interactive output

Usage of summary:
  -all-props
        print all available properties from java.lang.System
  -hprof string
        path to .hprof file (required)
  -no-color
        disable color output
  -non-interactive
        disable interactive output

Usage of objects:
  -hprof string
        path to .hprof file (required)
  -no-color
        disable color output
  -non-interactive
        disable interactive output
  -sort-by value
        Sort output by 'size' or 'count' (default)

```

There are three sub-command: `threads`, `summary` and `objects`.

### `threads`

`threads` sub-command outputs thread dump at the time .hprof file was created.

```sh
neojhat threads --hprof /path/to/hprof/file
```

```java
"main", ID=1, prio=5, status=TIMED_WAITING
    void java.lang.Thread.sleep(long) Thread.java:NativeMethod
    void Main.main(java.lang.String[]) Main.java:6
    java.lang.Object jdk.internal.reflect.NativeMethodAccessorImpl.invoke0(java.lang.reflect.Method, java.lang.Object, java.lang.Object[]) NativeMethodAccessorImpl.java:NativeMethod
    java.lang.Object jdk.internal.reflect.NativeMethodAccessorImpl.invoke(java.lang.Object, java.lang.Object[]) NativeMethodAccessorImpl.java:62
    java.lang.Object jdk.internal.reflect.DelegatingMethodAccessorImpl.invoke(java.lang.Object, java.lang.Object[]) DelegatingMethodAccessorImpl.java:43
    java.lang.Object java.lang.reflect.Method.invoke(java.lang.Object, java.lang.Object[]) Method.java:566
    void com.sun.tools.javac.launcher.Main.execute(java.lang.String, java.lang.String[], com.sun.tools.javac.launcher.Main$Context) Main.java:404
    void com.sun.tools.javac.launcher.Main.run(java.lang.String[], java.lang.String[]) Main.java:179
    void com.sun.tools.javac.launcher.Main.main(java.lang.String[]) Main.java:119

"Reference Handler", ID=2, prio=10, status=RUNNABLE (daemon)
    void java.lang.ref.Reference.waitForReferencePendingList() Reference.java:NativeMethod
    void java.lang.ref.Reference.processPendingReferences() Reference.java:241
    void java.lang.ref.Reference$ReferenceHandler.run() Reference.java:213

// ... full output omitted ...
```

It is also possible to print GC Roots in each stack frame with `--local-vars` option.

```sh
neojhat threads --hprof /path/to/hprof/file --local-vars
```

```java
"main", ID=1, prio=5, status=TIMED_WAITING
    void java.lang.Thread.sleep(long) Thread.java:NativeMethod
    void Main.main(java.lang.String[]) Main.java:6
        local java.lang.String[]
        local java.lang.String
    java.lang.Object jdk.internal.reflect.NativeMethodAccessorImpl.invoke0(java.lang.reflect.Method, java.lang.Object, java.lang.Object[]) NativeMethodAccessorImpl.java:NativeMethod
    java.lang.Object jdk.internal.reflect.NativeMethodAccessorImpl.invoke(java.lang.Object, java.lang.Object[]) NativeMethodAccessorImpl.java:62
        local jdk.internal.reflect.NativeMethodAccessorImpl
        local java.lang.Integer
        local java.lang.Object[]
    java.lang.Object jdk.internal.reflect.DelegatingMethodAccessorImpl.invoke(java.lang.Object, java.lang.Object[]) DelegatingMethodAccessorImpl.java:43
        local jdk.internal.reflect.DelegatingMethodAccessorImpl
        local java.lang.Integer
        local java.lang.Object[]
    java.lang.Object java.lang.reflect.Method.invoke(java.lang.Object, java.lang.Object[]) Method.java:566
        local java.lang.reflect.Method
        local java.lang.Integer
        local java.lang.Object[]
        local jdk.internal.reflect.DelegatingMethodAccessorImpl
    void com.sun.tools.javac.launcher.Main.execute(java.lang.String, java.lang.String[], com.sun.tools.javac.launcher.Main$Context) Main.java:404
        local com.sun.tools.javac.launcher.Main
        local java.lang.String
        local java.lang.String[]
        local com.sun.tools.javac.launcher.Main$Context
        local com.sun.tools.javac.launcher.Main$MemoryClassLoader
        local class Main
        local java.lang.reflect.Method
    void com.sun.tools.javac.launcher.Main.run(java.lang.String[], java.lang.String[]) Main.java:179
        local com.sun.tools.javac.launcher.Main
        local java.lang.String[]
        local java.lang.String[]
        local sun.nio.fs.UnixPath
        local com.sun.tools.javac.launcher.Main$Context
        local java.lang.String
        local java.lang.String[]
    void com.sun.tools.javac.launcher.Main.main(java.lang.String[]) Main.java:119
        local java.lang.String[]

"Reference Handler", ID=2, prio=10, status=RUNNABLE (daemon)
    void java.lang.ref.Reference.waitForReferencePendingList() Reference.java:NativeMethod

// ... full output omitted ...
```

### `summary`

`summary` prints some remarkable information about the program and JVM.

```sh
neojhat summary --hprof /path/to/hprof/file
```

```
- Environment
Architecture:          x86_64
JavaHome:              /Library/Java/JavaVirtualMachines/temurin-11.jdk/Contents/Home
JavaName:              OpenJDK 64-Bit Server VM (11.0.12+7, mixed mode)
JavaVendor:            Eclipse Foundation
JavaVersion:           11.0.12
System:                Mac OS X

- Heap
Classes:               3563
GC Roots:              2058
Heap Size:             2M
Instances:             73224

- System
JVM Uptime:            45.813s
```

The output consists of three sections: **Environment**, **Heap** and **System**.

#### Environment

The first **Environment** section is selected fields from the content of system
map of various properties stored inside `java.lang.System` class.

```java
                                  // selected properties are:
java.lang.System.getProperties(); //
                                  // os.arch, os.name, java.home, java.version,
                                  // java.vm.name, java.vm.version, java.vm.info,
                                  // java.vm.vendor
```

The `--all-props` switch could be added to the command to print full content of `Properties`.
It will add **Properties** section with the list of all properties and their values.

```sh
neojhat summary --hprof /path/to/hprof/file --all-props
```

```
- Properties
awt.toolkit:                         sun.lwawt.macosx.LWCToolkit
file.encoding:                       UTF-8
file.separator:                      /
ftp.nonProxyHosts:                   local|*.local|169.254/16|*.169.254/16
gopherProxySet:                      false
// ... full output omitted ...
```

#### Heap

**Heap** is the section with summed up statistics of the heap dump:
*number of loaded classes*, *number of GC Roots*, *heap size (in memory)*
and *number of allocated instances* appeared in heap dump.

#### System

**System** section shows uptime of the JVM. The duration computed as a
difference between the time then the heap dump was taken and the value of

```java
sun.management.ManagementFactoryHelper.runtimeMBean.vmStartupTime
```

which should contain approximate time when JVM was started.

### `objects`

This command can print the table with the list of classes in the
heap dump that sorted either by total size or by the number of instances.
Sorting order is controlled with the `--sort-by` option.

```sh
neojhat objects --hprof /path/to/hprof/file --sort-by count
```

```java
Instances: 73224
Total Suze: 2M

Class Name                                       |              Count ↓ |                Size |
-----------------------------------------------------------------------------------------------
byte[]                                           |          16558 (22%) |          571K (22%) |
java.lang.String                                 |          16004 (21%) |           203K (7%) |
java.util.HashMap$Node                           |            4851 (6%) |           132K (5%) |
java.util.concurrent.ConcurrentHashMap$Node      |            3746 (5%) |           102K (4%) |
java.lang.Object[]                               |            3141 (4%) |          340K (13%) |
// ... full output omitted ...
```

<br>

```sh
neojhat objects --hprof /path/to/hprof/file --sort-by size
```

```java
Instances: 73224
Total Suze: 2M

Class Name                                       |                Count |              Size ↓ |
-----------------------------------------------------------------------------------------------
byte[]                                           |          16558 (22%) |          571K (22%) |
java.lang.Object[]                               |            3141 (4%) |          340K (13%) |
java.lang.String                                 |          16004 (21%) |           203K (7%) |
java.lang.reflect.Method                         |            1119 (1%) |           142K (5%) |
java.util.HashMap$Node                           |            4851 (6%) |           132K (5%) |
// ... full output omitted ...
```
