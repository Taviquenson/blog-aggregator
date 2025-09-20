# gator CLI blog-aggregator
An RSS feed aggregator written in Go to follow your favorite RSS
and have them refresh on a terminal window. It supports multiple
users so everyone can follow their own preferred feeds.

## Before using
To use it you will need to install Postgres to manage the program's
database. As well as the programming language Go.

## Installing the gator CLI using `go install`
You can install `gator` as follows:
1. Navigate to the program's root directory.
2. Rung `go install`
This command will compile your program and place the executable binary in your GOBIN directory. 
By default, GOBIN is set to $GOPATH/bin or $HOME/go/bin if GOPATH is not explicitly set.

If necessary, add GOBIN to your PATH:
`export PATH=$PATH:$(go env GOBIN)`
(Add this to your shell's profile file like .bashrc or .zshrc for persistence.)

## Setting up the config file
You'll need a `.gatorconfig.json` file placed in your home directory.
Its contents should be as follows:
```
{"db_url":"postgres://<computer_username>:@localhost:5432/gator?sslmode=disable","current_user_name":"<computer_username>"}
```

##  Running the program
Begin by registering a user. For example:
<br>`gator register tavo`</br>

Now you should add a feed. The following should work. If not, try
with a different RSS feed url:
<br>`gator addfeed blog.boot.dev https://blog.boot.dev/index.xml`</br>

With a feed added, remember, you still won't have the feed's posts
in the database. To get the posts you need to aggregate them. Run
this in a new terminal window:
<br>`gator agg 5s`</br>
This will keep a loop running every 5 seconds to pull all the posts from all the added feeds

Now you can see some info about the posts like title, description and publishing date by running:
<br>`gator browse`</br>
This will show all the posts for the feeds that the current user is following.

To see said feeds, run:
<br>`gator following`<br>
(The currently logged in user will automatically follow the feed they add)


Later, once you add more users you can switch between them, like so:
<br>`gator login mike`</br>

