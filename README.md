# dockit

_Orchestrating docker with go_

This is not supposed to be a general-purpose tool, but an example of how wiring up Docker environments using go is...

- Testable: a testing framework is built in
- Easy to read: idiomatic style that is almost universally followed making for a terse, highly readable ecosystem
- Easy to maintain: I'm assuming you were doing this work in bash before
- Portable: really, we're only talking about Linux with Docker, so a single statically linked binary is as portable as you can get

So get some ideas here (or use this as a starting point) and go build your own toolkit to codify (no pun intended) the conventions, opinions, 
and nuances of your environment in a format which you can maintain.

## Usage

By default, `dockit` looks for a config.json file in your current directory, connects using unix:///var/run/docker.sock, and keeps track
of running containers with "pid" files in /var/run/dockit-containers.

You can define services for your environment in the config file, and specify ports, environment variables, and dependancy services 
(which are translated into links.)

### Example Environment

Included in this repo, is an example config to build an environment with a webapp (the included `webapp-ex/`) fronted by a reverse proxy (Hipache.)


The config specifies: 

- the webapp ("WebApp" service) container should link in the hipache/redis ("Hipache" service) container 
- which ports to expose (in the link and externally)
- an environment variable for the webapp to use as a host name when registering with Hipache

The webapp entry point script (`webapp-ex/start.sh`) uses the link and host env var to register the webapp with Hipache and to deregister on shutdown.

Here is a copy of `config.json` from our example environment:

	{"Hipache" : {
		"Image" : "stackbrew/hipache",
		"Ports" : {
			"80" : "80",
			"6379" : ""
		}
		
	}, "WebApp" : {
		"Image" : "benschw/go-webapp",
		"Deps" : [
			"Hipache"
		],
		"Env" : {
			"HOST" : "webapp.local"
		}
	}}

(the image `benschw/go-webapp` was built from the contents of the webapp-ex directory)

#### Run the example

	sudo ./dockit -service WebApp -start

This will start up the `WebApp` service (bringing `Hipache` up too as it is a dependency) and register the private ip:port of 
the `WebApp` container with Hipache (see `webapp-ex/start.sh`) under the name `webapp.local`.

add `127.0.0.1  webapp.local` to your `/etc/hosts` file and the example webapp should be available on at [http://webapp.local](http://webapp.local)


#### Stop the example

	sudo ./dockit -service WebApp -stop

This will only stop the `WebApp` container (and deregister from Hipache); Hipache is still running. To stop it too, run:

	sudo ./dockit -service Hipache -stop

Note the containers are still there in a "stopped" state, and that a subsequent `-start` will run new instances.
