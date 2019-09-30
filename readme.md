## 这是一个加载配置文件的模块

### 提供方法：

- New (path,note_by_slash) *Config
> 说明：实例化配置对象

| 参数 | 类型 | 备注|
|:----|:----|:----|
|path | string | 配置文件路径
| note_by_slash | bool | 是否使用双斜杠作为备注

- SetDefaultGroup (group) *Config
> 说明：设置默认获取的组

| 参数 | 类型 | 备注|
|:----|:----|:----|
|group | string | 组名，不采用.方法取值时,默认的取值方法

- Load () bool 
> 说明：加载并解析文件

| 参数 | 类型 | 备注|
|:----|:----|:----|
|path | string | 配置文件路径

- Get(key) interface{}
> 说明：获取参数,返回值为string | nil | []string

| 参数 | 类型 | 备注|
|:----|:----|:----|
|str | string | key的名称，或着key + "." + group （用点拼接的名称）

- GetGroup(group) (map[string]interface{},bool)
> 说明：获取组参数，第一个参数为值，第二个参数为是否有该group,true为有，false为无

| 参数 | 类型 | 备注|
|:----|:----|:----|
|group | string | group的名称



