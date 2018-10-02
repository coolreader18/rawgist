# RawGist

A simple http server to redirect to the correct cdn.rawgit.com url.

## Usage

The url is [https://rawgist.now.sh](https://rawgist.now.sh). The request format
is `/{gist_id}/{file}`.

### If running yourself

Build it, you can use `./build.sh` to build it if you don't have a go
environment set up. Put `GITHUB_TOKEN=xxxxx` in the `.env` file, and run it. It
defaults to port `3030`, but you can edit it in the code.

## License

This project is licensed under the MIT license. See the [LICENSE](LICENSE) file
for more details.
