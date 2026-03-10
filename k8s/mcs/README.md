## Usage
```
# Create
./mcs --name my-svc --namespace default

# Get
./mcs --name my-svc --namespace default --action=get

# Delete
./mcs --name my-svc --namespace default --action=delete

# Update labels
./mcs --name my-svc --namespace default --action=patch --labels='{"env":"prod","team":"platform"}'
```
