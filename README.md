User Microservice
=================

Use this service to do the following.

- Create a user with GitHub data.
- Manage user data, including GitHub data.

Hacking
-------

To get started with this service, you'll need to have some things set up beforehand.

Have these installed.

- Make
- Go
- Docker
- Fig

If you're on Mac, you can get these by running `brew install go boot2docker fig`;
boot2docker will also require VirtualBox (`brew cask install virtualbox`).
You'll also need to be able to cross-compile for 64-bit Linux: http://stackoverflow.com/questions/12168873/cross-compile-go-on-osx.

In order to make calls to GitHub, you'll need to obtain a GitHub client ID and secret (https://github.com/settings/applications).
With those in hand, make a Dockerfile in ./dist.

    FROM scratch
    ENV GITHUB_CLIENT_ID someid
    ENV GITHUB_CLIENT_SECRET somesecret
    ADD user-service user-service
    EXPOSE 8080
    ENTRYPOINT ["/user-service"]

Ensure your $GOPATH is set up and $GOPATH/bin is on your $PATH.
Then, you can run

    make install

to get all the Go packages you'll need to build and execute the service. Use

    make watch

to watch all Go source files and rebuild whenever one of those files changes.
Test it by running `curl`.

    curl <docker_host>:8080/users/ 2> /dev/null | python -m json.tool

Rationale
---------

This build process might seem cumbersome, but
surprisingly, the entire process is very quick.
Additionally, your development image can be the exact same image that is used for production.
Since the built Docker image is based on "scratch", the image is essentially the size of your service's binary.
