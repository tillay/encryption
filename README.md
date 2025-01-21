# encryption
Source code for all of the end to end encryption projects I have been doing with some friends

for everything you need: `pycryptodome`

for image.py you need `pillow`

for image.py, pgpcli, and encrypt.py you need `xclip` or `wl-clipboard` (depending whether you are on wayland or X11)

minimal.py only requires `pycryptodome`.

# General Usage
to install:

```git clone https://github.com/tillay8/encryption&&cd encryption```

to run a python script:

`python3 <script.py> <optional flag>`

running a go script (pgpcli) is explained later

note: support of encrypt.py and image.py are only for Linux with X11 or Wayland due to clipboard utility needs

# Usage (encrypt.py)

The default location for stored password is /tmp/key, which is wiped on reboot because its in ram.

This can be changed my modifying the source code for encrypt.py

The python script automatically checks for whats on your clipboard:

- if it starts with `&&` it parses it as encrypted text

- if it starts with `@@` it makes it the new password

# Usage (image.py)

The python script automatically checks for what to do:

- if a file ~/Downloads/message.txt exists, it parses that as an encrypted image

- if the clipboard starts with `££` it parses it as an encrypted image

- if the clipboard is an png image, it encryptes the image (Other image types probably coming soonish maybe)

- if it starts with none of these, it does nothing

# Flags 

these work on image.py, encrypt.py, and utils.py

- `-p`: create new password, overwrite current saved password with it, and copy the new one to the clipboard. Add a number with the length of the password after. Example: ```python3 utils.py -p 420```

- `-n`: manually type in a new password

# Usage (minimal.py)

minimal.py is designed to work on a lot more OS'es due to no clipboard integration.

to use, type your message into the text box. if its unencrypted, it will encrypt it, and if it is encrypted, it will decrypt it. 

# Usage (pgpcli)

For sharing passwords between people securely use pgpcli.

pgpcli uses the PGP asymetric standard for encryption and provides a wrapper for PGP encryption

to compile:

make sure you are in the pgpcli subdirectory, then run ```go build -o pgpcli```

this will create a pgpcli binary, which can be "dot slashed" or run by ```./pgpcli```

The menu is pretty straightfoward. By default, keys are stored in ~/wpgp, but that can be changed

pgpcli also checks your clipboard to see if you have a public or private key on your clipboard, and will run the proper steps if so. 

# Usage (discord.py)

Discord.py is designed for seamless encryption integration with discord. 

put your discord token into a file named `token` in the directory `~/tilcord`. This is the default configuration file directory

if you dont know how to get your discord token, google it

make the file `~/tilcord/channel` to be the channel id of where you want to send the messages

to set it from inside the script, type `set <channel id>`

you also needs a users.csv script at `~/tilcord/users.csv`. This script should be formatted like so:
```
user1, red
user2, blue
user3, teal
```

to add a user from within the script, type `adduser username, green` into the send console. If it worked, it will give an output, otherwise try again. 

if you run the script with flag `-listen` it will print out all decrypted messages in real time that are sent in that channel

if you run the script with flag `-send` it will prompt for a message to send to the channel, and then send it encrypted

make sure you have already set a password with encrypt, image, or utils scripts 
