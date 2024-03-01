# Bracket

## Build and run
```bash
docker build . -t bracket:dev
docker run --expose 8080 -p 8080:8080 -v $(pwd)/data:/data -v $(pwd)/templates:/templates bracket:dev
```


