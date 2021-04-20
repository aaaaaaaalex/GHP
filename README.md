# (G)o (h)ome(p)age
GHP is a FastCGI auto-indexer similar to the NGINX Auto-index
feature, but allows for fully customisable index pages via Go Templates

## Installation
This is a very dumb single-file program. Install it wherever you like :D
For use with Nginx, see the [sample nginx config](sample.nginx.conf)

## Usage
Just run `ghp`!

`ghp -h` for full list of config options


### Writing index files
Index files are written and rendered as [Go Templates](https://golang.org/pkg/text/template/).
The currently-served directory's contents are passed as data to the index with the type `[]os.FileInfo`, and can be referenced with `.`.
For more detail, see [an example](samples/index.gohtml)

### Additional util functions

 - `Request`: Get the current request (*http.Request)

 - `BaseURL( *http.Request )`: returns the current request's `URL.Path`, formatted for use with the HTML `<base>` tag

