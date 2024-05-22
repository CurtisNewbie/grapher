# grapher

Simple lib to draw directed graph. grapher depends on graphviz.

## cmd/mtree

mtree supports parsing output of `mvn dependency:tree` into a graph to display.

```sh
Usage of mtree:
  -file string
        mvn dependency:tree output file
  -filter string
        filter tree branches by label name for tree-shaking
  -pom string
        maven pom file
```

For example,

```sh
go build -o mtree cmd/mtree/main.go

# write to local file then parse it
mvn dependency:tree > tree.out && mtree -file tree.out

# pipe output to mtree
mvn dependency:tree | mtree -file tree.out

# let mtree obtain output of dependency:tree directly
mtree -pom myproject
```