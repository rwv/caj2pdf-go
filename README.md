# caj2pdf-go

A Go port for [caj2pdf/caj2pdf](https://github.com/caj2pdf/caj2pdf)

## Purpose

* ⚡️ Better Performance
* 💻 Cross Compile
* 📦 Single Binary

## Current Status

* ✅ `CAJ` format
* ✅ `PDF` format
* ❌ `HN` format
* ❌ `KDH` format

## Requirements

* `mutool`: Fix PDF structure
* `fntsample`: Add outlines to PDF
## TODO

* Add outlines using pure Go to remove `fntsample` dependency
* Implement `JBIG` and `JBIG2` support (UniPDF has `JBIG2 `support but there's a [license issue](https://github.com/unidoc/unipdf/blob/master/LICENSE.md))