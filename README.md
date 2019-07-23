shigoto
=======

[![GoDoc](http://www.godoc.org/github.com/DeedleFake/shigoto?status.svg)](http://www.godoc.org/github.com/DeedleFake/shigoto)
[![Go Report Card](https://goreportcard.com/badge/github.com/DeedleFake/shigoto)](https://goreportcard.com/report/github.com/DeedleFake/shigoto)

shigoto is a very simple static site generator. Hugo is great, but sometimes Hugo is just complete overkill. This is for those situations.

Features
--------

*Note: This section is currently a lie. Some of these features haven't actually been implemented yet. But this is the plan, at least.*

* No config file. None. All configuration is either figured out automatically from directory structure or is embedded directly into the actual content files.
* No tags. No categories. No archetypes. Just layout, drafts, and published pages. That's it.
* Automatic date insertion when publishing drafts.
* A [gh-pages][gh-pages]-style automatic deployment system.

[gh-pages]: https://www.npmjs.com/package/gh-pages
