<a name="unreleased"></a>
## [Unreleased]


<a name="1.1.1"></a>
## [1.1.1] - 2023-04-13
### Build
- **docker-compose:** add depends_on with conditions
- **docker-compose:** add kafka depends
- **docker-compose:** move envs to the file

### Docs
- **readme:** fix test config's path

### Feat
- **errors:** add postgres 23505 error handling

### Style
- remove legacy config


<a name="v1.1.0"></a>
## [v1.1.0] - 2023-04-13
### Build
- **task:** run unit tests before release

### Feat
- **companies:** send event to kafka

### Test
- **companies:** add uuid validation test


<a name="v1.0.0"></a>
## v1.0.0 - 2023-04-13
### Ci
- **github:** set go version
- **github:** install dependencies
- **github:** update action version
- **release:** fix tags
- **release:** build docker image
- **tests:** run only unit tests

### Docs
- update readme

### Feat
- **companies:** list is deprecated
- **fx:** shutdown migrations
- **rest:** use custom logger

### Style
- add gitignore
- remove idea
- remove idea

### Test
- add integration testing


[Unreleased]: /compare/1.1.1...HEAD
[1.1.1]: /compare/v1.1.0...1.1.1
[v1.1.0]: /compare/v1.0.0...v1.1.0
