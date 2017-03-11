# builder
通过向导来自动生成一个项目的初始代码

# 目录结构（初稿）

## 函数式
```
github.com/webx-top/{name}
│
├─ main.go
├─ application
│    ├─ init.go
│    ├─ handler
│    │     ├─ {group_1}
│    │     │     ├─ {handler_1}.go
│    │     │     ├─ {handler_2}.go
│    │     │     └─ {handler_N}.go
│    │     └─ {group_N}
│    │           ├─ {handler_1}.go
│    │           ├─ {handler_2}.go
│    │           └─ {handler_N}.go
│    ├─ dbschema
│    ├─ model
│    │     ├─ {model_1}.go
│    │     ├─ {model_2}.go
│    │     └─ {model_N}.go
│    ├─ libary
│    │     ├─ {package_1}
│    │     │     ├─ {file_1}.go
│    │     │     ├─ {file_2}.go
│    │     │     └─ {file_N}.go
│    │     └─ {package_N}
│    │           ├─ {file_1}.go
│    │           ├─ {file_2}.go
│    │           └─ {file_N}.go
│    └─ middleware
├─ data
│    ├─ config
│    │     └─ config.yaml
│    ├─ public
│    │     ├─ img
│    │     ├─ js
│    │     └─ css
│    ├─ theme
│    │     ├─ {theme_1}
│    │     │     ├─ {group_1}
│    │     │     │     ├─ {handler_1}.html
│    │     │     │     ├─ {handler_2}.html
│    │     │     │     └─ {handler_N}.html
│    │     │     └─ {group_N}
│    │     └─ {theme_N}
│    └─ upload
└─ tool
```


## HMVC模式
```
github.com/webx-top/{name}
│
├─ main.go
├─ application
│    ├─ init.go
│    ├─ dbschema
│    ├─ middleware
│    ├─ {module_1}
│    │     ├─ init.go
│    │     ├─ controller
│    │     │     ├─ {controller_1}.go
│    │     │     ├─ {controller_2}.go
│    │     │     └─ {controller_N}.go
│    │     ├─ model
│    │     │     ├─ {model_1}.go
│    │     │     ├─ {model_2}.go
│    │     │     └─ {model_N}.go
│    │     ├─ libary
│    │     │     ├─ {package_1}
│    │     │     │     ├─ {file_1}.go
│    │     │     │     ├─ {file_2}.go
│    │     │     │     └─ {file_N}.go
│    │     │     └─ {package_N}
│    │     │           ├─ {file_1}.go
│    │     │           ├─ {file_2}.go
│    │     │           └─ {file_N}.go
│    │     ├─ config
│    │     │     └─ config.yaml
│    │     └─ middleware
│    └─ {module_N}
│          ├─ init.go
│          ├─ controller
│          │     ├─ {controller_1}.go
│          │     ├─ {controller_2}.go
│          │     └─ {controller_N}.go
│          ├─ model
│          │     ├─ {model_1}.go
│          │     ├─ {model_2}.go
│          │     └─ {model_N}.go
│          ├─ libary
│          │     ├─ {package_1}
│          │     │     ├─ {file_1}.go
│          │     │     ├─ {file_2}.go
│          │     │     └─ {file_N}.go
│          │     └─ {package_N}
│          │           ├─ {file_1}.go
│          │           ├─ {file_2}.go
│          │           └─ {file_N}.go
│          ├─ config
│          │     └─ config.yaml
│          └─ middleware
├─ data
│    ├─ config
│    │     └─ config.yaml
│    ├─ public
│    │     ├─ img
│    │     ├─ js
│    │     └─ css
│    ├─ theme
│    │     ├─ {theme_1}
│    │     │     ├─ {module_1}
│    │     │     │     ├─ {controller_1}
│    │     │     │     │     ├─ {action_1}.html
│    │     │     │     │     ├─ {action_2}.html
│    │     │     │     │     └─ {action_N}.html
│    │     │     │     └─ {controller_N}
│    │     │     └─ {module_N}
│    │     └─ {theme_N}
│    └─ upload
└─ tool
```