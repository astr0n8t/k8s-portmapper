# k8s-portmapper

Run these commands to get started:
```
find . -type f -exec sed -i 's/k8s-portmapper/actual_app_name/g' {} +
go mod init github.com/astr0n8t/actual_app_name
go get github.com/spf13/cobra
go get github.com/spf13/viper 
```

See the config.yaml file for an example configuration
