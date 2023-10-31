# Building the Application for Different Platforms

To build the application for different platforms and architectures, you can use the following instructions. Make sure you have Go installed on your system.

### Build for Windows (64-bit)

```bash
GOOS=windows GOARCH=amd64 ./build
```

### Build for Windows (32-bit)

```bash
GOOS=windows GOARCH=386 ./build
```

### Build for Linux (64-bit)

```bash
GOOS=linux GOARCH=amd64 ./build
```

### Build for Linux (32-bit)

```bash
GOOS=linux GOARCH=386 ./build
```

### Build for Mac OS (64-bit)

```bash
GOOS=darwin GOARCH=amd64 ./build
```
