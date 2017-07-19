## Mockdemo  

This repo contains the source code for my blog post on [mocking the database in Go][blog]. 

> Blog post is WIP

Features 

- [x] Posts can be created.
- [x] Posts can be deleted.
- [x] A single post can be viewed.

The app uses an sqlite database but for the tests, the database interaction is mocked out... Run `cat db.sql | sqlite3 demo.sqlite` to get the sample database. 

### License
MIT

[blog]: https://lanreadelowo.com/blog/2017/07/06/mocking-go/
