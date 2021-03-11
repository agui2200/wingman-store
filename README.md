# Wingman-store

fast use entgo make a store

# Install

```shell
go get entgo.io/ent
go get github.com/agui2200/wingman-store
wingman-store init User
wingman-store generate

```

# Use

example store.yml
```yaml
schemapackage: schema
targetpackage: store
feature:
  {
    privacy: false,
    entql: false,
    snapshot: false,
    schemaconfig: false
  }
```