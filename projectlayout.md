# Go 工程项目文件布局  

## Go 目录 
## `/cmd`
项目主干 

    /cmd
        /myapp1
            /main.go
        /myapp2
            /main.go

## `/pkg`
这个目录，放置的是公开的库，可供别人导入使用 

## `/internal`
私有目录，不希望对外公开的库 


## 服务应用程序目录 
````
    /api
    /configs
    /test 
    /internal
    /cmd
````

https://blog.csdn.net/qq_41371349/article/details/106628372



## 通用应用目录

## 其它目录 

## 不该拥有的目录 

### `/src` 