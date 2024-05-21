# grapher

grapher depends on graphviz.

## build

```sh
go build -o mtree cmd/mtree/main.go

mvn dependency:tree > tree.out && mtree -file tree.out -filter "spring"
```