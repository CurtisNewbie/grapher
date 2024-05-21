# grapher

grapher depends on graphviz.

## cmd/mtree

mtree supports parsing output of `mvn dependency:tree` into a graph to display.

```sh
go build -o mtree cmd/mtree/main.go

# write to local file then parse it
mvn dependency:tree > tree.out && mtree -file tree.out -filter "spring"

# or pipe output redirectly to mtree
mvn dependency:tree | mtree -file tree.out -filter "spring"
```