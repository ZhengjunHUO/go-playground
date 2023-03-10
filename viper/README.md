```sh
$ go build
$ FOO_SECRET=foo_pwd BAR_SECRET=bar_pwd BAZ_SECRET=baz_pwd HUO_SECRET=huo_pwd ./viper --config config.yaml --env env.yaml -o delta1.yaml,delta2.yaml conf
```
